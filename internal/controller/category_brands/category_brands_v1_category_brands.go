package brands

import (
	"context"

	"gf-eshop/api/category_brands/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.CategoryBrandListReq) (res *v1.CategoryBrandListRes, err error) {
	return service.CategoryBrands().List(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.CategoryBrandUpdateReq) (res *v1.CategoryBrandUpdateRes, err error) {
	return service.CategoryBrands().Update(ctx, req)
}

