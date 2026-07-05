// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTemplates is the golang structure for table notification_templates.
type NotificationTemplates struct {
	Id              int64       `json:"id"               description:""`
	TemplateCode    string      `json:"template_code"    description:"模板代码（如 ORDER_PAID_SUCCESS）"`
	Channel         int         `json:"channel"          description:"渠道 1-站内 2-Push 3-短信 4-邮件"`
	TitleTemplate   string      `json:"title_template"   description:"标题模板（支持变量 {{.OrderID}}）"`
	ContentTemplate string      `json:"content_template" description:"内容模板"`
	Category        int         `json:"category"         description:"默认分类"`
	Priority        int         `json:"priority"         description:"默认优先级"`
	Status          int         `json:"status"           description:"1-启用 0-停用"`
	CreatedAt       *gtime.Time `json:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:""`
}
