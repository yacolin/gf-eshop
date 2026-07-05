package brands

import (
	"context"

	"gf-eshop/api/brands/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.BrandsListReq) (res *v1.BrandsListRes, err error) {
	return service.Brands().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.BrandsDetailReq) (res *v1.BrandsDetailRes, err error) {
	return service.Brands().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.BrandsCreateReq) (res *v1.BrandsCreateRes, err error) {
	return service.Brands().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.BrandsUpdateReq) (res *v1.BrandsUpdateRes, err error) {
	return service.Brands().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.BrandsDeleteReq) (res *v1.BrandsDeleteRes, err error) {
	return service.Brands().Delete(ctx, req)
}
