// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionUsageLogs is the golang structure of table mkt_promotion_usage_logs for DAO operations like Where/Data.
type PromotionUsageLogs struct {
	g.Meta          `orm:"table:mkt_promotion_usage_logs, do:true"`
	Id              interface{} // 主键ID
	PromotionId     interface{} // 促销ID
	UserPromotionId interface{} // 用户促销资产ID
	UserId          interface{} // 用户ID
	MerchantId      interface{} // 所属商家ID
	OrderId         interface{} // 使用的订单ID
	UsageType       interface{} // 1-下单使用 2-自动优惠
	DiscountAmount  interface{} // 优惠金额（分）
	CreatedAt       *gtime.Time //
}
