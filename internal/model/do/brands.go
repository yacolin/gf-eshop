// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Brands is the golang structure of table sp_brands for DAO operations like Where/Data.
type Brands struct {
	g.Meta      `orm:"table:sp_brands, do:true"`
	Id          interface{} //
	Name        interface{} // 品牌名称（如：苹果）
	EnglishName interface{} // 英文名
	LogoUrl     interface{} // 品牌Logo（CDN）
	FirstLetter interface{} // 首字母（A-Z，用于前台索引筛选）
	SortOrder   interface{} // 排序权重
	Status      interface{} // 1-启用 0-禁用
	Description interface{} // 品牌故事
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
