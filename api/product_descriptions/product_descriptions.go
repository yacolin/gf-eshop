package productDescriptions

import (
	"context"

	"gf-eshop/api/product_descriptions/v1"
)

type IProductDescriptionsV1 interface {
	Detail(ctx context.Context, req *v1.ProductDescriptionsDetailReq) (res *v1.ProductDescriptionsDetailRes, err error)
	Save(ctx context.Context, req *v1.ProductDescriptionsSaveReq) (res *v1.ProductDescriptionsSaveRes, err error)
}
