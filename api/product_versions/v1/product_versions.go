package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type ProductVersionsListReq struct {
	g.Meta `path:"/product_versions" tags:"ProductVersions" method:"get" summary:"商品版本列表"`

	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	ProductId int64 `json:"product_id"`
}
type ProductVersionsListRes struct {
	List  []*entity.ProductVersions `json:"list"`
	Total int                       `json:"total"`
}

type ProductVersionsDetailReq struct {
	g.Meta `path:"/product_versions/{id}" tags:"ProductVersions" method:"get" summary:"版本详情"`
	Id     int64 `json:"id"`
}
type ProductVersionsDetailRes struct {
	*entity.ProductVersions
}
