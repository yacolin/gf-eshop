package level_rules

import (
	"context"

	"gf-eshop/api/level_rules/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.LevelRulesListReq) (res *v1.LevelRulesListRes, err error) {
	return service.LevelRules().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.LevelRulesDetailReq) (res *v1.LevelRulesDetailRes, err error) {
	return service.LevelRules().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.LevelRulesCreateReq) (res *v1.LevelRulesCreateRes, err error) {
	return service.LevelRules().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.LevelRulesUpdateReq) (res *v1.LevelRulesUpdateRes, err error) {
	return service.LevelRules().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.LevelRulesDeleteReq) (res *v1.LevelRulesDeleteRes, err error) {
	return service.LevelRules().Delete(ctx, req)
}
