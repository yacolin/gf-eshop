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
	g.Meta            `orm:"table:mkt_promotion_usage_logs, do:true"`
	Id                interface{} // 主键ID
	PromotionId       interface{} // 促销ID
	UserPromotionId   interface{} // 用户促销资产ID
	UserId            interface{} // 用户ID
	OrderId           interface{} // 使用的订单ID
	DiscountAmount    interface{} // 优惠金额（分）
	PromotionSnapshot interface{} // 优惠快照（名称、规则等）
	CreatedAt         *gtime.Time //
}
