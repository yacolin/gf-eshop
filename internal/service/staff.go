package service

import (
	"context"
	"gf-eshop/api/staff/v1"
)

type IStaff interface {
	Login(ctx context.Context, req *v1.StaffLoginReq) (res *v1.StaffLoginRes, err error)
	RefreshToken(ctx context.Context, req *v1.StaffRefreshTokenReq) (res *v1.StaffRefreshTokenRes, err error)
	Logout(ctx context.Context, req *v1.StaffLogoutReq) (res *v1.StaffLogoutRes, err error)
	Profile(ctx context.Context, req *v1.StaffProfileReq) (res *v1.StaffProfileRes, err error)
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
