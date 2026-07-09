// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AttributeValues is the golang structure for table attribute_values.
type AttributeValues struct {
	Id           int64       `json:"id"            description:""`
	AttributeId  int64       `json:"attribute_id"  description:"关联属性ID"`
	Value        string      `json:"value"         description:"属性值（如：256G、红色）"`
	Alias        string      `json:"alias"         description:"别名列表，如[\"深空灰\",\"黑灰\"]，用于搜索纠错、模糊匹配"`
	SearchWeight int         `json:"search_weight" description:"搜索权重（值越大匹配优先级越高）"`
	NumericValue float64     `json:"numeric_value" description:"数值型值（用于区间筛选）"`
	ColorHex     string      `json:"color_hex"     description:"颜色色值（#FF0000）"`
	SortOrder    int         `json:"sort_order"    description:"排序权重"`
	Status       int         `json:"status"        description:"1-启用 0-禁用"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
}
