package carts

import (
	"context"

	"gf-eshop/api/carts/v1"
)

type ICartsV1 interface {
	GetCart(ctx context.Context, req *v1.CartsGetReq) (res *v1.CartsGetRes, err error)
	AddItem(ctx context.Context, req *v1.CartsAddItemReq) (res *v1.CartsAddItemRes, err error)
	UpdateItem(ctx context.Context, req *v1.CartsUpdateItemReq) (res *v1.CartsUpdateItemRes, err error)
	RemoveItem(ctx context.Context, req *v1.CartsRemoveItemReq) (res *v1.CartsRemoveItemRes, err error)
	ClearCart(ctx context.Context, req *v1.CartsClearReq) (res *v1.CartsClearRes, err error)
}
