// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BrandsDao is the data access object for table sp_brands.
type BrandsDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns BrandsColumns // columns contains all the column names of Table for convenient usage.
}

// BrandsColumns defines and stores column names for table sp_brands.
type BrandsColumns struct {
	Id          string //
	Name        string // 品牌名称（如：苹果）
	EnglishName string // 英文名
	LogoUrl     string // 品牌Logo（CDN）
	FirstLetter string // 首字母（A-Z，用于前台索引筛选）
	SortOrder   string // 排序权重
	Status      string // 1-启用 0-禁用
	Description string // 品牌故事
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// brandsColumns holds the columns for table sp_brands.
var brandsColumns = BrandsColumns{
	Id:          "id",
	Name:        "name",
	EnglishName: "english_name",
	LogoUrl:     "logo_url",
	FirstLetter: "first_letter",
	SortOrder:   "sort_order",
	Status:      "status",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewBrandsDao creates and returns a new DAO object for table data access.
func NewBrandsDao() *BrandsDao {
	return &BrandsDao{
		group:   "default",
		table:   "sp_brands",
		columns: brandsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *BrandsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *BrandsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *BrandsDao) Columns() BrandsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *BrandsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *BrandsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *BrandsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
