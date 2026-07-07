// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionProducts is the golang structure of table mkt_promotion_products for DAO operations like Where/Data.
type PromotionProducts struct {
	g.Meta      `orm:"table:mkt_promotion_products, do:true"`
	Id          interface{} // 主键ID
	PromotionId interface{} // 促销ID
	MerchantId  interface{} // 所属商家ID
	ProductType interface{} // 1-全站 2-指定分类 3-指定SPU 4-指定SKU
	ProductId   interface{} // 产品ID（product_type=3/4时使用）
	CategoryId  interface{} // 分类ID（product_type=2时使用）
	CreatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
