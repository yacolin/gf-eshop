package staff

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/staff/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sStaff struct{}

func init() {
	service.RegisterStaff(&sStaff{})
}

func refreshRedisKey(staffId int64, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", utility.RefreshRedisKey, staffId, tokenId)
}

func saveRefreshToken(ctx context.Context, staffId int64, tokenId string, expire time.Duration) {
	key := refreshRedisKey(staffId, tokenId)
	g.Redis().Do(ctx, "SETEX", key, int(expire.Seconds()), "1")
}

func deleteRefreshToken(ctx context.Context, staffId int64, tokenId string) {
	g.Redis().Do(ctx, "DEL", refreshRedisKey(staffId, tokenId))
}

func deleteAllRefreshTokens(ctx context.Context, staffId int64) {
	pattern := fmt.Sprintf("%s:%d:*", utility.RefreshRedisKey, staffId)
	keys, _ := g.Redis().Do(ctx, "KEYS", pattern)
	if !keys.IsNil() {
		for _, key := range keys.Vars() {
			g.Redis().Do(ctx, "DEL", key.String())
		}
	}
}

func (s *sStaff) Login(ctx context.Context, req *v1.StaffLoginReq) (res *v1.StaffLoginRes, err error) {
	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Username, req.Username).
		Scan(&staff)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errcode.ErrInvalidCredentials
	}
	if staff.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errcode.ErrInvalidCredentials
	}

	now := time.Now()
	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
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
	_, _ = dao.Staff.Ctx(ctx).Where(dao.Staff.Columns().Id, staff.Id).Update(do.Staff{
		LastLoginIp: ip,
		LastLoginAt: gtime.New(now),
	})
	if _, err := dao.SysLoginHistories.Ctx(ctx).Insert(do.SysLoginHistories{
		StaffId:     staff.Id,
		LoginIp:     ip,
		LoginDevice: device,
		LoginMethod: "password",
		LoginStatus: 1,
	}); err != nil {
		g.Log().Warning(ctx, "insert login history failed: %v", err)
	}

	refreshClaims, _ := utility.ParseStaffToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveRefreshToken(ctx, staff.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.StaffLoginRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
		StaffId:      staff.Id,
		Username:     staff.Username,
		RealName:     staff.RealName,
	}, nil
}

func (s *sStaff) RefreshToken(ctx context.Context, req *v1.StaffRefreshTokenReq) (res *v1.StaffRefreshTokenRes, err error) {
	claims, err := utility.ParseStaffToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	if claims.TokenType != utility.TokenTypeRefresh {
		return nil, errcode.ErrInvalidParams
	}

	key := refreshRedisKey(claims.StaffId, claims.TokenId)
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil || v.IsNil() {
		return nil, errcode.ErrInvalidToken
	}

	deleteRefreshToken(ctx, claims.StaffId, claims.TokenId)

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil || staff == nil {
		return nil, errcode.ErrUserNotFound
	}
	if staff.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	refreshClaims, _ := utility.ParseStaffToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveRefreshToken(ctx, staff.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.StaffRefreshTokenRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
	}, nil
}

func (s *sStaff) Logout(ctx context.Context, req *v1.StaffLogoutReq) (res *v1.StaffLogoutRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims != nil {
		deleteAllRefreshTokens(ctx, claims.StaffId)
	}
	return &v1.StaffLogoutRes{}, nil
}

func (s *sStaff) Profile(ctx context.Context, req *v1.StaffProfileReq) (res *v1.StaffProfileRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errcode.ErrUserNotFound
	}

	return &v1.StaffProfileRes{
		Id:          staff.Id,
		Username:    staff.Username,
		RealName:    staff.RealName,
		Email:       staff.Email,
		Phone:       staff.Phone,
		Avatar:      staff.Avatar,
		Status:      staff.Status,
		LastLoginIp: staff.LastLoginIp,
	}, nil
}

func (s *sStaff) Permissions(ctx context.Context, req *v1.StaffPermissionsReq) (res *v1.StaffPermissionsRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}
	perms, roles, err := service.Roles().GetPermissions(ctx, claims.StaffId)
	if err != nil {
		return nil, err
	}
	return &v1.StaffPermissionsRes{
		Roles:       roles,
		Permissions: perms,
	}, nil
}


