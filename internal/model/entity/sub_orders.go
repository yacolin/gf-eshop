// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubOrders is the golang structure for table sub_orders.
type SubOrders struct {
	Id             int64       `json:"id"              description:"子订单ID"`
	SubOrderNo     string      `json:"sub_order_no"    description:"子订单号（按商家拆单后的业务唯一键）"`
	ParentOrderId  int64       `json:"parent_order_id" description:"父订单ID"`
	ParentOrderNo  string      `json:"parent_order_no" description:"父订单号（冗余，便于查询）"`
	UserId         int64       `json:"user_id"         description:"用户ID"`
	MerchantId     int64       `json:"merchant_id"     description:"所属商家ID"`
	TotalAmount    int64       `json:"total_amount"    description:"商品总金额（分）"`
	DiscountAmount int64       `json:"discount_amount" description:"优惠金额（分）"`
	ShippingFee    int64       `json:"shipping_fee"    description:"运费（分）"`
	PayAmount      int64       `json:"pay_amount"      description:"子订单实付金额（分）"`
	Status         string      `json:"status"          description:"子订单状态：pending-待支付 paid-已支付 shipped-已发货 delivered-已签收 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款"`
	RefundStatus   string      `json:"refund_status"   description:"退款状态：none-无 partial_refunded-部分退款 refunded-已退款"`
	SellerRemark   string      `json:"seller_remark"   description:"卖家备注"`
	PaidAt         *gtime.Time `json:"paid_at"         description:"支付时间"`
	ShippedAt      *gtime.Time `json:"shipped_at"      description:"发货时间"`
	DeliveredAt    *gtime.Time `json:"delivered_at"    description:"签收时间"`
	CompletedAt    *gtime.Time `json:"completed_at"    description:"完成时间"`
	ClosedAt       *gtime.Time `json:"closed_at"       description:"关闭时间"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:""`
}
