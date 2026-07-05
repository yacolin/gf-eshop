package service

import (
	"context"

	"gf-eshop/api/payments/v1"
)

type IPayments interface {
	CreatePayment(ctx context.Context, req *v1.PaymentsCreateReq) (res *v1.PaymentsCreateRes, err error)
	HandleCallback(ctx context.Context, req *v1.PaymentsCallbackReq) (res *v1.PaymentsCallbackRes, err error)
	GetPayment(ctx context.Context, req *v1.PaymentsGetReq) (res *v1.PaymentsGetRes, err error)
	CreateRefund(ctx context.Context, req *v1.RefundsCreateReq) (res *v1.RefundsCreateRes, err error)
}

var localPayments IPayments

func Payments() IPayments {
	if localPayments == nil {
		panic("implement not found for interface IPayments, forgot register?")
	}
	return localPayments
}

func RegisterPayments(i IPayments) {
	localPayments = i
}