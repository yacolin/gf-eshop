// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CategoryBrands is the golang structure for table category_brands.
type CategoryBrands struct {
	Id         int64       `json:"id"         description:""`
	CategoryId int64       `json:"categoryId" description:"关联 categories.id"`
	BrandId    int64       `json:"brandId"    description:"关联 brands.id"`
	SortOrder  int         `json:"sortOrder"  description:"排序权重（越小越靠前，控制该类目下品牌的展示顺序）"`
	CreatedAt  *gtime.Time `json:"createdAt"  description:""`
}
