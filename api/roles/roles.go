package roles

import (
	"context"

	"gf-eshop/api/roles/v1"
)

type IRolesV1 interface {
	List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error)
	Detail(ctx context.Context, req *v1.RoleDetailReq) (res *v1.RoleDetailRes, err error)
	Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error)
	Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error)
	Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error)
}
