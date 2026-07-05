// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentLogs is the golang structure for table payment_logs.
type PaymentLogs struct {
	Id            int64       `json:"id"             description:"日志ID"`
	PaymentId     int64       `json:"payment_id"     description:"关联 tx_payments.id"`
	PaymentNo     string      `json:"payment_no"     description:"支付单号（冗余）"`
	Channel       string      `json:"channel"        description:"支付渠道"`
	TransactionId string      `json:"transaction_id" description:"渠道交易号"`
	Action        string      `json:"action"         description:"操作类型：create-创建 pay-支付回调 refund-退款 refund_callback-退款回调 close-关闭"`
	RequestBody   string      `json:"request_body"   description:"请求参数（渠道原始数据，用于对账排查）"`
	ResponseBody  string      `json:"response_body"  description:"响应结果（渠道原始数据）"`
	Status        string      `json:"status"         description:"操作结果状态"`
	CreatedAt     *gtime.Time `json:"created_at"     description:""`
}
