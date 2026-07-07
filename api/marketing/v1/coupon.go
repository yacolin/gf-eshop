package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type CouponClaimReq struct {
	g.Meta      `path:"/coupons/claim" tags:"Marketing" method:"post" summary:"领取优惠券"`
	PromotionId int64 `json:"promotion_id" v:"required"`
}
type CouponClaimRes struct{}

type CouponUseReq struct {
	g.Meta          `path:"/coupons/use" tags:"Marketing" method:"post" summary:"使用优惠券"`
	UserPromotionId int64 `json:"user_promotion_id" v:"required"`
	OrderId         int64 `json:"order_id"          v:"required"`
}
type CouponUseRes struct{}

type CouponListReq struct {
	g.Meta   `path:"/coupons/me" tags:"Marketing" method:"get" summary:"我的优惠券列表"`
	Page     int `json:"page"`      // 页码，默认1
	PageSize int `json:"page_size"` // 每页条数，默认10
	Status   int `json:"status"`    // 按状态筛选
}
type CouponListRes struct {
	List  []*entity.UserPromotions `json:"list"`
	Total int                      `json:"total"`
}
