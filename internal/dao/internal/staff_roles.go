// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// StaffRolesDao is the data access object for table sys_staff_roles.
type StaffRolesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns StaffRolesColumns // columns contains all the column names of Table for convenient usage.
}

// StaffRolesColumns defines and stores column names for table sys_staff_roles.
type StaffRolesColumns struct {
	Id        string // 主键
	StaffId   string // 员工ID（关联 sys_staff.id）
	RoleId    string // 角色ID（关联 sys_roles.id）
	CreatedAt string // 创建时间
	DeletedAt string // 删除时间
}

// staffRolesColumns holds the columns for table sys_staff_roles.
var staffRolesColumns = StaffRolesColumns{
	Id:        "id",
	StaffId:   "staff_id",
	RoleId:    "role_id",
	CreatedAt: "created_at",
	DeletedAt: "deleted_at",
}

// NewStaffRolesDao creates and returns a new DAO object for table data access.
func NewStaffRolesDao() *StaffRolesDao {
	return &StaffRolesDao{
		group:   "default",
		table:   "sys_staff_roles",
		columns: staffRolesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *StaffRolesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *StaffRolesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *StaffRolesDao) Columns() StaffRolesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *StaffRolesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *StaffRolesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *StaffRolesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
