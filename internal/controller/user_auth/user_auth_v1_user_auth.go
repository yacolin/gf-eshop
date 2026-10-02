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

func (c *ControllerV1) SendEmailCode(ctx context.Context, req *v1.UserSendEmailCodeReq) (res *v1.UserSendEmailCodeRes, err error) {
	return service.UserAuth().SendEmailCode(ctx, req)
}

func (c *ControllerV1) EmailLogin(ctx context.Context, req *v1.UserEmailLoginReq) (res *v1.UserEmailLoginRes, err error) {
	return service.UserAuth().EmailLogin(ctx, req)
}

func (c *ControllerV1) ResetPassword(ctx context.Context, req *v1.UserResetPasswordReq) (res *v1.UserResetPasswordRes, err error) {
	return service.UserAuth().ResetPassword(ctx, req)
}
