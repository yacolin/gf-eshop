// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Payments is the golang structure for table payments.
type Payments struct {
	Id              int64       `json:"id"               description:"支付单ID"`
	PaymentNo       string      `json:"payment_no"       description:"支付单号（业务唯一键）"`
	OrderNo         string      `json:"order_no"         description:"关联订单号"`
	OrderId         int64       `json:"order_id"         description:"关联 tx_orders.id"`
	MerchantId      int64       `json:"merchant_id"      description:"所属商家ID"`
	OrderType       string      `json:"order_type"       description:"订单类型：order-普通订单 flash-秒杀订单"`
	Amount          int64       `json:"amount"           description:"支付金额（分）"`
	Currency        string      `json:"currency"         description:""`
	PaymentMethod   string      `json:"payment_method"   description:"支付方式：wechat-微信 alipay-支付宝 wallet-余额"`
	Channel         string      `json:"channel"          description:"支付渠道（如 wechat_native-微信 native alipay_page-支付宝页面）"`
	TradeType       string      `json:"trade_type"       description:"交易类型：native-jsapi-app-h5-page"`
	TransactionId   string      `json:"transaction_id"   description:"支付渠道交易号（微信/支付宝订单号，用于对账）"`
	IdempotencyKey  string      `json:"idempotency_key"  description:"支付创建幂等键（防重复提交）"`
	Status          string      `json:"status"           description:"支付状态：pending-待支付 processing-处理中 success-已支付 failed-支付失败 refunding-退款中 refunded-已退款"`
	FailureReason   string      `json:"failure_reason"   description:"失败原因"`
	ClientIp        string      `json:"client_ip"        description:"客户端IP"`
	ExpireAt        *gtime.Time `json:"expire_at"        description:"支付过期时间"`
	PaidAt          *gtime.Time `json:"paid_at"          description:"支付成功时间"`
	NotifyAt        *gtime.Time `json:"notify_at"        description:"最近一次渠道回调时间"`
	ChannelResponse string      `json:"channel_response" description:"渠道最近一次响应/回调摘要"`
	CreatedAt       *gtime.Time `json:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:""`
	DeletedAt       *gtime.Time `json:"deleted_at"       description:""`
}
