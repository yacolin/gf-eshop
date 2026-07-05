package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
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

func (s *sDashboard) rebuildStats(ctx context.Context) (*v1.DashboardStatsRes, error) {
	g, egCtx := errgroup.WithContext(ctx)

	var summary v1.SummaryDTO
	g.Go(func() error {
		summary = s.computeSummary(egCtx)
		return nil
	})

	var orderTrend []v1.OrderTrendDTO
	g.Go(func() error {
		orderTrend = s.computeOrderTrend(egCtx)
		return nil
	})

	var orderStatusDist []v1.StatusDistDTO
	g.Go(func() error {
		orderStatusDist = s.computeOrderStatusDist(egCtx)
		return nil
	})

	var paymentMethodDist []v1.MethodDistDTO
	g.Go(func() error {
		paymentMethodDist = s.computePaymentMethodDist(egCtx)
		return nil
	})

	var categoryDist []v1.CategoryDistDTO
	g.Go(func() error {
		categoryDist = s.computeCategoryDist(egCtx)
		return nil
	})

	var inventoryStatusDist []v1.StatusDistDTO
	g.Go(func() error {
		inventoryStatusDist = s.computeInventoryStatusDist(egCtx)
		return nil
	})

	var topProducts []v1.TopProductDTO
	g.Go(func() error {
		topProducts = s.computeTopProducts(egCtx)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

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

func (s *sDashboard) computeSummary(ctx context.Context) v1.SummaryDTO {
	var summary v1.SummaryDTO

	// 商品总数
	n, _ := dao.Products.Ctx(ctx).Count()
	summary.TotalProducts = int64(n)

	// 库存预警数量（status=2 缺货, status=3 无货）
	var lowStock int
	dao.Inventories.Ctx(ctx).WhereIn(dao.Inventories.Columns().Status, g.Slice{2, 3}).Count(&lowStock)
	summary.LowStockCount = int64(lowStock)

	// 订单/支付表不存在时返回 0
	return summary
}

func (s *sDashboard) computeOrderTrend(ctx context.Context) []v1.OrderTrendDTO {
	// 订单表不存在，返回空
	return make([]v1.OrderTrendDTO, 0)
}

func (s *sDashboard) computeOrderStatusDist(ctx context.Context) []v1.StatusDistDTO {
	return make([]v1.StatusDistDTO, 0)
}

func (s *sDashboard) computePaymentMethodDist(ctx context.Context) []v1.MethodDistDTO {
	return make([]v1.MethodDistDTO, 0)
}

func (s *sDashboard) computeCategoryDist(ctx context.Context) []v1.CategoryDistDTO {
	type row struct {
		Category string `orm:"category"`
		Value    int64  `orm:"value"`
	}
	var rows []row
	err := g.DB().Model("sp_products p").
		Fields("COALESCE(c.name, '未分类') AS category", "COUNT(p.id) AS value").
		LeftJoin("sp_categories c", "c.id = p.category_id").
		Where("p.deleted_at IS NULL").
		Group("c.id").
		Order("value DESC").
		Limit(8).
		Scan(&rows)
	if err != nil {
		return make([]v1.CategoryDistDTO, 0)
	}
	result := make([]v1.CategoryDistDTO, len(rows))
	for i, r := range rows {
		result[i] = v1.CategoryDistDTO{Category: r.Category, Value: r.Value}
	}
	return result
}

func (s *sDashboard) computeInventoryStatusDist(ctx context.Context) []v1.StatusDistDTO {
	// 查询 sp_inventories 按 status 分组
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
		return make([]v1.StatusDistDTO, 0)
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
	return result
}

func (s *sDashboard) computeTopProducts(ctx context.Context) []v1.TopProductDTO {
	// 无订单明细表，返回空
	return make([]v1.TopProductDTO, 0)
}
