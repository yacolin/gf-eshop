package address

import (
	"context"

	"gf-eshop/api/address/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.AddressCreateReq) (res *v1.AddressCreateRes, err error) {
	return service.Address().Create(ctx, req)
}

func (c *ControllerV1) List(ctx context.Context, req *v1.AddressListReq) (res *v1.AddressListRes, err error) {
	return service.Address().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.AddressDetailReq) (res *v1.AddressDetailRes, err error) {
	return service.Address().Detail(ctx, req)
}

func (c *ControllerV1) GetDefault(ctx context.Context, req *v1.AddressGetDefaultReq) (res *v1.AddressGetDefaultRes, err error) {
	return service.Address().GetDefault(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.AddressUpdateReq) (res *v1.AddressUpdateRes, err error) {
	return service.Address().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.AddressDeleteReq) (res *v1.AddressDeleteRes, err error) {
	return service.Address().Delete(ctx, req)
}
