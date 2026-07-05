package productDescriptions

import (
	"context"

	"gf-eshop/api/product_descriptions/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Detail(ctx context.Context, req *v1.ProductDescriptionsDetailReq) (res *v1.ProductDescriptionsDetailRes, err error) {
	return service.ProductDescriptions().Detail(ctx, req)
}

func (c *ControllerV1) Save(ctx context.Context, req *v1.ProductDescriptionsSaveReq) (res *v1.ProductDescriptionsSaveRes, err error) {
	return service.ProductDescriptions().Save(ctx, req)
}
