package service

import (
	"context"

	"gf-eshop/api/marketing/v1"
)

type IMarketing interface {
	// Promotion
	PromotionList(ctx context.Context, req *v1.PromotionListReq) (res *v1.PromotionListRes, err error)
	PromotionDetail(ctx context.Context, req *v1.PromotionDetailReq) (res *v1.PromotionDetailRes, err error)
	PromotionFullDetail(ctx context.Context, req *v1.PromotionFullDetailReq) (res *v1.PromotionFullDetailRes, err error)
	PromotionCreate(ctx context.Context, req *v1.PromotionCreateReq) (res *v1.PromotionCreateRes, err error)
	PromotionUpdate(ctx context.Context, req *v1.PromotionUpdateReq) (res *v1.PromotionUpdateRes, err error)
	PromotionDelete(ctx context.Context, req *v1.PromotionDeleteReq) (res *v1.PromotionDeleteRes, err error)

	// Coupon
	CouponClaim(ctx context.Context, req *v1.CouponClaimReq) (res *v1.CouponClaimRes, err error)
	CouponUse(ctx context.Context, req *v1.CouponUseReq) (res *v1.CouponUseRes, err error)
	CouponList(ctx context.Context, req *v1.CouponListReq) (res *v1.CouponListRes, err error)

	// Flash
	FlashBuy(ctx context.Context, req *v1.FlashBuyReq) (res *v1.FlashBuyRes, err error)
	FlashConfirm(ctx context.Context, req *v1.FlashConfirmReq) (res *v1.FlashConfirmRes, err error)
}

var localMarketing IMarketing

func Marketing() IMarketing {
	if localMarketing == nil {
		panic("implement not found for interface IMarketing, forgot register?")
	}
	return localMarketing
}

func RegisterMarketing(i IMarketing) {
	localMarketing = i
}
