// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionRules is the golang structure of table mkt_promotion_rules for DAO operations like Where/Data.
type PromotionRules struct {
	g.Meta         `orm:"table:mkt_promotion_rules, do:true"`
	Id             interface{} // 规则ID
	PromotionId    interface{} // 所属促销ID
	MerchantId     interface{} // 所属商家ID
	RuleName       interface{} // 规则名称
	ConditionType  interface{} // 1-无门槛 2-满金额 3-满件数 4-指定用户等级
	ConditionValue interface{} // 门槛值（分）
	BenefitConfig  interface{} // 优惠配置JSON。例：{"type":1,"value":3000} 或 {"type":2,"steps":[{"limit":10000,"rate":900},{"limit":20000,"rate":800}]}
	IsStackable    interface{} // 是否可与其他促销叠加 0-否 1-是
	StackGroup     interface{} // 叠加组ID（同组内互斥，不同组可叠加）
	CreatedBy      interface{} // 创建人
	UpdatedBy      interface{} // 更新人
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
