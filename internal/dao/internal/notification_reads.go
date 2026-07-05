// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationReadsDao is the data access object for table base_notification_reads.
type NotificationReadsDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns NotificationReadsColumns // columns contains all the column names of Table for convenient usage.
}

// NotificationReadsColumns defines and stores column names for table base_notification_reads.
type NotificationReadsColumns struct {
	Id             string // 主键
	NotificationId string // 关联通知ID
	UserId         string // 用户ID
	ReadAt         string // 阅读时间
}

// notificationReadsColumns holds the columns for table base_notification_reads.
var notificationReadsColumns = NotificationReadsColumns{
	Id:             "id",
	NotificationId: "notification_id",
	UserId:         "user_id",
	ReadAt:         "read_at",
}

// NewNotificationReadsDao creates and returns a new DAO object for table data access.
func NewNotificationReadsDao() *NotificationReadsDao {
	return &NotificationReadsDao{
		group:   "default",
		table:   "base_notification_reads",
		columns: notificationReadsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *NotificationReadsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *NotificationReadsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *NotificationReadsDao) Columns() NotificationReadsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *NotificationReadsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *NotificationReadsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *NotificationReadsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
