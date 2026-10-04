package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"

	"gf-eshop/api/dashboard/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/service"
)

const (
	dashboardCacheKey        = "dashboard:stats"
	dashboardCacheTTL        = 20 * time.Minute
	dashboardRefreshInterval = 15 * time.Minute
)

type sDashboard struct {
	sf singleflight.Group
}

func init() {
	service.RegisterDashboard(&sDashboard{})
}

func (s *sDashboard) Stats(ctx context.Context, req *v1.DashboardStatsReq) (res *v1.DashboardStatsRes, err error) {
	v, err, _ := s.sf.Do(dashboardCacheKey, func() (interface{}, error) {
		// 尝试 Redis 缓存
		cached, err := g.Redis().Do(ctx, "GET", dashboardCacheKey)
		if err == nil && !cached.IsNil() {
			var resp v1.DashboardStatsRes
			if sonic.Unmarshal(cached.Bytes(), &resp) == nil {
				return &resp, nil
			}
		}

		// 缓存 miss，重建
		resp, err := s.rebuildStats(ctx)
		if err != nil {
			return nil, err
		}

		data, marshalErr := sonic.Marshal(resp)
		if marshalErr == nil {
			g.Redis().Do(ctx, "SETEX", dashboardCacheKey, int(dashboardCacheTTL.Seconds()), string(data))
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*v1.DashboardStatsRes), nil
}

func (s *sDashboard) StartPeriodicRefresh(ctx context.Context) {
	go func() {
		g.Log().Info(ctx, "Starting dashboard cache warmup...")
		if err := s.refreshCache(ctx); err != nil {
			g.Log().Warningf(ctx, "Dashboard cache warmup failed: %v", err)
		} else {
			g.Log().Info(ctx, "Dashboard cache warmup completed")
		}

		ticker := time.NewTicker(dashboardRefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.refreshCache(ctx); err != nil {
				g.Log().Warningf(ctx, "Dashboard cache refresh failed: %v", err)
			}
		}
	}()
}

func (s *sDashboard) refreshCache(ctx context.Context) error {
	resp, err := s.rebuildStats(ctx)
	if err != nil {
		return err
	}
	data, err := sonic.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", dashboardCacheKey, int(dashboardCacheTTL.Seconds()), string(data))
	return err
}

func (s *sDashboard) InvalidateCache(ctx context.Context) {
	_, err := g.Redis().Do(ctx, "DEL", dashboardCacheKey)
	if err != nil {
		g.Log().Warningf(ctx, "dashboard invalidate cache failed: %v", err)
	}
}

// rebuildStats 并行执行 7 个统计查询，出错时日志记录不阻断
func (s *sDashboard) rebuildStats(ctx context.Context) (*v1.DashboardStatsRes, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var summary v1.SummaryDTO
	eg.Go(func() error {
		var err error
		summary, err = s.computeSummary(egCtx)
		if err != nil {
			g.Log().Warningf(egCtx, "dashboard computeSummary failed: %v", err)
		}
		return nil // 不阻断，超时兜底
	})

	var orderTrend []v1.OrderTrendDTO
	eg.Go(func() error {
		orderTrend = s.computeOrderTrend(egCtx)
		return nil
	})

	var orderStatusDist []v1.StatusDistDTO
	eg.Go(func() error {
		orderStatusDist = s.computeOrderStatusDist(egCtx)
		return nil
	})

	var paymentMethodDist []v1.MethodDistDTO
	eg.Go(func() error {
		paymentMethodDist = s.computePaymentMethodDist(egCtx)
		return nil
	})

	var categoryDist []v1.CategoryDistDTO
	eg.Go(func() error {
		var err error
		categoryDist, err = s.computeCategoryDist(egCtx)
		if err != nil {
			g.Log().Warningf(egCtx, "dashboard computeCategoryDist failed: %v", err)
		}
		return nil
	})

	var inventoryStatusDist []v1.StatusDistDTO
	eg.Go(func() error {
		var err error
		inventoryStatusDist, err = s.computeInventoryStatusDist(egCtx)
		if err != nil {
			g.Log().Warningf(egCtx, "dashboard computeInventoryStatusDist failed: %v", err)
		}
		return nil
	})

	var topProducts []v1.TopProductDTO
	eg.Go(func() error {
		topProducts = s.computeTopProducts(egCtx)
		return nil
	})

	_ = eg.Wait() // 所有 goroutine 都返回 nil，不阻断
	return &v1.DashboardStatsRes{
		Summary:             summary,
		OrderTrend:          orderTrend,
		OrderStatusDist:     orderStatusDist,
		PaymentMethodDist:   paymentMethodDist,
		CategoryDist:        categoryDist,
		InventoryStatusDist: inventoryStatusDist,
		TopProducts:         topProducts,
	}, nil
}

// ── Compute Functions ──────────────────────────────────────────────────

func (s *sDashboard) computeSummary(ctx context.Context) (v1.SummaryDTO, error) {
	var summary v1.SummaryDTO

	// 商品总数
	n, err := dao.Products.Ctx(ctx).Count()
	if err != nil {
		return summary, fmt.Errorf("count products: %w", err)
	}
	summary.TotalProducts = int64(n)

	// 库存预警数量（status=2 缺货, status=3 无货）
	n, err = dao.Inventories.Ctx(ctx).
		WhereIn(dao.Inventories.Columns().Status, g.Slice{2, 3}).
		Count()
	if err != nil {
		return summary, fmt.Errorf("count low stock: %w", err)
	}
	summary.LowStockCount = int64(n)

	// 订单规模与营收：订单表只能经 service.Orders() 访问（见 docs/order-sharding-design.md），
	// 分表后这里会改走日汇总表，看板侧不需要跟着改。
	orderStats, err := service.Orders().StatsSummary(ctx)
	summary.TotalOrders = orderStats.TotalOrders
	summary.TotalRevenue = orderStats.TotalRevenue
	if err != nil {
		return summary, fmt.Errorf("order summary: %w", err)
	}

	return summary, nil
}

func (s *sDashboard) computeOrderTrend(ctx context.Context) []v1.OrderTrendDTO {
	sevenDaysAgo := gtime.Now().AddDate(0, 0, -6).Format("Y-m-d") + " 00:00:00"
	rows, err := service.Orders().DailyTrend(ctx, sevenDaysAgo)
	if err != nil {
		g.Log().Warningf(ctx, "query order trend: %v", err)
		return make([]v1.OrderTrendDTO, 0)
	}

	trendMap := make(map[string]v1.OrderTrendDTO, len(rows))
	for _, r := range rows {
		trendMap[r.Date] = v1.OrderTrendDTO{
			Date:   r.Date,
			Count:  r.Count,
			Amount: r.Amount,
		}
	}

	result := make([]v1.OrderTrendDTO, 0, 7)
	now := gtime.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("m-d")
		if t, ok := trendMap[date]; ok {
			result = append(result, t)
		} else {
			result = append(result, v1.OrderTrendDTO{Date: date, Count: 0, Amount: 0})
		}
	}
	return result
}

func (s *sDashboard) computeOrderStatusDist(ctx context.Context) []v1.StatusDistDTO {
	rows, err := service.Orders().StatusDistribution(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "query order status dist: %v", err)
		return make([]v1.StatusDistDTO, 0)
	}

	labelMap := map[string]string{
		"pending":   "待付款",
		"paid":      "已付款",
		"shipped":   "已发货",
		"delivered": "已送达",
		"cancelled": "已取消",
		"refunded":  "已退款",
	}
	result := make([]v1.StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = r.Status
		}
		result = append(result, v1.StatusDistDTO{
			Status: r.Status,
			Label:  label,
			Value:  r.Value,
		})
	}
	return result
}

func (s *sDashboard) computePaymentMethodDist(ctx context.Context) []v1.MethodDistDTO {
	type methodRow struct {
		Method string `orm:"method"`
		Value  int64  `orm:"value"`
	}
	var rows []methodRow
	err := g.DB().Model("tx_payments").
		Fields("payment_method AS method", "COUNT(*) AS value").
		WhereIn("status", g.Slice{"success", "paid"}).
		Where("deleted_at IS NULL").
		Group("method").
		Scan(&rows)
	if err != nil {
		g.Log().Warningf(ctx, "query payment method dist: %v", err)
		return make([]v1.MethodDistDTO, 0)
	}

	labelMap := map[string]string{
		"alipay": "支付宝",
		"wechat": "微信支付",
		"wallet": "余额",
		"bank":   "银行卡",
		"cash":   "现金",
	}
	result := make([]v1.MethodDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Method]
		if label == "" {
			label = r.Method
		}
		result = append(result, v1.MethodDistDTO{
			Method: r.Method,
			Label:  label,
			Value:  r.Value,
		})
	}
	return result
}

func (s *sDashboard) computeCategoryDist(ctx context.Context) ([]v1.CategoryDistDTO, error) {
	type row struct {
		Category string `orm:"category"`
		Value    int64  `orm:"value"`
	}
	var rows []row
	err := g.DB().Model("sp_products", "p").
		Fields("COALESCE(c.name, '未分类') AS category", "COUNT(p.id) AS value").
		LeftJoin("sp_categories c", "c.id = p.category_id").
		Where("p.deleted_at IS NULL").
		Group("c.id").
		Order("value DESC").
		Limit(8).
		Scan(&rows)
	if err != nil {
		return nil, fmt.Errorf("query category dist: %w", err)
	}
	result := make([]v1.CategoryDistDTO, len(rows))
	for i, r := range rows {
		result[i] = v1.CategoryDistDTO{Category: r.Category, Value: r.Value}
	}
	return result, nil
}

func (s *sDashboard) computeInventoryStatusDist(ctx context.Context) ([]v1.StatusDistDTO, error) {
	type row struct {
		Status int   `orm:"status"`
		Value  int64 `orm:"value"`
	}
	var rows []row
	err := g.DB().Model("sp_inventories").
		Fields("status", "COUNT(*) AS value").
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&rows)
	if err != nil {
		return nil, fmt.Errorf("query inventory status dist: %w", err)
	}
	labelMap := map[int]string{1: "库存充足", 2: "库存偏低", 3: "缺货"}
	result := make([]v1.StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = fmt.Sprintf("状态%d", r.Status)
		}
		result = append(result, v1.StatusDistDTO{
			Status: fmt.Sprintf("%d", r.Status),
			Label:  label,
			Value:  r.Value,
		})
	}
	return result, nil
}

func (s *sDashboard) computeTopProducts(ctx context.Context) []v1.TopProductDTO {
	rows, err := service.Orders().TopProducts(ctx, 10)
	if err != nil {
		g.Log().Warningf(ctx, "query top products: %v", err)
		return make([]v1.TopProductDTO, 0)
	}
	if len(rows) == 0 {
		return make([]v1.TopProductDTO, 0)
	}

	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.ProductId
	}
	type nameRow struct {
		Id   int64  `orm:"id"`
		Name string `orm:"name"`
	}
	var names []nameRow
	err = dao.Products.Ctx(ctx).
		Fields(dao.Products.Columns().Id, dao.Products.Columns().Name).
		WhereIn(dao.Products.Columns().Id, ids).
		Scan(&names)
	if err != nil {
		g.Log().Warningf(ctx, "query product names: %v", err)
	}
	nameMap := make(map[int64]string, len(names))
	for _, n := range names {
		nameMap[n.Id] = n.Name
	}

	result := make([]v1.TopProductDTO, len(rows))
	for i, r := range rows {
		name := nameMap[r.ProductId]
		if name == "" {
			name = fmt.Sprintf("已删除(ID:%d)", r.ProductId)
		}
		result[i] = v1.TopProductDTO{
			Name:   name,
			Count:  r.Count,
			Amount: r.Amount,
		}
	}
	return result
}
