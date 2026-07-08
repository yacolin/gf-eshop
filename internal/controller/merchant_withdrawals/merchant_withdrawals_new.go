package merchantWithdrawals

import "gf-eshop/api/merchant_withdrawals"

type ControllerV1 struct{}

func NewV1() merchantWithdrawals.IMerchantWithdrawalsV1 {
	return &ControllerV1{}
}
