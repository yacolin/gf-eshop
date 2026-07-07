package marketing

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/marketing/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/utility"
)

func (s *sMarketing) CouponClaim(ctx context.Context, req *v1.CouponClaimReq) (res *v1.CouponClaimRes, err error) {
	userID := utility.GetUserClaims(ctx).UserId

	p, err := s.promotionByID(ctx, req.PromotionId)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}
	if p.PromoType != 1 && p.PromoType != 2 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not a coupon type promotion")
	}
	if p.Status != 2 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "promotion not active")
	}

	now := gtime.Now()
	if now.Before(p.StartTime) || now.After(p.EndTime) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "promotion not in active time window")
	}
	if p.TotalQuantity > 0 && p.UsedQuantity >= p.TotalQuantity {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "coupon sold out")
	}

	if p.PerUserLimit > 0 {
		count, err := dao.UserPromotions.Ctx(ctx).
			Where(dao.UserPromotions.Columns().UserId, userID).
			Where(dao.UserPromotions.Columns().PromotionId, req.PromotionId).
			Count()
		if err != nil {
			return nil, err
		}
		if count >= p.PerUserLimit {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "already claimed")
		}
	}

	sno := generateUserPromotionNo(userID, req.PromotionId)
	upt := dao.UserPromotions.Table()
	pt := dao.Promotions.Table()

	err = dao.Promotions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(upt).Data(do.UserPromotions{
			UserPromotionNo: sno,
			UserId:          userID,
			PromotionId:     req.PromotionId,
			AcquireTime:     now,
			ExpireTime:      p.EndTime,
			Status:          1,
		}).Insert()
		if err != nil {
			return err
		}
		_, err = tx.Model(pt).
			Data(g.Map{"used_quantity": gdb.Raw("used_quantity + 1")}).
			Where("id", req.PromotionId).
			Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	return &v1.CouponClaimRes{}, nil
}

func (s *sMarketing) CouponUse(ctx context.Context, req *v1.CouponUseReq) (res *v1.CouponUseRes, err error) {
	userID := utility.GetUserClaims(ctx).UserId

	_, err = dao.UserPromotions.Ctx(ctx).
		Data(do.UserPromotions{
			Status:   2,
			UsedTime: gtime.Now(),
			OrderId:  req.OrderId,
		}).
		Where(dao.UserPromotions.Columns().Id, req.UserPromotionId).
		Where(dao.UserPromotions.Columns().UserId, userID).
		Where(dao.UserPromotions.Columns().Status, 1).
		Update()
	if err != nil {
		return nil, err
	}

	return &v1.CouponUseRes{}, nil
}

func (s *sMarketing) CouponList(ctx context.Context, req *v1.CouponListReq) (res *v1.CouponListRes, err error) {
	userID := utility.GetUserClaims(ctx).UserId

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.PageSize
	if size <= 0 {
		size = 10
	}

	m := dao.UserPromotions.Ctx(ctx).Where(dao.UserPromotions.Columns().UserId, userID)
	if req.Status != nil {
		m = m.Where(dao.UserPromotions.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.CouponListRes{
			List:  make([]*entity.UserPromotions, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.UserPromotions
	err = m.Page(page, size).OrderDesc(dao.UserPromotions.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.CouponListRes{List: list, Total: total}, nil
}

var promoNoCounter int64

func generateUserPromotionNo(userID, promotionID int64) string {
	promoNoCounter++
	return gtime.Now().Format("YmdHis") + "-" + itoa(promotionID) + "-" + itoa(userID) + "-" + itoa(promoNoCounter)
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
