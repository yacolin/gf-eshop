package skus

import (
	"context"

	"gf-eshop/api/skus/v1"
)

type ISkusV1 interface {
	List(ctx context.Context, req *v1.SkusListReq) (res *v1.SkusListRes, err error)
	Detail(ctx context.Context, req *v1.SkusDetailReq) (res *v1.SkusDetailRes, err error)
	GetByCode(ctx context.Context, req *v1.SkusGetByCodeReq) (res *v1.SkusGetByCodeRes, err error)
	Create(ctx context.Context, req *v1.SkusCreateReq) (res *v1.SkusCreateRes, err error)
	Update(ctx context.Context, req *v1.SkusUpdateReq) (res *v1.SkusUpdateRes, err error)
	Delete(ctx context.Context, req *v1.SkusDeleteReq) (res *v1.SkusDeleteRes, err error)
}
