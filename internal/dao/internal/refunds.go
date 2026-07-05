// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RefundsDao is the data access object for table tx_refunds.
type RefundsDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns RefundsColumns // columns contains all the column names of Table for convenient usage.
}

// RefundsColumns defines and stores column names for table tx_refunds.
type RefundsColumns struct {
	Id              string // 退款单ID
	RefundNo        string // 退款单号（业务唯一键）
	PaymentId       string // 关联 tx_payments.id
	PaymentNo       string // 关联支付单号
	OrderNo         string // 关联订单号
	OrderId         string // 关联 tx_orders.id
	MerchantId      string // 所属商家ID
	Amount          string // 退款金额（分）
	Reason          string // 退款原因
	Status          string // 退款状态：pending-待处理 processing-处理中 success-已退款 failed-退款失败 rejected-已拒绝
	ChannelRefundId string // 渠道退款交易号
	FailureReason   string // 失败原因
	ChannelResponse string // 渠道退款响应/回调摘要
	AppliedAt       string // 申请时间
	SuccessAt       string // 退款成功时间
	NotifyAt        string // 最近一次退款回调时间
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// refundsColumns holds the columns for table tx_refunds.
var refundsColumns = RefundsColumns{
	Id:              "id",
	RefundNo:        "refund_no",
	PaymentId:       "payment_id",
	PaymentNo:       "payment_no",
	OrderNo:         "order_no",
	OrderId:         "order_id",
	MerchantId:      "merchant_id",
	Amount:          "amount",
	Reason:          "reason",
	Status:          "status",
	ChannelRefundId: "channel_refund_id",
	FailureReason:   "failure_reason",
	ChannelResponse: "channel_response",
	AppliedAt:       "applied_at",
	SuccessAt:       "success_at",
	NotifyAt:        "notify_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewRefundsDao creates and returns a new DAO object for table data access.
func NewRefundsDao() *RefundsDao {
	return &RefundsDao{
		group:   "default",
		table:   "tx_refunds",
		columns: refundsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *RefundsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *RefundsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *RefundsDao) Columns() RefundsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *RefundsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *RefundsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *RefundsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
