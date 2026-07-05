// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderLogsDao is the data access object for table tx_order_logs.
type OrderLogsDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns OrderLogsColumns // columns contains all the column names of Table for convenient usage.
}

// OrderLogsColumns defines and stores column names for table tx_order_logs.
type OrderLogsColumns struct {
	Id           string // 日志ID
	OrderId      string // 关联 tx_orders.id
	OrderNo      string // 订单号（冗余，便于按号查日志）
	FromStatus   string // 变更前状态
	ToStatus     string // 变更后状态
	Operator     string // 操作人
	OperatorType string // 操作人类型：system-系统 user-用户 admin-管理员
	Note         string // 备注（如：支付成功、超时取消）
	CreatedAt    string //
}

// orderLogsColumns holds the columns for table tx_order_logs.
var orderLogsColumns = OrderLogsColumns{
	Id:           "id",
	OrderId:      "order_id",
	OrderNo:      "order_no",
	FromStatus:   "from_status",
	ToStatus:     "to_status",
	Operator:     "operator",
	OperatorType: "operator_type",
	Note:         "note",
	CreatedAt:    "created_at",
}

// NewOrderLogsDao creates and returns a new DAO object for table data access.
func NewOrderLogsDao() *OrderLogsDao {
	return &OrderLogsDao{
		group:   "default",
		table:   "tx_order_logs",
		columns: orderLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OrderLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OrderLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OrderLogsDao) Columns() OrderLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OrderLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OrderLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OrderLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
