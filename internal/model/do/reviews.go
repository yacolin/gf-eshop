// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Reviews is the golang structure of table rev_reviews for DAO operations like Where/Data.
type Reviews struct {
	g.Meta          `orm:"table:rev_reviews, do:true"`
	Id              interface{} // 评价ID
	UserId          interface{} // 用户ID（冗余，便于查询）
	OrderId         interface{} // 订单ID（校验必须已购）
	OrderItemId     interface{} // 订单明细ID（用于区分同订单多商品）
	SpuId           interface{} // 商品SPU ID
	SkuId           interface{} // 商品SKU ID（若评价具体规格则填）
	MerchantId      interface{} // 所属商家ID
	OverallRating   interface{} // 总体评分（1-5星）
	QualityRating   interface{} // 质量评分（1-5）
	LogisticsRating interface{} // 物流评分（1-5）
	ServiceRating   interface{} // 服务评分（1-5）
	Content         interface{} // 评价文字内容
	IsAnonymous     interface{} // 是否匿名 0-否 1-是
	Status          interface{} // 0-待审核 1-审核通过 2-审核拒绝 3-用户删除
	RejectReason    interface{} // 拒绝原因（审核不通过时填写）
	LatestReplyId   interface{} // 最新回复ID
	ReplyCount      interface{} // 回复总数
	LikeCount       interface{} // 点赞数
	HelpfulCount    interface{} // 有用数
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time // 软删除
}
