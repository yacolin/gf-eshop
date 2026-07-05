// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Roles is the golang structure of table sys_roles for DAO operations like Where/Data.
type Roles struct {
	g.Meta      `orm:"table:sys_roles, do:true"`
	Id          interface{} // 主键
	Name        interface{} // 角色名称（唯一，如 admin/editor/vip）
	DisplayName interface{} // 角色显示名称
	Description interface{} // 角色描述
	RoleType    interface{} // builtin-系统内置 custom-自定义
	SortOrder   interface{} // 排序值
	Status      interface{} // 1-启用 0-禁用
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
