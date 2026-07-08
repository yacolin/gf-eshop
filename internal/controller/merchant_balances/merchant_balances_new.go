package merchantBalances

import "gf-eshop/api/merchant_balances"

type ControllerV1 struct{}

func NewV1() merchantBalances.IMerchantBalancesV1 {
	return &ControllerV1{}
}
