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

func (c *ControllerV1) DetailPure(ctx context.Context, req *v1.ProductsDetailPureReq) (res *v1.ProductsDetailPureRes, err error) {
	return service.Products().DetailPure(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error) {
	return service.Products().Create(ctx, req)
}

func (c *ControllerV1) CreateFull(ctx context.Context, req *v1.ProductsCreateFullReq) (res *v1.ProductsCreateFullRes, err error) {
	return service.Products().CreateFull(ctx, req)
}


func (c *ControllerV1) GetAttributes(ctx context.Context, req *v1.ProductsGetAttributesReq) (res *v1.ProductsGetAttributesRes, err error) {
	return service.Products().GetAttributes(ctx, req)
}

func (c *ControllerV1) UpdateAttributes(ctx context.Context, req *v1.ProductsUpdateAttributesReq) (res *v1.ProductsUpdateAttributesRes, err error) {
	return service.Products().UpdateAttributes(ctx, req)
}

func (c *ControllerV1) EnrichedDetail(ctx context.Context, req *v1.ProductsEnrichedDetailReq) (res *v1.ProductsEnrichedDetailRes, err error) {
	return service.Products().EnrichedDetail(ctx, req)
}

func (c *ControllerV1) BatchCreateSKUs(ctx context.Context, req *v1.ProductsBatchCreateSKUsReq) (res *v1.ProductsBatchCreateSKUsRes, err error) {
	return service.Products().BatchCreateSKUs(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error) {
	return service.Products().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error) {
	return service.Products().Delete(ctx, req)
}
