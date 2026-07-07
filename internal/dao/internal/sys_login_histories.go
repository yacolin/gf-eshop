// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysLoginHistoriesDao is the data access object for table sys_login_histories.
type SysLoginHistoriesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns SysLoginHistoriesColumns // columns contains all the column names of Table for convenient usage.
}

// SysLoginHistoriesColumns defines and stores column names for table sys_login_histories.
type SysLoginHistoriesColumns struct {
	Id            string // 主键
	StaffId       string // 员工ID（关联 sys_staff.id）
	LoginIp       string // 登录IP
	LoginDevice   string // 登录设备信息（UA）
	LoginLocation string // 登录地点
	LoginMethod   string // 登录方式：password/sms/oauth
	LoginStatus   string // 1-成功 0-失败
	FailureReason string // 失败原因
	CreatedAt     string // 创建时间
}

// sysLoginHistoriesColumns holds the columns for table sys_login_histories.
var sysLoginHistoriesColumns = SysLoginHistoriesColumns{
	Id:            "id",
	StaffId:       "staff_id",
	LoginIp:       "login_ip",
	LoginDevice:   "login_device",
	LoginLocation: "login_location",
	LoginMethod:   "login_method",
	LoginStatus:   "login_status",
	FailureReason: "failure_reason",
	CreatedAt:     "created_at",
}

// NewSysLoginHistoriesDao creates and returns a new DAO object for table data access.
func NewSysLoginHistoriesDao() *SysLoginHistoriesDao {
	return &SysLoginHistoriesDao{
		group:   "default",
		table:   "sys_login_histories",
		columns: sysLoginHistoriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysLoginHistoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SysLoginHistoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SysLoginHistoriesDao) Columns() SysLoginHistoriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SysLoginHistoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysLoginHistoriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysLoginHistoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
