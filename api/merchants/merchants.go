package merchants

import (
	"context"

	"gf-eshop/api/merchants/v1"
)

type IMerchantsV1 interface {
	List(ctx context.Context, req *v1.MerchantsListReq) (res *v1.MerchantsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantsDetailReq) (res *v1.MerchantsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantsCreateReq) (res *v1.MerchantsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantsUpdateReq) (res *v1.MerchantsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantsDeleteReq) (res *v1.MerchantsDeleteRes, err error)
	
}