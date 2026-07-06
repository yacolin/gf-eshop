// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table usr_users for DAO operations like Where/Data.
type Users struct {
	g.Meta         `orm:"table:usr_users, do:true"`
	Id             interface{} // 主键
	Username       interface{} // 用户名（唯一，NULL表示未设置）
	PasswordHash   interface{} // bcrypt 密码哈希
	Email          interface{} // 邮箱（唯一，NULL表示未绑定）
	EmailVerified  interface{} // 邮箱是否已验证
	Phone          interface{} // 手机号（唯一，NULL表示未绑定）
	PhoneVerified  interface{} // 手机号是否已验证
	Avatar         interface{} // 头像URL
	Nickname       interface{} // 昵称
	Status         interface{} // 状态：1-正常 0-禁用 2-冻结
	RegisterIp     interface{} // 注册IP
	RegisterSource interface{} // 注册来源：web/ios/android/admin
	LastLoginIp    interface{} // 最后登录IP
	LastLoginAt    *gtime.Time // 最后登录时间
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 删除时间
}
