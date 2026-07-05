// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubOrders is the golang structure of table tx_sub_orders for DAO operations like Where/Data.
type SubOrders struct {
	g.Meta         `orm:"table:tx_sub_orders, do:true"`
	Id             interface{} // 子订单ID
	SubOrderNo     interface{} // 子订单号（按商家拆单后的业务唯一键）
	ParentOrderId  interface{} // 父订单ID
	ParentOrderNo  interface{} // 父订单号（冗余，便于查询）
	UserId         interface{} // 用户ID
	MerchantId     interface{} // 所属商家ID
	TotalAmount    interface{} // 商品总金额（分）
	DiscountAmount interface{} // 优惠金额（分）
	ShippingFee    interface{} // 运费（分）
	PayAmount      interface{} // 子订单实付金额（分）
	Status         interface{} // 子订单状态：pending-待支付 paid-已支付 shipped-已发货 delivered-已签收 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款
	RefundStatus   interface{} // 退款状态：none-无 partial_refunded-部分退款 refunded-已退款
	SellerRemark   interface{} // 卖家备注
	PaidAt         *gtime.Time // 支付时间
	ShippedAt      *gtime.Time // 发货时间
	DeliveredAt    *gtime.Time // 签收时间
	CompletedAt    *gtime.Time // 完成时间
	ClosedAt       *gtime.Time // 关闭时间
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
