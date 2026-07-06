package user_auth

import (
	"context"

	"gf-eshop/api/user_auth/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	return service.UserAuth().Login(ctx, req)
}

func (c *ControllerV1) Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error) {
	return service.UserAuth().Register(ctx, req)
}

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.UserRefreshTokenReq) (res *v1.UserRefreshTokenRes, err error) {
	return service.UserAuth().RefreshToken(ctx, req)
}
