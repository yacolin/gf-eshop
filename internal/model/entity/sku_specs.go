// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SkuSpecs is the golang structure for table sku_specs.
type SkuSpecs struct {
	Id               int64 `json:"id"                 description:""`
	SkuId            int64 `json:"sku_id"             description:"关联SKU ID"`
	AttributeId      int64 `json:"attribute_id"       description:"关联属性ID"`
	AttributeValueId int64 `json:"attribute_value_id" description:"关联属性值ID"`
	SortOrder        int   `json:"sort_order"         description:"展示顺序"`
}
