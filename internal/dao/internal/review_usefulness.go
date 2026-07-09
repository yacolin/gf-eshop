// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewUsefulnessDao is the data access object for table rev_review_usefulness.
type ReviewUsefulnessDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns ReviewUsefulnessColumns // columns contains all the column names of Table for convenient usage.
}

// ReviewUsefulnessColumns defines and stores column names for table rev_review_usefulness.
type ReviewUsefulnessColumns struct {
	Id        string // 主键
	ReviewId  string // 评价ID
	UserId    string // 用户ID
	CreatedAt string //
}

// reviewUsefulnessColumns holds the columns for table rev_review_usefulness.
var reviewUsefulnessColumns = ReviewUsefulnessColumns{
	Id:        "id",
	ReviewId:  "review_id",
	UserId:    "user_id",
	CreatedAt: "created_at",
}

// NewReviewUsefulnessDao creates and returns a new DAO object for table data access.
func NewReviewUsefulnessDao() *ReviewUsefulnessDao {
	return &ReviewUsefulnessDao{
		group:   "default",
		table:   "rev_review_usefulness",
		columns: reviewUsefulnessColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ReviewUsefulnessDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ReviewUsefulnessDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ReviewUsefulnessDao) Columns() ReviewUsefulnessColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ReviewUsefulnessDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ReviewUsefulnessDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ReviewUsefulnessDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
