// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationReads is the golang structure for table notification_reads.
type NotificationReads struct {
	Id             int64       `json:"id"              description:"主键"`
	NotificationId int64       `json:"notification_id" description:"关联通知ID"`
	UserId         int64       `json:"user_id"         description:"用户ID"`
	ReadAt         *gtime.Time `json:"read_at"         description:"阅读时间"`
}
