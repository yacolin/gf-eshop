package service

import (
	"context"
	"gf-eshop/api/user_auth/v1"
)

type IUserAuth interface {
	Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error)
	Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error)
	RefreshToken(ctx context.Context, req *v1.UserRefreshTokenReq) (res *v1.UserRefreshTokenRes, err error)
	SendEmailCode(ctx context.Context, req *v1.UserSendEmailCodeReq) (res *v1.UserSendEmailCodeRes, err error)
	EmailLogin(ctx context.Context, req *v1.UserEmailLoginReq) (res *v1.UserEmailLoginRes, err error)
	ResetPassword(ctx context.Context, req *v1.UserResetPasswordReq) (res *v1.UserResetPasswordRes, err error)
}

var localUserAuth IUserAuth

func UserAuth() IUserAuth {
	if localUserAuth == nil {
		panic("implement not found for interface IUserAuth, forgot register?")
	}
	return localUserAuth
}

func RegisterUserAuth(i IUserAuth) {
	localUserAuth = i
}
