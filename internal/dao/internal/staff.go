// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// StaffDao is the data access object for table sys_staff.
type StaffDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns StaffColumns // columns contains all the column names of Table for convenient usage.
}

// StaffColumns defines and stores column names for table sys_staff.
type StaffColumns struct {
	Id           string // 主键
	Username     string // 登录用户名（唯一）
	PasswordHash string // bcrypt 密码哈希
	RealName     string // 真实姓名
	Email        string // 邮箱
	Phone        string // 手机号
	Avatar       string // 头像URL
	Status       string // 1-正常 0-禁用
	LastLoginIp  string // 最后登录IP
	LastLoginAt  string // 最后登录时间
	CreatedBy    string // 创建人ID（0=系统）
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 删除时间
}

// staffColumns holds the columns for table sys_staff.
var staffColumns = StaffColumns{
	Id:           "id",
	Username:     "username",
	PasswordHash: "password_hash",
	RealName:     "real_name",
	Email:        "email",
	Phone:        "phone",
	Avatar:       "avatar",
	Status:       "status",
	LastLoginIp:  "last_login_ip",
	LastLoginAt:  "last_login_at",
	CreatedBy:    "created_by",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewStaffDao creates and returns a new DAO object for table data access.
func NewStaffDao() *StaffDao {
	return &StaffDao{
		group:   "default",
		table:   "sys_staff",
		columns: staffColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *StaffDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *StaffDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *StaffDao) Columns() StaffColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *StaffDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *StaffDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *StaffDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
