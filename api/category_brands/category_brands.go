package category_brands


import (
	"context"

	"gf-eshop/api/category_brands/v1"
)

type ICategoryBrandsV1 interface {
	List(ctx context.Context, req *v1.CategoryBrandListReq) (res *v1.CategoryBrandListRes, err error)
	Update(ctx context.Context, req *v1.CategoryBrandUpdateReq) (res *v1.CategoryBrandUpdateRes, err error)
}