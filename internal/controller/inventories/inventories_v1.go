package inventories

import (
	"context"

	"gf-eshop/api/inventories/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.InventoriesListReq) (res *v1.InventoriesListRes, err error) {
	return service.Inventories().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.InventoriesDetailReq) (res *v1.InventoriesDetailRes, err error) {
	return service.Inventories().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.InventoriesCreateReq) (res *v1.InventoriesCreateRes, err error) {
	return service.Inventories().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.InventoriesUpdateReq) (res *v1.InventoriesUpdateRes, err error) {
	return service.Inventories().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.InventoriesDeleteReq) (res *v1.InventoriesDeleteRes, err error) {
	return service.Inventories().Delete(ctx, req)
}

func (c *ControllerV1) Lock(ctx context.Context, req *v1.InventoriesLockReq) (res *v1.InventoriesLockRes, err error) {
	return service.Inventories().Lock(ctx, req)
}

func (c *ControllerV1) Unlock(ctx context.Context, req *v1.InventoriesUnlockReq) (res *v1.InventoriesUnlockRes, err error) {
	return service.Inventories().Unlock(ctx, req)
}

func (c *ControllerV1) Deduct(ctx context.Context, req *v1.InventoriesDeductReq) (res *v1.InventoriesDeductRes, err error) {
	return service.Inventories().Deduct(ctx, req)
}

func (c *ControllerV1) Restock(ctx context.Context, req *v1.InventoriesRestockReq) (res *v1.InventoriesRestockRes, err error) {
	return service.Inventories().Restock(ctx, req)
}

func (c *ControllerV1) GetStock(ctx context.Context, req *v1.InventoriesGetStockReq) (res *v1.InventoriesGetStockRes, err error) {
	return service.Inventories().GetStock(ctx, req)
}
