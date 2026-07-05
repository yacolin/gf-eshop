// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentsDao is the data access object for table tx_payments.
type PaymentsDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns PaymentsColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentsColumns defines and stores column names for table tx_payments.
type PaymentsColumns struct {
	Id              string // 支付单ID
	PaymentNo       string // 支付单号（业务唯一键）
	OrderNo         string // 关联订单号
	OrderId         string // 关联 tx_orders.id
	MerchantId      string // 所属商家ID
	OrderType       string // 订单类型：order-普通订单 flash-秒杀订单
	Amount          string // 支付金额（分）
	Currency        string //
	PaymentMethod   string // 支付方式：wechat-微信 alipay-支付宝 wallet-余额
	Channel         string // 支付渠道（如 wechat_native-微信 native alipay_page-支付宝页面）
	TradeType       string // 交易类型：native-jsapi-app-h5-page
	TransactionId   string // 支付渠道交易号（微信/支付宝订单号，用于对账）
	IdempotencyKey  string // 支付创建幂等键（防重复提交）
	Status          string // 支付状态：pending-待支付 processing-处理中 success-已支付 failed-支付失败 refunding-退款中 refunded-已退款
	FailureReason   string // 失败原因
	ClientIp        string // 客户端IP
	ExpireAt        string // 支付过期时间
	PaidAt          string // 支付成功时间
	NotifyAt        string // 最近一次渠道回调时间
	ChannelResponse string // 渠道最近一次响应/回调摘要
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// paymentsColumns holds the columns for table tx_payments.
var paymentsColumns = PaymentsColumns{
	Id:              "id",
	PaymentNo:       "payment_no",
	OrderNo:         "order_no",
	OrderId:         "order_id",
	MerchantId:      "merchant_id",
	OrderType:       "order_type",
	Amount:          "amount",
	Currency:        "currency",
	PaymentMethod:   "payment_method",
	Channel:         "channel",
	TradeType:       "trade_type",
	TransactionId:   "transaction_id",
	IdempotencyKey:  "idempotency_key",
	Status:          "status",
	FailureReason:   "failure_reason",
	ClientIp:        "client_ip",
	ExpireAt:        "expire_at",
	PaidAt:          "paid_at",
	NotifyAt:        "notify_at",
	ChannelResponse: "channel_response",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewPaymentsDao creates and returns a new DAO object for table data access.
func NewPaymentsDao() *PaymentsDao {
	return &PaymentsDao{
		group:   "default",
		table:   "tx_payments",
		columns: paymentsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentsDao) Columns() PaymentsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
