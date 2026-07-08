// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewRepliesDao is the data access object for table rev_review_replies.
type ReviewRepliesDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns ReviewRepliesColumns // columns contains all the column names of Table for convenient usage.
}

// ReviewRepliesColumns defines and stores column names for table rev_review_replies.
type ReviewRepliesColumns struct {
	Id         string // 回复ID
	ReviewId   string // 关联评价ID
	ParentId   string // 父级回复ID（支持多级回复）
	ReplyType  string // 1-商家回复 2-用户追问 3-平台回复
	Content    string // 回复内容
	OperatorId string // 操作人ID
	CreatedAt  string //
	DeletedAt  string //
}

// reviewRepliesColumns holds the columns for table rev_review_replies.
var reviewRepliesColumns = ReviewRepliesColumns{
	Id:         "id",
	ReviewId:   "review_id",
	ParentId:   "parent_id",
	ReplyType:  "reply_type",
	Content:    "content",
	OperatorId: "operator_id",
	CreatedAt:  "created_at",
	DeletedAt:  "deleted_at",
}

// NewReviewRepliesDao creates and returns a new DAO object for table data access.
func NewReviewRepliesDao() *ReviewRepliesDao {
	return &ReviewRepliesDao{
		group:   "default",
		table:   "rev_review_replies",
		columns: reviewRepliesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ReviewRepliesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ReviewRepliesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ReviewRepliesDao) Columns() ReviewRepliesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ReviewRepliesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ReviewRepliesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ReviewRepliesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
