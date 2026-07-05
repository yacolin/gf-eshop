// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductDescriptions is the golang structure of table sp_product_descriptions for DAO operations like Where/Data.
type ProductDescriptions struct {
	g.Meta            `orm:"table:sp_product_descriptions, do:true"`
	Id                interface{} //
	ProductId         interface{} // 关联 products.id
	Description       interface{} // 商品详情（富文本HTML）
	MobileDescription interface{} // 移动端详情（可选，适配手机展示）
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
