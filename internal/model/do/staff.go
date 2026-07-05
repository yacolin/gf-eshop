// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Staff is the golang structure of table sys_staff for DAO operations like Where/Data.
type Staff struct {
	g.Meta       `orm:"table:sys_staff, do:true"`
	Id           interface{} // 主键
	Username     interface{} // 登录用户名（唯一）
	PasswordHash interface{} // bcrypt 密码哈希
	RealName     interface{} // 真实姓名
	Email        interface{} // 邮箱
	Phone        interface{} // 手机号
	Avatar       interface{} // 头像URL
	Status       interface{} // 1-正常 0-禁用
	LastLoginIp  interface{} // 最后登录IP
	LastLoginAt  *gtime.Time // 最后登录时间
	CreatedBy    interface{} // 创建人ID（0=系统）
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 删除时间
}
