package points_rules

import (
	"context"

	"gf-eshop/api/points_rules/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.PointsRulesListReq) (res *v1.PointsRulesListRes, err error) {
	return service.PointsRules().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.PointsRulesDetailReq) (res *v1.PointsRulesDetailRes, err error) {
	return service.PointsRules().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.PointsRulesCreateReq) (res *v1.PointsRulesCreateRes, err error) {
	return service.PointsRules().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.PointsRulesUpdateReq) (res *v1.PointsRulesUpdateRes, err error) {
	return service.PointsRules().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.PointsRulesDeleteReq) (res *v1.PointsRulesDeleteRes, err error) {
	return service.PointsRules().Delete(ctx, req)
}
