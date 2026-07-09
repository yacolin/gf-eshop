package category_attributes

import (
	"context"

	"gf-eshop/api/category_attributes/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.CategoryAttributesListReq) (res *v1.CategoryAttributesListRes, err error) {
	return service.CategoryAttributes().List(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.CategoryAttributesCreateReq) (res *v1.CategoryAttributesCreateRes, err error) {
	return service.CategoryAttributes().Create(ctx, req)
}

func (c *ControllerV1) BatchCreate(ctx context.Context, req *v1.CategoryAttributesBatchCreateReq) (res *v1.CategoryAttributesBatchCreateRes, err error) {
	return service.CategoryAttributes().BatchCreate(ctx, req)
}

func (c *ControllerV1) ListByCat(ctx context.Context, req *v1.CategoryAttributesListByCatReq) (res *v1.CategoryAttributesListByCatRes, err error) {
	return service.CategoryAttributes().ListByCat(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.CategoryAttributesDeleteReq) (res *v1.CategoryAttributesDeleteRes, err error) {
	return service.CategoryAttributes().Delete(ctx, req)
}
