package productDescriptions

import (
	"gf-eshop/api/product_descriptions"
)

type ControllerV1 struct{}

func NewV1() productDescriptions.IProductDescriptionsV1 {
	return &ControllerV1{}
}
