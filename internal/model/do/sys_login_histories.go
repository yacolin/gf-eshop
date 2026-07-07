// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysLoginHistories is the golang structure of table sys_login_histories for DAO operations like Where/Data.
type SysLoginHistories struct {
	g.Meta        `orm:"table:sys_login_histories, do:true"`
	Id            interface{} // 主键
	StaffId       interface{} // 员工ID（关联 sys_staff.id）
	LoginIp       interface{} // 登录IP
	LoginDevice   interface{} // 登录设备信息（UA）
	LoginLocation interface{} // 登录地点
	LoginMethod   interface{} // 登录方式：password/sms/oauth
	LoginStatus   interface{} // 1-成功 0-失败
	FailureReason interface{} // 失败原因
	CreatedAt     *gtime.Time // 创建时间
}
