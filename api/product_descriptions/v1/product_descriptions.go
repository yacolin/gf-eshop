package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type ProductDescriptionsDetailReq struct {
	g.Meta `path:"/products/{product_id}/description" tags:"ProductDescriptions" method:"get" summary:"商品详情内容"`
	ProductId int64 `json:"product_id"`
}
type ProductDescriptionsDetailRes struct {
	*entity.ProductDescriptions
}

type ProductDescriptionsSaveReq struct {
	g.Meta `path:"/products/{product_id}/description" tags:"ProductDescriptions" method:"put" summary:"保存商品详情"`

	ProductId         int64  `json:"product_id"         v:"required"`
	Description       string `json:"description"        description:"商品详情(富文本HTML)"`
	MobileDescription string `json:"mobile_description" description:"移动端详情"`
}
type ProductDescriptionsSaveRes struct{}
