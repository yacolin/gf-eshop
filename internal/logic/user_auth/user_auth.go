package user_auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/user_auth/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sUserAuth struct{}

func init() {
	service.RegisterUserAuth(&sUserAuth{})
}

func userRefreshRedisKey(userId int64, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", utility.UserRefreshRedisKey, userId, tokenId)
}

func saveUserRefreshToken(ctx context.Context, userId int64, tokenId string, expire time.Duration) {
	g.Redis().Do(ctx, "SETEX", userRefreshRedisKey(userId, tokenId), int(expire.Seconds()), "1")
}

func deleteUserRefreshToken(ctx context.Context, userId int64, tokenId string) {
	g.Redis().Do(ctx, "DEL", userRefreshRedisKey(userId, tokenId))
}

func (s *sUserAuth) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	var user *entity.Users
	err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Username, req.Username).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errcode.ErrInvalidCredentials
	}
	if user.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errcode.ErrInvalidCredentials
	}

	now := time.Now()
	pair, err := utility.GenerateUserTokenPair(ctx, user.Id, user.Username)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	ip := ""
	device := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
		device = r.Header.Get("User-Agent")
		if len(device) > 100 {
			device = device[:100]
		}
	}
	_, _ = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, user.Id).Update(do.Users{
		LastLoginIp: ip,
		LastLoginAt: gtime.New(now),
	})
	if _, err := dao.UsrLoginHistories.Ctx(ctx).Insert(do.UsrLoginHistories{
		UserId:      user.Id,
		LoginIp:     ip,
		LoginDevice: device,
		LoginMethod: "password",
		LoginStatus: 1,
	}); err != nil {
		g.Log().Warning(ctx, "insert login history failed: %v", err)
	}

	refreshClaims, _ := utility.ParseUserToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveUserRefreshToken(ctx, user.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.UserLoginRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
		UserId:       user.Id,
		Username:     user.Username,
	}, nil
}

func (s *sUserAuth) Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error) {
	count, err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Username, req.Username).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errcode.ErrUsernameAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	ip := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}

	userId, err := dao.Users.Ctx(ctx).InsertAndGetId(do.Users{
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       1,
		RegisterIp:   ip,
	})
	if err != nil {
		return nil, err
	}

	_, err = dao.Infos.Ctx(ctx).Insert(do.Infos{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	pair, err := utility.GenerateUserTokenPair(ctx, userId, req.Username)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	refreshClaims, _ := utility.ParseUserToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveUserRefreshToken(ctx, userId, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.UserRegisterRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
		UserId:       userId,
		Username:     req.Username,
	}, nil
}

func (s *sUserAuth) RefreshToken(ctx context.Context, req *v1.UserRefreshTokenReq) (res *v1.UserRefreshTokenRes, err error) {
	claims, err := utility.ParseUserToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	if claims.TokenType != utility.TokenTypeRefresh {
		return nil, errcode.ErrInvalidParams
	}

	key := userRefreshRedisKey(claims.UserId, claims.TokenId)
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil || v.IsNil() {
		return nil, errcode.ErrInvalidToken
	}

	deleteUserRefreshToken(ctx, claims.UserId, claims.TokenId)

	var user *entity.Users
	err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, claims.UserId).Scan(&user)
	if err != nil || user == nil {
		return nil, errcode.ErrUserNotFound
	}
	if user.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	pair, err := utility.GenerateUserTokenPair(ctx, user.Id, user.Username)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	refreshClaims, _ := utility.ParseUserToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveUserRefreshToken(ctx, user.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.UserRefreshTokenRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
	}, nil
}
