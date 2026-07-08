// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewsDao is the data access object for table rev_reviews.
type ReviewsDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns ReviewsColumns // columns contains all the column names of Table for convenient usage.
}

// ReviewsColumns defines and stores column names for table rev_reviews.
type ReviewsColumns struct {
	Id              string // 评价ID
	UserId          string // 用户ID（冗余，便于查询）
	OrderId         string // 订单ID（校验必须已购）
	OrderItemId     string // 订单明细ID（用于区分同订单多商品）
	SpuId           string // 商品SPU ID
	SkuId           string // 商品SKU ID（若评价具体规格则填）
	MerchantId      string // 所属商家ID
	OverallRating   string // 总体评分（1-5星）
	QualityRating   string // 质量评分（1-5）
	LogisticsRating string // 物流评分（1-5）
	ServiceRating   string // 服务评分（1-5）
	Content         string // 评价文字内容
	IsAnonymous     string // 是否匿名 0-否 1-是
	Status          string // 0-待审核 1-审核通过 2-审核拒绝 3-用户删除
	RejectReason    string // 拒绝原因（审核不通过时填写）
	LatestReplyId   string // 最新回复ID
	ReplyCount      string // 回复总数
	LikeCount       string // 点赞数
	HelpfulCount    string // 有用数
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string // 软删除
}

// reviewsColumns holds the columns for table rev_reviews.
var reviewsColumns = ReviewsColumns{
	Id:              "id",
	UserId:          "user_id",
	OrderId:         "order_id",
	OrderItemId:     "order_item_id",
	SpuId:           "spu_id",
	SkuId:           "sku_id",
	MerchantId:      "merchant_id",
	OverallRating:   "overall_rating",
	QualityRating:   "quality_rating",
	LogisticsRating: "logistics_rating",
	ServiceRating:   "service_rating",
	Content:         "content",
	IsAnonymous:     "is_anonymous",
	Status:          "status",
	RejectReason:    "reject_reason",
	LatestReplyId:   "latest_reply_id",
	ReplyCount:      "reply_count",
	LikeCount:       "like_count",
	HelpfulCount:    "helpful_count",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewReviewsDao creates and returns a new DAO object for table data access.
func NewReviewsDao() *ReviewsDao {
	return &ReviewsDao{
		group:   "default",
		table:   "rev_reviews",
		columns: reviewsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ReviewsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ReviewsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ReviewsDao) Columns() ReviewsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ReviewsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ReviewsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ReviewsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
