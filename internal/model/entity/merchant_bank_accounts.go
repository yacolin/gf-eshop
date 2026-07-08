// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBankAccounts is the golang structure for table merchant_bank_accounts.
type MerchantBankAccounts struct {
	Id          int64       `json:"id"           description:"主键"`
	MerchantId  int64       `json:"merchant_id"  description:"商家ID"`
	BankName    string      `json:"bank_name"    description:"开户行"`
	BankBranch  string      `json:"bank_branch"  description:"开户支行"`
	AccountName string      `json:"account_name" description:"开户名"`
	AccountNo   string      `json:"account_no"   description:"银行账号"`
	AccountType int         `json:"account_type" description:"1-对公账户 2-对私账户"`
	IsDefault   int         `json:"is_default"   description:"是否默认结算账户 0-否 1-是"`
	Status      int         `json:"status"       description:"1-正常 2-禁用"`
	CreatedAt   *gtime.Time `json:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:""`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:""`
}
