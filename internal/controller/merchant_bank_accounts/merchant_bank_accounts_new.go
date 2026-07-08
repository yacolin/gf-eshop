package merchantBankAccounts

import "gf-eshop/api/merchant_bank_accounts"

type ControllerV1 struct{}

func NewV1() merchantBankAccounts.IMerchantBankAccountsV1 {
	return &ControllerV1{}
}
