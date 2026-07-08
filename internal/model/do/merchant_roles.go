// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRoles is the golang structure of table mch_merchant_roles for DAO operations like Where/Data.
type MerchantRoles struct {
	g.Meta      `orm:"table:mch_merchant_roles, do:true"`
	Id          interface{} // 主键
	MerchantId  interface{} // 商家ID（0=平台预置角色）
	Name        interface{} // 角色名称（如店长/运营/财务）
	DisplayName interface{} // 角色显示名称
	Description interface{} // 角色描述
	RoleType    interface{} // builtin-系统预置 custom-商家自定义
	SortOrder   interface{} // 排序值
	Status      interface{} // 1-启用 0-禁用
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
