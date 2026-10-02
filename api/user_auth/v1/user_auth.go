package v1

import "github.com/gogf/gf/v2/frame/g"

type UserLoginReq struct {
	g.Meta   `path:"/user/auth/login" tags:"UserAuth" method:"post" summary:"用户密码登录"`
	Username string `json:"username" v:"required" description:"用户名"`
	Password string `json:"password" v:"required" description:"密码"`
}
type UserLoginRes struct {
	AccessToken  string `json:"access_token"  description:"访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
	UserId       int64  `json:"user_id"       description:"用户ID"`
	Username     string `json:"username"      description:"用户名"`
}

type UserRegisterReq struct {
	g.Meta   `path:"/user/auth/register" tags:"UserAuth" method:"post" summary:"用户注册"`
	Username string `json:"username" v:"required" description:"用户名"`
	Password string `json:"password" v:"required" description:"密码"`
	Email    string `json:"email"    description:"邮箱"`
	Phone    string `json:"phone"    description:"手机号"`
	// EmailCode 可选：填写后会校验注册场景的邮箱验证码，并把 email_verified 置为 1；
	// 不填写则保持旧行为（邮箱未验证）。
	EmailCode string `json:"email_code" description:"邮箱验证码（scene=register）"`
}
type UserRegisterRes struct {
	AccessToken  string `json:"access_token"  description:"访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
	UserId       int64  `json:"user_id"       description:"用户ID"`
	Username     string `json:"username"      description:"用户名"`
}

type UserRefreshTokenReq struct {
	g.Meta       `path:"/user/auth/refresh" tags:"UserAuth" method:"post" summary:"刷新令牌"`
	RefreshToken string `json:"refresh_token" v:"required" description:"刷新令牌"`
}
type UserRefreshTokenRes struct {
	AccessToken  string `json:"access_token"  description:"新的访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"新的刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
}

// 验证码用途（与投递渠道无关）：login=验证码登录，register=注册前验证收件人归属，
// reset=忘记密码重置。不同用途的验证码互相隔离，不可跨场景使用。
const (
	VerifySceneLogin    = "login"
	VerifySceneRegister = "register"
	VerifySceneReset    = "reset"
)

type UserSendEmailCodeReq struct {
	g.Meta `path:"/user/auth/email/code" tags:"UserAuth" method:"post" summary:"发送邮箱验证码"`
	Email  string `json:"email" v:"required|email" description:"接收验证码的邮箱"`
	// Scene 验证码用途，默认 login。不同用途的验证码互相隔离，不可跨场景使用。
	Scene string `json:"scene" description:"验证码用途：login-登录 register-注册 reset-重置密码，默认 login"`
}
type UserSendEmailCodeRes struct {
	Success bool `json:"success" description:"是否已受理（为防邮箱枚举，未注册邮箱同样返回 true）"`
	// ExpireIn 秒；ResendAfter 秒内不允许重复发送
	ExpireIn    int `json:"expire_in"    description:"验证码有效期（秒）"`
	ResendAfter int `json:"resend_after" description:"多少秒后可重新发送"`
}

type UserEmailLoginReq struct {
	g.Meta `path:"/user/auth/email/login" tags:"UserAuth" method:"post" summary:"邮箱验证码登录"`
	Email  string `json:"email" v:"required|email" description:"邮箱"`
	Code   string `json:"code"  v:"required" description:"邮箱验证码"`
}
type UserEmailLoginRes struct {
	AccessToken  string `json:"access_token"  description:"访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
	UserId       int64  `json:"user_id"       description:"用户ID"`
	Username     string `json:"username"      description:"用户名"`
}

type UserResetPasswordReq struct {
	g.Meta `path:"/user/auth/password/reset" tags:"UserAuth" method:"post" summary:"邮箱验证码重置密码"`
	Email  string `json:"email" v:"required|email" description:"注册邮箱"`
	// Code 必须是 scene=reset 的验证码；login / register 场景的码不可用于重置。
	Code        string `json:"code"         v:"required" description:"邮箱验证码（scene=reset）"`
	NewPassword string `json:"new_password" v:"required|length:8,64" description:"新密码（8-64 位）"`
}
type UserResetPasswordRes struct {
	Success bool `json:"success" description:"是否重置成功"`
	// RevokedSessions 重置后已吊销的刷新令牌数量（所有设备需重新登录）
	RevokedSessions int `json:"revoked_sessions" description:"已吊销的刷新令牌数量"`
}
