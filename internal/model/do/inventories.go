// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Inventories is the golang structure of table sp_inventories for DAO operations like Where/Data.
type Inventories struct {
	g.Meta        `orm:"table:sp_inventories, do:true"`
	Id            interface{} // 库存记录ID
	SkuId         interface{} // 关联 skus.id
	MerchantId    interface{} // 所属商家ID
	WarehouseId   interface{} // 仓库ID（关联 sp_warehouses.id）
	Quantity      interface{} // 物理库存总量（含预占）
	Reserved      interface{} // 预占库存（下单未支付）
	InTransit     interface{} // 在途库存（采购中/调拨中）
	Threshold     interface{} // 安全库存预警阈值（低于此值触发告警）
	MaxThreshold  interface{} // 最大库存上限（入库不能超过此值）
	Status        interface{} // 1-充足 2-缺货 3-无货
	LastCountedAt *gtime.Time // 最后盘点时间
	LastCountedBy interface{} // 最后盘点人
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
