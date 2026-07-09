// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewMediaDao is the data access object for table rev_review_media.
type ReviewMediaDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns ReviewMediaColumns // columns contains all the column names of Table for convenient usage.
}

// ReviewMediaColumns defines and stores column names for table rev_review_media.
type ReviewMediaColumns struct {
	Id        string // 媒体ID
	ReviewId  string // 关联评价ID
	MediaType string // 1-图片 2-视频
	MediaUrl  string // 媒体文件URL
	FileSize  string // 文件大小（字节）
	Width     string // 宽度（图片/视频）
	Height    string // 高度（图片/视频）
	Duration  string // 时长（视频，秒）
	SortOrder string // 排序
	CreatedAt string //
}

// reviewMediaColumns holds the columns for table rev_review_media.
var reviewMediaColumns = ReviewMediaColumns{
	Id:        "id",
	ReviewId:  "review_id",
	MediaType: "media_type",
	MediaUrl:  "media_url",
	FileSize:  "file_size",
	Width:     "width",
	Height:    "height",
	Duration:  "duration",
	SortOrder: "sort_order",
	CreatedAt: "created_at",
}

// NewReviewMediaDao creates and returns a new DAO object for table data access.
func NewReviewMediaDao() *ReviewMediaDao {
	return &ReviewMediaDao{
		group:   "default",
		table:   "rev_review_media",
		columns: reviewMediaColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ReviewMediaDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ReviewMediaDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ReviewMediaDao) Columns() ReviewMediaColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ReviewMediaDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ReviewMediaDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ReviewMediaDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
