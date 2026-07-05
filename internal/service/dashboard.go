package service

import (
	"context"

	"gf-eshop/api/dashboard/v1"
)

type IDashboard interface {
	Stats(ctx context.Context, req *v1.DashboardStatsReq) (res *v1.DashboardStatsRes, err error)
	StartPeriodicRefresh(ctx context.Context)
}

var localDashboard IDashboard

func Dashboard() IDashboard {
	if localDashboard == nil {
		panic("implement not found for interface IDashboard, forgot register?")
	}
	return localDashboard
}

func RegisterDashboard(i IDashboard) {
	localDashboard = i
}
