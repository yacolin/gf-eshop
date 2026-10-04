package service

import (
	"context"

	"gf-eshop/api/orders/v1"
	"gf-eshop/internal/model/entity"
)

type IOrders interface {
	Create(ctx context.Context, req *v1.OrdersCreateReq) (res *v1.OrdersCreateRes, err error)
	List(ctx context.Context, req *v1.OrdersListReq) (res *v1.OrdersListRes, err error)
	Detail(ctx context.Context, req *v1.OrdersDetailReq) (res *v1.OrdersDetailRes, err error)
	UpdateStatus(ctx context.Context, req *v1.OrdersUpdateStatusReq) (res *v1.OrdersUpdateStatusRes, err error)

	// ── 跨模块能力 ────────────────────────────────────────────────────────
	//
	// 订单表对其他模块只有这一个访问出口：payments 与 dashboard 通过下列方法读写订单，
	// 不再直接引用 dao.Orders / tx_orders 等表名。
	// 这样做的直接收益是订单分表（见 docs/order-sharding-design.md）落地时，
	// 跨分片逻辑只需要改 orders 模块一处。

	// GetByOrderNo 按业务单号查询订单；不存在时返回 (nil, nil)。
	GetByOrderNo(ctx context.Context, orderNo string) (*entity.Orders, error)

	// MarkPaidByOrderNo 支付成功后回写主订单与子订单状态。
	// 覆盖该订单号下的全部子订单；在调用方的事务中执行时自动加入该事务。
	// orderID 供分表双写镜像 tx_order_logs 使用（该表按 order_id 定位）。
	MarkPaidByOrderNo(ctx context.Context, orderNo string, orderID int64) error

	// StatsSummary 订单总数与已支付金额合计。
	StatsSummary(ctx context.Context) (OrderStatsSummary, error)

	// DailyTrend 统计 since（含）之后的按日订单数与金额。
	DailyTrend(ctx context.Context, since string) ([]OrderDayTrend, error)

	// StatusDistribution 按状态统计订单数。
	StatusDistribution(ctx context.Context) ([]OrderStatusCount, error)

	// TopProducts 按订单明细聚合商品销量与金额 TOP N。
	TopProducts(ctx context.Context, limit int) ([]OrderProductRank, error)
}

// OrderStatsSummary 订单规模与营收概览。
type OrderStatsSummary struct {
	TotalOrders  int64 // 订单总数
	TotalRevenue int64 // 已支付订单实付金额合计（分）
}

// OrderDayTrend 单日订单量与金额；Date 形如 08-20。
type OrderDayTrend struct {
	Date   string
	Count  int64
	Amount int64
}

// OrderStatusCount 单个状态下的订单数。
type OrderStatusCount struct {
	Status string
	Value  int64
}

// OrderProductRank 商品维度的销量与金额。
type OrderProductRank struct {
	ProductId int64
	Count     int64
	Amount    int64
}

var localOrders IOrders

func Orders() IOrders {
	if localOrders == nil {
		panic("implement not found for interface IOrders, forgot register?")
	}
	return localOrders
}

func RegisterOrders(i IOrders) {
	localOrders = i
}
