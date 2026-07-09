package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type CategoryAttributesListReq struct {
	g.Meta `path:"/category_attributes" tags:"CategoryAttributes" method:"get" summary:"类目属性关联列表"`

	CategoryId int64 `json:"category_id" v:"required" description:"类目ID"`
}
type CategoryAttributesListRes struct {
	List []*entity.CategoryAttributes `json:"list"`
}

type CategoryAttributesCreateReq struct {
	g.Meta `path:"/category_attributes" tags:"CategoryAttributes" method:"post" summary:"新增类目属性关联"`

	CategoryId      int64 `json:"category_id"       v:"required"           description:"类目ID"`
	AttributeId     int64 `json:"attribute_id"      v:"required"           description:"属性ID"`
	Required        int   `json:"required"          description:"该类目下是否必填"`
	IsDefaultFilter int   `json:"is_default_filter" description:"是否作为前台默认筛选项"`
	SortOrder       int   `json:"sort_order"        description:"排序"`
}
type CategoryAttributesCreateRes struct {
	Id int64 `json:"id"`
}

type CategoryAttributesBatchCreateReq struct {
	g.Meta `path:"/category_attributes/batch" tags:"CategoryAttributes" method:"post" summary:"批量新增类目属性关联"`

	CategoryId int64                             `json:"category_id" v:"required" description:"类目ID"`
	Items      []CategoryAttributesBatchItem     `json:"items"       v:"required" description:"属性列表"`
}
type CategoryAttributesBatchItem struct {
	AttributeId     int64 `json:"attribute_id"      v:"required" description:"属性ID"`
	Required        int   `json:"required"          description:"该类目下是否必填"`
	IsDefaultFilter int   `json:"is_default_filter" description:"是否作为前台默认筛选项"`
	SortOrder       int   `json:"sort_order"        description:"排序"`
}
type CategoryAttributesBatchCreateRes struct{}

type CategoryAttributesListByCatReq struct {
	g.Meta `path:"/categories/{id}/attributes" tags:"Categories" method:"get" summary:"获取类目推荐属性（按类目ID）"`
	Id     int64 `json:"id" v:"required"`
}
type CategoryAttributesListByCatRes struct {
	List []*entity.CategoryAttributes `json:"list"`
}

type CategoryAttributesDeleteReq struct {
	g.Meta `path:"/category_attributes/{id}" tags:"CategoryAttributes" method:"delete" summary:"删除类目属性关联"`

	Id int64 `json:"id" v:"required"`
}
type CategoryAttributesDeleteRes struct{}
