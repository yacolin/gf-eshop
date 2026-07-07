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
	MerchantId     interface{} // 所属商家ID（0表示平台级规则）
	RuleName       interface{} // 规则名称（便于理解）
	ConditionType  interface{} // 1-无门槛 2-满金额 3-满件数 4-指定用户等级
	ConditionValue interface{} // 门槛值（分，满20000则存20000）
	BenefitType    interface{} // 1-减固定金额 2-打折扣 3-赠品 4-免运费 5-送积分
	BenefitValue   interface{} // 优惠值（减固定金额填分如3000；打折扣填千分比如800=8折）
	IsStackable    interface{} // 是否可与其他促销叠加 0-否 1-是
	StackPriority  interface{} // 叠加优先级（数字越小越优先计算）
	CreatedBy      interface{} // 创建人
	UpdatedBy      interface{} // 更新人
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 软删除时间
}
