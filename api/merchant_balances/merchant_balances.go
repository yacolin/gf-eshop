package merchantBalances

import (
	"context"

	"gf-eshop/api/merchant_balances/v1"
)

type IMerchantBalancesV1 interface {
	GetByMerchant(ctx context.Context, req *v1.GetByMerchantReq) (res *v1.GetByMerchantRes, err error)
}
