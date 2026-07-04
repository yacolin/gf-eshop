package brands

import (
	"gf-eshop/api/brands"
)

type ControllerV1 struct{}

func NewV1() brands.IBrandsV1 {
	return &ControllerV1{}
}
