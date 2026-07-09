// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItems is the golang structure of table tx_order_items for DAO operations like Where/Data.
type OrderItems struct {
	g.Meta         `orm:"table:tx_order_items, do:true"`
	Id             interface{} // 订单项ID
	OrderId        interface{} // 关联 tx_orders.id
	SubOrderId     interface{} // 子订单ID
	MerchantId     interface{} // 所属商家ID
	OrderNo        interface{} // 订单号（冗余，方便按订单号查）
	SubOrderNo     interface{} // 子订单号（冗余，方便按商家订单查）
	SkuId          interface{} // 关联 skus.id(sp_skus)
	ProductId      interface{} // 关联 products.id(sp_products)
	SkuCode        interface{} // 商家编码（冗余快照）
	ProductName    interface{} // 商品名（冗余快照）
	SkuSpecSummary interface{} // 规格摘要（如：红色 / 256G）
	SkuSpec        interface{} // 规格JSON快照
	Image          interface{} // 商品图（冗余快照）
	Price          interface{} // 单价（分，下单时价格）
	Quantity       interface{} // 购买数量
	Subtotal       interface{} // 小计（分 = price * quantity）
	RefundStatus   interface{} // 退款状态：none-无 refunding-退款中 refunded-已退款
	RefundAmount   interface{} // 已退款金额（分）
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
