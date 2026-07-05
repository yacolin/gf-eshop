package skus

import (
	"context"

	"gf-eshop/api/skus/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.SkusListReq) (res *v1.SkusListRes, err error) {
	return service.Skus().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.SkusDetailReq) (res *v1.SkusDetailRes, err error) {
	return service.Skus().Detail(ctx, req)
}

func (c *ControllerV1) GetByCode(ctx context.Context, req *v1.SkusGetByCodeReq) (res *v1.SkusGetByCodeRes, err error) {
	return service.Skus().GetByCode(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.SkusCreateReq) (res *v1.SkusCreateRes, err error) {
	return service.Skus().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.SkusUpdateReq) (res *v1.SkusUpdateRes, err error) {
	return service.Skus().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.SkusDeleteReq) (res *v1.SkusDeleteRes, err error) {
	return service.Skus().Delete(ctx, req)
}
