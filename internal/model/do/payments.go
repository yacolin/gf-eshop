// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Payments is the golang structure of table tx_payments for DAO operations like Where/Data.
type Payments struct {
	g.Meta          `orm:"table:tx_payments, do:true"`
	Id              interface{} // 支付单ID
	PaymentNo       interface{} // 支付单号（业务唯一键）
	OrderNo         interface{} // 关联订单号
	OrderId         interface{} // 关联 tx_orders.id
	MerchantId      interface{} // 所属商家ID
	OrderType       interface{} // 订单类型：order-普通订单 flash-秒杀订单
	Amount          interface{} // 支付金额（分）
	Currency        interface{} //
	PaymentMethod   interface{} // 支付方式：wechat-微信 alipay-支付宝 wallet-余额
	Channel         interface{} // 支付渠道（如 wechat_native-微信 native alipay_page-支付宝页面）
	TradeType       interface{} // 交易类型：native-jsapi-app-h5-page
	TransactionId   interface{} // 支付渠道交易号（微信/支付宝订单号，用于对账）
	IdempotencyKey  interface{} // 支付创建幂等键（防重复提交）
	Status          interface{} // 支付状态：pending-待支付 processing-处理中 success-已支付 failed-支付失败 refunding-退款中 refunded-已退款
	FailureReason   interface{} // 失败原因
	ClientIp        interface{} // 客户端IP
	ExpireAt        *gtime.Time // 支付过期时间
	PaidAt          *gtime.Time // 支付成功时间
	NotifyAt        *gtime.Time // 最近一次渠道回调时间
	ChannelResponse interface{} // 渠道最近一次响应/回调摘要
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
