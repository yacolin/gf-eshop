package merchants

import (
	"context"

	"gf-eshop/api/merchants/v1"
	"gf-eshop/internal/service"
)



func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantsListReq) (res *v1.MerchantsListRes, err error) {
	return service.Merchants().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.MerchantsDetailReq) (res *v1.MerchantsDetailRes, err error) {
	return service.Merchants().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantsCreateReq) (res *v1.MerchantsCreateRes, err error) {
	return service.Merchants().Create(ctx, req)
}


func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantsUpdateReq) (res *v1.MerchantsUpdateRes, err error) {
	return service.Merchants().Update(ctx, req)
}


func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantsDeleteReq) (res *v1.MerchantsDeleteRes, err error) {
	return service.Merchants().Delete(ctx, req)
}