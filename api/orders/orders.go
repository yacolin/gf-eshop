package orders

import (
	"context"

	"gf-eshop/api/orders/v1"
)

type IOrdersV1 interface {
	Create(ctx context.Context, req *v1.OrdersCreateReq) (res *v1.OrdersCreateRes, err error)
	List(ctx context.Context, req *v1.OrdersListReq) (res *v1.OrdersListRes, err error)
	Detail(ctx context.Context, req *v1.OrdersDetailReq) (res *v1.OrdersDetailRes, err error)
	UpdateStatus(ctx context.Context, req *v1.OrdersUpdateStatusReq) (res *v1.OrdersUpdateStatusRes, err error)
}