package service

import (
	"context"

	"gf-eshop/api/inventories/v1"
)

type IWarehouses interface {
	List(ctx context.Context, req *v1.WarehousesListReq) (res *v1.WarehousesListRes, err error)
	Detail(ctx context.Context, req *v1.WarehousesDetailReq) (res *v1.WarehousesDetailRes, err error)
	Create(ctx context.Context, req *v1.WarehousesCreateReq) (res *v1.WarehousesCreateRes, err error)
	Update(ctx context.Context, req *v1.WarehousesUpdateReq) (res *v1.WarehousesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.WarehousesDeleteReq) (res *v1.WarehousesDeleteRes, err error)
}

type IInventories interface {
	List(ctx context.Context, req *v1.InventoriesListReq) (res *v1.InventoriesListRes, err error)
	Detail(ctx context.Context, req *v1.InventoriesDetailReq) (res *v1.InventoriesDetailRes, err error)
	Create(ctx context.Context, req *v1.InventoriesCreateReq) (res *v1.InventoriesCreateRes, err error)
	Update(ctx context.Context, req *v1.InventoriesUpdateReq) (res *v1.InventoriesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.InventoriesDeleteReq) (res *v1.InventoriesDeleteRes, err error)
	Lock(ctx context.Context, req *v1.InventoriesLockReq) (res *v1.InventoriesLockRes, err error)
	Unlock(ctx context.Context, req *v1.InventoriesUnlockReq) (res *v1.InventoriesUnlockRes, err error)
	Deduct(ctx context.Context, req *v1.InventoriesDeductReq) (res *v1.InventoriesDeductRes, err error)
	Restock(ctx context.Context, req *v1.InventoriesRestockReq) (res *v1.InventoriesRestockRes, err error)
	GetStock(ctx context.Context, req *v1.InventoriesGetStockReq) (res *v1.InventoriesGetStockRes, err error)
	Alerts(ctx context.Context, req *v1.InventoriesAlertsReq) (res *v1.InventoriesAlertsRes, err error)
	AlertResolve(ctx context.Context, req *v1.InventoriesAlertResolveReq) (res *v1.InventoriesAlertResolveRes, err error)
	Export(ctx context.Context, req *v1.InventoriesExportReq) (res *v1.InventoriesExportRes, err error)
}

var (
	localInventories IInventories
	localWarehouses  IWarehouses
)

func Inventories() IInventories {
	if localInventories == nil {
		panic("implement not found for interface IInventories, forgot register?")
	}
	return localInventories
}

func RegisterInventories(i IInventories) {
	localInventories = i
}

func Warehouses() IWarehouses {
	if localWarehouses == nil {
		panic("implement not found for interface IWarehouses, forgot register?")
	}
	return localWarehouses
}

func RegisterWarehouses(i IWarehouses) {
	localWarehouses = i
}
