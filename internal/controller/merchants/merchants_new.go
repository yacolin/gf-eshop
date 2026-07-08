package merchants

import "gf-eshop/api/merchants"

type ControllerV1 struct{}

func NewV1() merchants.IMerchantsV1 {
	return &ControllerV1{}
}