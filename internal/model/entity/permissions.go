// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Permissions is the golang structure for table permissions.
type Permissions struct {
	Id          int64       `json:"id"           description:"主键"`
	Name        string      `json:"name"         description:"权限标识（唯一，如 order:create）"`
	DisplayName string      `json:"display_name" description:"权限显示名称"`
	Description string      `json:"description"  description:"权限描述"`
	Resource    string      `json:"resource"     description:"资源（如 order/product/user）"`
	Action      string      `json:"action"       description:"操作（如 create/read/update/delete）"`
	ParentId    int64       `json:"parent_id"    description:"父级ID（0=根节点，支持菜单/按钮树形层级）"`
	Category    string      `json:"category"     description:"分类（如 business/system/admin）"`
	SortOrder   int         `json:"sort_order"   description:"排序值"`
	Status      int         `json:"status"       description:"1-启用 0-禁用"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:"删除时间"`
}
