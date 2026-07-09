// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewStatisticsDao is the data access object for table rev_review_statistics.
type ReviewStatisticsDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns ReviewStatisticsColumns // columns contains all the column names of Table for convenient usage.
}

// ReviewStatisticsColumns defines and stores column names for table rev_review_statistics.
type ReviewStatisticsColumns struct {
	Id              string //
	TargetType      string // 统计目标类型 1-SPU 2-商家
	TargetId        string // 目标ID（SPU_ID或Merchant_ID）
	Rating1Count    string // 1星数量
	Rating2Count    string // 2星数量
	Rating3Count    string // 3星数量
	Rating4Count    string // 4星数量
	Rating5Count    string // 5星数量
	TotalCount      string // 总评价数
	AvgRating       string // 平均评分
	GoodRate        string // 好评率（%）
	HasMediaCount   string // 带图评价数
	HasContentCount string // 有内容评价数
	LastUpdatedAt   string //
}

// reviewStatisticsColumns holds the columns for table rev_review_statistics.
var reviewStatisticsColumns = ReviewStatisticsColumns{
	Id:              "id",
	TargetType:      "target_type",
	TargetId:        "target_id",
	Rating1Count:    "rating_1_count",
	Rating2Count:    "rating_2_count",
	Rating3Count:    "rating_3_count",
	Rating4Count:    "rating_4_count",
	Rating5Count:    "rating_5_count",
	TotalCount:      "total_count",
	AvgRating:       "avg_rating",
	GoodRate:        "good_rate",
	HasMediaCount:   "has_media_count",
	HasContentCount: "has_content_count",
	LastUpdatedAt:   "last_updated_at",
}

// NewReviewStatisticsDao creates and returns a new DAO object for table data access.
func NewReviewStatisticsDao() *ReviewStatisticsDao {
	return &ReviewStatisticsDao{
		group:   "default",
		table:   "rev_review_statistics",
		columns: reviewStatisticsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ReviewStatisticsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ReviewStatisticsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ReviewStatisticsDao) Columns() ReviewStatisticsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ReviewStatisticsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ReviewStatisticsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ReviewStatisticsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
