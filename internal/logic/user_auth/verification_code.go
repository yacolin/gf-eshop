package user_auth

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/user_auth/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/verifycode"
	"gf-eshop/utility"
)

// 本文件是「验证码」在用户认证场景下的业务胶水层：
// 与渠道无关的生成/风控/校验逻辑都在 internal/verifycode，
// 这里只负责三件业务事：场景合法性、收件人与账号的对应关系、以及登录态签发。
//
// 想接入短信（手机号验证码）时，只需在 api 层加一个 phone/code 接口、
// 这里按 channel 分派即可，风控与校验完全复用。

func isValidVerifyScene(scene string) bool {
	switch scene {
	case v1.VerifySceneLogin, v1.VerifySceneRegister, v1.VerifySceneReset:
		return true
	}
	return false
}

// findUserByTarget 按渠道把收件人映射到账号。
// 接短信时在这里加 ChannelSMS → users.phone 的分支即可。
func findUserByTarget(ctx context.Context, ch verifycode.Channel, target string) (*entity.Users, error) {
	var (
		user *entity.Users
		err  error
	)
	switch ch {
	case verifycode.ChannelSMS:
		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Phone, target).Scan(&user)
	default:
		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, target).Scan(&user)
	}
	return user, err
}

// sendUserAuthCode 是发送验证码的渠道无关实现。
//
// 防枚举的关键：收件人不存在时不是「直接返回」，而是把 deliver=false 交给
// verifycode 服务 —— 服务仍会照常扣额度、写冷却、落验证码，两条路径不可区分。
func sendUserAuthCode(ctx context.Context, ch verifycode.Channel, scene, target string) (*v1.UserSendEmailCodeRes, error) {
	svc, err := verifycode.Default(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "初始化验证码服务失败：%v", err)
		return nil, errcode.ErrVerifyChannelNotReady
	}

	var (
		deliver bool
		reason  string
	)
	policy := func(ctx context.Context) (bool, error) {
		user, err := findUserByTarget(ctx, ch, target)
		if err != nil {
			return false, err
		}
		switch scene {
		case v1.VerifySceneRegister:
			if user != nil {
				// 注册场景需要明确告知占用情况（注册接口本身也会暴露这一点）
				return false, errcode.ErrEmailAlreadyExists
			}
			return true, nil
		default:
			// 登录 / 重置场景：未注册则静默跳过投递，但风控副作用照常
			if user == nil {
				reason = "收件人未注册"
				return false, nil
			}
			return true, nil
		}
	}

	result, err := svc.Send(ctx, ch, scene, target, policy)
	if err != nil {
		return nil, err
	}
	if !deliver && reason != "" {
		g.Log().Infof(ctx, "验证码按防枚举策略跳过投递（channel=%s scene=%s 原因=%s）", ch, scene, reason)
	}

	return &v1.UserSendEmailCodeRes{
		Success:     true,
		ExpireIn:    result.ExpireIn,
		ResendAfter: result.ResendAfter,
	}, nil
}

// SendEmailCode 发送邮箱验证码（邮箱渠道入口）。
func (s *sUserAuth) SendEmailCode(ctx context.Context, req *v1.UserSendEmailCodeReq) (res *v1.UserSendEmailCodeRes, err error) {
	email := utility.NormalizeEmail(req.Email)
	if email == "" {
		return nil, errcode.ErrInvalidParams
	}
	scene := strings.TrimSpace(req.Scene)
	if scene == "" {
		scene = v1.VerifySceneLogin
	}
	if !isValidVerifyScene(scene) {
		return nil, errcode.ErrInvalidParams
	}
	return sendUserAuthCode(ctx, verifycode.ChannelEmail, scene, email)
}

// EmailLogin 使用邮箱 + 验证码登录。
func (s *sUserAuth) EmailLogin(ctx context.Context, req *v1.UserEmailLoginReq) (res *v1.UserEmailLoginRes, err error) {
	email := utility.NormalizeEmail(req.Email)
	if email == "" {
		return nil, errcode.ErrInvalidParams
	}
	svc, err := verifycode.Default(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "初始化验证码服务失败：%v", err)
		return nil, errcode.ErrVerifyChannelNotReady
	}

	// 先验码、再查账号：未注册邮箱本来就不可能有验证码，因此「邮箱不存在」与
	// 「验证码错误」会返回完全相同的错误码（1017），不会因为错误码差异泄露注册状态。
	if err = svc.Consume(ctx, verifycode.ChannelEmail, v1.VerifySceneLogin, email, strings.TrimSpace(req.Code)); err != nil {
		return nil, err
	}

	user, err := findUserByTarget(ctx, verifycode.ChannelEmail, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 验码已通过却查不到账号，正常不会发生（验证码只可能发给已绑定邮箱），防御性返回同一错误
		return nil, errcode.ErrVerifyCodeInvalid
	}
	// 走到这里说明验证码已校验通过（只有邮箱持有人才能拿到码），
	// 因此后续针对账号状态的错误码不构成枚举风险
	if user.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	pair, err := issueUserSession(ctx, user, "email_code")
	if err != nil {
		return nil, err
	}

	return &v1.UserEmailLoginRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
		UserId:       user.Id,
		Username:     user.Username,
	}, nil
}

// issueUserSession 生成令牌对、更新登录信息、写登录历史与刷新白名单。
// 密码登录与验证码登录共用，保证两种入口的会话行为完全一致。
func issueUserSession(ctx context.Context, user *entity.Users, loginMethod string) (*utility.TokenPair, error) {
	pair, err := utility.GenerateUserTokenPair(ctx, user.Id, user.Username)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	ip := clientIp(ctx)
	device := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		device = r.Header.Get("User-Agent")
		if len(device) > 100 {
			device = device[:100]
		}
	}
	if _, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, user.Id).Update(do.Users{
		LastLoginIp: ip,
		LastLoginAt: gtime.New(time.Now()),
	}); err != nil {
		g.Log().Warningf(ctx, "更新最后登录信息失败 user_id=%d: %v", user.Id, err)
	}
	if _, err = dao.UsrLoginHistories.Ctx(ctx).Insert(do.UsrLoginHistories{
		UserId:      user.Id,
		LoginIp:     ip,
		LoginDevice: device,
		LoginMethod: loginMethod,
		LoginStatus: 1,
	}); err != nil {
		g.Log().Warning(ctx, "insert login history failed: %v", err)
	}

	refreshClaims, _ := utility.ParseUserToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveUserRefreshToken(ctx, user.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}
	return pair, nil
}

func clientIp(ctx context.Context) string {
	if r := g.RequestFromCtx(ctx); r != nil {
		return r.GetClientIp()
	}
	return ""
}
