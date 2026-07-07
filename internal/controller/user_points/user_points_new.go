package user_points

import "gf-eshop/api/user_points"

type ControllerV1 struct{}

func NewV1() user_points.IUserPointsV1 {
	return &ControllerV1{}
}