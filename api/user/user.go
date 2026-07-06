package user

import (
	"context"
	"gf-eshop/api/user/v1"
)

type IUserV1 interface {
	Profile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error)
	UpdateInfo(ctx context.Context, req *v1.UserUpdateInfoReq) (res *v1.UserUpdateInfoRes, err error)
}
