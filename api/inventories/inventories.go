package inventories

import (
	"context"

	"gf-eshop/api/inventories/v1"
)

type IInventoriesV1 interface {
	List(ctx context.Context, req *v1.InventoriesListReq) (res *v1.InventoriesListRes, err error)
	Detail(ctx context.Context, req *v1.InventoriesDetailReq) (res *v1.InventoriesDetailRes, err error)
	Create(ctx context.Context, req *v1.InventoriesCreateReq) (res *v1.InventoriesCreateRes, err error)
	Update(ctx context.Context, req *v1.InventoriesUpdateReq) (res *v1.InventoriesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.InventoriesDeleteReq) (res *v1.InventoriesDeleteRes, err error)
}

type Warehouses interface {
	List(ctx context.Context, req *v1.WarehousesListReq) (res *v1.WarehousesListRes, err error)
	Detail(ctx context.Context, req *v1.WarehousesDetailReq) (res *v1.WarehousesDetailRes, err error)
	Create(ctx context.Context, req *v1.WarehousesCreateReq) (res *v1.WarehousesCreateRes, err error)
	Update(ctx context.Context, req *v1.WarehousesUpdateReq) (res *v1.WarehousesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.WarehousesDeleteReq) (res *v1.WarehousesDeleteRes, err error)
}
