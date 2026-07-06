package orders

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/api/orders/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

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

func generateOrderNo() string {
	return fmt.Sprintf("ORD%s%04d", gtime.Now().Format("YmdHis"), grand.Intn(10000))
}

func generateSubOrderNo() string {
	return fmt.Sprintf("SUB%s%04d", gtime.Now().Format("YmdHis"), grand.Intn(10000))
}

// resolveUserID 从 ctx 或 req 获取用户ID
func resolveUserID(ctx context.Context, userID int64) int64 {
	if userID > 0 {
		return userID
	}
	if claims := utility.GetStaffClaims(ctx); claims != nil {
		return claims.StaffId
	}
	return 0
}

func (s *sOrders) Create(ctx context.Context, req *v1.OrdersCreateReq) (res *v1.OrdersCreateRes, err error) {
	orderNo := generateOrderNo()
	userID := resolveUserID(ctx, req.UserId)

	var order *entity.Orders
	err = dao.Orders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("SKU不存在: %d", ri.SkuID))
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
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("SKU库存不足: %d", ri.SkuID))
			}
			availableQty := inv.Quantity - inv.Reserved
			if availableQty < int64(ri.Quantity) {
				return gerror.NewCode(gcode.CodeValidationFailed,
					fmt.Sprintf("SKU库存不足: %d (可用%d, 需要%d)", ri.SkuID, availableQty, ri.Quantity))
			}

			// 获取商品名称（从 SPU 表）
			var product entity.Products
			err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, sku.ProductId).Scan(&product)
			if err != nil {
				return err
			}
			productName := sku.Spec
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
				SkuSpec:     sku.Spec,
				Image:       sku.Image,
				Price:       sku.Price,
				Quantity:    ri.Quantity,
				Subtotal:    subtotal,
			})
		}

		// 写入主订单
		orderId, err := tx.Model("tx_orders").InsertAndGetId(g.Map{
			"order_no":       orderNo,
			"user_id":        userID,
			"total_amount":   totalAmount,
			"pay_amount":     payAmount,
			"shipping_fee":   0,
			"discount_amount": 0,
			"status":         "pending",
			"payment_status": "unpaid",
			"consignee":      req.Consignee,
			"phone":          req.Phone,
			"province":       req.Province,
			"city":           req.City,
			"district":       req.District,
			"detail_addr":    req.DetailAddr,
			"zip_code":       req.ZipCode,
			"coupon_id":      req.CouponID,
			"buyer_remark":   req.BuyerRemark,
			"source":         req.Source,
			"created_at":     gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 创建子订单
		subOrderNo := generateSubOrderNo()
		subOrderId, err := tx.Model("tx_sub_orders").InsertAndGetId(g.Map{
			"sub_order_no":    subOrderNo,
			"parent_order_id": orderId,
			"parent_order_no": orderNo,
			"user_id":         userID,
			"merchant_id":     0,
			"total_amount":    totalAmount,
			"discount_amount": 0,
			"pay_amount":      payAmount,
			"shipping_fee":  0,
			"status":        "pending",
			"created_at":    gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 写入订单项
		for _, item := range items {
			_, err = tx.Model("tx_order_items").Insert(g.Map{
				"sub_order_id": subOrderId,
				"order_id":     orderId,
				"order_no":     orderNo,
				"sub_order_no": subOrderNo,
				"merchant_id":  0,
				"sku_id":       item.SkuId,
				"product_id":   item.ProductId,
				"sku_code":     item.SkuCode,
				"product_name": item.ProductName,
				"sku_spec":     item.SkuSpec,
				"image":        item.Image,
				"price":        item.Price,
				"quantity":     item.Quantity,
				"subtotal":     item.Subtotal,
				"refund_status": "none",
				"created_at":   gtime.Now(),
			})
			if err != nil {
				return err
			}
		}

		// 创建订单日志
		_, err = tx.Model("tx_order_logs").Insert(g.Map{
			"order_id":      orderId,
			"order_no":      orderNo,
			"from_status":   "",
			"to_status":     "pending",
			"operator":      "system",
			"operator_type": "system",
			"note":          "订单创建",
			"created_at":    gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 重新读取完整订单
		err = tx.Model("tx_orders").Where("id", orderId).Scan(&order)
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

	m := dao.Orders.Ctx(ctx)
	if req.UserID > 0 {
		m = m.Where(dao.Orders.Columns().UserId, req.UserID)
	}
	if req.Status != "" {
		m = m.Where(dao.Orders.Columns().Status, req.Status)
	}
	if req.PaymentStatus != "" {
		m = m.Where(dao.Orders.Columns().PaymentStatus, req.PaymentStatus)
	}
	if req.OrderNo != "" {
		m = m.Where(dao.Orders.Columns().OrderNo, req.OrderNo)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.OrdersListRes{
			List:  make([]*entity.Orders, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Orders
	err = m.Page(page, size).OrderDesc(dao.Orders.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.OrdersListRes{
		List:  list,
		Total: int(total),
	}, nil
}

func (s *sOrders) Detail(ctx context.Context, req *v1.OrdersDetailReq) (res *v1.OrdersDetailRes, err error) {
	var order *entity.Orders
	err = dao.Orders.Ctx(ctx).Where(dao.Orders.Columns().OrderNo, req.OrderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "订单不存在")
	}

	var subOrders []*entity.SubOrders
	err = dao.SubOrders.Ctx(ctx).Where(dao.SubOrders.Columns().ParentOrderId, order.Id).Scan(&subOrders)
	if err != nil {
		return nil, err
	}
	if subOrders == nil {
		subOrders = make([]*entity.SubOrders, 0)
	}

	var items []*entity.OrderItems
	err = dao.OrderItems.Ctx(ctx).Where(dao.OrderItems.Columns().OrderId, order.Id).Scan(&items)
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
	// 查询订单
	var order *entity.Orders
	err = dao.Orders.Ctx(ctx).Where(dao.Orders.Columns().OrderNo, req.OrderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "订单不存在")
	}

	// 校验状态流转
	allowed, ok := validTransitions[order.Status]
	if !ok {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "当前状态不允许变更")
	}
	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, gerror.NewCode(gcode.CodeValidationFailed,
			fmt.Sprintf("状态不允许从 %s 变更为 %s", order.Status, req.Status))
	}

	// 在事务中更新状态
	err = dao.Orders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		updateData := g.Map{
			"status":     req.Status,
			"updated_at": gtime.Now(),
		}

		// 根据目标状态设置对应的时间戳
		switch req.Status {
		case "paid":
			updateData["payment_status"] = "paid"
			updateData["paid_at"] = gtime.Now()
		case "cancelled":
			updateData["closed_at"] = gtime.Now()
			updateData["payment_status"] = "refunded"
		case "shipped":
			updateData["shipped_at"] = gtime.Now()
		case "delivered":
			updateData["delivered_at"] = gtime.Now()
		case "completed":
			updateData["completed_at"] = gtime.Now()
		}
		_, err = tx.Model("tx_orders").Where("id", order.Id).Data(updateData).Update()
		if err != nil {
			return err
		}

		// 同步更新子订单状态
		subUpdateData := g.Map{
			"status":     req.Status,
			"updated_at": gtime.Now(),
		}
		switch req.Status {
		case "paid":
			subUpdateData["paid_at"] = gtime.Now()
		case "cancelled":
			subUpdateData["closed_at"] = gtime.Now()
		case "shipped":
			subUpdateData["shipped_at"] = gtime.Now()
		case "delivered":
			subUpdateData["delivered_at"] = gtime.Now()
		case "completed":
			subUpdateData["completed_at"] = gtime.Now()
		}
		_, err = tx.Model("tx_sub_orders").Where("parent_order_id", order.Id).Data(subUpdateData).Update()
		if err != nil {
			return err
		}

		// 创建订单日志
		_, err = tx.Model("tx_order_logs").Insert(g.Map{
			"order_id":      order.Id,
			"order_no":      order.OrderNo,
			"from_status":   order.Status,
			"to_status":     req.Status,
			"operator":      "system",
			"operator_type": "system",
			"note":          req.Note,
			"created_at":    gtime.Now(),
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	// 重新读取更新后的订单
	var updatedOrder *entity.Orders
	err = dao.Orders.Ctx(ctx).Where(dao.Orders.Columns().OrderNo, req.OrderNo).Scan(&updatedOrder)
	if err != nil {
		return nil, err
	}
	return &v1.OrdersUpdateStatusRes{Orders: updatedOrder}, nil
}