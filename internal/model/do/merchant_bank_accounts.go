// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBankAccounts is the golang structure of table mch_merchant_bank_accounts for DAO operations like Where/Data.
type MerchantBankAccounts struct {
	g.Meta      `orm:"table:mch_merchant_bank_accounts, do:true"`
	Id          interface{} // 主键
	MerchantId  interface{} // 商家ID
	BankName    interface{} // 开户行
	BankBranch  interface{} // 开户支行
	AccountName interface{} // 开户名
	AccountNo   interface{} // 银行账号
	AccountType interface{} // 1-对公账户 2-对私账户
	IsDefault   interface{} // 是否默认结算账户（NULL=非默认, 1=默认）
	Status      interface{} // 1-正常 2-禁用
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
