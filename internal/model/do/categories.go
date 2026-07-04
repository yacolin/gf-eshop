// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure of table sp_categories for DAO operations like Where/Data.
type Categories struct {
	g.Meta    `orm:"table:sp_categories, do:true"`
	Id        interface{} //
	Name      interface{} // 类目名称（如：手机）
	ParentId  interface{} // 父级ID（0表示根节点）
	Level     interface{} // 层级（1-3级）
	Path      interface{} // 路径（如：1/23/45/）
	IconUrl   interface{} // 类目图标
	SortOrder interface{} // 排序
	Status    interface{} // 1-启用 0-禁用
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
