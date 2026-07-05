package inventories

import (
	"context"

	"gf-eshop/api/inventories/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerWarehousesV1) List(ctx context.Context, req *v1.WarehousesListReq) (res *v1.WarehousesListRes, err error) {
	return service.Warehouses().List(ctx, req)
}

func (c *ControllerWarehousesV1) Detail(ctx context.Context, req *v1.WarehousesDetailReq) (res *v1.WarehousesDetailRes, err error) {
	return service.Warehouses().Detail(ctx, req)
}

func (c *ControllerWarehousesV1) Create(ctx context.Context, req *v1.WarehousesCreateReq) (res *v1.WarehousesCreateRes, err error) {
	return service.Warehouses().Create(ctx, req)
}

func (c *ControllerWarehousesV1) Update(ctx context.Context, req *v1.WarehousesUpdateReq) (res *v1.WarehousesUpdateRes, err error) {
	return service.Warehouses().Update(ctx, req)
}

func (c *ControllerWarehousesV1) Delete(ctx context.Context, req *v1.WarehousesDeleteReq) (res *v1.WarehousesDeleteRes, err error) {
	return service.Warehouses().Delete(ctx, req)
}
