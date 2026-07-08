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
	UserId          int64       `json:"user_id"          description:"用户ID（冗余，便于查询）"`
	OrderId         int64       `json:"order_id"         description:"订单ID（校验必须已购）"`
	OrderItemId     int64       `json:"order_item_id"    description:"订单明细ID（用于区分同订单多商品）"`
	SpuId           int64       `json:"spu_id"           description:"商品SPU ID"`
	SkuId           int64       `json:"sku_id"           description:"商品SKU ID（若评价具体规格则填）"`
	MerchantId      int64       `json:"merchant_id"      description:"所属商家ID"`
	OverallRating   int         `json:"overall_rating"   description:"总体评分（1-5星）"`
	QualityRating   int         `json:"quality_rating"   description:"质量评分（1-5）"`
	LogisticsRating int         `json:"logistics_rating" description:"物流评分（1-5）"`
	ServiceRating   int         `json:"service_rating"   description:"服务评分（1-5）"`
	Content         string      `json:"content"          description:"评价文字内容"`
	IsAnonymous     int         `json:"is_anonymous"     description:"是否匿名 0-否 1-是"`
	Status          int         `json:"status"           description:"0-待审核 1-审核通过 2-审核拒绝 3-用户删除"`
	RejectReason    string      `json:"reject_reason"    description:"拒绝原因（审核不通过时填写）"`
	LatestReplyId   int64       `json:"latest_reply_id"  description:"最新回复ID"`
	ReplyCount      int         `json:"reply_count"      description:"回复总数"`
	LikeCount       int         `json:"like_count"       description:"点赞数"`
	HelpfulCount    int         `json:"helpful_count"    description:"有用数"`
	CreatedAt       *gtime.Time `json:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:""`
	DeletedAt       *gtime.Time `json:"deleted_at"       description:"软删除"`
}
