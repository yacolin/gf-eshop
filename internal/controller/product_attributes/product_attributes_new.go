package productAttributes

import (
	"gf-eshop/api/product_attributes"
)

type ControllerV1 struct{}

func NewV1() productAttributes.IProductAttributesV1 {
	return &ControllerV1{}
}
