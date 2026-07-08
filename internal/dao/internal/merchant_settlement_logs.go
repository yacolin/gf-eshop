// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantSettlementLogsDao is the data access object for table mch_merchant_settlement_logs.
type MerchantSettlementLogsDao struct {
	table   string                        // table is the underlying table name of the DAO.
	group   string                        // group is the database configuration group name of current DAO.
	columns MerchantSettlementLogsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantSettlementLogsColumns defines and stores column names for table mch_merchant_settlement_logs.
type MerchantSettlementLogsColumns struct {
	Id               string // 主键
	MerchantId       string // 商家ID
	SettlementNo     string // 结算单号
	SettlementCycle  string // 结算周期（如 2026-07-01~2026-07-15）
	TotalAmount      string // 期内总金额（分）
	CommissionAmount string // 平台佣金（分）
	SettlementAmount string // 应结算金额（分）
	Status           string // 0-待结算 1-已结算 2-已打款
	SettledAt        string // 结算时间
	PaidAt           string // 打款时间
	Remark           string // 备注
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// merchantSettlementLogsColumns holds the columns for table mch_merchant_settlement_logs.
var merchantSettlementLogsColumns = MerchantSettlementLogsColumns{
	Id:               "id",
	MerchantId:       "merchant_id",
	SettlementNo:     "settlement_no",
	SettlementCycle:  "settlement_cycle",
	TotalAmount:      "total_amount",
	CommissionAmount: "commission_amount",
	SettlementAmount: "settlement_amount",
	Status:           "status",
	SettledAt:        "settled_at",
	PaidAt:           "paid_at",
	Remark:           "remark",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewMerchantSettlementLogsDao creates and returns a new DAO object for table data access.
func NewMerchantSettlementLogsDao() *MerchantSettlementLogsDao {
	return &MerchantSettlementLogsDao{
		group:   "default",
		table:   "mch_merchant_settlement_logs",
		columns: merchantSettlementLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantSettlementLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantSettlementLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantSettlementLogsDao) Columns() MerchantSettlementLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantSettlementLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantSettlementLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantSettlementLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
