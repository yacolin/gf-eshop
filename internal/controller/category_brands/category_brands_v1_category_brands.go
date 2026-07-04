package brands

import (
	"context"

	"gf-eshop/api/category_brands/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	return service.CategoryBrands().List(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	return service.CategoryBrands().Update(ctx, req)
}

