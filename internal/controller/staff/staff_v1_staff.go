package staff

import (
	"context"

	"gf-eshop/api/staff/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.StaffLoginReq) (res *v1.StaffLoginRes, err error) {
	return service.Staff().Login(ctx, req)
}

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.StaffRefreshTokenReq) (res *v1.StaffRefreshTokenRes, err error) {
	return service.Staff().RefreshToken(ctx, req)
}

func (c *ControllerV1) Logout(ctx context.Context, req *v1.StaffLogoutReq) (res *v1.StaffLogoutRes, err error) {
	return service.Staff().Logout(ctx, req)
}

func (c *ControllerV1) Profile(ctx context.Context, req *v1.StaffProfileReq) (res *v1.StaffProfileRes, err error) {
	return service.Staff().Profile(ctx, req)
}

func (c *ControllerV1) Permissions(ctx context.Context, req *v1.StaffPermissionsReq) (res *v1.StaffPermissionsRes, err error) {
	return service.Staff().Permissions(ctx, req)
}

func (c *ControllerV1) List(ctx context.Context, req *v1.StaffListReq) (res *v1.StaffListRes, err error) {
	return service.Staff().List(ctx, req)
}

func (c *ControllerV1) AssignRoles(ctx context.Context, req *v1.StaffAssignRolesReq) (res *v1.StaffAssignRolesRes, err error) {
	return service.Staff().AssignRoles(ctx, req)
}
