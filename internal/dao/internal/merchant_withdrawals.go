// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantWithdrawalsDao is the data access object for table mch_merchant_withdrawals.
type MerchantWithdrawalsDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns MerchantWithdrawalsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantWithdrawalsColumns defines and stores column names for table mch_merchant_withdrawals.
type MerchantWithdrawalsColumns struct {
	Id            string // 主键
	MerchantId    string // 商家ID
	WithdrawNo    string // 提现单号
	Amount        string // 提现金额（分）
	BankAccountId string // 结算账户ID
	Status        string // 0-待审核 1-审核通过 2-已打款 3-拒绝
	AuditRemark   string // 审批备注
	AppliedAt     string // 申请时间
	ApprovedAt    string // 审批时间
	PaidAt        string // 打款时间
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// merchantWithdrawalsColumns holds the columns for table mch_merchant_withdrawals.
var merchantWithdrawalsColumns = MerchantWithdrawalsColumns{
	Id:            "id",
	MerchantId:    "merchant_id",
	WithdrawNo:    "withdraw_no",
	Amount:        "amount",
	BankAccountId: "bank_account_id",
	Status:        "status",
	AuditRemark:   "audit_remark",
	AppliedAt:     "applied_at",
	ApprovedAt:    "approved_at",
	PaidAt:        "paid_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewMerchantWithdrawalsDao creates and returns a new DAO object for table data access.
func NewMerchantWithdrawalsDao() *MerchantWithdrawalsDao {
	return &MerchantWithdrawalsDao{
		group:   "default",
		table:   "mch_merchant_withdrawals",
		columns: merchantWithdrawalsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantWithdrawalsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantWithdrawalsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantWithdrawalsDao) Columns() MerchantWithdrawalsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantWithdrawalsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantWithdrawalsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantWithdrawalsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
