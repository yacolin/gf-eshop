// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantContactsDao is the data access object for table mch_merchant_contacts.
type MerchantContactsDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns MerchantContactsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantContactsColumns defines and stores column names for table mch_merchant_contacts.
type MerchantContactsColumns struct {
	Id           string // 主键
	MerchantId   string // 商家ID
	ContactName  string // 联系人姓名
	ContactPhone string // 联系电话
	ContactRole  string // 联系人角色：finance-财务 legal-法人 operation-运营
	IsPrimary    string // 是否主要联系人 0-否 1-是
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
}

// merchantContactsColumns holds the columns for table mch_merchant_contacts.
var merchantContactsColumns = MerchantContactsColumns{
	Id:           "id",
	MerchantId:   "merchant_id",
	ContactName:  "contact_name",
	ContactPhone: "contact_phone",
	ContactRole:  "contact_role",
	IsPrimary:    "is_primary",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewMerchantContactsDao creates and returns a new DAO object for table data access.
func NewMerchantContactsDao() *MerchantContactsDao {
	return &MerchantContactsDao{
		group:   "default",
		table:   "mch_merchant_contacts",
		columns: merchantContactsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantContactsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantContactsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantContactsDao) Columns() MerchantContactsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantContactsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantContactsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantContactsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
