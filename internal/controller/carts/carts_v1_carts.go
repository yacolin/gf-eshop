package carts

import (
	"context"

	"gf-eshop/api/carts/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) GetCart(ctx context.Context, req *v1.CartsGetReq) (res *v1.CartsGetRes, err error) {
	return service.Carts().GetCart(ctx, req)
}

func (c *ControllerV1) AddItem(ctx context.Context, req *v1.CartsAddItemReq) (res *v1.CartsAddItemRes, err error) {
	return service.Carts().AddItem(ctx, req)
}

func (c *ControllerV1) UpdateItem(ctx context.Context, req *v1.CartsUpdateItemReq) (res *v1.CartsUpdateItemRes, err error) {
	return service.Carts().UpdateItem(ctx, req)
}

func (c *ControllerV1) RemoveItem(ctx context.Context, req *v1.CartsRemoveItemReq) (res *v1.CartsRemoveItemRes, err error) {
	return service.Carts().RemoveItem(ctx, req)
}

func (c *ControllerV1) ClearCart(ctx context.Context, req *v1.CartsClearReq) (res *v1.CartsClearRes, err error) {
	return service.Carts().ClearCart(ctx, req)
}
