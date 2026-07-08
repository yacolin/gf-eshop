package merchantWithdrawals

import (
	"context"

	"gf-eshop/api/merchant_withdrawals/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	return service.MerchantWithdrawals().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	return service.MerchantWithdrawals().Detail(ctx, req)
}

func (c *ControllerV1) Apply(ctx context.Context, req *v1.ApplyReq) (res *v1.ApplyRes, err error) {
	return service.MerchantWithdrawals().Apply(ctx, req)
}

func (c *ControllerV1) Approve(ctx context.Context, req *v1.ApproveReq) (res *v1.ApproveRes, err error) {
	return service.MerchantWithdrawals().Approve(ctx, req)
}

func (c *ControllerV1) Reject(ctx context.Context, req *v1.RejectReq) (res *v1.RejectRes, err error) {
	return service.MerchantWithdrawals().Reject(ctx, req)
}
