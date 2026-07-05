package skus

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/skus/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sSkus struct{}

func init() {
	service.RegisterSkus(&sSkus{})
}

func (s *sSkus) List(ctx context.Context, req *v1.SkusListReq) (res *v1.SkusListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.Skus.Ctx(ctx)
	if req.ProductId > 0 {
		m = m.Where(dao.Skus.Columns().ProductId, req.ProductId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.SkusListRes{
			List:  make([]*entity.Skus, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Skus
	err = m.Page(page, size).OrderAsc(dao.Skus.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.SkusListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sSkus) Detail(ctx context.Context, req *v1.SkusDetailReq) (res *v1.SkusDetailRes, err error) {
	var entity *entity.Skus
	err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "SKU不存在")
	}
	return &v1.SkusDetailRes{Skus: entity}, nil
}

func (s *sSkus) GetByCode(ctx context.Context, req *v1.SkusGetByCodeReq) (res *v1.SkusGetByCodeRes, err error) {
	var entity *entity.Skus
	err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().SkuCode, req.SkuCode).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "SKU不存在")
	}
	return &v1.SkusGetByCodeRes{Skus: entity}, nil
}

func (s *sSkus) Create(ctx context.Context, req *v1.SkusCreateReq) (res *v1.SkusCreateRes, err error) {
	result, err := dao.Skus.Ctx(ctx).Insert(do.Skus{
		ProductId:    req.ProductId,
		SkuCode:      req.SkuCode,
		Barcode:      req.Barcode,
		Spec:         req.Spec,
		Price:        req.Price,
		MarketPrice:  req.MarketPrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
		Volume:       req.Volume,
		Length:       req.Length,
		Width:        req.Width,
		Height:       req.Height,
		MinPurchaseQty: req.MinPurchaseQty,
		MaxPurchaseQty: req.MaxPurchaseQty,
		Image:        req.Image,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.SkusCreateRes{Id: id}, nil
}

func (s *sSkus) Update(ctx context.Context, req *v1.SkusUpdateReq) (res *v1.SkusUpdateRes, err error) {
	count, err := dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "SKU不存在")
	}

	_, err = dao.Skus.Ctx(ctx).Data(do.Skus{
		SkuCode:      req.SkuCode,
		Barcode:      req.Barcode,
		Price:        req.Price,
		MarketPrice:  req.MarketPrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
		Volume:       req.Volume,
		Length:       req.Length,
		Width:        req.Width,
		Height:       req.Height,
		MinPurchaseQty: req.MinPurchaseQty,
		MaxPurchaseQty: req.MaxPurchaseQty,
		Image:        req.Image,
		Status:       req.Status,
	}).Where(dao.Skus.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.SkusUpdateRes{}, nil
}

func (s *sSkus) Delete(ctx context.Context, req *v1.SkusDeleteReq) (res *v1.SkusDeleteRes, err error) {
	_, err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.SkusDeleteRes{}, nil
}
