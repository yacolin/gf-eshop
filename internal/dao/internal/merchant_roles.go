// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantRolesDao is the data access object for table mch_merchant_roles.
type MerchantRolesDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns MerchantRolesColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantRolesColumns defines and stores column names for table mch_merchant_roles.
type MerchantRolesColumns struct {
	Id          string // 主键
	MerchantId  string // 商家ID（0=平台预置角色）
	Name        string // 角色名称（如店长/运营/财务）
	DisplayName string // 角色显示名称
	Description string // 角色描述
	RoleType    string // builtin-系统预置 custom-商家自定义
	SortOrder   string // 排序值
	Status      string // 1-启用 0-禁用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// merchantRolesColumns holds the columns for table mch_merchant_roles.
var merchantRolesColumns = MerchantRolesColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
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

// NewMerchantRolesDao creates and returns a new DAO object for table data access.
func NewMerchantRolesDao() *MerchantRolesDao {
	return &MerchantRolesDao{
		group:   "default",
		table:   "mch_merchant_roles",
		columns: merchantRolesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantRolesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantRolesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantRolesDao) Columns() MerchantRolesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantRolesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantRolesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantRolesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
