package attributes

import (
	"gf-eshop/api/attributes"
)

type ControllerV1 struct{}

func NewV1() attributes.IAttributesV1 {
	return &ControllerV1{}
}
