// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSales is the golang structure for table after_sales.
type AfterSales struct {
	Id                 int64       `json:"id"                   description:"主键"`
	AfterSaleNo        string      `json:"after_sale_no"        description:"售后单号"`
	OrderId            int64       `json:"order_id"             description:"订单ID"`
	OrderItemId        int64       `json:"order_item_id"        description:"订单明细ID"`
	MerchantId         int64       `json:"merchant_id"          description:"商家ID"`
	UserId             int64       `json:"user_id"              description:"用户ID"`
	AfterSaleType      int         `json:"after_sale_type"      description:"1-退款 2-退货 3-换货"`
	ApplyQuantity      int         `json:"apply_quantity"       description:"申请售后数量"`
	Reason             string      `json:"reason"               description:"申请原因"`
	Amount             int64       `json:"amount"               description:"退款金额（分）"`
	RefundId           int64       `json:"refund_id"            description:"关联退款单ID"`
	RefundNo           string      `json:"refund_no"            description:"关联退款单号"`
	Status             int         `json:"status"               description:"0-待审核 1-审核通过(待退货) 2-退货中 3-待退款 4-已完成 5-已拒绝 6-已取消"`
	ReturnCarrier      string      `json:"return_carrier"       description:"退货物流商"`
	ReturnTrackingNo   string      `json:"return_tracking_no"   description:"退货运单号"`
	ReturnShippedAt    *gtime.Time `json:"return_shipped_at"    description:"买家退货发出时间"`
	MerchantReceivedAt *gtime.Time `json:"merchant_received_at" description:"商家收货时间"`
	ApplyAt            *gtime.Time `json:"apply_at"             description:"申请时间"`
	AuditedAt          *gtime.Time `json:"audited_at"           description:"审核时间"`
	CompletedAt        *gtime.Time `json:"completed_at"         description:"完成时间"`
	CreatedAt          *gtime.Time `json:"created_at"           description:""`
	UpdatedAt          *gtime.Time `json:"updated_at"           description:""`
	DeletedAt          *gtime.Time `json:"deleted_at"           description:""`
}
