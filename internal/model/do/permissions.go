// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Permissions is the golang structure of table sys_permissions for DAO operations like Where/Data.
type Permissions struct {
	g.Meta      `orm:"table:sys_permissions, do:true"`
	Id          interface{} // 主键
	Name        interface{} // 权限标识（唯一，如 order:create）
	DisplayName interface{} // 权限显示名称
	Description interface{} // 权限描述
	Resource    interface{} // 资源（如 order/product/user）
	Action      interface{} // 操作（如 create/read/update/delete）
	ParentId    interface{} // 父级ID（0=根节点，支持菜单/按钮树形层级）
	Category    interface{} // 分类（如 business/system/admin）
	SortOrder   interface{} // 排序值
	Status      interface{} // 1-启用 0-禁用
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
