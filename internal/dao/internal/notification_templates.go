// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationTemplatesDao is the data access object for table base_notification_templates.
type NotificationTemplatesDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns NotificationTemplatesColumns // columns contains all the column names of Table for convenient usage.
}

// NotificationTemplatesColumns defines and stores column names for table base_notification_templates.
type NotificationTemplatesColumns struct {
	Id              string //
	TemplateCode    string // 模板代码（如 ORDER_PAID_SUCCESS）
	Channel         string // 渠道 1-站内 2-Push 3-短信 4-邮件
	TitleTemplate   string // 标题模板（支持变量 {{.OrderID}}）
	ContentTemplate string // 内容模板
	Category        string // 默认分类
	Priority        string // 默认优先级
	Status          string // 1-启用 0-停用
	CreatedAt       string //
	UpdatedAt       string //
}

// notificationTemplatesColumns holds the columns for table base_notification_templates.
var notificationTemplatesColumns = NotificationTemplatesColumns{
	Id:              "id",
	TemplateCode:    "template_code",
	Channel:         "channel",
	TitleTemplate:   "title_template",
	ContentTemplate: "content_template",
	Category:        "category",
	Priority:        "priority",
	Status:          "status",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewNotificationTemplatesDao creates and returns a new DAO object for table data access.
func NewNotificationTemplatesDao() *NotificationTemplatesDao {
	return &NotificationTemplatesDao{
		group:   "default",
		table:   "base_notification_templates",
		columns: notificationTemplatesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *NotificationTemplatesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *NotificationTemplatesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *NotificationTemplatesDao) Columns() NotificationTemplatesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *NotificationTemplatesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *NotificationTemplatesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *NotificationTemplatesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
