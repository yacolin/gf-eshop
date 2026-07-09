// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SkuSpecs is the golang structure of table sp_sku_specs for DAO operations like Where/Data.
type SkuSpecs struct {
	g.Meta           `orm:"table:sp_sku_specs, do:true"`
	Id               interface{} //
	SkuId            interface{} // 关联SKU ID
	AttributeId      interface{} // 关联属性ID
	AttributeValueId interface{} // 关联属性值ID
	SortOrder        interface{} // 展示顺序
}
