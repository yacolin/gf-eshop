package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type NotificationListReq struct {
	g.Meta   `path:"/notification" tags:"Notification" method:"get" summary:"通知列表"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
type NotificationListItem struct {
	Id          int64       `json:"id"`
	UserId      int64       `json:"user_id"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Channel     int         `json:"channel"`
	Category    int         `json:"category"`
	TargetType  string      `json:"target_type"`
	TargetId    int64       `json:"target_id"`
	RedirectUrl string      `json:"redirect_url"`
	IconUrl     string      `json:"icon_url"`
	Priority    int         `json:"priority"`
	IsRead      bool        `json:"is_read"`
	CreatedBy   int64       `json:"created_by"`
	CreatedAt   *gtime.Time `json:"created_at"`
}
type NotificationListRes struct {
	List  []*NotificationListItem `json:"list"`
	Total int                     `json:"total"`
}

type NotificationUnreadCountReq struct {
	g.Meta `path:"/notification/unread" tags:"Notification" method:"get" summary:"未读通知数"`
}
type NotificationUnreadCountRes struct {
	Count int64 `json:"count"`
}

type NotificationMarkAsReadReq struct {
	g.Meta `path:"/notification/{id}/read" tags:"Notification" method:"put" summary:"标记已读"`
	Id     int64 `json:"id"`
}
type NotificationMarkAsReadRes struct{}

type NotificationMarkAllAsReadReq struct {
	g.Meta `path:"/notification/readall" tags:"Notification" method:"put" summary:"全部已读"`
}
type NotificationMarkAllAsReadRes struct{}

type NotificationDeleteReq struct {
	g.Meta `path:"/notification/{id}" tags:"Notification" method:"delete" summary:"删除通知"`
	Id     int64 `json:"id"`
}
type NotificationDeleteRes struct{}

type NotificationSendSystemReq struct {
	g.Meta       `path:"/notification/system" tags:"Notification" method:"post" summary:"发送系统通知（admin）"`
	UserId       int64  `json:"user_id" v:"required|min:0" description:"0=全体用户"`
	TemplateCode string `json:"template_code" description:"模板代码（优先级高于 title/content）"`
	Title        string `json:"title" description:"通知标题"`
	Content      string `json:"content" description:"通知内容"`
}
type NotificationSendSystemRes struct {
	Id int64 `json:"id"`
}

type NotificationListTemplatesReq struct {
	g.Meta `path:"/notification/templates" tags:"Notification" method:"get" summary:"通知模板列表"`
}
type NotificationListTemplatesItem struct {
	Id              int64  `json:"id"`
	TemplateCode    string `json:"template_code"`
	Channel         int    `json:"channel"`
	TitleTemplate   string `json:"title_template"`
	ContentTemplate string `json:"content_template"`
	Category        int    `json:"category"`
	Priority        int    `json:"priority"`
	Status          int    `json:"status"`
}
type NotificationListTemplatesRes struct {
	List []*NotificationListTemplatesItem `json:"list"`
}
