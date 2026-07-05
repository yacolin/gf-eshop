// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id              int64       `json:"id"               description:"通知ID"`
	UserId          int64       `json:"user_id"          description:"0-全体用户 >0-指定用户"`
	MerchantId      int64       `json:"merchant_id"      description:"所属商家ID（0表示平台）"`
	Title           string      `json:"title"            description:"通知标题"`
	Content         string      `json:"content"          description:"通知内容（最终渲染后的文本）"`
	ContentTemplate string      `json:"content_template" description:"内容模板ID（用于统计/调试）"`
	TemplateParams  string      `json:"template_params"  description:"模板参数（存原始变量，便于回溯）"`
	Channel         int         `json:"channel"          description:"渠道: 1-站内消息 2-App Push 3-短信 4-邮件 5-微信模板消息"`
	Category        int         `json:"category"         description:"分类: 1-系统公告 2-订单通知 3-营销推广 4-互动通知 5-安全提醒"`
	TargetType      string      `json:"target_type"      description:"关联业务类型: order/coupon/review/activity 等"`
	TargetId        int64       `json:"target_id"        description:"关联业务ID"`
	RedirectUrl     string      `json:"redirect_url"     description:"跳转链接（优先级高于 target_type+id）"`
	IconUrl         string      `json:"icon_url"         description:"通知图标（如订单图标、优惠券图标）"`
	Priority        int         `json:"priority"         description:"优先级: 0-高 1-中 2-低（用于排序展示）"`
	CreatedBy       int64       `json:"created_by"       description:"创建人（0表示系统自动）"`
	CreatedAt       *gtime.Time `json:"created_at"       description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:"更新时间"`
	DeletedAt       *gtime.Time `json:"deleted_at"       description:"系统软删除（数据清理用）"`
}
