package orders

import (
	"context"

	"gf-eshop/api/orders/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.OrdersCreateReq) (res *v1.OrdersCreateRes, err error) {
	return service.Orders().Create(ctx, req)
}

func (c *ControllerV1) List(ctx context.Context, req *v1.OrdersListReq) (res *v1.OrdersListRes, err error) {
	return service.Orders().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.OrdersDetailReq) (res *v1.OrdersDetailRes, err error) {
	return service.Orders().Detail(ctx, req)
}

func (c *ControllerV1) UpdateStatus(ctx context.Context, req *v1.OrdersUpdateStatusReq) (res *v1.OrdersUpdateStatusRes, err error) {
	return service.Orders().UpdateStatus(ctx, req)
}