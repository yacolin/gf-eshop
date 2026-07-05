package service

import (
	"context"

	"gf-eshop/api/skus/v1"
)

type ISkus interface {
	List(ctx context.Context, req *v1.SkusListReq) (res *v1.SkusListRes, err error)
	Detail(ctx context.Context, req *v1.SkusDetailReq) (res *v1.SkusDetailRes, err error)
	GetByCode(ctx context.Context, req *v1.SkusGetByCodeReq) (res *v1.SkusGetByCodeRes, err error)
	Create(ctx context.Context, req *v1.SkusCreateReq) (res *v1.SkusCreateRes, err error)
	Update(ctx context.Context, req *v1.SkusUpdateReq) (res *v1.SkusUpdateRes, err error)
	Delete(ctx context.Context, req *v1.SkusDeleteReq) (res *v1.SkusDeleteRes, err error)
}

var localSkus ISkus

func Skus() ISkus {
	if localSkus == nil {
		panic("implement not found for interface ISkus, forgot register?")
	}
	return localSkus
}

func RegisterSkus(i ISkus) {
	localSkus = i
}
