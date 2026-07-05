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

func (s *sBrands) List(ctx context.Context, req *v1.BrandsListReq) (res *v1.BrandsListRes, err error) {
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

	// 无筛选条件时走 Lua 脚本（一次 Redis 往返完成 ZCARD+ZRANGE+MGET）
	if req.Name == "" && req.FirstLetter == "" && req.Status == 0 {
		list, total, err := getBrandPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.BrandsListRes{List: list, Total: total}, nil
		}
		// 缓存不完整，单次重建
		if ctx.Err() != nil {
			return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
		}
		ensureBrandCache(ctx)
		// 重建后重试一次
		list, total, err = getBrandPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.BrandsListRes{List: list, Total: total}, nil
		}
		// 最终兜底：直接查库
		var dbAll []*entity.Brands
		if err := dao.Brands.Ctx(ctx).OrderAsc(dao.Brands.Columns().SortOrder).Scan(&dbAll); err != nil {
			return nil, err
		}
		if len(dbAll) == 0 {
			return &v1.BrandsListRes{List: make([]*entity.Brands, 0), Total: 0}, nil
		}
		return paginateBrands(page, size, dbAll), nil
	}

	// 有筛选条件时直接查询数据库
	var (
		m    = dao.Brands.Ctx(ctx)
		list []*entity.Brands
	)
	if req.Name != "" {
		m = m.WhereLike(dao.Brands.Columns().Name, "%"+req.Name+"%")
	}
	if req.FirstLetter != "" {
		m = m.Where(dao.Brands.Columns().FirstLetter, req.FirstLetter)
	}
	if req.Status > 0 {
		m = m.Where(dao.Brands.Columns().Status, req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.BrandsListRes{
			List:  make([]*entity.Brands, 0),
			Total: 0,
		}, nil
	}

	err = m.Page(page, size).OrderAsc(dao.Brands.Columns().SortOrder).OrderDesc(dao.Brands.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.BrandsListRes{
		List:  list,
		Total: total,
	}, nil
}

func paginateBrands(page, size int, all []*entity.Brands) *v1.BrandsListRes {
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return &v1.BrandsListRes{List: make([]*entity.Brands, 0), Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return &v1.BrandsListRes{List: all[start:end], Total: total}
}

func (s *sBrands) Detail(ctx context.Context, req *v1.BrandsDetailReq) (res *v1.BrandsDetailRes, err error) {
	cached, err := getBrandEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.BrandsDetailRes{Brands: cached}, nil
	}
	if ctx.Err() != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
	}

	var entity *entity.Brands
	err = dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "品牌不存在")
	}
	if err := setBrandEntityCache(context.Background(), entity); err != nil {
		g.Log().Warning(ctx, "setBrandEntityCache failed: %v", err)
	}
	return &v1.BrandsDetailRes{Brands: entity}, nil
}

func (s *sBrands) Create(ctx context.Context, req *v1.BrandsCreateReq) (res *v1.BrandsCreateRes, err error) {
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
	id, _ := result.LastInsertId()
	addBrandToIndex(context.Background(), id, req.SortOrder)
	return &v1.BrandsCreateRes{Id: id}, nil
}

func (s *sBrands) Update(ctx context.Context, req *v1.BrandsUpdateReq) (res *v1.BrandsUpdateRes, err error) {
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
	g.Redis().Do(context.Background(), "ZADD", brandIdsKey, encodeBrandScore(req.SortOrder, req.Id), req.Id)
	delBrandEntityCache(context.Background(), req.Id)
	return &v1.BrandsUpdateRes{}, nil
}

func (s *sBrands) Delete(ctx context.Context, req *v1.BrandsDeleteReq) (res *v1.BrandsDeleteRes, err error) {
	_, err = dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	removeBrandFromIndex(context.Background(), req.Id)
	delBrandEntityCache(context.Background(), req.Id)
	return &v1.BrandsDeleteRes{}, nil
}
