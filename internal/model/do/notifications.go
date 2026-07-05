// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure of table base_notifications for DAO operations like Where/Data.
type Notifications struct {
	g.Meta          `orm:"table:base_notifications, do:true"`
	Id              interface{} // 通知ID
	UserId          interface{} // 0-全体用户 >0-指定用户
	MerchantId      interface{} // 所属商家ID（0表示平台）
	Title           interface{} // 通知标题
	Content         interface{} // 通知内容（最终渲染后的文本）
	ContentTemplate interface{} // 内容模板ID（用于统计/调试）
	TemplateParams  interface{} // 模板参数（存原始变量，便于回溯）
	Channel         interface{} // 渠道: 1-站内消息 2-App Push 3-短信 4-邮件 5-微信模板消息
	Category        interface{} // 分类: 1-系统公告 2-订单通知 3-营销推广 4-互动通知 5-安全提醒
	TargetType      interface{} // 关联业务类型: order/coupon/review/activity 等
	TargetId        interface{} // 关联业务ID
	RedirectUrl     interface{} // 跳转链接（优先级高于 target_type+id）
	IconUrl         interface{} // 通知图标（如订单图标、优惠券图标）
	Priority        interface{} // 优先级: 0-高 1-中 2-低（用于排序展示）
	CreatedBy       interface{} // 创建人（0表示系统自动）
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
	DeletedAt       *gtime.Time // 系统软删除（数据清理用）
}
