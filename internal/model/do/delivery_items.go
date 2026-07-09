// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeliveryItems is the golang structure of table tx_delivery_items for DAO operations like Where/Data.
type DeliveryItems struct {
	g.Meta      `orm:"table:tx_delivery_items, do:true"`
	Id          interface{} // 主键
	DeliveryId  interface{} // 关联 tx_deliveries.id
	OrderItemId interface{} // 关联 tx_order_items.id
	SkuId       interface{} // SKU ID（冗余）
	ProductName interface{} // 商品名（冗余快照）
	Quantity    interface{} // 本次发货数量
	CreatedAt   *gtime.Time //
}
