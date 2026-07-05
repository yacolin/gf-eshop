package payments

import (
	"context"

	"gf-eshop/api/payments/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) CreatePayment(ctx context.Context, req *v1.PaymentsCreateReq) (res *v1.PaymentsCreateRes, err error) {
	return service.Payments().CreatePayment(ctx, req)
}

func (c *ControllerV1) Callback(ctx context.Context, req *v1.PaymentsCallbackReq) (res *v1.PaymentsCallbackRes, err error) {
	return service.Payments().HandleCallback(ctx, req)
}

func (c *ControllerV1) GetPayment(ctx context.Context, req *v1.PaymentsGetReq) (res *v1.PaymentsGetRes, err error) {
	return service.Payments().GetPayment(ctx, req)
}

func (c *ControllerV1) CreateRefund(ctx context.Context, req *v1.RefundsCreateReq) (res *v1.RefundsCreateRes, err error) {
	return service.Payments().CreateRefund(ctx, req)
}