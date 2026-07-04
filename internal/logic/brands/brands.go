package brands

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/brands/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sBrands struct{}

func init() {
	service.RegisterBrands(&sBrands{})
}

func (s *sBrands) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	var (
		m      = dao.Brands.Ctx(ctx)
		list   []*entity.Brands
		total  int
		page   = req.Page
		size   = req.PageSize
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	// Conditions
	if req.Name != "" {
		m = m.WhereLike(dao.Brands.Columns().Name, "%"+req.Name+"%")
	}
	if req.FirstLetter != "" {
		m = m.Where(dao.Brands.Columns().FirstLetter, req.FirstLetter)
	}
	if req.Status > 0 {
		m = m.Where(dao.Brands.Columns().Status, req.Status)
	}

	total, err = m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ListRes{
			List:  make([]*entity.Brands, 0),
			Total: 0,
		}, nil
	}

	err = m.Page(page, size).OrderAsc(dao.Brands.Columns().SortOrder).OrderDesc(dao.Brands.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sBrands) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	// 尝试从缓存读取
	cached, err := getBrandEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.DetailRes{Brands: cached}, nil
	}

	var entity *entity.Brands
	err = dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "品牌不存在")
	}

	// 回写缓存
	if err := setBrandEntityCache(ctx, entity); err != nil {
		g.Log().Warning(ctx, "setBrandEntityCache failed: %v", err)
	}

	return &v1.DetailRes{
		Brands: entity,
	}, nil
}

func (s *sBrands) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	result, err := dao.Brands.Ctx(ctx).Insert(do.Brands{
		Name:        req.Name,
		EnglishName: req.EnglishName,
		LogoUrl:     req.LogoUrl,
		FirstLetter: req.FirstLetter,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	// 缓存失效
	delBrandAllCache(ctx)
	return &v1.CreateRes{
		Id: lastInsertId,
	}, nil
}

func (s *sBrands) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	count, err := dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "品牌不存在")
	}

	_, err = dao.Brands.Ctx(ctx).Data(do.Brands{
		Name:        req.Name,
		EnglishName: req.EnglishName,
		LogoUrl:     req.LogoUrl,
		FirstLetter: req.FirstLetter,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		Description: req.Description,
	}).Where(dao.Brands.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	// 缓存失效
	delBrandEntityCache(ctx, req.Id)
	delBrandAllCache(ctx)
	return &v1.UpdateRes{}, nil
}

func (s *sBrands) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	_, err = dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	// 缓存失效
	delBrandEntityCache(ctx, req.Id)
	delBrandAllCache(ctx)
	return &v1.DeleteRes{}, nil
}
