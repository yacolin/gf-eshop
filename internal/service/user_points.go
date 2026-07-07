package service

import (
	"context"
	"gf-eshop/api/user_points/v1"
)

type IUserPoints interface {
	List(ctx context.Context, req *v1.PointsListReq) (res *v1.PointsListRes, err error)
	Balance(ctx context.Context, req *v1.PointsBalanceReq) (res *v1.PointsBalanceRes, err error)
	Trend(ctx context.Context, req *v1.PointsTrendReq) (res *v1.PointsTrendRes, err error)
	Adjust(ctx context.Context, req *v1.PointsAdjustReq) (res *v1.PointsAdjustRes, err error)
}

var localUserPoints IUserPoints

func UserPoints() IUserPoints {
	if localUserPoints == nil {
		panic("implement not found for interface IUserPoints, forgot register?")
	}

	return localUserPoints
}

func RegisterUserPoints(i IUserPoints) {
	localUserPoints = i
}
