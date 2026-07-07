package marketing

import (
	"context"

	"gf-eshop/api/marketing/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) CouponClaim(ctx context.Context, req *v1.CouponClaimReq) (res *v1.CouponClaimRes, err error) {
	return service.Marketing().CouponClaim(ctx, req)
}

func (c *ControllerV1) CouponUse(ctx context.Context, req *v1.CouponUseReq) (res *v1.CouponUseRes, err error) {
	return service.Marketing().CouponUse(ctx, req)
}

func (c *ControllerV1) CouponList(ctx context.Context, req *v1.CouponListReq) (res *v1.CouponListRes, err error) {
	return service.Marketing().CouponList(ctx, req)
}
