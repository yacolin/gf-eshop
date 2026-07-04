// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CategoriesDao is the data access object for table sp_categories.
type CategoriesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns CategoriesColumns // columns contains all the column names of Table for convenient usage.
}

// CategoriesColumns defines and stores column names for table sp_categories.
type CategoriesColumns struct {
	Id        string //
	Name      string // 类目名称（如：手机）
	ParentId  string // 父级ID（0表示根节点）
	Level     string // 层级（1-3级）
	Path      string // 路径（如：1/23/45/）
	IconUrl   string // 类目图标
	SortOrder string // 排序
	Status    string // 1-启用 0-禁用
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// categoriesColumns holds the columns for table sp_categories.
var categoriesColumns = CategoriesColumns{
	Id:        "id",
	Name:      "name",
	ParentId:  "parent_id",
	Level:     "level",
	Path:      "path",
	IconUrl:   "icon_url",
	SortOrder: "sort_order",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewCategoriesDao creates and returns a new DAO object for table data access.
func NewCategoriesDao() *CategoriesDao {
	return &CategoriesDao{
		group:   "default",
		table:   "sp_categories",
		columns: categoriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CategoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CategoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CategoriesDao) Columns() CategoriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CategoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CategoriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CategoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
