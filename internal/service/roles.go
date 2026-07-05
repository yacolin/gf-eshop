package service

import (
	"context"

	"gf-eshop/api/roles/v1"
)

type IRoles interface {
	List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error)
	Detail(ctx context.Context, req *v1.RoleDetailReq) (res *v1.RoleDetailRes, err error)
	Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error)
	Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error)
	Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error)
	IsAdmin(ctx context.Context, staffId int64) (bool, error)
	HasPermission(ctx context.Context, staffId int64, permission string) (bool, error)
	GetPermissions(ctx context.Context, staffId int64) (permissions []string, roles []string, err error)
}

var localRoles IRoles

func Roles() IRoles {
	if localRoles == nil {
		panic("implement not found for interface IRoles, forgot register?")
	}
	return localRoles
}

func RegisterRoles(i IRoles) {
	localRoles = i
}
