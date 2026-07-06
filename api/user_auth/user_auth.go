package user_auth

import (
	"context"
	"gf-eshop/api/user_auth/v1"
)

type IUserAuthV1 interface {
	Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error)
	Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error)
	RefreshToken(ctx context.Context, req *v1.UserRefreshTokenReq) (res *v1.UserRefreshTokenRes, err error)
}
