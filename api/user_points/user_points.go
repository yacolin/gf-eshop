package user_points

import (
	"context"
	"gf-eshop/api/user_points/v1"
)

type IUserPointsV1 interface {
	List(ctx context.Context, req *v1.PointsListReq) (res *v1.PointsListRes, err error)
	Balance(ctx context.Context, req *v1.PointsBalanceReq) (res *v1.PointsBalanceRes, err error)
	Trend(ctx context.Context, req *v1.PointsTrendReq) (res *v1.PointsTrendRes, err error)
	Adjust(ctx context.Context, req *v1.PointsAdjustReq) (res *v1.PointsAdjustRes, err error)
}
