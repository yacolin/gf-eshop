package service

import (
	"context"

	"gf-eshop/api/inventory_logs/v1"
)

type IInventoryLogs interface {
	List(ctx context.Context, req *v1.InventoryLogsListReq) (res *v1.InventoryLogsListRes, err error)
	Detail(ctx context.Context, req *v1.InventoryLogsDetailReq) (res *v1.InventoryLogsDetailRes, err error)
}

var localInventoryLogs IInventoryLogs

func InventoryLogs() IInventoryLogs {
	if localInventoryLogs == nil {
		panic("implement not found for interface IInventoryLogs, forgot register?")
	}
	return localInventoryLogs
}

func RegisterInventoryLogs(i IInventoryLogs) {
	localInventoryLogs = i
}
