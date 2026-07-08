// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DepartmentsDao is the data access object for table sys_departments.
type DepartmentsDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns DepartmentsColumns // columns contains all the column names of Table for convenient usage.
}

// DepartmentsColumns defines and stores column names for table sys_departments.
type DepartmentsColumns struct {
	Id        string // 主键
	Name      string // 部门名称
	ParentId  string // 上级部门ID（0=根部门）
	SortOrder string // 排序值
	Status    string // 1-启用 0-禁用
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	DeletedAt string // 删除时间
}

// departmentsColumns holds the columns for table sys_departments.
var departmentsColumns = DepartmentsColumns{
	Id:        "id",
	Name:      "name",
	ParentId:  "parent_id",
	SortOrder: "sort_order",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewDepartmentsDao creates and returns a new DAO object for table data access.
func NewDepartmentsDao() *DepartmentsDao {
	return &DepartmentsDao{
		group:   "default",
		table:   "sys_departments",
		columns: departmentsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DepartmentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DepartmentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DepartmentsDao) Columns() DepartmentsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DepartmentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DepartmentsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DepartmentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
