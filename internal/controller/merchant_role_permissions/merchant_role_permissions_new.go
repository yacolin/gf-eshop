package merchantRolePermissions

import "gf-eshop/api/merchant_role_permissions"

type ControllerV1 struct{}

func NewV1() merchantRolePermissions.IMerchantRolePermissionsV1 {
	return &ControllerV1{}
}
