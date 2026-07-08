package service

import (
	"context"

	"gf-eshop/api/merchant_role_permissions/v1"
)

type IMerchantRolePermissions interface {
	List(ctx context.Context, req *v1.MerchantRolePermissionsListReq) (res *v1.MerchantRolePermissionsListRes, err error)
	Create(ctx context.Context, req *v1.MerchantRolePermissionsCreateReq) (res *v1.MerchantRolePermissionsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantRolePermissionsUpdateReq) (res *v1.MerchantRolePermissionsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantRolePermissionsDeleteReq) (res *v1.MerchantRolePermissionsDeleteRes, err error)
}

var localMerchantRolePermissions IMerchantRolePermissions

func MerchantRolePermissions() IMerchantRolePermissions {
	if localMerchantRolePermissions == nil {
		panic("implement not found for interface IMerchantRolePermissions, forgot register?")
	}
	return localMerchantRolePermissions
}

func RegisterMerchantRolePermissions(i IMerchantRolePermissions) {
	localMerchantRolePermissions = i
}
