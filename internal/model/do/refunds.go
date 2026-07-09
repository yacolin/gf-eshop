// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Refunds is the golang structure of table tx_refunds for DAO operations like Where/Data.
type Refunds struct {
	g.Meta          `orm:"table:tx_refunds, do:true"`
	Id              interface{} // 退款单ID
	RefundNo        interface{} // 退款单号（业务唯一键）
	PaymentId       interface{} // 关联 tx_payments.id
	PaymentNo       interface{} // 关联支付单号
	OrderNo         interface{} // 关联订单号
	OrderId         interface{} // 关联 tx_orders.id
	MerchantId      interface{} // 所属商家ID
	Amount          interface{} // 退款金额（分）
	Reason          interface{} // 退款原因
	Status          interface{} // 退款状态：pending-待处理 processing-处理中 success-已退款 failed-退款失败 rejected-已拒绝
	ChannelRefundId interface{} // 渠道退款交易号
	FailureReason   interface{} // 失败原因
	ChannelResponse interface{} // 渠道退款响应/回调摘要
	IdempotencyKey  interface{} // 退款幂等键
	AppliedAt       *gtime.Time // 申请时间
	SuccessAt       *gtime.Time // 退款成功时间
	NotifyAt        *gtime.Time // 最近一次退款回调时间
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
