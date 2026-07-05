package roles

import (
	"context"

	"gf-eshop/api/roles/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	return service.Roles().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.RoleDetailReq) (res *v1.RoleDetailRes, err error) {
	return service.Roles().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	return service.Roles().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	return service.Roles().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	return service.Roles().Delete(ctx, req)
}
