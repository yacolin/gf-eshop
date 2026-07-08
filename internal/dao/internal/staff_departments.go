// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// StaffDepartmentsDao is the data access object for table sys_staff_departments.
type StaffDepartmentsDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns StaffDepartmentsColumns // columns contains all the column names of Table for convenient usage.
}

// StaffDepartmentsColumns defines and stores column names for table sys_staff_departments.
type StaffDepartmentsColumns struct {
	Id           string // 主键
	StaffId      string // 员工ID（关联 sys_staff.id）
	DepartmentId string // 部门ID（关联 sys_departments.id）
	IsPrimary    string // 是否主部门
	CreatedAt    string // 创建时间
	DeletedAt    string // 删除时间
}

// staffDepartmentsColumns holds the columns for table sys_staff_departments.
var staffDepartmentsColumns = StaffDepartmentsColumns{
	Id:           "id",
	StaffId:      "staff_id",
	DepartmentId: "department_id",
	IsPrimary:    "is_primary",
	CreatedAt:    "created_at",
	DeletedAt:    "deleted_at",
}

// NewStaffDepartmentsDao creates and returns a new DAO object for table data access.
func NewStaffDepartmentsDao() *StaffDepartmentsDao {
	return &StaffDepartmentsDao{
		group:   "default",
		table:   "sys_staff_departments",
		columns: staffDepartmentsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *StaffDepartmentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *StaffDepartmentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *StaffDepartmentsDao) Columns() StaffDepartmentsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *StaffDepartmentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *StaffDepartmentsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *StaffDepartmentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
