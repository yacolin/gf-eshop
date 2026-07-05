package categories

import (
	"context"

	"gf-eshop/api/categories/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	return service.Categories().List(ctx, req)
}

func (c *ControllerV1) All(ctx context.Context, req *v1.CategoryAllReq) (res *v1.CategoryAllRes, err error) {
	return service.Categories().All(ctx, req)
}

func (c *ControllerV1) NonRoot(ctx context.Context, req *v1.CategoryNonRootReq) (res *v1.CategoryNonRootRes, err error) {
	return service.Categories().NonRoot(ctx, req)
}

func (c *ControllerV1) Root(ctx context.Context, req *v1.CategoryRootReq) (res *v1.CategoryRootRes, err error) {
	return service.Categories().Root(ctx, req)
}

func (c *ControllerV1) Children(ctx context.Context, req *v1.CategoryChildrenReq) (res *v1.CategoryChildrenRes, err error) {
	return service.Categories().Children(ctx, req)
}

func (c *ControllerV1) Level(ctx context.Context, req *v1.CategoryLevelReq) (res *v1.CategoryLevelRes, err error) {
	return service.Categories().Level(ctx, req)
}
func (c *ControllerV1) Tree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	return service.Categories().Tree(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.CategoryDetailReq) (res *v1.CategoryDetailRes, err error) {
	return service.Categories().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error) {
	return service.Categories().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error) {
	return service.Categories().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error) {
	return service.Categories().Delete(ctx, req)
}