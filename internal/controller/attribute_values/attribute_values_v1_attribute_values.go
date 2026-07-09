package attribute_values

import (
	"context"

	"gf-eshop/api/attribute_values/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.AttributeValuesListReq) (res *v1.AttributeValuesListRes, err error) {
	return service.AttributeValues().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.AttributeValuesDetailReq) (res *v1.AttributeValuesDetailRes, err error) {
	return service.AttributeValues().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.AttributeValuesCreateReq) (res *v1.AttributeValuesCreateRes, err error) {
	return service.AttributeValues().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.AttributeValuesUpdateReq) (res *v1.AttributeValuesUpdateRes, err error) {
	return service.AttributeValues().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.AttributeValuesDeleteReq) (res *v1.AttributeValuesDeleteRes, err error) {
	return service.AttributeValues().Delete(ctx, req)
}
