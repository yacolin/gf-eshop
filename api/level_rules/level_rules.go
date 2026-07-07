package level_rules

import (
	"context"
	"gf-eshop/api/level_rules/v1"
)

type ILevelRulesV1 interface {
	List(ctx context.Context, req *v1.LevelRulesListReq) (res *v1.LevelRulesListRes, err error)
	Detail(ctx context.Context, req *v1.LevelRulesDetailReq) (res *v1.LevelRulesDetailRes, err error)
	Create(ctx context.Context, req *v1.LevelRulesCreateReq) (res *v1.LevelRulesCreateRes, err error)
	Update(ctx context.Context, req *v1.LevelRulesUpdateReq) (res *v1.LevelRulesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.LevelRulesDeleteReq) (res *v1.LevelRulesDeleteRes, err error)
}
