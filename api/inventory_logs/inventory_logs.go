package inventoryLogs

import (
	"context"

	"gf-eshop/api/inventory_logs/v1"
)

type IInventoryLogsV1 interface {
	List(ctx context.Context, req *v1.InventoryLogsListReq) (res *v1.InventoryLogsListRes, err error)
	Detail(ctx context.Context, req *v1.InventoryLogsDetailReq) (res *v1.InventoryLogsDetailRes, err error)
}
