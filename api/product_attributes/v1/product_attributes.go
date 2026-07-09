package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type ProductAttributesListReq struct {
	g.Meta `path:"/product_attributes" tags:"ProductAttributes" method:"get" summary:"商品属性值列表"`

	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	ProductId int64 `json:"product_id"`
}
type ProductAttributesListRes struct {
	List  []*entity.ProductAttributes `json:"list"`
	Total int                         `json:"total"`
}

type ProductAttributesCreateReq struct {
	g.Meta `path:"/product_attributes" tags:"ProductAttributes" method:"post" summary:"新增商品属性值"`

	ProductId        int64  `json:"product_id"         v:"required" description:"商品ID"`
	AttributeId      int64  `json:"attribute_id"       v:"required" description:"属性ID"`
	AttributeValueId int64  `json:"attribute_value_id"               description:"引用属性值字典ID（可选），优先使用"`
	Value            string `json:"value"              v:"required" description:"属性值。有字典值时冗余存储展示文本，无字典值时存自由文本"`
}
type ProductAttributesCreateRes struct {
	Id int64 `json:"id"`
}

type ProductAttributesDeleteReq struct {
	g.Meta `path:"/product_attributes/{id}" tags:"ProductAttributes" method:"delete" summary:"删除商品属性值"`
	Id     int64 `json:"id"`
}
type ProductAttributesDeleteRes struct{}
