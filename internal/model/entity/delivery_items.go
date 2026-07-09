// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeliveryItems is the golang structure for table delivery_items.
type DeliveryItems struct {
	Id          int64       `json:"id"            description:"主键"`
	DeliveryId  int64       `json:"delivery_id"   description:"关联 tx_deliveries.id"`
	OrderItemId int64       `json:"order_item_id" description:"关联 tx_order_items.id"`
	SkuId       int64       `json:"sku_id"        description:"SKU ID（冗余）"`
	ProductName string      `json:"product_name"  description:"商品名（冗余快照）"`
	Quantity    int         `json:"quantity"      description:"本次发货数量"`
	CreatedAt   *gtime.Time `json:"created_at"    description:""`
}
