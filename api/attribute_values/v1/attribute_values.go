package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type AttributeValuesListReq struct {
	g.Meta `path:"/attribute_values" tags:"AttributeValues" method:"get" summary:"属性值列表"`

	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	AttributeId int64 `json:"attribute_id"`
}
type AttributeValuesListRes struct {
	List  []*entity.AttributeValues `json:"list"`
	Total int                       `json:"total"`
}

type AttributeValuesDetailReq struct {
	g.Meta `path:"/attribute_values/{id}" tags:"AttributeValues" method:"get" summary:"属性值详情"`
	Id     int64 `json:"id"`
}
type AttributeValuesDetailRes struct {
	*entity.AttributeValues
}

type AttributeValuesCreateReq struct {
	g.Meta `path:"/attribute_values" tags:"AttributeValues" method:"post" summary:"新增属性值"`

	AttributeId  int64    `json:"attribute_id"  v:"required"              description:"属性ID"`
	Value        string   `json:"value"         v:"required|length:1,200" description:"属性值"`
	Alias        []string `json:"alias"         description:"别名列表"`
	SearchWeight int      `json:"search_weight" description:"搜索权重"`
	NumericValue float64  `json:"numeric_value" description:"数值型值"`
	ColorHex     string   `json:"color_hex"     description:"颜色色值"`
	SortOrder    int      `json:"sort_order"    description:"排序"`
	Status       int      `json:"status"        description:"状态"`
}
type AttributeValuesCreateRes struct {
	Id int64 `json:"id"`
}

type AttributeValuesUpdateReq struct {
	g.Meta `path:"/attribute_values/{id}" tags:"AttributeValues" method:"put" summary:"更新属性值"`

	Id           int64    `json:"id"            v:"required"`
	Value        string   `json:"value"         v:"length:1,200" description:"属性值"`
	Alias        []string `json:"alias"         description:"别名列表"`
	SearchWeight int      `json:"search_weight" description:"搜索权重"`
	NumericValue float64  `json:"numeric_value" description:"数值型值"`
	ColorHex     string   `json:"color_hex"     description:"颜色色值"`
	SortOrder    int      `json:"sort_order"    description:"排序"`
	Status       int      `json:"status"        description:"状态"`
}
type AttributeValuesUpdateRes struct{}

type AttributeValuesListByAttrReq struct {
	g.Meta `path:"/attributes/{id}/values" tags:"Attributes" method:"get" summary:"获取属性值列表（按属性ID）"`
	Id     int64 `json:"id" v:"required"`
}
type AttributeValuesListByAttrRes struct {
	List []*entity.AttributeValues `json:"list"`
}

type AttributeValuesDeleteReq struct {
	g.Meta `path:"/attribute_values/{id}" tags:"AttributeValues" method:"delete" summary:"删除属性值"`
	Id     int64 `json:"id" v:"required"`
}
type AttributeValuesDeleteRes struct{}
