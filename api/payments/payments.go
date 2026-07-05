package payments

import (
	"context"

	"gf-eshop/api/payments/v1"
)

type IPaymentsV1 interface {
	CreatePayment(ctx context.Context, req *v1.PaymentsCreateReq) (res *v1.PaymentsCreateRes, err error)
	Callback(ctx context.Context, req *v1.PaymentsCallbackReq) (res *v1.PaymentsCallbackRes, err error)
	GetPayment(ctx context.Context, req *v1.PaymentsGetReq) (res *v1.PaymentsGetRes, err error)
	CreateRefund(ctx context.Context, req *v1.RefundsCreateReq) (res *v1.RefundsCreateRes, err error)
}