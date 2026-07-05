package productAttributes

import (
	"context"

	"gf-eshop/api/product_attributes/v1"
)

type IProductAttributesV1 interface {
	List(ctx context.Context, req *v1.ProductAttributesListReq) (res *v1.ProductAttributesListRes, err error)
	Create(ctx context.Context, req *v1.ProductAttributesCreateReq) (res *v1.ProductAttributesCreateRes, err error)
	Delete(ctx context.Context, req *v1.ProductAttributesDeleteReq) (res *v1.ProductAttributesDeleteRes, err error)
}
