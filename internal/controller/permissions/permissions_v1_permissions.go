package permissions

import (
	"context"

	"gf-eshop/api/permissions/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.PermissionListReq) (res *v1.PermissionListRes, err error) {
	return service.Permissions().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.PermissionDetailReq) (res *v1.PermissionDetailRes, err error) {
	return service.Permissions().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.PermissionCreateReq) (res *v1.PermissionCreateRes, err error) {
	return service.Permissions().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.PermissionUpdateReq) (res *v1.PermissionUpdateRes, err error) {
	return service.Permissions().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.PermissionDeleteReq) (res *v1.PermissionDeleteRes, err error) {
	return service.Permissions().Delete(ctx, req)
}

func (c *ControllerV1) Check(ctx context.Context, req *v1.PermissionCheckReq) (res *v1.PermissionCheckRes, err error) {
	return service.Permissions().Check(ctx, req)
}

func (c *ControllerV1) RolePermissionList(ctx context.Context, req *v1.RolePermissionListReq) (res *v1.RolePermissionListRes, err error) {
	return service.Permissions().RolePermissionList(ctx, req)
}

func (c *ControllerV1) RolePermissionUpdate(ctx context.Context, req *v1.RolePermissionUpdateReq) (res *v1.RolePermissionUpdateRes, err error) {
	return service.Permissions().RolePermissionUpdate(ctx, req)
}
