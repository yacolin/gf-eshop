package dashboard

import (
	"context"

	"gf-eshop/api/dashboard/v1"
)

type IDashboardV1 interface {
	Stats(ctx context.Context, req *v1.DashboardStatsReq) (res *v1.DashboardStatsRes, err error)
}
