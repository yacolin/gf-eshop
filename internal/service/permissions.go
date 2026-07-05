package service

import (
	"context"

	"gf-eshop/api/permissions/v1"
)

type IPermissions interface {
	List(ctx context.Context, req *v1.PermissionListReq) (res *v1.PermissionListRes, err error)
	Detail(ctx context.Context, req *v1.PermissionDetailReq) (res *v1.PermissionDetailRes, err error)
	Create(ctx context.Context, req *v1.PermissionCreateReq) (res *v1.PermissionCreateRes, err error)
	Update(ctx context.Context, req *v1.PermissionUpdateReq) (res *v1.PermissionUpdateRes, err error)
	Delete(ctx context.Context, req *v1.PermissionDeleteReq) (res *v1.PermissionDeleteRes, err error)
	Check(ctx context.Context, req *v1.PermissionCheckReq) (res *v1.PermissionCheckRes, err error)
	RolePermissionList(ctx context.Context, req *v1.RolePermissionListReq) (res *v1.RolePermissionListRes, err error)
	RolePermissionUpdate(ctx context.Context, req *v1.RolePermissionUpdateReq) (res *v1.RolePermissionUpdateRes, err error)
}

var localPermissions IPermissions

func Permissions() IPermissions {
	if localPermissions == nil {
		panic("implement not found for interface IPermissions, forgot register?")
	}
	return localPermissions
}

func RegisterPermissions(i IPermissions) {
	localPermissions = i
}
