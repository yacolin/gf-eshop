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
	ReviewNo        interface{} // 评价业务单号（幂等键）
	UserId          interface{} // 用户ID
	OrderId         interface{} // 订单ID（校验必须已购）
	OrderItemId     interface{} // 订单明细ID（用于区分同订单多商品）
	SpuId           interface{} // 商品SPU ID
	SkuId           interface{} // 商品SKU ID
	MerchantId      interface{} // 所属商家ID
	OverallRating   interface{} // 总体评分（1-5星）
	QualityRating   interface{} // 质量评分
	LogisticsRating interface{} // 物流评分
	ServiceRating   interface{} // 服务评分
	Content         interface{} // 评价文字内容
	ContentLength   interface{} // 内容长度（冗余，用于筛选优质评价）
	IsAnonymous     interface{} // 是否匿名 0-否 1-是
	HasMedia        interface{} // 是否包含媒体 0-否 1-是
	Status          interface{} // 0-待审核 1-审核通过 2-审核拒绝 3-用户删除 4-平台屏蔽
	RiskLevel       interface{} // 风险等级 0-正常 1-低风险 2-高风险
	RejectReason    interface{} // 拒绝原因
	AuditedBy       interface{} // 审核人ID
	AuditedAt       *gtime.Time // 审核时间
	LikeCount       interface{} // 点赞数（异步校准）
	HelpfulCount    interface{} // 有用数（异步校准）
	ReplyCount      interface{} // 回复总数（异步校准）
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
