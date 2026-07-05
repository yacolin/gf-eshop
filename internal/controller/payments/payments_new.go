package payments

import (
	"gf-eshop/api/payments"
)

type ControllerV1 struct{}

func NewV1() payments.IPaymentsV1 {
	return &ControllerV1{}
}