package merchantRoles

import (
	"context"

	"gf-eshop/api/merchant_roles/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantRolesListReq) (res *v1.MerchantRolesListRes, err error) {
	return service.MerchantRoles().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.MerchantRolesDetailReq) (res *v1.MerchantRolesDetailRes, err error) {
	return service.MerchantRoles().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantRolesCreateReq) (res *v1.MerchantRolesCreateRes, err error) {
	return service.MerchantRoles().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantRolesUpdateReq) (res *v1.MerchantRolesUpdateRes, err error) {
	return service.MerchantRoles().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantRolesDeleteReq) (res *v1.MerchantRolesDeleteRes, err error) {
	return service.MerchantRoles().Delete(ctx, req)
}
