package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type ListReq struct {
	g.Meta `path:"/brands" tags:"Brands" method:"get" summary:"品牌列表"`

	Page       int    `json:"page"`        // 页码，默认1
	PageSize   int    `json:"page_size"`    // 每页条数，默认20
	Name       string `json:"name"`        // 按名称模糊搜索
	FirstLetter string `json:"first_letter"` // 按首字母筛选
	Status     int    `json:"status"`      // 按状态筛选
}
type ListRes struct {
	List  []*entity.Brands `json:"list"`
	Total int              `json:"total"`
}

// ---------- Detail ----------
type DetailReq struct {
	g.Meta `path:"/brands/{id}" tags:"Brands" method:"get" summary:"品牌详情"`
	Id     int64 `json:"id"`
}
type DetailRes struct {
	*entity.Brands
}

// ---------- Create ----------
type CreateReq struct {
	g.Meta `path:"/brands" tags:"Brands" method:"post" summary:"新增品牌"`

	Name        string `json:"name"        v:"required|length:1,100" description:"品牌名称"`
	EnglishName string `json:"english_name" description:"英文名"`
	LogoUrl     string `json:"logo_url"     description:"品牌Logo"`
	FirstLetter string `json:"first_letter" v:"length:1,1" description:"首字母"`
	SortOrder   int    `json:"sort_order"   description:"排序权重"`
	Status      int    `json:"status"      description:"状态"`
	Description string `json:"description" description:"品牌故事"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type UpdateReq struct {
	g.Meta `path:"/brands/{id}" tags:"Brands" method:"put" summary:"更新品牌"`

	Id          int64  `json:"id"          v:"required"`
	Name        string `json:"name"        v:"length:1,100" description:"品牌名称"`
	EnglishName string `json:"english_name" description:"英文名"`
	LogoUrl     string `json:"logo_url"     description:"品牌Logo"`
	FirstLetter string `json:"first_letter" v:"length:1,1" description:"首字母"`
	SortOrder   int    `json:"sort_order"   description:"排序权重"`
	Status      int    `json:"status"      description:"状态"`
	Description string `json:"description" description:"品牌故事"`
}
type UpdateRes struct{}

// ---------- Delete ----------
type DeleteReq struct {
	g.Meta `path:"/brands/{id}" tags:"Brands" method:"delete" summary:"删除品牌"`
	Id     int64 `json:"id"`
}
type DeleteRes struct{}
