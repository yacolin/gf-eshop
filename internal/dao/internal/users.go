// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersDao is the data access object for table usr_users.
type UsersDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns UsersColumns // columns contains all the column names of Table for convenient usage.
}

// UsersColumns defines and stores column names for table usr_users.
type UsersColumns struct {
	Id             string // 主键
	Username       string // 用户名（唯一，NULL表示未设置）
	PasswordHash   string // bcrypt 密码哈希
	Email          string // 邮箱（唯一，NULL表示未绑定）
	EmailVerified  string // 邮箱是否已验证
	Phone          string // 手机号（唯一，NULL表示未绑定）
	PhoneVerified  string // 手机号是否已验证
	Avatar         string // 头像URL
	Nickname       string // 昵称
	Status         string // 状态：1-正常 0-禁用 2-冻结
	RegisterIp     string // 注册IP
	RegisterSource string // 注册来源：web/ios/android/admin
	LastLoginIp    string // 最后登录IP
	LastLoginAt    string // 最后登录时间
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 删除时间
}

// usersColumns holds the columns for table usr_users.
var usersColumns = UsersColumns{
	Id:             "id",
	Username:       "username",
	PasswordHash:   "password_hash",
	Email:          "email",
	EmailVerified:  "email_verified",
	Phone:          "phone",
	PhoneVerified:  "phone_verified",
	Avatar:         "avatar",
	Nickname:       "nickname",
	Status:         "status",
	RegisterIp:     "register_ip",
	RegisterSource: "register_source",
	LastLoginIp:    "last_login_ip",
	LastLoginAt:    "last_login_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewUsersDao creates and returns a new DAO object for table data access.
func NewUsersDao() *UsersDao {
	return &UsersDao{
		group:   "default",
		table:   "usr_users",
		columns: usersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersDao) Columns() UsersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
