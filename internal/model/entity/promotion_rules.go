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
	MerchantId     int64       `json:"merchant_id"     description:"所属商家ID"`
	RuleName       string      `json:"rule_name"       description:"规则名称"`
	ConditionType  int         `json:"condition_type"  description:"1-无门槛 2-满金额 3-满件数 4-指定用户等级"`
	ConditionValue int64       `json:"condition_value" description:"门槛值（分）"`
	BenefitConfig  string      `json:"benefit_config"  description:"优惠配置JSON。例：{\"type\":1,\"value\":3000} 或 {\"type\":2,\"steps\":[{\"limit\":10000,\"rate\":900},{\"limit\":20000,\"rate\":800}]}"`
	IsStackable    int         `json:"is_stackable"    description:"是否可与其他促销叠加 0-否 1-是"`
	StackGroup     int         `json:"stack_group"     description:"叠加组ID（同组内互斥，不同组可叠加）"`
	CreatedBy      int64       `json:"created_by"      description:"创建人"`
	UpdatedBy      int64       `json:"updated_by"      description:"更新人"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:""`
}
