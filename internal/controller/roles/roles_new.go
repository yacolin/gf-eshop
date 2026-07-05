package roles

import (
	"gf-eshop/api/roles"
)

type ControllerV1 struct{}

func NewV1() roles.IRolesV1 {
	return &ControllerV1{}
}
