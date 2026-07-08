// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PermissionsDao is the data access object for table sys_permissions.
type PermissionsDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PermissionsColumns // columns contains all the column names of Table for convenient usage.
}

// PermissionsColumns defines and stores column names for table sys_permissions.
type PermissionsColumns struct {
	Id          string // 主键
	Name        string // 权限标识（唯一，如 order:create）
	DisplayName string // 权限显示名称
	Description string // 权限描述
	Resource    string // 资源（如 order/product/user）
	Action      string // 操作（如 create/read/update/delete）
	ParentId    string // 父级ID（0=根节点，支持菜单/按钮树形层级）
	Category    string // 分类（如 business/system/admin）
	SortOrder   string // 排序值
	Status      string // 1-启用 0-禁用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// permissionsColumns holds the columns for table sys_permissions.
var permissionsColumns = PermissionsColumns{
	Id:          "id",
	Name:        "name",
	DisplayName: "display_name",
	Description: "description",
	Resource:    "resource",
	Action:      "action",
	ParentId:    "parent_id",
	Category:    "category",
	SortOrder:   "sort_order",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewPermissionsDao creates and returns a new DAO object for table data access.
func NewPermissionsDao() *PermissionsDao {
	return &PermissionsDao{
		group:   "default",
		table:   "sys_permissions",
		columns: permissionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PermissionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PermissionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PermissionsDao) Columns() PermissionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PermissionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PermissionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PermissionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
