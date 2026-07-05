package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type ProductsListReq struct {
	g.Meta `path:"/products" tags:"Products" method:"get" summary:"商品列表"`

	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Name       string `json:"name"`
	CategoryId int64  `json:"category_id"`
	BrandId    int64  `json:"brand_id"`
	Status     int    `json:"status"`
	PriceMin   int64  `json:"price_min"`
	PriceMax   int64  `json:"price_max"`
}
type ProductsListRes struct {
	List  []*entity.Products `json:"list"`
	Total int                `json:"total"`
}

type ProductsDetailReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"get" summary:"商品详情"`
	Id     int64 `json:"id"`
}
type ProductsDetailRes struct {
	*entity.Products
}

type ProductsCreateReq struct {
	g.Meta `path:"/products" tags:"Products" method:"post" summary:"新增商品"`

	Name       string `json:"name"        v:"required|length:1,200" description:"商品名称"`
	Subtitle   string `json:"subtitle"    description:"副标题"`
	CategoryId int64  `json:"category_id" v:"required"             description:"类目ID"`
	BrandId    int64  `json:"brand_id"    description:"品牌ID"`
	Unit       string `json:"unit"        description:"单位"`
	MainImage  string `json:"main_image"  v:"required"             description:"主图"`
	Images     string `json:"images"      description:"附图JSON"`
	VideoUrl   string `json:"video_url"   description:"视频URL"`
	SortOrder  int    `json:"sort_order"  description:"排序权重"`
	Status     int    `json:"status"      description:"状态"`
	CreatedBy  string `json:"created_by"  description:"创建人"`
}
type ProductsCreateRes struct {
	Id int64 `json:"id"`
}

type ProductsUpdateReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"put" summary:"更新商品"`

	Id         int64  `json:"id"          v:"required"`
	Name       string `json:"name"        v:"length:1,200" description:"商品名称"`
	Subtitle   string `json:"subtitle"    description:"副标题"`
	CategoryId int64  `json:"category_id" description:"类目ID"`
	BrandId    int64  `json:"brand_id"    description:"品牌ID"`
	Unit       string `json:"unit"        description:"单位"`
	MainImage  string `json:"main_image"  description:"主图"`
	Images     string `json:"images"      description:"附图JSON"`
	VideoUrl   string `json:"video_url"   description:"视频URL"`
	SortOrder  int    `json:"sort_order"  description:"排序权重"`
	Status     int    `json:"status"      description:"状态"`
	UpdatedBy  string `json:"updated_by"  description:"更新人"`
}
type ProductsUpdateRes struct{}

type ProductsDeleteReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"delete" summary:"删除商品"`
	Id     int64 `json:"id"`
}
type ProductsDeleteRes struct{}
