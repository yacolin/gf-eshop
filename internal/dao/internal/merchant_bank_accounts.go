// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantBankAccountsDao is the data access object for table mch_merchant_bank_accounts.
type MerchantBankAccountsDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns MerchantBankAccountsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantBankAccountsColumns defines and stores column names for table mch_merchant_bank_accounts.
type MerchantBankAccountsColumns struct {
	Id          string // 主键
	MerchantId  string // 商家ID
	BankName    string // 开户行
	BankBranch  string // 开户支行
	AccountName string // 开户名
	AccountNo   string // 银行账号
	AccountType string // 1-对公账户 2-对私账户
	IsDefault   string // 是否默认结算账户（NULL=非默认, 1=默认）
	Status      string // 1-正常 2-禁用
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// merchantBankAccountsColumns holds the columns for table mch_merchant_bank_accounts.
var merchantBankAccountsColumns = MerchantBankAccountsColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	BankName:    "bank_name",
	BankBranch:  "bank_branch",
	AccountName: "account_name",
	AccountNo:   "account_no",
	AccountType: "account_type",
	IsDefault:   "is_default",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewMerchantBankAccountsDao creates and returns a new DAO object for table data access.
func NewMerchantBankAccountsDao() *MerchantBankAccountsDao {
	return &MerchantBankAccountsDao{
		group:   "default",
		table:   "mch_merchant_bank_accounts",
		columns: merchantBankAccountsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantBankAccountsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantBankAccountsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantBankAccountsDao) Columns() MerchantBankAccountsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantBankAccountsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantBankAccountsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantBankAccountsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
