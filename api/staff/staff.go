package staff

import (
	"context"
	"gf-eshop/api/staff/v1"
)

type IStaffV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error)
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error)
}
