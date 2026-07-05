package productVersions

import (
	"gf-eshop/api/product_versions"
)

type ControllerV1 struct{}

func NewV1() productVersions.IProductVersionsV1 {
	return &ControllerV1{}
}
