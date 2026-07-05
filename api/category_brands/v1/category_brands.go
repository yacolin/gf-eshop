package v1


import (
	"github.com/gogf/gf/v2/frame/g"
)



// ---------- List ----------
type CategoryBrandListReq struct {
	g.Meta `path:"/categories/{id}/brands" tags:"Categories" method:"get" summary:"类目下品牌"`
	Id     int64 `json:"id"`
}

type CategoryBrandListRes struct {
	List []*CategoryBrandItem `json:"list"`
}

// CategoryBrandItem 关联关系 + 品牌详情
type CategoryBrandItem struct {
	Id         int64  `json:"id"`
	CategoryId int64  `json:"category_id"`
	BrandId    int64  `json:"brand_id"`
	SortOrder  int    `json:"sort_order"`

	BrandName   string `json:"brand_name"`
	EnglishName string `json:"english_name"`
	LogoUrl     string `json:"logo_url"`
	FirstLetter string `json:"first_letter"`
}


// ---------- Update ----------
type CategoryBrandUpdateReq struct {
	g.Meta `path:"/categories/{id}/brands" tags:"Categories" method:"put" summary:"类目关联品牌"`

	Id        int64   `json:"id"`
	BrandIDs  []int64 `json:"brand_ids"  v:"required"`
	SortOrder int     `json:"sort_order" description:"排序权重"`
}
type CategoryBrandUpdateRes struct{}