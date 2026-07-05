package permissions

import (
	"gf-eshop/api/permissions"
)

type ControllerV1 struct{}

func NewV1() permissions.IPermissionsV1 {
	return &ControllerV1{}
}
