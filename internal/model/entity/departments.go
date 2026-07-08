// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Departments is the golang structure for table departments.
type Departments struct {
	Id        int64       `json:"id"         description:"主键"`
	Name      string      `json:"name"       description:"部门名称"`
	ParentId  int64       `json:"parent_id"  description:"上级部门ID（0=根部门）"`
	SortOrder int         `json:"sort_order" description:"排序值"`
	Status    int         `json:"status"     description:"1-启用 0-禁用"`
	CreatedAt *gtime.Time `json:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at" description:"更新时间"`
	DeletedAt *gtime.Time `json:"deleted_at" description:"删除时间"`
}
