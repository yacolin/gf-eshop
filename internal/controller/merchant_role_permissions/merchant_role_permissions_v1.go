package merchantRolePermissions

import (
	"context"

	"gf-eshop/api/merchant_role_permissions/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantRolePermissionsListReq) (res *v1.MerchantRolePermissionsListRes, err error) {
	return service.MerchantRolePermissions().List(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantRolePermissionsCreateReq) (res *v1.MerchantRolePermissionsCreateRes, err error) {
	return service.MerchantRolePermissions().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantRolePermissionsUpdateReq) (res *v1.MerchantRolePermissionsUpdateRes, err error) {
	return service.MerchantRolePermissions().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantRolePermissionsDeleteReq) (res *v1.MerchantRolePermissionsDeleteRes, err error) {
	return service.MerchantRolePermissions().Delete(ctx, req)
}
