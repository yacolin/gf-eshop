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
