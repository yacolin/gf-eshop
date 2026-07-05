// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTemplates is the golang structure of table base_notification_templates for DAO operations like Where/Data.
type NotificationTemplates struct {
	g.Meta          `orm:"table:base_notification_templates, do:true"`
	Id              interface{} //
	TemplateCode    interface{} // 模板代码（如 ORDER_PAID_SUCCESS）
	Channel         interface{} // 渠道 1-站内 2-Push 3-短信 4-邮件
	TitleTemplate   interface{} // 标题模板（支持变量 {{.OrderID}}）
	ContentTemplate interface{} // 内容模板
	Category        interface{} // 默认分类
	Priority        interface{} // 默认优先级
	Status          interface{} // 1-启用 0-停用
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
