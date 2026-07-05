// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Attributes is the golang structure for table attributes.
type Attributes struct {
	Id         int64       `json:"id"          description:""`
	Name       string      `json:"name"        description:"属性名称（如：处理器、屏幕尺寸）"`
	CategoryId int64       `json:"category_id" description:"所属类目ID（该属性只出现在这个类目下）"`
	InputType  int         `json:"input_type"  description:"1-文本输入 2-单选 3-多选 4-数字"`
	Values     string      `json:"values"      description:"可选值列表（如[\"A15\",\"A16\"]，仅单选/多选时使用）"`
	Unit       string      `json:"unit"        description:"单位（如：英寸、GB）"`
	Required   int         `json:"required"    description:"1-必填（该属性在该类目下创建商品时必须填写）"`
	Searchable int         `json:"searchable"  description:"1-作为前台筛选条件（列表页筛选项来源）"`
	IsSkuSpec  int         `json:"is_sku_spec" description:"1-是SKU规格（如颜色、内存） 0-仅SPU属性（如上市时间）"`
	SortOrder  int         `json:"sort_order"  description:""`
	Status     int         `json:"status"      description:""`
	CreatedAt  *gtime.Time `json:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updated_at"  description:""`
	DeletedAt  *gtime.Time `json:"deleted_at"  description:""`
}
