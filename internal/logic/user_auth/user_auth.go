package user_auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/user_auth/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/internal/verifycode"
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

	pair, err := issueUserSession(ctx, user, "password")
	if err != nil {
		return nil, err
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

	email := utility.NormalizeEmail(req.Email)
	if email != "" {
		emailCount, err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, email).Count()
		if err != nil {
			return nil, err
		}
		if emailCount > 0 {
			return nil, errcode.ErrEmailAlreadyExists
		}
	}

	// 可选：携带邮箱验证码时校验邮箱归属，通过后把 email_verified 置为 1。
	// 先做完占用检查再消费验证码，避免注册失败白白作废一个验证码。
	emailVerified := 0
	if emailCode := strings.TrimSpace(req.EmailCode); emailCode != "" {
		if email == "" {
			return nil, errcode.ErrInvalidParams
		}
		svc, svcErr := verifycode.Default(ctx)
		if svcErr != nil {
			g.Log().Errorf(ctx, "初始化验证码服务失败：%v", svcErr)
			return nil, errcode.ErrVerifyChannelNotReady
		}
		if err = svc.Consume(ctx, verifycode.ChannelEmail, v1.VerifySceneRegister, email, emailCode); err != nil {
			return nil, err
		}
		emailVerified = 1
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	ip := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}

	// email / phone 上有唯一索引且允许为 NULL，而 GoFrame 的 DO 结构体会自动 OmitNilData。
	// 因此未填写时必须留 nil（落 DEFAULT NULL），写空串会撞唯一索引，
	// 导致「第一个不带手机号的用户注册成功后，后续同类注册全部失败」。
	var emailValue, phoneValue interface{}
	if email != "" {
		emailValue = email
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		phoneValue = phone
	}

	userId, err := dao.Users.Ctx(ctx).InsertAndGetId(do.Users{
		Username:      req.Username,
		PasswordHash:  string(hash),
		Email:         emailValue,
		EmailVerified: emailVerified,
		Phone:         phoneValue,
		Status:        1,
		RegisterIp:    ip,
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
