package service

import (
	"context"

	"gf-eshop/api/category_brands/v1"
)

type ICategoryBrands interface {
	List(ctx context.Context, req *v1.CategoryBrandListReq) (res *v1.CategoryBrandListRes, err error)
	Update(ctx context.Context, req *v1.CategoryBrandUpdateReq) (res *v1.CategoryBrandUpdateRes, err error)
}

var localCategoryBrands ICategoryBrands

func CategoryBrands() ICategoryBrands {
	if localCategoryBrands == nil {
		panic("implement not found for interface ICategoryBrands, forgot register?")
	}
	return localCategoryBrands
}

func RegisterCategoryBrands(i ICategoryBrands) {
	localCategoryBrands = i
}
