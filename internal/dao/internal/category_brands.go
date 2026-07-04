// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CategoryBrandsDao is the data access object for table sp_category_brands.
type CategoryBrandsDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns CategoryBrandsColumns // columns contains all the column names of Table for convenient usage.
}

// CategoryBrandsColumns defines and stores column names for table sp_category_brands.
type CategoryBrandsColumns struct {
	Id         string //
	CategoryId string // 关联 categories.id
	BrandId    string // 关联 brands.id
	SortOrder  string // 排序权重（越小越靠前，控制该类目下品牌的展示顺序）
	CreatedAt  string //
}

// categoryBrandsColumns holds the columns for table sp_category_brands.
var categoryBrandsColumns = CategoryBrandsColumns{
	Id:         "id",
	CategoryId: "category_id",
	BrandId:    "brand_id",
	SortOrder:  "sort_order",
	CreatedAt:  "created_at",
}

// NewCategoryBrandsDao creates and returns a new DAO object for table data access.
func NewCategoryBrandsDao() *CategoryBrandsDao {
	return &CategoryBrandsDao{
		group:   "default",
		table:   "sp_category_brands",
		columns: categoryBrandsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CategoryBrandsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CategoryBrandsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CategoryBrandsDao) Columns() CategoryBrandsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CategoryBrandsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CategoryBrandsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CategoryBrandsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
