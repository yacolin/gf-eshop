package merchantRoles

import "gf-eshop/api/merchant_roles"

type ControllerV1 struct{}

func NewV1() merchantRoles.IMerchantRolesV1 {
	return &ControllerV1{}
}
