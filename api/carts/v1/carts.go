package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CartItemResponse struct {
	SkuID       int64  `json:"sku_id"`
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	SkuSpec     string `json:"sku_spec,omitempty"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}

type CartResponse struct {
	ID          int64              `json:"id"`
	ItemCount   int                `json:"item_count"`
	TotalAmount int64              `json:"total_amount"`
	Items       []CartItemResponse `json:"items"`
}

type CartsGetReq struct {
	g.Meta `path:"" tags:"Carts" method:"get" summary:"获取购物车"`
}

type CartsGetRes struct {
	*CartResponse
}

type CartsAddItemReq struct {
	g.Meta  `path:"/items" tags:"Carts" method:"post" summary:"添加商品到购物车"`
	SkuID   int64 `json:"sku_id"   v:"required" description:"SKU ID"`
	Quantity int  `json:"quantity" v:"required|min:1|max:99" description:"数量"`
}

type CartsAddItemRes struct {
	*CartResponse
}

type CartsUpdateItemReq struct {
	g.Meta   `path:"/items" tags:"Carts" method:"put" summary:"更新购物车商品数量"`
	SkuID    int64 `json:"sku_id"   v:"required" description:"SKU ID"`
	Quantity int   `json:"quantity" v:"min:0|max:99" description:"数量（为0时删除）"`
}

type CartsUpdateItemRes struct {
	*CartResponse
}

type CartsRemoveItemReq struct {
	g.Meta `path:"/items/{sku_id}" tags:"Carts" method:"delete" summary:"删除购物车商品"`
	SkuID  int64 `json:"sku_id" v:"required" description:"SKU ID"`
}

type CartsRemoveItemRes struct {
	*CartResponse
}

type CartsClearReq struct {
	g.Meta `path:"/clear" tags:"Carts" method:"post" summary:"清空购物车"`
}

type CartsClearRes struct {
}
