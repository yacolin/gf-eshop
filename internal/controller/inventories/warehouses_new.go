package inventories

import (
	"gf-eshop/api/inventories"
)

type ControllerWarehousesV1 struct{}

func NewWarehousesV1() inventories.Warehouses {
	return &ControllerWarehousesV1{}
}
