// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Reviews is the golang structure for table reviews.
type Reviews struct {
	Id              int64       `json:"id"               description:"评价ID"`
	ReviewNo        string      `json:"review_no"        description:"评价业务单号（幂等键）"`
	UserId          int64       `json:"user_id"          description:"用户ID"`
	OrderId         int64       `json:"order_id"         description:"订单ID（校验必须已购）"`
	OrderItemId     int64       `json:"order_item_id"    description:"订单明细ID（用于区分同订单多商品）"`
	SpuId           int64       `json:"spu_id"           description:"商品SPU ID"`
	SkuId           int64       `json:"sku_id"           description:"商品SKU ID"`
	MerchantId      int64       `json:"merchant_id"      description:"所属商家ID"`
	OverallRating   int         `json:"overall_rating"   description:"总体评分（1-5星）"`
	QualityRating   int         `json:"quality_rating"   description:"质量评分"`
	LogisticsRating int         `json:"logistics_rating" description:"物流评分"`
	ServiceRating   int         `json:"service_rating"   description:"服务评分"`
	Content         string      `json:"content"          description:"评价文字内容"`
	ContentLength   int         `json:"content_length"   description:"内容长度（冗余，用于筛选优质评价）"`
	IsAnonymous     int         `json:"is_anonymous"     description:"是否匿名 0-否 1-是"`
	HasMedia        int         `json:"has_media"        description:"是否包含媒体 0-否 1-是"`
	Status          int         `json:"status"           description:"0-待审核 1-审核通过 2-审核拒绝 3-用户删除 4-平台屏蔽"`
	RiskLevel       int         `json:"risk_level"       description:"风险等级 0-正常 1-低风险 2-高风险"`
	RejectReason    string      `json:"reject_reason"    description:"拒绝原因"`
	AuditedBy       int64       `json:"audited_by"       description:"审核人ID"`
	AuditedAt       *gtime.Time `json:"audited_at"       description:"审核时间"`
	LikeCount       int         `json:"like_count"       description:"点赞数（异步校准）"`
	HelpfulCount    int         `json:"helpful_count"    description:"有用数（异步校准）"`
	ReplyCount      int         `json:"reply_count"      description:"回复总数（异步校准）"`
	CreatedAt       *gtime.Time `json:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:""`
}
