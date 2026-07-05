// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Staff is the golang structure for table staff.
type Staff struct {
	Id           int64       `json:"id"            description:"主键"`
	Username     string      `json:"username"      description:"登录用户名（唯一）"`
	PasswordHash string      `json:"password_hash" description:"bcrypt 密码哈希"`
	RealName     string      `json:"real_name"     description:"真实姓名"`
	Email        string      `json:"email"         description:"邮箱"`
	Phone        string      `json:"phone"         description:"手机号"`
	Avatar       string      `json:"avatar"        description:"头像URL"`
	Status       int         `json:"status"        description:"1-正常 0-禁用"`
	LastLoginIp  string      `json:"last_login_ip" description:"最后登录IP"`
	LastLoginAt  *gtime.Time `json:"last_login_at" description:"最后登录时间"`
	CreatedBy    int64       `json:"created_by"    description:"创建人ID（0=系统）"`
	CreatedAt    *gtime.Time `json:"created_at"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updated_at"    description:"更新时间"`
	DeletedAt    *gtime.Time `json:"deleted_at"    description:"删除时间"`
}
