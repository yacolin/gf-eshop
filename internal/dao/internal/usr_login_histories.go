// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsrLoginHistoriesDao is the data access object for table usr_login_histories.
type UsrLoginHistoriesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns UsrLoginHistoriesColumns // columns contains all the column names of Table for convenient usage.
}

// UsrLoginHistoriesColumns defines and stores column names for table usr_login_histories.
type UsrLoginHistoriesColumns struct {
	Id            string // 主键
	UserId        string // 用户ID
	LoginIp       string // 登录IP
	LoginDevice   string // 登录设备信息（UA）
	LoginLocation string // 登录地点
	LoginMethod   string // 登录方式：password/sms/oauth
	LoginStatus   string // 登录结果：1-成功 0-失败
	FailureReason string // 失败原因
	CreatedAt     string // 创建时间
}

// usrLoginHistoriesColumns holds the columns for table usr_login_histories.
var usrLoginHistoriesColumns = UsrLoginHistoriesColumns{
	Id:            "id",
	UserId:        "user_id",
	LoginIp:       "login_ip",
	LoginDevice:   "login_device",
	LoginLocation: "login_location",
	LoginMethod:   "login_method",
	LoginStatus:   "login_status",
	FailureReason: "failure_reason",
	CreatedAt:     "created_at",
}

// NewUsrLoginHistoriesDao creates and returns a new DAO object for table data access.
func NewUsrLoginHistoriesDao() *UsrLoginHistoriesDao {
	return &UsrLoginHistoriesDao{
		group:   "default",
		table:   "usr_login_histories",
		columns: usrLoginHistoriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsrLoginHistoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsrLoginHistoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsrLoginHistoriesDao) Columns() UsrLoginHistoriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsrLoginHistoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsrLoginHistoriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsrLoginHistoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
