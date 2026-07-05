package products

import (
	"context"

	"gf-eshop/api/products/v1"
)

type IProductsV1 interface {
	List(ctx context.Context, req *v1.ProductsListReq) (res *v1.ProductsListRes, err error)
	Detail(ctx context.Context, req *v1.ProductsDetailReq) (res *v1.ProductsDetailRes, err error)
	DetailPure(ctx context.Context, req *v1.ProductsDetailPureReq) (res *v1.ProductsDetailPureRes, err error)
	Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error)
	CreateFull(ctx context.Context, req *v1.ProductsCreateFullReq) (res *v1.ProductsCreateFullRes, err error)
	Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error)
}
