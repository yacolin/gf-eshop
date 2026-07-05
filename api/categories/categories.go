package categories

import (
	"context"

	"gf-eshop/api/categories/v1"
)

type ICategoriesV1 interface {
	List(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error)
	All(ctx context.Context, req *v1.CategoryAllReq) (res *v1.CategoryAllRes, err error)
	Root(ctx context.Context, req *v1.CategoryRootReq) (res *v1.CategoryRootRes, err error)
	Children(ctx context.Context, req *v1.CategoryChildrenReq) (res *v1.CategoryChildrenRes, err error)
	Level(ctx context.Context, req *v1.CategoryLevelReq) (res *v1.CategoryLevelRes, err error)
	Tree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error)
	Detail(ctx context.Context, req *v1.CategoryDetailReq) (res *v1.CategoryDetailRes, err error)
	Create(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error)
	Update(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error)
	Delete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error)
}
