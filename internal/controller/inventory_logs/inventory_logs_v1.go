package inventoryLogs

import (
	"context"

	"gf-eshop/api/inventory_logs/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.InventoryLogsListReq) (res *v1.InventoryLogsListRes, err error) {
	return service.InventoryLogs().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.InventoryLogsDetailReq) (res *v1.InventoryLogsDetailRes, err error) {
	return service.InventoryLogs().Detail(ctx, req)
}
