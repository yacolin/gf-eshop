package products

import (
	"context"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ProductsListReq) (res *v1.ProductsListRes, err error) {
	return service.Products().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.ProductsDetailReq) (res *v1.ProductsDetailRes, err error) {
	return service.Products().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error) {
	return service.Products().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error) {
	return service.Products().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error) {
	return service.Products().Delete(ctx, req)
}
