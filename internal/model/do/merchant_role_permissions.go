// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRolePermissions is the golang structure of table mch_merchant_role_permissions for DAO operations like Where/Data.
type MerchantRolePermissions struct {
	g.Meta         `orm:"table:mch_merchant_role_permissions, do:true"`
	Id             interface{} // 主键
	MerchantId     interface{} // 商家ID
	RoleId         interface{} // 角色ID（关联 mch_roles.id）
	PermissionName interface{} // 权限标识（对应 sys_permissions.name）
	CreatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
