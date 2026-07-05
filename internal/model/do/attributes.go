// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Attributes is the golang structure of table sp_attributes for DAO operations like Where/Data.
type Attributes struct {
	g.Meta     `orm:"table:sp_attributes, do:true"`
	Id         interface{} //
	Name       interface{} // 属性名称（如：处理器、屏幕尺寸）
	CategoryId interface{} // 所属类目ID（该属性只出现在这个类目下）
	InputType  interface{} // 1-文本输入 2-单选 3-多选 4-数字
	Values     interface{} // 可选值列表（如["A15","A16"]，仅单选/多选时使用）
	Unit       interface{} // 单位（如：英寸、GB）
	Required   interface{} // 1-必填（该属性在该类目下创建商品时必须填写）
	Searchable interface{} // 1-作为前台筛选条件（列表页筛选项来源）
	IsSkuSpec  interface{} // 1-是SKU规格（如颜色、内存） 0-仅SPU属性（如上市时间）
	SortOrder  interface{} //
	Status     interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
