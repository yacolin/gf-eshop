package service

import (
	"context"

	"gf-eshop/api/level_rules/v1"
)

type ILevelRules interface {
	List(ctx context.Context, req *v1.LevelRulesListReq) (res *v1.LevelRulesListRes, err error)
	Detail(ctx context.Context, req *v1.LevelRulesDetailReq) (res *v1.LevelRulesDetailRes, err error)
	Create(ctx context.Context, req *v1.LevelRulesCreateReq) (res *v1.LevelRulesCreateRes, err error)
	Update(ctx context.Context, req *v1.LevelRulesUpdateReq) (res *v1.LevelRulesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.LevelRulesDeleteReq) (res *v1.LevelRulesDeleteRes, err error)
}

var localLevelRules ILevelRules

func LevelRules() ILevelRules {
	if localLevelRules == nil {
		panic("implement not found for interface ILevelRules, forgot register?")
	}
	return localLevelRules
}

func RegisterLevelRules(i ILevelRules) {
	localLevelRules = i
}
