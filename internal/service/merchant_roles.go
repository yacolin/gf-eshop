package service

import (
	"context"

	"gf-eshop/api/merchant_roles/v1"
)

type IMerchantRoles interface {
	List(ctx context.Context, req *v1.MerchantRolesListReq) (res *v1.MerchantRolesListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantRolesDetailReq) (res *v1.MerchantRolesDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantRolesCreateReq) (res *v1.MerchantRolesCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantRolesUpdateReq) (res *v1.MerchantRolesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantRolesDeleteReq) (res *v1.MerchantRolesDeleteRes, err error)
}

var localMerchantRoles IMerchantRoles

func MerchantRoles() IMerchantRoles {
	if localMerchantRoles == nil {
		panic("implement not found for interface IMerchantRoles, forgot register?")
	}
	return localMerchantRoles
}

func RegisterMerchantRoles(i IMerchantRoles) {
	localMerchantRoles = i
}
