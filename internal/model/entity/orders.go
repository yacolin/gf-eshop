// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Orders is the golang structure for table orders.
type Orders struct {
	Id             int64       `json:"id"              description:"自增主键"`
	OrderNo        string      `json:"order_no"        description:"父订单号（业务唯一键，如 202612010001）"`
	UserId         int64       `json:"user_id"         description:"用户ID"`
	TotalAmount    int64       `json:"total_amount"    description:"商品总金额（分）"`
	DiscountAmount int64       `json:"discount_amount" description:"优惠金额（分，含优惠券/满减）"`
	ShippingFee    int64       `json:"shipping_fee"    description:"运费（分）"`
	PayAmount      int64       `json:"pay_amount"      description:"实付金额（分 = total - discount + shipping）"`
	Status         string      `json:"status"          description:"父订单状态：pending-待支付 paid-已支付 partial_shipped-部分发货 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款"`
	PaymentStatus  string      `json:"payment_status"  description:"支付状态：unpaid-未支付 paying-支付中 paid-已支付 refunding-退款中 refunded-已退款"`
	PaymentMethod  string      `json:"payment_method"  description:"支付方式：wechat-微信 alipay-支付宝 wallet-余额"`
	Consignee      string      `json:"consignee"       description:"收货人"`
	Phone          string      `json:"phone"           description:"联系电话"`
	Province       string      `json:"province"        description:"省"`
	City           string      `json:"city"            description:"市"`
	District       string      `json:"district"        description:"区"`
	DetailAddr     string      `json:"detail_addr"     description:"详细地址"`
	ZipCode        string      `json:"zip_code"        description:"邮编"`
	CouponId       int64       `json:"coupon_id"       description:"使用的优惠券ID"`
	CouponSnapshot string      `json:"coupon_snapshot" description:"优惠券快照（名称/面值等，便于售后追溯）"`
	BuyerRemark    string      `json:"buyer_remark"    description:"买家备注"`
	SellerRemark   string      `json:"seller_remark"   description:"卖家备注"`
	Source         string      `json:"source"          description:"订单来源：pc-电脑端 app-APP miniapp-小程序 h5-H5"`
	PaidAt         *gtime.Time `json:"paid_at"         description:"支付时间"`
	ShippedAt      *gtime.Time `json:"shipped_at"      description:"发货时间"`
	DeliveredAt    *gtime.Time `json:"delivered_at"    description:"签收时间"`
	CompletedAt    *gtime.Time `json:"completed_at"    description:"完成时间"`
	ClosedAt       *gtime.Time `json:"closed_at"       description:"关闭时间"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:""`
}
