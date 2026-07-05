// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Orders is the golang structure of table tx_orders for DAO operations like Where/Data.
type Orders struct {
	g.Meta         `orm:"table:tx_orders, do:true"`
	Id             interface{} // 自增主键
	OrderNo        interface{} // 父订单号（业务唯一键，如 202612010001）
	UserId         interface{} // 用户ID
	TotalAmount    interface{} // 商品总金额（分）
	DiscountAmount interface{} // 优惠金额（分，含优惠券/满减）
	ShippingFee    interface{} // 运费（分）
	PayAmount      interface{} // 实付金额（分 = total - discount + shipping）
	Status         interface{} // 父订单状态：pending-待支付 paid-已支付 partial_shipped-部分发货 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款
	PaymentStatus  interface{} // 支付状态：unpaid-未支付 paying-支付中 paid-已支付 refunding-退款中 refunded-已退款
	PaymentMethod  interface{} // 支付方式：wechat-微信 alipay-支付宝 wallet-余额
	Consignee      interface{} // 收货人
	Phone          interface{} // 联系电话
	Province       interface{} // 省
	City           interface{} // 市
	District       interface{} // 区
	DetailAddr     interface{} // 详细地址
	ZipCode        interface{} // 邮编
	CouponId       interface{} // 使用的优惠券ID
	CouponSnapshot interface{} // 优惠券快照（名称/面值等，便于售后追溯）
	BuyerRemark    interface{} // 买家备注
	SellerRemark   interface{} // 卖家备注
	Source         interface{} // 订单来源：pc-电脑端 app-APP miniapp-小程序 h5-H5
	PaidAt         *gtime.Time // 支付时间
	ShippedAt      *gtime.Time // 发货时间
	DeliveredAt    *gtime.Time // 签收时间
	CompletedAt    *gtime.Time // 完成时间
	ClosedAt       *gtime.Time // 关闭时间
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
