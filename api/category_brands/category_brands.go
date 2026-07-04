package category_brands


import (
	"context"

	"gf-eshop/api/category_brands/v1"
)

type ICategoryBrandsV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
}