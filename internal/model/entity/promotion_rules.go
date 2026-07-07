// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionRules is the golang structure for table promotion_rules.
type PromotionRules struct {
	Id             int64       `json:"id"              description:"规则ID"`
	PromotionId    int64       `json:"promotion_id"    description:"所属促销ID"`
	MerchantId     int64       `json:"merchant_id"     description:"所属商家ID（0表示平台级规则）"`
	RuleName       string      `json:"rule_name"       description:"规则名称（便于理解）"`
	ConditionType  int         `json:"condition_type"  description:"1-无门槛 2-满金额 3-满件数 4-指定用户等级"`
	ConditionValue int64       `json:"condition_value" description:"门槛值（分，满20000则存20000）"`
	BenefitType    int         `json:"benefit_type"    description:"1-减固定金额 2-打折扣 3-赠品 4-免运费 5-送积分"`
	BenefitValue   int64       `json:"benefit_value"   description:"优惠值（减固定金额填分如3000；打折扣填千分比如800=8折）"`
	IsStackable    int         `json:"is_stackable"    description:"是否可与其他促销叠加 0-否 1-是"`
	StackPriority  int         `json:"stack_priority"  description:"叠加优先级（数字越小越优先计算）"`
	CreatedBy      int64       `json:"created_by"      description:"创建人"`
	UpdatedBy      int64       `json:"updated_by"      description:"更新人"`
	CreatedAt      *gtime.Time `json:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:"更新时间"`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:"软删除时间"`
}
