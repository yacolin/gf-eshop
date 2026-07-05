package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// --- Inventories ---

type InventoriesListReq struct {
	g.Meta `path:"/inventories" tags:"Inventories" method:"get" summary:"库存列表"`

	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	SkuId       int64 `json:"sku_id"`
	WarehouseId int64 `json:"warehouse_id"`
}
type InventoriesListRes struct {
	List  []*entity.Inventories `json:"list"`
	Total int                  `json:"total"`
}

type InventoriesDetailReq struct {
	g.Meta `path:"/inventories/{id}" tags:"Inventories" method:"get" summary:"库存详情"`
	Id     int64 `json:"id"`
}
type InventoriesDetailRes struct {
	*entity.Inventories
}

type InventoriesCreateReq struct {
	g.Meta `path:"/inventories" tags:"Inventories" method:"post" summary:"新增库存记录"`

	SkuId       int64 `json:"sku_id"       v:"required" description:"SKU ID"`
	WarehouseId int64 `json:"warehouse_id" v:"required" description:"仓库ID"`
	Quantity    int64 `json:"quantity"     v:"required" description:"物理库存"`
	Threshold   int64 `json:"threshold"    description:"安全库存阈值"`
}
type InventoriesCreateRes struct {
	Id int64 `json:"id"`
}

type InventoriesUpdateReq struct {
	g.Meta `path:"/inventories/{id}" tags:"Inventories" method:"put" summary:"更新库存"`

	Id          int64 `json:"id"          v:"required"`
	Quantity    int64 `json:"quantity"    description:"物理库存"`
	Reserved    int64 `json:"reserved"    description:"预占库存"`
	Threshold   int64 `json:"threshold"   description:"安全库存阈值"`
	MaxThreshold int64 `json:"max_threshold" description:"最大库存上限"`
}
type InventoriesUpdateRes struct{}

type InventoriesDeleteReq struct {
	g.Meta `path:"/inventories/{id}" tags:"Inventories" method:"delete" summary:"删除库存记录"`
	Id     int64 `json:"id"`
}
type InventoriesDeleteRes struct{}

// --- Warehouses ---

type WarehousesListReq struct {
	g.Meta `path:"/warehouses" tags:"Warehouses" method:"get" summary:"仓库列表"`

	Page   int `json:"page"`
	PageSize int `json:"page_size"`
}
type WarehousesListRes struct {
	List  []*entity.Warehouses `json:"list"`
	Total int                  `json:"total"`
}

type WarehousesDetailReq struct {
	g.Meta `path:"/warehouses/{id}" tags:"Warehouses" method:"get" summary:"仓库详情"`
	Id     int64 `json:"id"`
}
type WarehousesDetailRes struct {
	*entity.Warehouses
}

type WarehousesCreateReq struct {
	g.Meta `path:"/warehouses" tags:"Warehouses" method:"post" summary:"新增仓库"`

	WarehouseName string `json:"warehouse_name" v:"required" description:"仓库名称"`
	WarehouseType int    `json:"warehouse_type" description:"仓库类型"`
	Province      string `json:"province"       description:"省"`
	City          string `json:"city"           description:"市"`
	District      string `json:"district"       description:"区"`
	Address       string `json:"address"        description:"详细地址"`
	Status        int    `json:"status"         description:"状态"`
}
type WarehousesCreateRes struct {
	Id int64 `json:"id"`
}

type WarehousesUpdateReq struct {
	g.Meta `path:"/warehouses/{id}" tags:"Warehouses" method:"put" summary:"更新仓库"`

	Id            int64  `json:"id"            v:"required"`
	WarehouseName string `json:"warehouse_name" description:"仓库名称"`
	WarehouseType int    `json:"warehouse_type" description:"仓库类型"`
	Province      string `json:"province"       description:"省"`
	City          string `json:"city"           description:"市"`
	District      string `json:"district"       description:"区"`
	Address       string `json:"address"        description:"详细地址"`
	Status        int    `json:"status"         description:"状态"`
}
type WarehousesUpdateRes struct{}

type WarehousesDeleteReq struct {
	g.Meta `path:"/warehouses/{id}" tags:"Warehouses" method:"delete" summary:"删除仓库"`
	Id     int64 `json:"id"`
}
type WarehousesDeleteRes struct{}
