package attributes

import (
	"context"

	"gf-eshop/api/attributes/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.AttributesListReq) (res *v1.AttributesListRes, err error) {
	return service.Attributes().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.AttributesDetailReq) (res *v1.AttributesDetailRes, err error) {
	return service.Attributes().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.AttributesCreateReq) (res *v1.AttributesCreateRes, err error) {
	return service.Attributes().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.AttributesUpdateReq) (res *v1.AttributesUpdateRes, err error) {
	return service.Attributes().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.AttributesDeleteReq) (res *v1.AttributesDeleteRes, err error) {
	return service.Attributes().Delete(ctx, req)
}
