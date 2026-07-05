// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RolesDao is the data access object for table sys_roles.
type RolesDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns RolesColumns // columns contains all the column names of Table for convenient usage.
}

// RolesColumns defines and stores column names for table sys_roles.
type RolesColumns struct {
	Id          string // 主键
	Name        string // 角色名称（唯一，如 admin/editor/vip）
	DisplayName string // 角色显示名称
	Description string // 角色描述
	RoleType    string // builtin-系统内置 custom-自定义
	SortOrder   string // 排序值
	Status      string // 1-启用 0-禁用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// rolesColumns holds the columns for table sys_roles.
var rolesColumns = RolesColumns{
	Id:          "id",
	Name:        "name",
	DisplayName: "display_name",
	Description: "description",
	RoleType:    "role_type",
	SortOrder:   "sort_order",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewRolesDao creates and returns a new DAO object for table data access.
func NewRolesDao() *RolesDao {
	return &RolesDao{
		group:   "default",
		table:   "sys_roles",
		columns: rolesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *RolesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *RolesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *RolesDao) Columns() RolesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *RolesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *RolesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *RolesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
