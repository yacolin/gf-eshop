// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffRoles is the golang structure of table sys_staff_roles for DAO operations like Where/Data.
type StaffRoles struct {
	g.Meta    `orm:"table:sys_staff_roles, do:true"`
	Id        interface{} // 主键
	StaffId   interface{} // 员工ID（关联 sys_staff.id）
	RoleId    interface{} // 角色ID（关联 sys_roles.id）
	CreatedAt *gtime.Time // 创建时间
	DeletedAt *gtime.Time // 删除时间
}
