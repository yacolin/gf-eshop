// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionUsageLogs is the golang structure for table promotion_usage_logs.
type PromotionUsageLogs struct {
	Id              int64       `json:"id"                description:"主键ID"`
	PromotionId     int64       `json:"promotion_id"      description:"促销ID"`
	UserPromotionId int64       `json:"user_promotion_id" description:"用户促销资产ID"`
	UserId          int64       `json:"user_id"           description:"用户ID"`
	MerchantId      int64       `json:"merchant_id"       description:"所属商家ID"`
	OrderId         int64       `json:"order_id"          description:"使用的订单ID"`
	UsageType       int         `json:"usage_type"        description:"1-下单使用 2-自动优惠"`
	DiscountAmount  int64       `json:"discount_amount"   description:"优惠金额（分）"`
	CreatedAt       *gtime.Time `json:"created_at"        description:""`
}
