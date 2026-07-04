package brands

import (
	"gf-eshop/api/category_brands"
)

type ControllerV1 struct{}

func NewV1() category_brands.ICategoryBrandsV1 {
	return &ControllerV1{}
}