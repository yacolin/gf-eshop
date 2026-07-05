// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffRoles is the golang structure for table staff_roles.
type StaffRoles struct {
	Id        int64       `json:"id"         description:"主键"`
	StaffId   int64       `json:"staff_id"   description:"员工ID（关联 sys_staff.id）"`
	RoleId    int64       `json:"role_id"    description:"角色ID（关联 sys_roles.id）"`
	CreatedAt *gtime.Time `json:"created_at" description:"创建时间"`
	DeletedAt *gtime.Time `json:"deleted_at" description:"删除时间"`
}
