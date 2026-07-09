// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AttributeValues is the golang structure of table sp_attribute_values for DAO operations like Where/Data.
type AttributeValues struct {
	g.Meta       `orm:"table:sp_attribute_values, do:true"`
	Id           interface{} //
	AttributeId  interface{} // 关联属性ID
	Value        interface{} // 属性值（如：256G、红色）
	NumericValue interface{} // 数值型值（用于区间筛选）
	ColorHex     interface{} // 颜色色值（#FF0000）
	SortOrder    interface{} // 排序权重
	Status       interface{} // 1-启用 0-禁用
	CreatedAt    *gtime.Time //
}
