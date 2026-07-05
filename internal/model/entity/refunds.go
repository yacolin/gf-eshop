// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Refunds is the golang structure for table refunds.
type Refunds struct {
	Id              int64       `json:"id"                description:"退款单ID"`
	RefundNo        string      `json:"refund_no"         description:"退款单号（业务唯一键）"`
	PaymentId       int64       `json:"payment_id"        description:"关联 tx_payments.id"`
	PaymentNo       string      `json:"payment_no"        description:"关联支付单号"`
	OrderNo         string      `json:"order_no"          description:"关联订单号"`
	OrderId         int64       `json:"order_id"          description:"关联 tx_orders.id"`
	MerchantId      int64       `json:"merchant_id"       description:"所属商家ID"`
	Amount          int64       `json:"amount"            description:"退款金额（分）"`
	Reason          string      `json:"reason"            description:"退款原因"`
	Status          string      `json:"status"            description:"退款状态：pending-待处理 processing-处理中 success-已退款 failed-退款失败 rejected-已拒绝"`
	ChannelRefundId string      `json:"channel_refund_id" description:"渠道退款交易号"`
	FailureReason   string      `json:"failure_reason"    description:"失败原因"`
	ChannelResponse string      `json:"channel_response"  description:"渠道退款响应/回调摘要"`
	AppliedAt       *gtime.Time `json:"applied_at"        description:"申请时间"`
	SuccessAt       *gtime.Time `json:"success_at"        description:"退款成功时间"`
	NotifyAt        *gtime.Time `json:"notify_at"         description:"最近一次退款回调时间"`
	CreatedAt       *gtime.Time `json:"created_at"        description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"        description:""`
	DeletedAt       *gtime.Time `json:"deleted_at"        description:""`
}
