package categories

import (
	"gf-eshop/api/categories"
)

type ControllerV1 struct{}

func NewV1() categories.ICategoriesV1 {
	return &ControllerV1{}
}
