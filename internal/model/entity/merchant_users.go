// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUsers is the golang structure for table merchant_users.
type MerchantUsers struct {
	Id          int64       `json:"id"            description:"主键"`
	MerchantId  int64       `json:"merchant_id"   description:"商家ID"`
	StaffId     int64       `json:"staff_id"      description:"员工ID（关联 sys_staff.id）"`
	RoleId      int64       `json:"role_id"       description:"商家角色ID（关联 mch_roles.id，与平台RBAC隔离）"`
	Status      int         `json:"status"        description:"1-正常 2-禁用"`
	InvitedAt   *gtime.Time `json:"invited_at"    description:"邀请时间"`
	LastLoginAt *gtime.Time `json:"last_login_at" description:"最后登录时间"`
	CreatedAt   *gtime.Time `json:"created_at"    description:""`
	UpdatedAt   *gtime.Time `json:"updated_at"    description:""`
	DeletedAt   *gtime.Time `json:"deleted_at"    description:""`
}
