package products

import (
	"gf-eshop/api/products"
)

type ControllerV1 struct{}

func NewV1() products.IProductsV1 {
	return &ControllerV1{}
}
