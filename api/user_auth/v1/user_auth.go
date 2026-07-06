package v1

import "github.com/gogf/gf/v2/frame/g"

type UserLoginReq struct {
	g.Meta   `path:"/login" tags:"用户认证" method:"post" summary:"用户密码登录"`
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
	g.Meta   `path:"/register" tags:"用户认证" method:"post" summary:"用户注册"`
	Username string `json:"username" v:"required" description:"用户名"`
	Password string `json:"password" v:"required" description:"密码"`
	Email    string `json:"email"    description:"邮箱"`
	Phone    string `json:"phone"    description:"手机号"`
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
	g.Meta       `path:"/refresh" tags:"用户认证" method:"post" summary:"刷新令牌"`
	RefreshToken string `json:"refresh_token" v:"required" description:"刷新令牌"`
}
type UserRefreshTokenRes struct {
	AccessToken  string `json:"access_token"  description:"新的访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"新的刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
}
