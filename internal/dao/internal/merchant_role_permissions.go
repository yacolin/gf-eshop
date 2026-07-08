// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantRolePermissionsDao is the data access object for table mch_merchant_role_permissions.
type MerchantRolePermissionsDao struct {
	table   string                         // table is the underlying table name of the DAO.
	group   string                         // group is the database configuration group name of current DAO.
	columns MerchantRolePermissionsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantRolePermissionsColumns defines and stores column names for table mch_merchant_role_permissions.
type MerchantRolePermissionsColumns struct {
	Id             string // 主键
	MerchantId     string // 商家ID
	RoleId         string // 角色ID（关联 mch_roles.id）
	PermissionName string // 权限标识（对应 sys_permissions.name）
	CreatedAt      string //
	DeletedAt      string //
}

// merchantRolePermissionsColumns holds the columns for table mch_merchant_role_permissions.
var merchantRolePermissionsColumns = MerchantRolePermissionsColumns{
	Id:             "id",
	MerchantId:     "merchant_id",
	RoleId:         "role_id",
	PermissionName: "permission_name",
	CreatedAt:      "created_at",
	DeletedAt:      "deleted_at",
}

// NewMerchantRolePermissionsDao creates and returns a new DAO object for table data access.
func NewMerchantRolePermissionsDao() *MerchantRolePermissionsDao {
	return &MerchantRolePermissionsDao{
		group:   "default",
		table:   "mch_merchant_role_permissions",
		columns: merchantRolePermissionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantRolePermissionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantRolePermissionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantRolePermissionsDao) Columns() MerchantRolePermissionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantRolePermissionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantRolePermissionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantRolePermissionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
