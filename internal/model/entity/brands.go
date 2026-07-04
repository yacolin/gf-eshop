// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Brands is the golang structure for table brands.
type Brands struct {
	Id          int64       `json:"id"          description:""`
	Name        string      `json:"name"        description:"品牌名称（如：苹果）"`
	EnglishName string      `json:"englishName" description:"英文名"`
	LogoUrl     string      `json:"logoUrl"     description:"品牌Logo（CDN）"`
	FirstLetter string      `json:"firstLetter" description:"首字母（A-Z，用于前台索引筛选）"`
	SortOrder   int         `json:"sortOrder"   description:"排序权重"`
	Status      int         `json:"status"      description:"1-启用 0-禁用"`
	Description string      `json:"description" description:"品牌故事"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   description:""`
}
