package inventoryLogs

import (
	"gf-eshop/api/inventory_logs"
)

type ControllerV1 struct{}

func NewV1() inventoryLogs.IInventoryLogsV1 {
	return &ControllerV1{}
}
