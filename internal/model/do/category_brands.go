// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CategoryBrands is the golang structure of table sp_category_brands for DAO operations like Where/Data.
type CategoryBrands struct {
	g.Meta     `orm:"table:sp_category_brands, do:true"`
	Id         interface{} //
	CategoryId interface{} // 关联 categories.id
	BrandId    interface{} // 关联 brands.id
	SortOrder  interface{} // 排序权重（越小越靠前，控制该类目下品牌的展示顺序）
	CreatedAt  *gtime.Time //
}
