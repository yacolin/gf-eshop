package user_levels

import (
	"context"
	"gf-eshop/api/user_levels/v1"
)

type IUserLevelsV1 interface {
	List(ctx context.Context, req *v1.LevelsListReq) (res *v1.LevelsListRes, err error)
	Detail(ctx context.Context, req *v1.LevelsDetailReq) (res *v1.LevelsDetailRes, err error)
	Create(ctx context.Context, req *v1.LevelsCreateReq) (res *v1.LevelsCreateRes, err error)
	Update(ctx context.Context, req *v1.LevelsUpdateReq) (res *v1.LevelsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.LevelsDeleteReq) (res *v1.LevelsDeleteRes, err error)
}
