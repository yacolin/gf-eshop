package user

import (
	"context"

	"gf-eshop/api/user/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Profile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	return service.User().Profile(ctx, req)
}

func (c *ControllerV1) UpdateInfo(ctx context.Context, req *v1.UserUpdateInfoReq) (res *v1.UserUpdateInfoRes, err error) {
	return service.User().UpdateInfo(ctx, req)
}
