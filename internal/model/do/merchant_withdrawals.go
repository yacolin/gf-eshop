// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWithdrawals is the golang structure of table mch_merchant_withdrawals for DAO operations like Where/Data.
type MerchantWithdrawals struct {
	g.Meta        `orm:"table:mch_merchant_withdrawals, do:true"`
	Id            interface{} // 主键
	MerchantId    interface{} // 商家ID
	WithdrawNo    interface{} // 提现单号
	Amount        interface{} // 提现金额（分）
	BankAccountId interface{} // 结算账户ID
	Status        interface{} // 0-待审核 1-审核通过 2-已打款 3-拒绝
	AuditRemark   interface{} // 审批备注
	AppliedAt     *gtime.Time // 申请时间
	ApprovedAt    *gtime.Time // 审批时间
	PaidAt        *gtime.Time // 打款时间
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
