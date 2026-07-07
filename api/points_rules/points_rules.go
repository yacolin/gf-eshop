package points_rules

import (
	"context"
	"gf-eshop/api/points_rules/v1"
)

type IPointsRulesV1 interface {
	List(ctx context.Context, req *v1.PointsRulesListReq) (res *v1.PointsRulesListRes, err error)
	Detail(ctx context.Context, req *v1.PointsRulesDetailReq) (res *v1.PointsRulesDetailRes, err error)
	Create(ctx context.Context, req *v1.PointsRulesCreateReq) (res *v1.PointsRulesCreateRes, err error)
	Update(ctx context.Context, req *v1.PointsRulesUpdateReq) (res *v1.PointsRulesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.PointsRulesDeleteReq) (res *v1.PointsRulesDeleteRes, err error)
}
