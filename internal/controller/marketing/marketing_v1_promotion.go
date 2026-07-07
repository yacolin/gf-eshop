package marketing

import (
	"context"

	"gf-eshop/api/marketing/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) PromotionList(ctx context.Context, req *v1.PromotionListReq) (res *v1.PromotionListRes, err error) {
	return service.Marketing().PromotionList(ctx, req)
}

func (c *ControllerV1) PromotionDetail(ctx context.Context, req *v1.PromotionDetailReq) (res *v1.PromotionDetailRes, err error) {
	return service.Marketing().PromotionDetail(ctx, req)
}

func (c *ControllerV1) PromotionFullDetail(ctx context.Context, req *v1.PromotionFullDetailReq) (res *v1.PromotionFullDetailRes, err error) {
	return service.Marketing().PromotionFullDetail(ctx, req)
}

func (c *ControllerV1) PromotionCreate(ctx context.Context, req *v1.PromotionCreateReq) (res *v1.PromotionCreateRes, err error) {
	return service.Marketing().PromotionCreate(ctx, req)
}

func (c *ControllerV1) PromotionUpdate(ctx context.Context, req *v1.PromotionUpdateReq) (res *v1.PromotionUpdateRes, err error) {
	return service.Marketing().PromotionUpdate(ctx, req)
}

func (c *ControllerV1) PromotionDelete(ctx context.Context, req *v1.PromotionDeleteReq) (res *v1.PromotionDeleteRes, err error) {
	return service.Marketing().PromotionDelete(ctx, req)
}
