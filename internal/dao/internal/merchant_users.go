// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantUsersDao is the data access object for table mch_merchant_users.
type MerchantUsersDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns MerchantUsersColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantUsersColumns defines and stores column names for table mch_merchant_users.
type MerchantUsersColumns struct {
	Id          string // 主键
	MerchantId  string // 商家ID
	StaffId     string // 员工ID（关联 sys_staff.id）
	RoleId      string // 商家角色ID（关联 mch_roles.id，与平台RBAC隔离）
	Status      string // 1-正常 2-禁用
	InvitedAt   string // 邀请时间
	LastLoginAt string // 最后登录时间
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// merchantUsersColumns holds the columns for table mch_merchant_users.
var merchantUsersColumns = MerchantUsersColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	StaffId:     "staff_id",
	RoleId:      "role_id",
	Status:      "status",
	InvitedAt:   "invited_at",
	LastLoginAt: "last_login_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewMerchantUsersDao creates and returns a new DAO object for table data access.
func NewMerchantUsersDao() *MerchantUsersDao {
	return &MerchantUsersDao{
		group:   "default",
		table:   "mch_merchant_users",
		columns: merchantUsersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantUsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantUsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantUsersDao) Columns() MerchantUsersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantUsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantUsersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantUsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
