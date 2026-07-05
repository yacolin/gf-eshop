// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RolePermissions is the golang structure for table role_permissions.
type RolePermissions struct {
	Id           int64       `json:"id"            description:"主键"`
	RoleId       int64       `json:"role_id"       description:"角色ID（关联 sys_roles.id）"`
	PermissionId int64       `json:"permission_id" description:"权限ID（关联 sys_permissions.id）"`
	ScopeType    string      `json:"scope_type"    description:"范围：platform-平台 merchant-商家"`
	ScopeId      int64       `json:"scope_id"      description:"范围ID（平台0/商家ID）"`
	CreatedAt    *gtime.Time `json:"created_at"    description:"创建时间"`
	DeletedAt    *gtime.Time `json:"deleted_at"    description:"删除时间"`
}
