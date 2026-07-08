package service


import (
	"context"

	"gf-eshop/api/merchants/v1"
)


type IMerchants interface {
	List(ctx context.Context, req *v1.MerchantsListReq) (res *v1.MerchantsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantsDetailReq) (res *v1.MerchantsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantsCreateReq) (res *v1.MerchantsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantsUpdateReq) (res *v1.MerchantsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantsDeleteReq) (res *v1.MerchantsDeleteRes, err error)
}

var localMerchants IMerchants

func Merchants() IMerchants {
	if localMerchants == nil {
		panic("implement not found for interface IMerchants, forgot register?")
	}
	return localMerchants
}

func RegisterMerchants(i IMerchants) {
	localMerchants = i
}