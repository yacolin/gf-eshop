// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentLogsDao is the data access object for table tx_payment_logs.
type PaymentLogsDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PaymentLogsColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentLogsColumns defines and stores column names for table tx_payment_logs.
type PaymentLogsColumns struct {
	Id            string // 日志ID
	PaymentId     string // 关联 tx_payments.id
	PaymentNo     string // 支付单号（冗余）
	Channel       string // 支付渠道
	TransactionId string // 渠道交易号
	Action        string // 操作类型：create-创建 pay-支付回调 refund-退款 refund_callback-退款回调 close-关闭
	RequestBody   string // 请求参数（渠道原始数据，用于对账排查）
	ResponseBody  string // 响应结果（渠道原始数据）
	Status        string // 操作结果状态
	CreatedAt     string //
}

// paymentLogsColumns holds the columns for table tx_payment_logs.
var paymentLogsColumns = PaymentLogsColumns{
	Id:            "id",
	PaymentId:     "payment_id",
	PaymentNo:     "payment_no",
	Channel:       "channel",
	TransactionId: "transaction_id",
	Action:        "action",
	RequestBody:   "request_body",
	ResponseBody:  "response_body",
	Status:        "status",
	CreatedAt:     "created_at",
}

// NewPaymentLogsDao creates and returns a new DAO object for table data access.
func NewPaymentLogsDao() *PaymentLogsDao {
	return &PaymentLogsDao{
		group:   "default",
		table:   "tx_payment_logs",
		columns: paymentLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentLogsDao) Columns() PaymentLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
