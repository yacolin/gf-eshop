package service

import (
	"context"
	"gf-eshop/api/user_levels/v1"
)

type IUserLevels interface {
	List(ctx context.Context, req *v1.LevelsListReq) (res *v1.LevelsListRes, err error)
	Detail(ctx context.Context, req *v1.LevelsDetailReq) (res *v1.LevelsDetailRes, err error)
	Create(ctx context.Context, req *v1.LevelsCreateReq) (res *v1.LevelsCreateRes, err error)
	Update(ctx context.Context, req *v1.LevelsUpdateReq) (res *v1.LevelsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.LevelsDeleteReq) (res *v1.LevelsDeleteRes, err error)
}

var localUserLevels IUserLevels

func UserLevels() IUserLevels {
	if localUserLevels == nil {
		panic("implement not found for interface IUserLevels, forgot register?")
	}
	return localUserLevels
}

func RegisterUserLevels(i IUserLevels) {
	localUserLevels = i
}
