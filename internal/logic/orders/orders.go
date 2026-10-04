package orders

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/api/orders/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

// 本文件只做业务编排：所有订单表（tx_orders / tx_sub_orders / tx_order_items / tx_order_logs）
// 的 SQL 都在同目录的 repo_*.go 里，分片路由在同目录的 shard.go 里。
// 商品 / SKU / 库存不属于订单域，仍直接走各自 DAO。

type sOrders struct{}

func init() {
	service.RegisterOrders(&sOrders{})
}

// validTransitions 定义订单状态流转规则
var validTransitions = map[string][]string{
	"pending":   {"paid", "cancelled"},
	"paid":      {"shipped", "cancelled"},
	"shipped":   {"delivered"},
	"delivered": {"completed"},
}

// generateOrderNo 生成父订单号：ORD + 创建时间 + 随机后缀。
// 时间由调用方以「单一时间基准」传入，保证单号与 created_at 落在同一个月（分片一致）。
func generateOrderNo(now *gtime.Time) string {
	return fmt.Sprintf("ORD%s%04d", now.Format("YmdHis"), grand.Intn(10000))
}

// generateSubOrderNo 生成子订单号：SUB + 创建时间 + 随机后缀，时间基准同父订单。
func generateSubOrderNo(now *gtime.Time) string {
	return fmt.Sprintf("SUB%s%04d", now.Format("YmdHis"), grand.Intn(10000))
}

// resolveUserID 从 ctx 或 req 获取用户ID
func resolveUserID(ctx context.Context, userID int64) int64 {
	if userID > 0 {
		return userID
	}
	if claims := utility.GetUserClaims(ctx); claims != nil {
		return claims.UserId
	}
	return 0
}

func (s *sOrders) Create(ctx context.Context, req *v1.OrdersCreateReq) (res *v1.OrdersCreateRes, err error) {
	// 单一时间基准：这一个 now 同时派生单号、created_at 与分片。
	// 分表后单号/created_at 必须落在同一个月，否则读路径与写路径会命中不同分片。
	now := gtime.Now()
	orderNo := generateOrderNo(now)
	subOrderNo := generateSubOrderNo(now)
	sh := shardFromCreatedAt(ctx, now)
	userID := resolveUserID(ctx, req.UserId)

	var order *entity.Orders
	err = withOrdersTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		var totalAmount, payAmount int64
		type itemData struct {
			SkuId       int64
			ProductId   int64
			SkuCode     string
			ProductName string
			SkuSpec     string
			Image       string
			Price       int64
			Quantity    int
			Subtotal    int64
		}
		var items []itemData
		for _, ri := range req.Items {
			var sku entity.Skus
			err := dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, ri.SkuID).Scan(&sku)
			if err != nil {
				return err
			}
			if sku.Id == 0 {
				return errcode.Newf(errcode.CodeSKUNotFound, "SKU不存在: %d", ri.SkuID)
			}

			// 锁定库存：检查可用库存 >= 需求数量
			var inv entity.Inventories
			err = dao.Inventories.Ctx(ctx).
				Where(dao.Inventories.Columns().SkuId, ri.SkuID).
				Where(dao.Inventories.Columns().WarehouseId, 0).
				Scan(&inv)
			if err != nil {
				return err
			}
			if inv.Id == 0 {
				return errcode.Newf(errcode.CodeInsufficientStock, "SKU库存不足: %d", ri.SkuID)
			}
			availableQty := inv.Quantity - inv.Reserved
			if availableQty < int64(ri.Quantity) {
				return errcode.Newf(errcode.CodeInsufficientStock, "SKU库存不足: %d (可用%d, 需要%d)", ri.SkuID, availableQty, ri.Quantity)
			}

			// 获取商品名称（从 SPU 表）
			var product entity.Products
			err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, sku.ProductId).Scan(&product)
			if err != nil {
				return err
			}
			productName := sku.SpecSummary
			if product.Id > 0 && product.Name != "" {
				productName = product.Name
			}

			subtotal := sku.Price * int64(ri.Quantity)
			totalAmount += subtotal
			payAmount += subtotal

			// 扣减预留库存
			_, err = tx.Model("sp_inventories").
				Where("id", inv.Id).
				Where("quantity - reserved >= ?", ri.Quantity).
				Increment("reserved", ri.Quantity)
			if err != nil {
				return err
			}

			items = append(items, itemData{
				SkuId:       sku.Id,
				ProductId:   sku.ProductId,
				SkuCode:     sku.SkuCode,
				ProductName: productName,
				SkuSpec:     sku.SpecSummary,
				Image:       sku.Image,
				Price:       sku.Price,
				Quantity:    ri.Quantity,
				Subtotal:    subtotal,
			})
		}

		// 写入主订单
		orderId, err := insertOrder(ctx, sh, g.Map{
			"order_no":        orderNo,
			"user_id":         userID,
			"total_amount":    totalAmount,
			"pay_amount":      payAmount,
			"shipping_fee":    0,
			"discount_amount": 0,
			"status":          "pending",
			"payment_status":  "unpaid",
			"consignee":       req.Consignee,
			"phone":           req.Phone,
			"province":        req.Province,
			"city":            req.City,
			"district":        req.District,
			"detail_addr":     req.DetailAddr,
			"zip_code":        req.ZipCode,
			"coupon_id":       req.CouponID,
			"buyer_remark":    req.BuyerRemark,
			"source":          req.Source,
			"created_at":      now,
		})
		if err != nil {
			return err
		}

		// 创建子订单
		subOrderId, err := insertSubOrder(ctx, sh, g.Map{
			"sub_order_no":    subOrderNo,
			"parent_order_id": orderId,
			"parent_order_no": orderNo,
			"user_id":         userID,
			"merchant_id":     0,
			"total_amount":    totalAmount,
			"discount_amount": 0,
			"pay_amount":      payAmount,
			"shipping_fee":    0,
			"status":          "pending",
			"created_at":      now,
		})
		if err != nil {
			return err
		}

		// 写入订单项
		for _, item := range items {
			err = insertOrderItem(ctx, sh, g.Map{
				"sub_order_id":  subOrderId,
				"order_id":      orderId,
				"order_no":      orderNo,
				"sub_order_no":  subOrderNo,
				"merchant_id":   0,
				"sku_id":        item.SkuId,
				"product_id":    item.ProductId,
				"sku_code":      item.SkuCode,
				"product_name":  item.ProductName,
				"sku_spec":      item.SkuSpec,
				"image":         item.Image,
				"price":         item.Price,
				"quantity":      item.Quantity,
				"subtotal":      item.Subtotal,
				"refund_status": "none",
				"created_at":    now,
			})
			if err != nil {
				return err
			}
		}

		// 创建订单日志
		err = insertOrderLog(ctx, sh, g.Map{
			"order_id":      orderId,
			"order_no":      orderNo,
			"from_status":   "",
			"to_status":     "pending",
			"operator":      "system",
			"operator_type": "system",
			"note":          "订单创建",
			"created_at":    now,
		})
		if err != nil {
			return err
		}

		// 重新读取完整订单（分片已知，无需二次推导）
		order, err = findOrderByIDIn(ctx, sh, orderId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.OrdersCreateRes{Orders: order}, nil
}

func (s *sOrders) List(ctx context.Context, req *v1.OrdersListReq) (res *v1.OrdersListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	filter := orderListFilter{
		UserID:        req.UserID,
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		OrderNo:       req.OrderNo,
	}

	total, err := countOrdersByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.OrdersListRes{
			List:  make([]*entity.Orders, 0),
			Total: 0,
		}, nil
	}

	list, err := pageOrders(ctx, filter, page, size)
	if err != nil {
		return nil, err
	}
	return &v1.OrdersListRes{
		List:  list,
		Total: int(total),
	}, nil
}

func (s *sOrders) Detail(ctx context.Context, req *v1.OrdersDetailReq) (res *v1.OrdersDetailRes, err error) {
	// 由单号定位分片，子表查询复用同一个分片
	order, sh, err := findOrderWithShardByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errcode.ErrOrderNotFound
	}

	subOrders, err := listSubOrdersByParentOrderID(ctx, sh, order.Id)
	if err != nil {
		return nil, err
	}
	if subOrders == nil {
		subOrders = make([]*entity.SubOrders, 0)
	}

	items, err := listItemsByOrderID(ctx, sh, order.Id)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]*entity.OrderItems, 0)
	}

	// 回填 product_name（兼容旧数据将 sku_spec 误存为 product_name）
	var productIDs []int64
	for _, item := range items {
		if item.ProductId > 0 {
			productIDs = append(productIDs, item.ProductId)
		}
	}
	if len(productIDs) > 0 {
		var products []*entity.Products
		_ = dao.Products.Ctx(ctx).WhereIn(dao.Products.Columns().Id, productIDs).Scan(&products)
		nameByID := make(map[int64]string, len(products))
		for _, p := range products {
			if p.Name != "" {
				nameByID[p.Id] = p.Name
			}
		}
		for _, item := range items {
			if name, ok := nameByID[item.ProductId]; ok {
				item.ProductName = name
			}
		}
	}

	return &v1.OrdersDetailRes{
		Order:     order,
		SubOrders: subOrders,
		Items:     items,
	}, nil
}

func (s *sOrders) UpdateStatus(ctx context.Context, req *v1.OrdersUpdateStatusReq) (res *v1.OrdersUpdateStatusRes, err error) {
	// 查询订单（同时拿到分片，后续子表与更新都在同一分片内）
	order, sh, err := findOrderWithShardByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errcode.ErrOrderNotFound
	}

	// 校验状态流转
	allowed, ok := validTransitions[order.Status]
	if !ok {
		return nil, errcode.ErrInvalidOrderStatus
	}
	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errcode.Newf(errcode.CodeInvalidOrderStatus, "状态不允许从 %s 变更为 %s", order.Status, req.Status)
	}

	// 在事务中更新状态
	err = withOrdersTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := gtime.Now()
		updateData := g.Map{
			"status":     req.Status,
			"updated_at": now,
		}

		// 根据目标状态设置对应的时间戳
		switch req.Status {
		case "paid":
			updateData["payment_status"] = "paid"
			updateData["paid_at"] = now
		case "cancelled":
			updateData["closed_at"] = now
			updateData["payment_status"] = "refunded"
		case "shipped":
			updateData["shipped_at"] = now
		case "delivered":
			updateData["delivered_at"] = now
		case "completed":
			updateData["completed_at"] = now
		}
		err := updateOrderByID(ctx, sh, order.Id, updateData)
		if err != nil {
			return err
		}

		// 同步更新子订单状态
		subUpdateData := g.Map{
			"status":     req.Status,
			"updated_at": now,
		}
		switch req.Status {
		case "paid":
			subUpdateData["paid_at"] = now
		case "cancelled":
			subUpdateData["closed_at"] = now
		case "shipped":
			subUpdateData["shipped_at"] = now
		case "delivered":
			subUpdateData["delivered_at"] = now
		case "completed":
			subUpdateData["completed_at"] = now
		}
		err = updateSubOrdersByParentOrderID(ctx, sh, order.Id, subUpdateData)
		if err != nil {
			return err
		}

		// 创建订单日志
		return insertOrderLog(ctx, sh, g.Map{
			"order_id":      order.Id,
			"order_no":      order.OrderNo,
			"from_status":   order.Status,
			"to_status":     req.Status,
			"operator":      "system",
			"operator_type": "system",
			"note":          req.Note,
			"created_at":    now,
		})
	})
	if err != nil {
		return nil, err
	}

	// 重新读取更新后的订单（分片已知）
	updatedOrder, err := findOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	return &v1.OrdersUpdateStatusRes{Orders: updatedOrder}, nil
}

// ── 跨模块能力（供 payments / dashboard 使用） ─────────────────────────────
//
// 这些方法是订单表对其他模块**唯一**的访问入口：payments 不再直接 dao.Orders / tx_orders，
// dashboard 不再直接聚合订单表。分表落地时，跨片逻辑只在这一个模块里改。

// GetByOrderNo 按业务单号查询订单；不存在时返回 (nil, nil)。
func (s *sOrders) GetByOrderNo(ctx context.Context, orderNo string) (*entity.Orders, error) {
	return findOrderByOrderNo(ctx, orderNo)
}

// MarkPaidByOrderNo 支付成功后回写主订单与子订单状态（按业务单号，覆盖该单号下全部子订单）。
// 由调用方在事务中调用时，会自动加入该事务（GoFrame 从事务 ctx 中取 tx）。
func (s *sOrders) MarkPaidByOrderNo(ctx context.Context, orderNo string) error {
	return markOrderPaidByOrderNo(ctx, orderNo)
}

// StatsSummary 订单总数与已支付金额合计。
func (s *sOrders) StatsSummary(ctx context.Context) (service.OrderStatsSummary, error) {
	var out service.OrderStatsSummary

	total, err := countAllOrders(ctx)
	if err != nil {
		return out, fmt.Errorf("count orders: %w", err)
	}
	out.TotalOrders = total

	revenue, err := sumPaidAmount(ctx)
	if err != nil {
		return out, fmt.Errorf("sum revenue: %w", err)
	}
	out.TotalRevenue = revenue
	return out, nil
}

// DailyTrend 统计 since（含）之后的按日订单数与金额，date 形如 08-20。
func (s *sOrders) DailyTrend(ctx context.Context, since string) ([]service.OrderDayTrend, error) {
	rows, err := dailyOrderTrend(ctx, since)
	if err != nil {
		return nil, err
	}
	out := make([]service.OrderDayTrend, len(rows))
	for i, r := range rows {
		out[i] = service.OrderDayTrend{Date: r.Date, Count: r.Count, Amount: r.Amount}
	}
	return out, nil
}

// StatusDistribution 按状态统计订单数。
func (s *sOrders) StatusDistribution(ctx context.Context) ([]service.OrderStatusCount, error) {
	rows, err := orderStatusDistribution(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]service.OrderStatusCount, len(rows))
	for i, r := range rows {
		out[i] = service.OrderStatusCount{Status: r.Status, Value: r.Value}
	}
	return out, nil
}

// TopProducts 按订单明细聚合商品销量与金额 TOP N。
func (s *sOrders) TopProducts(ctx context.Context, limit int) ([]service.OrderProductRank, error) {
	rows, err := topProductsByOrderItems(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]service.OrderProductRank, len(rows))
	for i, r := range rows {
		out[i] = service.OrderProductRank{ProductId: r.ProductId, Count: r.Count, Amount: r.Amount}
	}
	return out, nil
}
