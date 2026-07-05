package orders

import (
	"gf-eshop/api/orders"
)

type ControllerV1 struct{}

func NewV1() orders.IOrdersV1 {
	return &ControllerV1{}
}