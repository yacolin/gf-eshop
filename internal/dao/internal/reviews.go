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
	ReviewNo        string // 评价业务单号（幂等键）
	UserId          string // 用户ID
	OrderId         string // 订单ID（校验必须已购）
	OrderItemId     string // 订单明细ID（用于区分同订单多商品）
	SpuId           string // 商品SPU ID
	SkuId           string // 商品SKU ID
	MerchantId      string // 所属商家ID
	OverallRating   string // 总体评分（1-5星）
	QualityRating   string // 质量评分
	LogisticsRating string // 物流评分
	ServiceRating   string // 服务评分
	Content         string // 评价文字内容
	ContentLength   string // 内容长度（冗余，用于筛选优质评价）
	IsAnonymous     string // 是否匿名 0-否 1-是
	HasMedia        string // 是否包含媒体 0-否 1-是
	Status          string // 0-待审核 1-审核通过 2-审核拒绝 3-用户删除 4-平台屏蔽
	RiskLevel       string // 风险等级 0-正常 1-低风险 2-高风险
	RejectReason    string // 拒绝原因
	AuditedBy       string // 审核人ID
	AuditedAt       string // 审核时间
	LikeCount       string // 点赞数（异步校准）
	HelpfulCount    string // 有用数（异步校准）
	ReplyCount      string // 回复总数（异步校准）
	CreatedAt       string //
	UpdatedAt       string //
}

// reviewsColumns holds the columns for table rev_reviews.
var reviewsColumns = ReviewsColumns{
	Id:              "id",
	ReviewNo:        "review_no",
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
	ContentLength:   "content_length",
	IsAnonymous:     "is_anonymous",
	HasMedia:        "has_media",
	Status:          "status",
	RiskLevel:       "risk_level",
	RejectReason:    "reject_reason",
	AuditedBy:       "audited_by",
	AuditedAt:       "audited_at",
	LikeCount:       "like_count",
	HelpfulCount:    "helpful_count",
	ReplyCount:      "reply_count",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
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
