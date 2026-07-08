package merchantBalances

import (
	"context"

	"gf-eshop/api/merchant_balances/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) GetByMerchant(ctx context.Context, req *v1.GetByMerchantReq) (res *v1.GetByMerchantRes, err error) {
	return service.MerchantBalances().GetByMerchant(ctx, req)
}
