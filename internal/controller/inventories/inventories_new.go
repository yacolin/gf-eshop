package inventories

import (
	"gf-eshop/api/inventories"
)

type ControllerV1 struct{}

func NewV1() inventories.IInventoriesV1 {
	return &ControllerV1{}
}
