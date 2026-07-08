// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRolePermissions is the golang structure for table merchant_role_permissions.
type MerchantRolePermissions struct {
	Id             int64       `json:"id"              description:"主键"`
	MerchantId     int64       `json:"merchant_id"     description:"商家ID"`
	RoleId         int64       `json:"role_id"         description:"角色ID（关联 mch_roles.id）"`
	PermissionName string      `json:"permission_name" description:"权限标识（对应 sys_permissions.name）"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:""`
}
