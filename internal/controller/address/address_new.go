package address

import (
	"gf-eshop/api/address"
)

type ControllerV1 struct{}

func NewV1() address.IAddressV1 {
	return &ControllerV1{}
}
