// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure for table categories.
type Categories struct {
	Id        int64       `json:"id"         description:""`
	Name      string      `json:"name"       description:"类目名称（如：手机）"`
	ParentId  int64       `json:"parent_id"  description:"父级ID（0表示根节点）"`
	Level     int         `json:"level"      description:"层级（1-3级）"`
	Path      string      `json:"path"       description:"路径（如：1/23/45/）"`
	IconUrl   string      `json:"icon_url"   description:"类目图标"`
	SortOrder int         `json:"sort_order" description:"排序"`
	Status    int         `json:"status"     description:"1-启用 0-禁用"`
	CreatedAt *gtime.Time `json:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updated_at" description:""`
	DeletedAt *gtime.Time `json:"deleted_at" description:""`
}
