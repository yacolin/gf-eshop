// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWithdrawals is the golang structure for table merchant_withdrawals.
type MerchantWithdrawals struct {
	Id            int64       `json:"id"              description:"主键"`
	MerchantId    int64       `json:"merchant_id"     description:"商家ID"`
	WithdrawNo    string      `json:"withdraw_no"     description:"提现单号"`
	Amount        int64       `json:"amount"          description:"提现金额（分）"`
	BankAccountId int64       `json:"bank_account_id" description:"结算账户ID"`
	Status        int         `json:"status"          description:"0-待审核 1-审核通过 2-已打款 3-拒绝"`
	AuditRemark   string      `json:"audit_remark"    description:"审批备注"`
	AppliedAt     *gtime.Time `json:"applied_at"      description:"申请时间"`
	ApprovedAt    *gtime.Time `json:"approved_at"     description:"审批时间"`
	PaidAt        *gtime.Time `json:"paid_at"         description:"打款时间"`
	CreatedAt     *gtime.Time `json:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updated_at"      description:""`
	DeletedAt     *gtime.Time `json:"deleted_at"      description:""`
}
