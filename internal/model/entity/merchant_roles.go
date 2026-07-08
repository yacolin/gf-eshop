// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRoles is the golang structure for table merchant_roles.
type MerchantRoles struct {
	Id          int64       `json:"id"           description:"主键"`
	MerchantId  int64       `json:"merchant_id"  description:"商家ID（0=平台预置角色）"`
	Name        string      `json:"name"         description:"角色名称（如店长/运营/财务）"`
	DisplayName string      `json:"display_name" description:"角色显示名称"`
	Description string      `json:"description"  description:"角色描述"`
	RoleType    string      `json:"role_type"    description:"builtin-系统预置 custom-商家自定义"`
	SortOrder   int         `json:"sort_order"   description:"排序值"`
	Status      int         `json:"status"       description:"1-启用 0-禁用"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:"删除时间"`
}
