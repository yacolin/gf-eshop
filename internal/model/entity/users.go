// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id             int64       `json:"id"              description:"主键"`
	Username       string      `json:"username"        description:"用户名（唯一，NULL表示未设置）"`
	PasswordHash   string      `json:"password_hash"   description:"bcrypt 密码哈希"`
	Email          string      `json:"email"           description:"邮箱（唯一，NULL表示未绑定）"`
	EmailVerified  int         `json:"email_verified"  description:"邮箱是否已验证"`
	Phone          string      `json:"phone"           description:"手机号（唯一，NULL表示未绑定）"`
	PhoneVerified  int         `json:"phone_verified"  description:"手机号是否已验证"`
	Avatar         string      `json:"avatar"          description:"头像URL"`
	Nickname       string      `json:"nickname"        description:"昵称"`
	Status         int         `json:"status"          description:"状态：1-正常 0-禁用 2-冻结"`
	RegisterIp     string      `json:"register_ip"     description:"注册IP"`
	RegisterSource string      `json:"register_source" description:"注册来源：web/ios/android/admin"`
	LastLoginIp    string      `json:"last_login_ip"   description:"最后登录IP"`
	LastLoginAt    *gtime.Time `json:"last_login_at"   description:"最后登录时间"`
	CreatedAt      *gtime.Time `json:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:"更新时间"`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:"删除时间"`
}
