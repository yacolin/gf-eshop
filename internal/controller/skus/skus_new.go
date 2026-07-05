package skus

import (
	"gf-eshop/api/skus"
)

type ControllerV1 struct{}

func NewV1() skus.ISkusV1 {
	return &ControllerV1{}
}
