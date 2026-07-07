package user_admin

import (
	"context"

	"gf-eshop/api/user_admin/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	return service.User().List(ctx, req)
}
