package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type UserListReq struct {
	g.Meta   `path:"/users" tags:"User" method:"get" summary:"用户列表"`
	Page     int    `json:"page"      description:"页码"`
	PageSize int    `json:"page_size" description:"每页条数"`
	Keyword  string `json:"keyword"   description:"搜索关键词（用户名/昵称/邮箱/手机号）"`
	Status   *int   `json:"status"    description:"状态：1-正常 0-禁用 2-冻结"`
}
type UserListItem struct {
	Id             int64       `json:"id"              description:"主键"`
	Username       string      `json:"username"        description:"用户名"`
	Nickname       string      `json:"nickname"        description:"昵称"`
	Email          string      `json:"email"           description:"邮箱"`
	EmailVerified  int         `json:"email_verified"  description:"邮箱是否已验证"`
	Phone          string      `json:"phone"           description:"手机号"`
	PhoneVerified  int         `json:"phone_verified"  description:"手机号是否已验证"`
	Avatar         string      `json:"avatar"          description:"头像URL"`
	Status         int         `json:"status"          description:"状态"`
	RegisterIp     string      `json:"register_ip"     description:"注册IP"`
	RegisterSource string      `json:"register_source" description:"注册来源"`
	LastLoginIp    string      `json:"last_login_ip"   description:"最后登录IP"`
	LastLoginAt    *gtime.Time `json:"last_login_at"   description:"最后登录时间"`
	CreatedAt      *gtime.Time `json:"created_at"      description:"创建时间"`
}
type UserListRes struct {
	List  []*UserListItem `json:"list"`
	Total int             `json:"total"`
}
