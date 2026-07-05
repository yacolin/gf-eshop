package service

import (
	"context"

	"gf-eshop/api/product_versions/v1"
)

type IProductVersions interface {
	List(ctx context.Context, req *v1.ProductVersionsListReq) (res *v1.ProductVersionsListRes, err error)
	Detail(ctx context.Context, req *v1.ProductVersionsDetailReq) (res *v1.ProductVersionsDetailRes, err error)
}

var localProductVersions IProductVersions

func ProductVersions() IProductVersions {
	if localProductVersions == nil {
		panic("implement not found for interface IProductVersions, forgot register?")
	}
	return localProductVersions
}

func RegisterProductVersions(i IProductVersions) {
	localProductVersions = i
}
