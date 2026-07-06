package user_auth

import (
	"gf-eshop/api/user_auth"
)

type ControllerV1 struct{}

func NewV1() user_auth.IUserAuthV1 {
	return &ControllerV1{}
}
