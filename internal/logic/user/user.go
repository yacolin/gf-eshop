package user

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/user/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) Profile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	claims := utility.GetUserClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}

	var user *entity.Users
	err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, claims.UserId).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errcode.ErrUserNotFound
	}

	var info *entity.Infos
	err = dao.Infos.Ctx(ctx).Where(dao.Infos.Columns().UserId, claims.UserId).Scan(&info)
	if err != nil {
		return nil, err
	}

	return &v1.UserProfileRes{
		Users: user,
		Info:  info,
	}, nil
}

func (s *sUser) UpdateInfo(ctx context.Context, req *v1.UserUpdateInfoReq) (res *v1.UserUpdateInfoRes, err error) {
	claims := utility.GetUserClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}

	if req.Nickname != "" {
		_, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, claims.UserId).
			Update(do.Users{Nickname: req.Nickname, Avatar: req.Avatar})
		if err != nil {
			return nil, err
		}
	}

	var info *entity.Infos
	err = dao.Infos.Ctx(ctx).Where(dao.Infos.Columns().UserId, claims.UserId).Scan(&info)
	if err != nil {
		return nil, err
	}

	updateData := do.Infos{
		Gender:   req.Gender,
		Bio:      req.Bio,
		Country:  req.Country,
		Province: req.Province,
		City:     req.City,
		ZipCode:  req.ZipCode,
		Language: req.Language,
		Timezone: req.Timezone,
	}
	if req.Birthday != "" {
		t, err := gtime.StrToTime(req.Birthday)
		if err == nil {
			updateData.Birthday = t
		}
	}

	if info == nil {
		updateData.UserId = claims.UserId
		_, err = dao.Infos.Ctx(ctx).Insert(updateData)
	} else {
		_, err = dao.Infos.Ctx(ctx).Where(dao.Infos.Columns().UserId, claims.UserId).Update(updateData)
	}
	if err != nil {
		return nil, err
	}

	return &v1.UserUpdateInfoRes{}, nil
}
