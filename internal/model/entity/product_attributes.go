// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductAttributes is the golang structure for table product_attributes.
type ProductAttributes struct {
	Id               int64       `json:"id"                 description:""`
	ProductId        int64       `json:"product_id"         description:"关联 products.id"`
	AttributeId      int64       `json:"attribute_id"       description:"关联 attributes.id"`
	AttributeValueId int64       `json:"attribute_value_id" description:"引用属性值字典ID（可选，关联 attribute_values.id）"`
	Value            string      `json:"value"              description:"属性值（如：A16）。有字典值时冗余存储便于展示，无字典值时存自由文本"`
	SortOrder        int         `json:"sort_order"         description:"排序权重（越小越靠前，用于控制前台展示顺序）"`
	CreatedAt        *gtime.Time `json:"created_at"         description:""`
	UpdatedAt        *gtime.Time `json:"updated_at"         description:""`
	DeletedAt        *gtime.Time `json:"deleted_at"         description:"软删除"`
}
