package service

import (
	"context"

	"gf-eshop/api/products/v1"
)

type IProducts interface {
	List(ctx context.Context, req *v1.ProductsListReq) (res *v1.ProductsListRes, err error)
	Detail(ctx context.Context, req *v1.ProductsDetailReq) (res *v1.ProductsDetailRes, err error)
	DetailPure(ctx context.Context, req *v1.ProductsDetailPureReq) (res *v1.ProductsDetailPureRes, err error)
	Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error)
	CreateFull(ctx context.Context, req *v1.ProductsCreateFullReq) (res *v1.ProductsCreateFullRes, err error)
	Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error)
}

var localProducts IProducts

func Products() IProducts {
	if localProducts == nil {
		panic("implement not found for interface IProducts, forgot register?")
	}
	return localProducts
}

func RegisterProducts(i IProducts) {
	localProducts = i
}
