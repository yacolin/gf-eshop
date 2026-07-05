// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductDescriptionsDao is the data access object for table sp_product_descriptions.
type ProductDescriptionsDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns ProductDescriptionsColumns // columns contains all the column names of Table for convenient usage.
}

// ProductDescriptionsColumns defines and stores column names for table sp_product_descriptions.
type ProductDescriptionsColumns struct {
	Id                string //
	ProductId         string // 关联 products.id
	Description       string // 商品详情（富文本HTML）
	MobileDescription string // 移动端详情（可选，适配手机展示）
	CreatedAt         string //
	UpdatedAt         string //
}

// productDescriptionsColumns holds the columns for table sp_product_descriptions.
var productDescriptionsColumns = ProductDescriptionsColumns{
	Id:                "id",
	ProductId:         "product_id",
	Description:       "description",
	MobileDescription: "mobile_description",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewProductDescriptionsDao creates and returns a new DAO object for table data access.
func NewProductDescriptionsDao() *ProductDescriptionsDao {
	return &ProductDescriptionsDao{
		group:   "default",
		table:   "sp_product_descriptions",
		columns: productDescriptionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductDescriptionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductDescriptionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductDescriptionsDao) Columns() ProductDescriptionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductDescriptionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductDescriptionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductDescriptionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
