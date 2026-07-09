package marketing

import (
	"context"
	"encoding/json"
	"sync"
	"time"

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
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sMarketing struct {
	cache     *multiLevelCache
	jobChan   chan FlashBuyJob
	closeOnce sync.Once
	stopCh    chan struct{}
}

func init() {
	campaignCache = initCache()
	svc := &sMarketing{
		cache:   campaignCache,
		jobChan: make(chan FlashBuyJob, flashJobChanSize),
		stopCh:  make(chan struct{}),
	}
	service.RegisterMarketing(svc)
	go svc.asyncWorker()
}

func (s *sMarketing) PromotionList(ctx context.Context, req *v1.PromotionListReq) (res *v1.PromotionListRes, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.PageSize
	if size <= 0 {
		size = 10
	}

	m := dao.Promotions.Ctx(ctx)
	if req.Status != nil {
		m = m.Where(dao.Promotions.Columns().Status, *req.Status)
	}
	if req.PromoType > 0 {
		m = m.Where(dao.Promotions.Columns().PromoType, req.PromoType)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.PromotionListRes{
			List:  make([]*entity.Promotions, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Promotions
	err = m.Page(page, size).OrderDesc(dao.Promotions.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.PromotionListRes{List: list, Total: total}, nil
}

func (s *sMarketing) promotionByID(ctx context.Context, id int64) (*entity.Promotions, error) {
	if !s.cache.bloomMayExist(id) {
		return nil, nil
	}

	if p, ok := s.cache.getLocal(id); ok {
		return p, nil
	}

	p, err := getRedisEntity(ctx, id)
	if err == nil && p != nil {
		s.cache.setLocal(id, p)
		return p, nil
	}

	var promo *entity.Promotions
	err = dao.Promotions.Ctx(ctx).Where(dao.Promotions.Columns().Id, id).Scan(&promo)
	if err != nil {
		return nil, err
	}
	if promo == nil {
		return nil, nil
	}

	s.cache.bloomAdd(id)
	s.cache.setLocal(id, promo)
	if err := setRedisEntity(context.Background(), promo); err != nil {
		g.Log().Warning(ctx, "setPromotionCache failed: %v", err)
	}
	return promo, nil
}

func (s *sMarketing) PromotionDetail(ctx context.Context, req *v1.PromotionDetailReq) (res *v1.PromotionDetailRes, err error) {
	p, err := s.promotionByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}
	return &v1.PromotionDetailRes{Promotions: p}, nil
}

func (s *sMarketing) PromotionFullDetail(ctx context.Context, req *v1.PromotionFullDetailReq) (res *v1.PromotionFullDetailRes, err error) {
	cached, err := getRedisDetail(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.PromotionFullDetailRes{
			Promotion: cached.Promotion,
			Rule:      cached.Rule,
			Products:  s.enrichProducts(ctx, cached.Products),
		}, nil
	}

	p, err := s.promotionByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}

	var rule *entity.PromotionRules
	if p.RuleId > 0 {
		err = dao.PromotionRules.Ctx(ctx).Where(dao.PromotionRules.Columns().Id, p.RuleId).Scan(&rule)
		if err != nil {
			return nil, err
		}
	}

	var products []*entity.PromotionProducts
	err = dao.PromotionProducts.Ctx(ctx).
		Where(dao.PromotionProducts.Columns().PromotionId, req.Id).
		OrderAsc(dao.PromotionProducts.Columns().ProductType).
		OrderAsc(dao.PromotionProducts.Columns().Id).
		Scan(&products)
	if err != nil {
		return nil, err
	}

	detail := &promotionDetail{
		Promotion: p,
		Rule:      rule,
		Products:  products,
	}
	if err := setRedisDetail(context.Background(), req.Id, detail); err != nil {
		g.Log().Warning(ctx, "setPromotionDetailCache failed: %v", err)
	}

	return &v1.PromotionFullDetailRes{
		Promotion: p,
		Rule:      rule,
		Products:  s.enrichProducts(ctx, products),
	}, nil
}

func (s *sMarketing) enrichProducts(ctx context.Context, products []*entity.PromotionProducts) []*v1.PromotionProductItem {
	if len(products) == 0 {
		return make([]*v1.PromotionProductItem, 0)
	}

	spuIDs := make([]int64, 0, len(products))
	for _, pp := range products {
		if pp.ProductType == 3 && pp.TargetId > 0 {
			spuIDs = append(spuIDs, pp.TargetId)
		}
	}

	spuMap := make(map[int64]*entity.Products)
	if len(spuIDs) > 0 {
		var spus []*entity.Products
		err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, spuIDs).Scan(&spus)
		if err == nil {
			for _, spu := range spus {
				spuMap[spu.Id] = spu
			}
		}
	}

	// 批量查询 SKU 价格区间
	type priceRow struct {
		ProductId int64 `orm:"product_id"`
		MinPrice  int64 `orm:"min_price"`
		MaxPrice  int64 `orm:"max_price"`
	}
	priceMap := make(map[int64]struct{ min, max int64 })
	if len(spuIDs) > 0 {
		var rows []priceRow
		err := g.DB().Model("sp_skus").
			Fields("product_id", "MIN(price) AS min_price", "MAX(price) AS max_price").
			Where("product_id IN (?)", spuIDs).
			Where("status", 1).
			Where("deleted_at IS NULL").
			Group("product_id").
			Scan(&rows)
		if err == nil {
			for _, r := range rows {
				priceMap[r.ProductId] = struct{ min, max int64 }{min: r.MinPrice, max: r.MaxPrice}
			}
		}
	}

	items := make([]*v1.PromotionProductItem, 0, len(products))
	for _, pp := range products {
		item := &v1.PromotionProductItem{
			Id:          pp.Id,
			ProductType: pp.ProductType,
			ProductId:   pp.TargetId,
		}
		if spu, ok := spuMap[pp.TargetId]; ok {
			item.SpuName = spu.Name
			item.Subtitle = spu.Subtitle
			item.MainImage = spu.MainImage
			item.Unit = spu.Unit
			item.SalesCount = spu.SalesCount
			item.SpuStatus = spu.Status
		}
		if pr, ok := priceMap[pp.TargetId]; ok {
			item.MinPrice = pr.min
			item.MaxPrice = pr.max
		}
		items = append(items, item)
	}
	return items
}

func (s *sMarketing) PromotionCreate(ctx context.Context, req *v1.PromotionCreateReq) (res *v1.PromotionCreateRes, err error) {
	st, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid start_time format")
	}
	et, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid end_time format")
	}

	userID := utility.GetUserClaims(ctx).UserId

	var promoID int64
	pt := dao.Promotions.Table()
	prt := dao.PromotionRules.Table()
	pp := dao.PromotionProducts.Table()

	err = dao.Promotions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(pt).Data(do.Promotions{
			PromoName:     req.PromoName,
			PromoType:     req.PromoType,
			PromoCode:     promoCodeOrNil(req.PromoCode),
			StartTime:     gtime.New(st),
			EndTime:       gtime.New(et),
			TotalQuantity: req.TotalQuantity,
			PerUserLimit:  req.PerUserLimit,
			CreatedBy:     userID,
			UpdatedBy:     userID,
		}).Insert()
		if err != nil {
			return err
		}
		promoID, _ = result.LastInsertId()

		if req.BenefitType > 0 {
			ruleResult, err := tx.Model(prt).Data(do.PromotionRules{
				PromotionId:    promoID,
				RuleName:       req.RuleName,
				ConditionType:  req.ConditionType,
				ConditionValue: req.ConditionValue,
				BenefitConfig:  benefitConfigJSON(req.BenefitType, req.BenefitValue),
				IsStackable:    req.IsStackable,
				StackGroup:     req.StackPriority,
				CreatedBy:      userID,
				UpdatedBy:      userID,
			}).Insert()
			if err != nil {
				return err
			}
			ruleID, _ := ruleResult.LastInsertId()
			_, err = tx.Model(pt).
				Data(g.Map{"rule_id": ruleID}).
				Where("id", promoID).
				Update()
			if err != nil {
				return err
			}
		}

		if len(req.ProductIDs) > 0 {
			for _, pid := range req.ProductIDs {
				_, err = tx.Model(pp).Data(do.PromotionProducts{
					PromotionId: promoID,
					ProductType: 3,
					TargetId:    pid,
				}).Insert()
				if err != nil {
					return err
				}
			}
		}

		s.cache.bloomAdd(promoID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.PromotionCreateRes{Id: promoID}, nil
}

func (s *sMarketing) PromotionUpdate(ctx context.Context, req *v1.PromotionUpdateReq) (res *v1.PromotionUpdateRes, err error) {
	p, err := s.promotionByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}

	userID := utility.GetUserClaims(ctx).UserId

	updates := g.Map{}
	if req.PromoName != "" {
		updates["promo_name"] = req.PromoName
	}
	if req.StartTime != "" {
		pt, e := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if e != nil {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid start_time format")
		}
		updates["start_time"] = gtime.New(pt)
	}
	if req.EndTime != "" {
		pt, e := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if e != nil {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid end_time format")
		}
		updates["end_time"] = gtime.New(pt)
	}
	if req.TotalQuantity > 0 {
		updates["total_quantity"] = req.TotalQuantity
	}
	if req.PerUserLimit > 0 {
		updates["per_user_limit"] = req.PerUserLimit
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	updates["updated_by"] = userID

	_, err = dao.Promotions.Ctx(ctx).Data(updates).Where(dao.Promotions.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}

	if req.BenefitType > 0 {
		ruleUpdates := g.Map{"updated_by": userID}
		ruleUpdates["is_stackable"] = req.IsStackable
		if req.RuleName != "" {
			ruleUpdates["rule_name"] = req.RuleName
		}
		if req.ConditionType > 0 {
			ruleUpdates["condition_type"] = req.ConditionType
		}
		ruleUpdates["condition_value"] = req.ConditionValue
		ruleUpdates["benefit_config"] = benefitConfigJSON(req.BenefitType, req.BenefitValue)
		if req.StackPriority > 0 {
			ruleUpdates["stack_group"] = req.StackPriority
		}
		_, err = dao.PromotionRules.Ctx(ctx).
			Data(ruleUpdates).
			Where(dao.PromotionRules.Columns().PromotionId, req.Id).
			Update()
		if err != nil {
			return nil, err
		}
	}

	delRedisEntity(context.Background(), req.Id)
	delRedisDetail(context.Background(), req.Id)
	s.cache.delLocal(req.Id)
	time.AfterFunc(500*time.Millisecond, func() {
		delRedisEntity(context.Background(), req.Id)
		delRedisDetail(context.Background(), req.Id)
	})

	return &v1.PromotionUpdateRes{}, nil
}

func (s *sMarketing) PromotionDelete(ctx context.Context, req *v1.PromotionDeleteReq) (res *v1.PromotionDeleteRes, err error) {
	p, err := s.promotionByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}

	pt := dao.Promotions.Table()
	prt := dao.PromotionRules.Table()
	pp := dao.PromotionProducts.Table()

	err = dao.Promotions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(prt).Where("promotion_id", req.Id).Delete()
		if err != nil {
			return err
		}
		_, err = tx.Model(pp).Where("promotion_id", req.Id).Delete()
		if err != nil {
			return err
		}
		_, err = tx.Model(pt).Where("id", req.Id).Delete()
		return err
	})
	if err != nil {
		return nil, err
	}

	delRedisEntity(context.Background(), req.Id)
	delRedisDetail(context.Background(), req.Id)
	s.cache.delLocal(req.Id)
	return &v1.PromotionDeleteRes{}, nil
}

func promoCodeOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

type benefitCfg struct {
	Type  int   `json:"type"`
	Value int64 `json:"value"`
}

func benefitConfigJSON(benefitType int, benefitValue int64) string {
	if benefitType == 0 {
		return ""
	}
	b, _ := json.Marshal(benefitCfg{Type: benefitType, Value: benefitValue})
	return string(b)
}
