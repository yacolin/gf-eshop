// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentLogs is the golang structure of table tx_payment_logs for DAO operations like Where/Data.
type PaymentLogs struct {
	g.Meta        `orm:"table:tx_payment_logs, do:true"`
	Id            interface{} // 日志ID
	PaymentId     interface{} // 关联 tx_payments.id
	PaymentNo     interface{} // 支付单号（冗余）
	Channel       interface{} // 支付渠道
	TransactionId interface{} // 渠道交易号
	Action        interface{} // 操作类型：create-创建 pay-支付回调 refund-退款 refund_callback-退款回调 close-关闭
	RequestBody   interface{} // 请求参数（渠道原始数据，用于对账排查）
	ResponseBody  interface{} // 响应结果（渠道原始数据）
	Status        interface{} // 操作结果状态
	CreatedAt     *gtime.Time //
}
