package productAttributes

import (
	"context"

	"gf-eshop/api/product_attributes/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ProductAttributesListReq) (res *v1.ProductAttributesListRes, err error) {
	return service.ProductAttributes().List(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.ProductAttributesCreateReq) (res *v1.ProductAttributesCreateRes, err error) {
	return service.ProductAttributes().Create(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.ProductAttributesDeleteReq) (res *v1.ProductAttributesDeleteRes, err error) {
	return service.ProductAttributes().Delete(ctx, req)
}
