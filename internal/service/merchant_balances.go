package service

import (
	"context"

	"gf-eshop/api/merchant_balances/v1"
)

type IMerchantBalances interface {
	GetByMerchant(ctx context.Context, req *v1.GetByMerchantReq) (res *v1.GetByMerchantRes, err error)
}

var localMerchantBalances IMerchantBalances

func MerchantBalances() IMerchantBalances {
	if localMerchantBalances == nil {
		panic("implement not found for interface IMerchantBalances, forgot register?")
	}
	return localMerchantBalances
}

func RegisterMerchantBalances(i IMerchantBalances) {
	localMerchantBalances = i
}
