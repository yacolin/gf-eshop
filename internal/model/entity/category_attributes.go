// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CategoryAttributes is the golang structure for table category_attributes.
type CategoryAttributes struct {
	Id              int64       `json:"id"                description:""`
	CategoryId      int64       `json:"category_id"       description:"类目ID"`
	AttributeId     int64       `json:"attribute_id"      description:"属性ID"`
	Required        int         `json:"required"          description:"该类目下是否必填（仅提示，非强校验）"`
	IsDefaultFilter int         `json:"is_default_filter" description:"是否作为前台默认筛选项"`
	SortOrder       int         `json:"sort_order"        description:""`
	CreatedAt       *gtime.Time `json:"created_at"        description:""`
}
