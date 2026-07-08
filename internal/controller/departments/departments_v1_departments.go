package departments

import (
	"context"

	"gf-eshop/api/departments/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	return service.Departments().List(ctx, req)
}

func (c *ControllerV1) All(ctx context.Context, req *v1.AllReq) (res *v1.AllRes, err error) {
	return service.Departments().All(ctx, req)
}

func (c *ControllerV1) Children(ctx context.Context, req *v1.ChildrenReq) (res *v1.ChildrenRes, err error) {
	return service.Departments().Children(ctx, req)
}

func (c *ControllerV1) Tree(ctx context.Context, req *v1.TreeReq) (res *v1.TreeRes, err error) {
	return service.Departments().Tree(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	return service.Departments().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	return service.Departments().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	return service.Departments().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	return service.Departments().Delete(ctx, req)
}
