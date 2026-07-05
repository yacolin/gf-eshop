package brands

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/category_brands/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/service"
)

type sCategoryBrands struct{}

func init() {
	service.RegisterCategoryBrands(&sCategoryBrands{})
}

// List 获取类目下的品牌关联列表（含品牌详情）
func (s *sCategoryBrands) List(ctx context.Context, req *v1.CategoryBrandListReq) (res *v1.CategoryBrandListRes, err error) {
	var (
		m    = dao.CategoryBrands.Ctx(ctx)
		list []*v1.CategoryBrandItem
	)
	err = m.LeftJoin("sp_brands", "sp_brands.id = sp_category_brands.brand_id").
		Fields(
			"sp_category_brands.id",
			"sp_category_brands.category_id",
			"sp_category_brands.brand_id",
			"sp_category_brands.sort_order",
			"sp_brands.name AS brand_name",
			"sp_brands.english_name",
			"sp_brands.logo_url",
			"sp_brands.first_letter",
		).
		Where(dao.CategoryBrands.Columns().CategoryId, req.Id).
		OrderAsc("sp_category_brands.sort_order").
		OrderDesc("sp_category_brands.id").
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*v1.CategoryBrandItem, 0)
	}
	return &v1.CategoryBrandListRes{List: list}, nil
}

// Update 批量替换类目下的品牌关联（先删后插）
func (s *sCategoryBrands) Update(ctx context.Context, req *v1.CategoryBrandUpdateReq) (res *v1.CategoryBrandUpdateRes, err error) {
	// 前置校验：类目是否存在
	categoryCount, err := dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if categoryCount == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "类目不存在")
	}
	// 前置校验：品牌是否存在
	if len(req.BrandIDs) > 0 {
		brandCount, err := dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, req.BrandIDs).Count()
		if err != nil {
			return nil, err
		}
		if brandCount != len(req.BrandIDs) {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "部分品牌不存在")
		}
	}

	err = dao.CategoryBrands.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 删除该类目下所有现有品牌关联
		_, err := dao.CategoryBrands.Ctx(ctx).TX(tx).Where(dao.CategoryBrands.Columns().CategoryId, req.Id).Delete()
		if err != nil {
			return err
		}
		// 2. 插入新的品牌关联
		for _, brandId := range req.BrandIDs {
			_, err = tx.Model("sp_category_brands").Insert(do.CategoryBrands{
				CategoryId: req.Id,
				BrandId:    brandId,
				SortOrder:  req.SortOrder,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.CategoryBrandUpdateRes{}, nil
}
