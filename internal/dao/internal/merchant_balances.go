// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantBalancesDao is the data access object for table mch_merchant_balances.
type MerchantBalancesDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns MerchantBalancesColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantBalancesColumns defines and stores column names for table mch_merchant_balances.
type MerchantBalancesColumns struct {
	Id               string // 主键
	MerchantId       string // 商家ID
	AvailableBalance string // 可提现余额（分）
	FreezeBalance    string // 冻结余额（分）
	Currency         string // 币种
	Version          string // 版本号（并发控制）
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// merchantBalancesColumns holds the columns for table mch_merchant_balances.
var merchantBalancesColumns = MerchantBalancesColumns{
	Id:               "id",
	MerchantId:       "merchant_id",
	AvailableBalance: "available_balance",
	FreezeBalance:    "freeze_balance",
	Currency:         "currency",
	Version:          "version",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewMerchantBalancesDao creates and returns a new DAO object for table data access.
func NewMerchantBalancesDao() *MerchantBalancesDao {
	return &MerchantBalancesDao{
		group:   "default",
		table:   "mch_merchant_balances",
		columns: merchantBalancesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantBalancesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantBalancesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantBalancesDao) Columns() MerchantBalancesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantBalancesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantBalancesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantBalancesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
