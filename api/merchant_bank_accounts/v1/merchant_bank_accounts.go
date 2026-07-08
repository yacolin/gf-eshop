package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type MerchantBankAccountsListReq struct {
	g.Meta `path:"/merchant-bank-accounts" tags:"MerchantBankAccounts" method:"get" summary:"银行账户列表"`

	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	MerchantId int64 `json:"merchant_id" description:"商家ID，可选筛选"`
}

type MerchantBankAccountsListRes struct {
	List  []*entity.MerchantBankAccounts `json:"list"`
	Total int                            `json:"total"`
}

// ---------- Detail ----------
type MerchantBankAccountsDetailReq struct {
	g.Meta `path:"/merchant-bank-accounts/{id}" tags:"MerchantBankAccounts" method:"get" summary:"银行账户详情"`
	Id     int64 `json:"id"`
}
type MerchantBankAccountsDetailRes struct {
	*entity.MerchantBankAccounts
}

// ---------- Create ----------
type MerchantBankAccountsCreateReq struct {
	g.Meta      `path:"/merchant-bank-accounts" tags:"MerchantBankAccounts" method:"post" summary:"新增银行账户"`

	MerchantId  int64  `json:"merchant_id"  v:"required" description:"商家ID"`
	BankName    string `json:"bank_name"    v:"required" description:"开户行"`
	BankBranch  string `json:"bank_branch"  description:"开户支行"`
	AccountName string `json:"account_name" v:"required" description:"开户名"`
	AccountNo   string `json:"account_no"   v:"required" description:"银行账号"`
	AccountType int    `json:"account_type" description:"1-对公账户 2-对私账户"`
	IsDefault   int    `json:"is_default"   description:"是否默认结算账户 0-否 1-是"`
	Status      int    `json:"status"       description:"1-正常 2-禁用"`
}

type MerchantBankAccountsCreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type MerchantBankAccountsUpdateReq struct {
	g.Meta      `path:"/merchant-bank-accounts/{id}" tags:"MerchantBankAccounts" method:"put" summary:"更新银行账户"`

	Id          int64  `json:"id"           v:"required"`
	BankName    string `json:"bank_name"    description:"开户行"`
	BankBranch  string `json:"bank_branch"  description:"开户支行"`
	AccountName string `json:"account_name" description:"开户名"`
	AccountNo   string `json:"account_no"   description:"银行账号"`
	AccountType int    `json:"account_type" description:"1-对公账户 2-对私账户"`
	IsDefault   int    `json:"is_default"   description:"是否默认结算账户 0-否 1-是"`
	Status      int    `json:"status"       description:"1-正常 2-禁用"`
}

type MerchantBankAccountsUpdateRes struct{}

// ---------- Delete ----------
type MerchantBankAccountsDeleteReq struct {
	g.Meta `path:"/merchant-bank-accounts/{id}" tags:"MerchantBankAccounts" method:"delete" summary:"删除银行账户"`
	Id     int64 `json:"id"`
}
type MerchantBankAccountsDeleteRes struct{}
