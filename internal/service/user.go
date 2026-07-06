package service

import (
	"context"
	"gf-eshop/api/user/v1"
)

type IUser interface {
	Profile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error)
	UpdateInfo(ctx context.Context, req *v1.UserUpdateInfoReq) (res *v1.UserUpdateInfoRes, err error)
}

var localUser IUser

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
