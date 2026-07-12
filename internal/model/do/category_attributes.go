// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CategoryAttributes is the golang structure of table sp_category_attributes for DAO operations like Where/Data.
type CategoryAttributes struct {
	g.Meta          `orm:"table:sp_category_attributes, do:true"`
	Id              interface{} //
	CategoryId      interface{} // 类目ID
	AttributeId     interface{} // 属性ID
	Required        interface{} // 该类目下是否必填（仅提示，非强校验）
	IsDefaultFilter interface{} // 是否作为前台默认筛选项
	SortOrder       interface{} // 排序
	CreatedAt       *gtime.Time // 创建时间
}
