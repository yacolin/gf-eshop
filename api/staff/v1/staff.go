package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type StaffLoginReq struct {
	g.Meta   `path:"/login" tags:"Staff" method:"post" summary:"系统用户登录"`
	Username string `json:"username" v:"required" description:"用户名"`
	Password string `json:"password" v:"required" description:"密码"`
}
type StaffLoginRes struct {
	AccessToken  string `json:"access_token"  description:"访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
	StaffId      int64  `json:"staff_id"      description:"用户ID"`
	Username     string `json:"username"      description:"用户名"`
	RealName     string `json:"real_name"     description:"真实姓名"`
}

type StaffRefreshTokenReq struct {
	g.Meta       `path:"/refresh" tags:"Staff" method:"post" summary:"刷新令牌"`
	RefreshToken string `json:"refresh_token" v:"required" description:"刷新令牌"`
}
type StaffRefreshTokenRes struct {
	AccessToken  string `json:"access_token"  description:"新的访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"新的刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
}

type StaffLogoutReq struct {
	g.Meta `path:"/logout" tags:"Staff" method:"post" summary:"系统用户退出登录"`
}
type StaffLogoutRes struct{}

type StaffProfileReq struct {
	g.Meta `path:"/profile" tags:"Staff" method:"get" summary:"获取当前用户信息"`
}
type StaffProfileRes struct {
	Id          int64  `json:"id"          description:"用户ID"`
	Username    string `json:"username"    description:"用户名"`
	RealName    string `json:"real_name"   description:"真实姓名"`
	Email       string `json:"email"       description:"邮箱"`
	Phone       string `json:"phone"       description:"手机号"`
	Avatar      string `json:"avatar"      description:"头像URL"`
	Status      int    `json:"status"      description:"状态"`
	LastLoginIp string `json:"last_login_ip" description:"最后登录IP"`
}
