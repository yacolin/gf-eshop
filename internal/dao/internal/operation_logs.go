// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OperationLogsDao is the data access object for table sys_operation_logs.
type OperationLogsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns OperationLogsColumns // columns contains all the column names of Table for convenient usage.
}

// OperationLogsColumns defines and stores column names for table sys_operation_logs.
type OperationLogsColumns struct {
	Id            string // 主键
	StaffId       string // 操作员工ID（关联 sys_staff.id）
	StaffName     string // 操作员工姓名（冗余，便于查询）
	Operation     string // 操作类型：update_price/disable_user/create_coupon/...
	Resource      string // 操作资源：product/order/user/coupon/...
	ResourceId    string // 资源ID
	Detail        string // 操作详情JSON（记录变更前后快照）
	Result        string // 1-成功 0-失败
	FailureReason string // 失败原因
	Ip            string // 操作IP
	CreatedAt     string // 创建时间
}

// operationLogsColumns holds the columns for table sys_operation_logs.
var operationLogsColumns = OperationLogsColumns{
	Id:            "id",
	StaffId:       "staff_id",
	StaffName:     "staff_name",
	Operation:     "operation",
	Resource:      "resource",
	ResourceId:    "resource_id",
	Detail:        "detail",
	Result:        "result",
	FailureReason: "failure_reason",
	Ip:            "ip",
	CreatedAt:     "created_at",
}

// NewOperationLogsDao creates and returns a new DAO object for table data access.
func NewOperationLogsDao() *OperationLogsDao {
	return &OperationLogsDao{
		group:   "default",
		table:   "sys_operation_logs",
		columns: operationLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OperationLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OperationLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OperationLogsDao) Columns() OperationLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OperationLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OperationLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OperationLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
