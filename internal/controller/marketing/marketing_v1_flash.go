package marketing

import (
	"context"

	"gf-eshop/api/marketing/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) FlashBuy(ctx context.Context, req *v1.FlashBuyReq) (res *v1.FlashBuyRes, err error) {
	return service.Marketing().FlashBuy(ctx, req)
}

func (c *ControllerV1) FlashConfirm(ctx context.Context, req *v1.FlashConfirmReq) (res *v1.FlashConfirmRes, err error) {
	return service.Marketing().FlashConfirm(ctx, req)
}
