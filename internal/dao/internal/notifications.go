// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationsDao is the data access object for table base_notifications.
type NotificationsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns NotificationsColumns // columns contains all the column names of Table for convenient usage.
}

// NotificationsColumns defines and stores column names for table base_notifications.
type NotificationsColumns struct {
	Id              string // 通知ID
	UserId          string // 0-全体用户 >0-指定用户
	MerchantId      string // 所属商家ID（0表示平台）
	Title           string // 通知标题
	Content         string // 通知内容（最终渲染后的文本）
	ContentTemplate string // 内容模板ID（用于统计/调试）
	TemplateParams  string // 模板参数（存原始变量，便于回溯）
	Channel         string // 渠道: 1-站内消息 2-App Push 3-短信 4-邮件 5-微信模板消息
	Category        string // 分类: 1-系统公告 2-订单通知 3-营销推广 4-互动通知 5-安全提醒
	TargetType      string // 关联业务类型: order/coupon/review/activity 等
	TargetId        string // 关联业务ID
	RedirectUrl     string // 跳转链接（优先级高于 target_type+id）
	IconUrl         string // 通知图标（如订单图标、优惠券图标）
	Priority        string // 优先级: 0-高 1-中 2-低（用于排序展示）
	CreatedBy       string // 创建人（0表示系统自动）
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	DeletedAt       string // 系统软删除（数据清理用）
}

// notificationsColumns holds the columns for table base_notifications.
var notificationsColumns = NotificationsColumns{
	Id:              "id",
	UserId:          "user_id",
	MerchantId:      "merchant_id",
	Title:           "title",
	Content:         "content",
	ContentTemplate: "content_template",
	TemplateParams:  "template_params",
	Channel:         "channel",
	Category:        "category",
	TargetType:      "target_type",
	TargetId:        "target_id",
	RedirectUrl:     "redirect_url",
	IconUrl:         "icon_url",
	Priority:        "priority",
	CreatedBy:       "created_by",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewNotificationsDao creates and returns a new DAO object for table data access.
func NewNotificationsDao() *NotificationsDao {
	return &NotificationsDao{
		group:   "default",
		table:   "base_notifications",
		columns: notificationsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *NotificationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *NotificationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *NotificationsDao) Columns() NotificationsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *NotificationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *NotificationsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *NotificationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
