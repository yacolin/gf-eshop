package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type FlashBuyReq struct {
	g.Meta      `path:"/flash/buy" tags:"Marketing" method:"post" summary:"秒杀下单"`
	PromotionId int64 `json:"promotion_id" v:"required"`
	ProductId   int64 `json:"product_id"   v:"required"`
	SkuId       int64 `json:"sku_id"       v:"required"`
	Quantity    int   `json:"quantity"     v:"required|between:1,99"`
}
type FlashBuyRes struct {
	Token string `json:"token"`
}

type FlashConfirmReq struct {
	g.Meta    `path:"/flash/confirm" tags:"Marketing" method:"post" summary:"确认秒杀订单"`
	Token     string `json:"token"      v:"required"`
	AddressId int64  `json:"address_id" v:"required"`
}
type FlashConfirmRes struct{}
