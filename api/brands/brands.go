package brands

import (
	"context"

	"gf-eshop/api/brands/v1"
)

type IBrandsV1 interface {
	List(ctx context.Context, req *v1.BrandsListReq) (res *v1.BrandsListRes, err error)
	Detail(ctx context.Context, req *v1.BrandsDetailReq) (res *v1.BrandsDetailRes, err error)
	Create(ctx context.Context, req *v1.BrandsCreateReq) (res *v1.BrandsCreateRes, err error)
	Update(ctx context.Context, req *v1.BrandsUpdateReq) (res *v1.BrandsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.BrandsDeleteReq) (res *v1.BrandsDeleteRes, err error)
}
