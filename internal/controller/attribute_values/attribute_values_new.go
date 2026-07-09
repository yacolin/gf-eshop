package attribute_values

import (
	"gf-eshop/api/attribute_values"
)

type ControllerV1 struct{}

func NewV1() attribute_values.IAttributeValuesV1 {
	return &ControllerV1{}
}
