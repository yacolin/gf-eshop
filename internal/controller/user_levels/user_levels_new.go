package user_levels

import "gf-eshop/api/user_levels"

type ControllerV1 struct{}

func NewV1() user_levels.IUserLevelsV1 {
	return &ControllerV1{}
}
