// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Inventories is the golang structure for table inventories.
type Inventories struct {
	Id            int64       `json:"id"              description:"库存记录ID"`
	SkuId         int64       `json:"sku_id"          description:"关联 skus.id"`
	MerchantId    int64       `json:"merchant_id"     description:"所属商家ID"`
	WarehouseId   int64       `json:"warehouse_id"    description:"仓库ID（关联 sp_warehouses.id）"`
	Quantity      int64       `json:"quantity"        description:"物理库存总量（含预占）"`
	Reserved      int64       `json:"reserved"        description:"预占库存（下单未支付）"`
	InTransit     int64       `json:"in_transit"      description:"在途库存（采购中/调拨中）"`
	Threshold     int64       `json:"threshold"       description:"安全库存预警阈值（低于此值触发告警）"`
	MaxThreshold  int64       `json:"max_threshold"   description:"最大库存上限（入库不能超过此值）"`
	Status        int         `json:"status"          description:"1-充足 2-缺货 3-无货"`
	LastCountedAt *gtime.Time `json:"last_counted_at" description:"最后盘点时间"`
	LastCountedBy string      `json:"last_counted_by" description:"最后盘点人"`
	CreatedAt     *gtime.Time `json:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updated_at"      description:""`
	DeletedAt     *gtime.Time `json:"deleted_at"      description:""`
}
