package category_attributes

import (
	"gf-eshop/api/category_attributes"
)

type ControllerV1 struct{}

func NewV1() category_attributes.ICategoryAttributesV1 {
	return &ControllerV1{}
}
