package staff

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/staff/v1"
	"gf-eshop/internal/dao"
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
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "用户名或密码错误")
	}
	if staff.Status != 1 {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "账号已被禁用")
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "用户名或密码错误")
	}

	now := time.Now()
	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "生成Token失败")
	}

	ip := g.RequestFromCtx(ctx).GetClientIp()
	_, _ = dao.Staff.Ctx(ctx).Where(dao.Staff.Columns().Id, staff.Id).Update(do.Staff{
		LastLoginIp: ip,
		LastLoginAt: gtime.New(now),
	})

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
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "RefreshToken无效或已过期")
	}
	if claims.TokenType != utility.TokenTypeRefresh {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "无效的Token类型")
	}

	key := refreshRedisKey(claims.StaffId, claims.TokenId)
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil || v.IsNil() {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "RefreshToken已失效")
	}

	deleteRefreshToken(ctx, claims.StaffId, claims.TokenId)

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil || staff == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "用户不存在")
	}
	if staff.Status != 1 {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "账号已被禁用")
	}

	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "生成Token失败")
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
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "用户不存在")
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


