// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AfterSaleLogsDao is the data access object for table tx_after_sale_logs.
type AfterSaleLogsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns AfterSaleLogsColumns // columns contains all the column names of Table for convenient usage.
}

// AfterSaleLogsColumns defines and stores column names for table tx_after_sale_logs.
type AfterSaleLogsColumns struct {
	Id           string // 主键
	AfterSaleId  string // 售后单ID
	OperatorId   string // 操作人ID
	OperatorType string // operator/user/merchant/admin
	Action       string // 动作
	Remark       string // 备注
	CreatedAt    string //
}

// afterSaleLogsColumns holds the columns for table tx_after_sale_logs.
var afterSaleLogsColumns = AfterSaleLogsColumns{
	Id:           "id",
	AfterSaleId:  "after_sale_id",
	OperatorId:   "operator_id",
	OperatorType: "operator_type",
	Action:       "action",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewAfterSaleLogsDao creates and returns a new DAO object for table data access.
func NewAfterSaleLogsDao() *AfterSaleLogsDao {
	return &AfterSaleLogsDao{
		group:   "default",
		table:   "tx_after_sale_logs",
		columns: afterSaleLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AfterSaleLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AfterSaleLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AfterSaleLogsDao) Columns() AfterSaleLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AfterSaleLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AfterSaleLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AfterSaleLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
