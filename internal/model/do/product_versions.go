// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductVersions is the golang structure of table sp_product_versions for DAO operations like Where/Data.
type ProductVersions struct {
	g.Meta        `orm:"table:sp_product_versions, do:true"`
	Id            interface{} // 主键
	ProductId     interface{} // 关联 sp_products.id
	Version       interface{} // 版本号（从1递增）
	Diff          interface{} // 变更JSON（{"before": {...}, "after": {...}}）
	ChangedFields interface{} // 变更字段列表（如：["name", "price", "status"]）
	Operator      interface{} // 操作人
	OperatorId    interface{} // 操作人ID
	Reason        interface{} // 变更原因
	CreatedAt     *gtime.Time //
}
