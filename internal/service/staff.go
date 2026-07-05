package service

import (
	"context"
	"gf-eshop/api/staff/v1"
)

type IStaff interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error)
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error)
}

var localStaff IStaff

func Staff() IStaff {
	if localStaff == nil {
		panic("implement not found for interface IStaff, forgot register?")
	}
	return localStaff
}

func RegisterStaff(i IStaff) {
	localStaff = i
}
