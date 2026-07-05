// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductAttributes is the golang structure of table sp_product_attributes for DAO operations like Where/Data.
type ProductAttributes struct {
	g.Meta      `orm:"table:sp_product_attributes, do:true"`
	Id          interface{} //
	ProductId   interface{} // 关联 products.id
	AttributeId interface{} // 关联 attributes.id
	Value       interface{} // 属性值（如：A16）
	SortOrder   interface{} // 排序权重（越小越靠前，用于控制前台展示顺序）
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time // 软删除
}
