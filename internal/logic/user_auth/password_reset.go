package user_auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/user_auth/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/verifycode"
	"gf-eshop/utility"
)

// ResetPassword 用邮箱验证码重置密码（忘记密码）。
//
// 流程：校验 scene=reset 的验证码 → 改密码 → 吊销该用户全部 refresh token（强制所有设备重新登录）。
//
// 安全约定：
//   - 账号不存在与验证码错误返回同一个错误码，避免邮箱枚举；
//   - 先认证（验证码）再鉴权（账号状态），与邮箱验证码登录保持一致；
//   - 验证码一次性消费：即使后续写库失败也不复用，防止重放。
func (s *sUserAuth) ResetPassword(ctx context.Context, req *v1.UserResetPasswordReq) (res *v1.UserResetPasswordRes, err error) {
	email := utility.NormalizeEmail(req.Email)
	if email == "" {
		return nil, errcode.ErrInvalidParams
	}
	svc, err := verifycode.Default(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "初始化验证码服务失败：%v", err)
		return nil, errcode.ErrVerifyChannelNotReady
	}

	user, err := findUserByTarget(ctx, verifycode.ChannelEmail, email)
	if err != nil {
		return nil, err
	}

	// 先验码、再判账号是否存在：未注册邮箱不可能持有验证码，
	// 因此「邮箱不存在」与「验证码错误」返回同一个错误码（1017），不会泄露注册状态。
	if err = svc.Consume(ctx, verifycode.ChannelEmail, v1.VerifySceneReset, email, strings.TrimSpace(req.Code)); err != nil {
		return nil, err
	}
	if user == nil {
		// 验码已通过却查不到账号，正常不会发生，防御性返回同一错误
		return nil, errcode.ErrVerifyCodeInvalid
	}
	// 走到这里说明验证码已校验通过（只有邮箱持有人才能拿到码），
	// 后续针对账号状态的错误码不构成枚举风险
	if user.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if _, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, user.Id).Update(do.Users{
		PasswordHash: string(hash),
	}); err != nil {
		return nil, err
	}

	// 密码已改成功；吊销失败只记日志，不能让用户看到「重置失败」而重复操作。
	// 注意：access_token 是无状态 JWT，无法回收，最长会在 expire_in（默认 30 分钟）后自然失效。
	revoked, revokeErr := revokeAllUserRefreshTokens(ctx, user.Id)
	if revokeErr != nil {
		g.Log().Errorf(ctx, "重置密码后吊销 refresh token 失败 user_id=%d: %v", user.Id, revokeErr)
	}
	g.Log().Infof(ctx, "用户通过邮箱验证码重置密码成功 user_id=%d，已吊销 %d 个刷新令牌", user.Id, revoked)

	return &v1.UserResetPasswordRes{Success: true, RevokedSessions: revoked}, nil
}

// revokeAllUserRefreshTokens 删除该用户全部 refresh 令牌白名单条目，使其所有会话下线。
//
// 令牌键形如 user:refresh:<userId>:<tokenId>（见 utility.UserRefreshRedisKey）。
// 用 SCAN 游标遍历而不是 KEYS，避免在线上大 key 空间上阻塞 Redis。
func revokeAllUserRefreshTokens(ctx context.Context, userId int64) (int, error) {
	pattern := fmt.Sprintf("%s:%d:*", utility.UserRefreshRedisKey, userId)
	cursor := "0"
	revoked := 0
	for {
		v, err := g.Redis().Do(ctx, "SCAN", cursor, "MATCH", pattern, "COUNT", 200)
		if err != nil {
			return revoked, err
		}
		if v == nil {
			return revoked, nil
		}
		parts := v.Vars()
		if len(parts) < 2 {
			return revoked, nil
		}
		cursor = parts[0].String()
		keys := parts[1].Strings()
		if len(keys) > 0 {
			args := make([]interface{}, 0, len(keys))
			for _, key := range keys {
				args = append(args, key)
			}
			if _, err = g.Redis().Do(ctx, "DEL", args...); err != nil {
				return revoked, err
			}
			revoked += len(keys)
		}
		// SCAN 游标回到 0 表示遍历结束
		if cursor == "0" {
			return revoked, nil
		}
	}
}
