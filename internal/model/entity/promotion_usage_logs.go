// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionUsageLogs is the golang structure for table promotion_usage_logs.
type PromotionUsageLogs struct {
	Id                int64       `json:"id"                 description:"主键ID"`
	PromotionId       int64       `json:"promotion_id"       description:"促销ID"`
	UserPromotionId   int64       `json:"user_promotion_id"  description:"用户促销资产ID"`
	UserId            int64       `json:"user_id"            description:"用户ID"`
	OrderId           int64       `json:"order_id"           description:"使用的订单ID"`
	DiscountAmount    int64       `json:"discount_amount"    description:"优惠金额（分）"`
	PromotionSnapshot string      `json:"promotion_snapshot" description:"优惠快照（名称、规则等）"`
	CreatedAt         *gtime.Time `json:"created_at"         description:""`
}
