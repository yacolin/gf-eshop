package carts

import (
	"gf-eshop/api/carts"
)

type ControllerV1 struct{}

func NewV1() carts.ICartsV1 {
	return &ControllerV1{}
}
