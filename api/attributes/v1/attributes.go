package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type AttributesListReq struct {
	g.Meta `path:"/attributes" tags:"Attributes" method:"get" summary:"属性列表"`

	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	CategoryId int64 `json:"category_id"`
	Searchable int   `json:"searchable"`
	IsSkuSpec  int   `json:"is_sku_spec"`
}
type AttributesListRes struct {
	List  []*entity.Attributes `json:"list"`
	Total int                  `json:"total"`
}

type AttributesDetailReq struct {
	g.Meta `path:"/attributes/{id}" tags:"Attributes" method:"get" summary:"属性详情"`
	Id     int64 `json:"id"`
}
type AttributesDetailRes struct {
	*entity.Attributes
}

type AttributesCreateReq struct {
	g.Meta `path:"/attributes" tags:"Attributes" method:"post" summary:"新增属性"`

	Name       string `json:"name"        v:"required|length:1,100" description:"属性名称"`
	CategoryId int64  `json:"category_id" v:"required"              description:"所属类目ID"`
	InputType  int    `json:"input_type"  v:"required|in:1,2,3,4"   description:"输入类型"`
	Values     string `json:"values"      description:"可选值列表JSON"`
	Unit       string `json:"unit"        description:"单位"`
	Required   int    `json:"required"    description:"是否必填"`
	Searchable int    `json:"searchable"  description:"是否可搜索"`
	IsSkuSpec  int    `json:"is_sku_spec" description:"是否SKU规格"`
	SortOrder  int    `json:"sort_order"  description:"排序"`
	Status     int    `json:"status"      description:"状态"`
}
type AttributesCreateRes struct {
	Id int64 `json:"id"`
}

type AttributesUpdateReq struct {
	g.Meta `path:"/attributes/{id}" tags:"Attributes" method:"put" summary:"更新属性"`

	Id         int64  `json:"id"          v:"required"`
	Name       string `json:"name"        v:"length:1,100" description:"属性名称"`
	CategoryId int64  `json:"category_id" description:"所属类目ID"`
	InputType  int    `json:"input_type"  v:"in:1,2,3,4"   description:"输入类型"`
	Values     string `json:"values"      description:"可选值列表JSON"`
	Unit       string `json:"unit"        description:"单位"`
	Required   int    `json:"required"    description:"是否必填"`
	Searchable int    `json:"searchable"  description:"是否可搜索"`
	IsSkuSpec  int    `json:"is_sku_spec" description:"是否SKU规格"`
	SortOrder  int    `json:"sort_order"  description:"排序"`
	Status     int    `json:"status"      description:"状态"`
}
type AttributesUpdateRes struct{}

type AttributesDeleteReq struct {
	g.Meta `path:"/attributes/{id}" tags:"Attributes" method:"delete" summary:"删除属性"`
	Id     int64 `json:"id"`
}
type AttributesDeleteRes struct{}
