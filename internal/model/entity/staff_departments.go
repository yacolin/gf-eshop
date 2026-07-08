// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffDepartments is the golang structure for table staff_departments.
type StaffDepartments struct {
	Id           int64       `json:"id"            description:"主键"`
	StaffId      int64       `json:"staff_id"      description:"员工ID（关联 sys_staff.id）"`
	DepartmentId int64       `json:"department_id" description:"部门ID（关联 sys_departments.id）"`
	IsPrimary    int         `json:"is_primary"    description:"是否主部门"`
	CreatedAt    *gtime.Time `json:"created_at"    description:"创建时间"`
	DeletedAt    *gtime.Time `json:"deleted_at"    description:"删除时间"`
}
