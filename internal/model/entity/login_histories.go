// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LoginHistories is the golang structure for table login_histories.
type LoginHistories struct {
	Id            int64       `json:"id"             description:"主键"`
	UserId        int64       `json:"user_id"        description:"用户ID"`
	LoginIp       string      `json:"login_ip"       description:"登录IP"`
	LoginDevice   string      `json:"login_device"   description:"登录设备信息（UA）"`
	LoginLocation string      `json:"login_location" description:"登录地点"`
	LoginMethod   string      `json:"login_method"   description:"登录方式：password/sms/oauth"`
	LoginStatus   int         `json:"login_status"   description:"登录结果：1-成功 0-失败"`
	FailureReason string      `json:"failure_reason" description:"失败原因"`
	CreatedAt     *gtime.Time `json:"created_at"     description:"创建时间"`
}
