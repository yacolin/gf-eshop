// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RolePermissions is the golang structure of table sys_role_permissions for DAO operations like Where/Data.
type RolePermissions struct {
	g.Meta       `orm:"table:sys_role_permissions, do:true"`
	Id           interface{} // 主键
	RoleId       interface{} // 角色ID（关联 sys_roles.id）
	PermissionId interface{} // 权限ID（关联 sys_permissions.id）
	ScopeType    interface{} // 范围：platform-平台 merchant-商家
	ScopeId      interface{} // 范围ID（平台0/商家ID）
	CreatedAt    *gtime.Time // 创建时间
	DeletedAt    *gtime.Time // 删除时间
}
