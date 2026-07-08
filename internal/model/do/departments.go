// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Departments is the golang structure of table sys_departments for DAO operations like Where/Data.
type Departments struct {
	g.Meta    `orm:"table:sys_departments, do:true"`
	Id        interface{} // 主键
	Name      interface{} // 部门名称
	ParentId  interface{} // 上级部门ID（0=根部门）
	SortOrder interface{} // 排序值
	Status    interface{} // 1-启用 0-禁用
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
}
