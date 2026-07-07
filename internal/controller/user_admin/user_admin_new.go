package user_admin

import (
	"gf-eshop/api/user_admin"
)

type ControllerV1 struct{}

func NewV1() user_admin.IUserAdminV1 {
	return &ControllerV1{}
}
