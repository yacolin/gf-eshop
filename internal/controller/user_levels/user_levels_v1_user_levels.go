package user_levels

import (
	"context"

	"gf-eshop/api/user_levels/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.LevelsListReq) (res *v1.LevelsListRes, err error) {
	return service.UserLevels().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.LevelsDetailReq) (res *v1.LevelsDetailRes, err error) {
	return service.UserLevels().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.LevelsCreateReq) (res *v1.LevelsCreateRes, err error) {
	return service.UserLevels().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.LevelsUpdateReq) (res *v1.LevelsUpdateRes, err error) {
	return service.UserLevels().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.LevelsDeleteReq) (res *v1.LevelsDeleteRes, err error) {
	return service.UserLevels().Delete(ctx, req)
}
