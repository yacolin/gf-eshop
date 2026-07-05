package dashboard

import (
	"context"

	"gf-eshop/api/dashboard/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.DashboardStatsReq) (res *v1.DashboardStatsRes, err error) {
	return service.Dashboard().Stats(ctx, req)
}
