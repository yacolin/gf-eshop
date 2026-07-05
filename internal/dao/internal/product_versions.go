// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductVersionsDao is the data access object for table sp_product_versions.
type ProductVersionsDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns ProductVersionsColumns // columns contains all the column names of Table for convenient usage.
}

// ProductVersionsColumns defines and stores column names for table sp_product_versions.
type ProductVersionsColumns struct {
	Id            string // 主键
	ProductId     string // 关联 sp_products.id
	Version       string // 版本号（从1递增）
	Diff          string // 变更JSON（{"before": {...}, "after": {...}}）
	ChangedFields string // 变更字段列表（如：["name", "price", "status"]）
	Operator      string // 操作人
	OperatorId    string // 操作人ID
	Reason        string // 变更原因
	CreatedAt     string //
}

// productVersionsColumns holds the columns for table sp_product_versions.
var productVersionsColumns = ProductVersionsColumns{
	Id:            "id",
	ProductId:     "product_id",
	Version:       "version",
	Diff:          "diff",
	ChangedFields: "changed_fields",
	Operator:      "operator",
	OperatorId:    "operator_id",
	Reason:        "reason",
	CreatedAt:     "created_at",
}

// NewProductVersionsDao creates and returns a new DAO object for table data access.
func NewProductVersionsDao() *ProductVersionsDao {
	return &ProductVersionsDao{
		group:   "default",
		table:   "sp_product_versions",
		columns: productVersionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductVersionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductVersionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductVersionsDao) Columns() ProductVersionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductVersionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductVersionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductVersionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
