package user_points

import (
	"context"

	"gf-eshop/api/user_points/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.PointsListReq) (res *v1.PointsListRes, err error) {
	return service.UserPoints().List(ctx, req)
}

func (c *ControllerV1) Balance(ctx context.Context, req *v1.PointsBalanceReq) (res *v1.PointsBalanceRes, err error) {
	return service.UserPoints().Balance(ctx, req)
}

func (c *ControllerV1) Trend(ctx context.Context, req *v1.PointsTrendReq) (res *v1.PointsTrendRes, err error) {
	return service.UserPoints().Trend(ctx, req)
}

func (c *ControllerV1) Adjust(ctx context.Context, req *v1.PointsAdjustReq) (res *v1.PointsAdjustRes, err error) {
	return service.UserPoints().Adjust(ctx, req)
}
