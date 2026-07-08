package merchantRolePermissions

import (
	"context"

	"gf-eshop/api/merchant_role_permissions/v1"
)

type IMerchantRolePermissionsV1 interface {
	List(ctx context.Context, req *v1.MerchantRolePermissionsListReq) (res *v1.MerchantRolePermissionsListRes, err error)
	Create(ctx context.Context, req *v1.MerchantRolePermissionsCreateReq) (res *v1.MerchantRolePermissionsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantRolePermissionsUpdateReq) (res *v1.MerchantRolePermissionsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantRolePermissionsDeleteReq) (res *v1.MerchantRolePermissionsDeleteRes, err error)
}
