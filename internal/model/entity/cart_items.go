// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItems is the golang structure for table cart_items.
type CartItems struct {
	Id          int64       `json:"id"           description:"购物车项ID"`
	CartId      int64       `json:"cart_id"      description:"关联 tx_carts.id"`
	SkuId       int64       `json:"sku_id"       description:"关联 skus.id(sp_skus)"`
	ProductId   int64       `json:"product_id"   description:"关联 products.id(sp_products)，冗余用于展示"`
	ProductName string      `json:"product_name" description:"商品名（冗余快照）"`
	SkuSpec     string      `json:"sku_spec"     description:"规格JSON快照（如{\"颜色\":\"红色\",\"内存\":\"256G\"}）"`
	Image       string      `json:"image"        description:"商品图（冗余快照）"`
	Price       int64       `json:"price"        description:"加入时的价格（分，防止下单时价格变动导致纠纷）"`
	Quantity    int         `json:"quantity"     description:"数量"`
	CreatedAt   *gtime.Time `json:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:""`
}
