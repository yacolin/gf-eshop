// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffDepartments is the golang structure of table sys_staff_departments for DAO operations like Where/Data.
type StaffDepartments struct {
	g.Meta       `orm:"table:sys_staff_departments, do:true"`
	Id           interface{} // 主键
	StaffId      interface{} // 员工ID（关联 sys_staff.id）
	DepartmentId interface{} // 部门ID（关联 sys_departments.id）
	IsPrimary    interface{} // 是否主部门
	CreatedAt    *gtime.Time // 创建时间
	DeletedAt    *gtime.Time // 删除时间
}
