package staff

import (
	"context"

	"gf-eshop/api/staff/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	return service.Staff().Login(ctx, req)
}

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	return service.Staff().RefreshToken(ctx, req)
}

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	return service.Staff().Logout(ctx, req)
}

func (c *ControllerV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	return service.Staff().Profile(ctx, req)
}
