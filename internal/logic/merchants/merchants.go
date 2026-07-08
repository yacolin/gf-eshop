package merchants

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/api/merchants/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchants struct{}

func init() {
	service.RegisterMerchants(&sMerchants{})
}

func (s *sMerchants) List(ctx context.Context, req *v1.MerchantsListReq) (res *v1.MerchantsListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.Merchants.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.MerchantsListRes{
			List:  make([]*entity.Merchants, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.Merchants
	err = m.Page(page, size).OrderDesc(dao.Merchants.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantsListRes{List: list, Total: total}, nil
}

func (s *sMerchants) Detail(ctx context.Context, req *v1.MerchantsDetailReq) (res *v1.MerchantsDetailRes, err error) {
	cached, err := getMerchantEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.MerchantsDetailRes{Merchants: cached}, nil
	}

	var entity *entity.Merchants
	err = dao.Merchants.Ctx(ctx).Where(dao.Merchants.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrMerchantsNotFound
	}
	if err := setMerchantEntityCache(context.Background(), entity); err != nil {
		g.Log().Warning(ctx, "setMerchantEntityCache failed: %v", err)
	}
	return &v1.MerchantsDetailRes{Merchants: entity}, nil
}

func (s *sMerchants) Create(ctx context.Context, req *v1.MerchantsCreateReq) (res *v1.MerchantsCreateRes, err error) {
	merchantCode := fmt.Sprintf("M%s%s", gtime.Now().Format("YmdHis"), grand.Letters(4))

	result, err := dao.Merchants.Ctx(ctx).Insert(do.Merchants{
		MerchantName:    req.MerchantName,
		MerchantCode:    merchantCode,
		MerchantType:    req.MerchantType,
		MerchantLevel:   req.MerchantLevel,
		BusinessScope:   req.BusinessScope,
		BusinessYears:   req.BusinessYears,
		ContactPerson:   req.ContactPerson,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		LogoUrl:         req.LogoUrl,
		BannerUrl:       req.BannerUrl,
		ShopDescription: req.ShopDescription,
		Status:          req.Status,
		CommissionRate:  req.CommissionRate,
		SettlementCycle: req.SettlementCycle,
		SettledAt:       req.SettledAt,
		ExpireAt:        req.ExpireAt,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MerchantsCreateRes{Id: id}, nil
}

func (s *sMerchants) Update(ctx context.Context, req *v1.MerchantsUpdateReq) (res *v1.MerchantsUpdateRes, err error) {
	count, err := dao.Merchants.Ctx(ctx).Where(dao.Merchants.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantsNotFound
	}

	_, err = dao.Merchants.Ctx(ctx).Data(do.Merchants{
		MerchantName:    req.MerchantName,
		MerchantType:    req.MerchantType,
		MerchantLevel:   req.MerchantLevel,
		BusinessScope:   req.BusinessScope,
		BusinessYears:   req.BusinessYears,
		ContactPerson:   req.ContactPerson,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		LogoUrl:         req.LogoUrl,
		BannerUrl:       req.BannerUrl,
		ShopDescription: req.ShopDescription,
		Status:          req.Status,
		AuditStatus:     req.AuditStatus,
		AuditReason:     req.AuditReason,
		AuditedAt:       req.AuditedAt,
		FrozenReason:    req.FrozenReason,
		CommissionRate:  req.CommissionRate,
		SettlementCycle: req.SettlementCycle,
		ExpireAt:        req.ExpireAt,
	}).Where(dao.Merchants.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	delMerchantEntityCache(context.Background(), req.Id)
	return &v1.MerchantsUpdateRes{}, nil
}

func (s *sMerchants) Delete(ctx context.Context, req *v1.MerchantsDeleteReq) (res *v1.MerchantsDeleteRes, err error) {
	_, err = dao.Merchants.Ctx(ctx).Where(dao.Merchants.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	delMerchantEntityCache(context.Background(), req.Id)
	return &v1.MerchantsDeleteRes{}, nil
}