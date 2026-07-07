package user_admin

import (
	"context"
	"gf-eshop/api/user_admin/v1"
)

type IUserAdminV1 interface {
	List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)
}
