// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItems is the golang structure of table tx_cart_items for DAO operations like Where/Data.
type CartItems struct {
	g.Meta      `orm:"table:tx_cart_items, do:true"`
	Id          interface{} // 购物车项ID
	CartId      interface{} // 关联 tx_carts.id
	SkuId       interface{} // 关联 skus.id(sp_skus)
	ProductId   interface{} // 关联 products.id(sp_products)，冗余用于展示
	ProductName interface{} // 商品名（冗余快照）
	SkuSpec     interface{} // 规格JSON快照（如{"颜色":"红色","内存":"256G"}）
	Image       interface{} // 商品图（冗余快照）
	Price       interface{} // 加入时的价格（分，防止下单时价格变动导致纠纷）
	Quantity    interface{} // 数量
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
