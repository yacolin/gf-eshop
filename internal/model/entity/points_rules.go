// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRules is the golang structure for table points_rules.
type PointsRules struct {
	Id          int         `json:"id"          description:""`
	Name        string      `json:"name"        description:"规则名称"`
	RuleKey     string      `json:"rule_key"    description:"规则键名：earn_rate-消费返积分比例 expire_days-积分过期天数 signin_points-签到奖励积分 review_points-评价奖励积分"`
	RuleValue   string      `json:"rule_value"  description:"规则值"`
	Description string      `json:"description" description:"规则说明"`
	SortOrder   int         `json:"sort_order"  description:"排序"`
	Status      int         `json:"status"      description:"状态：0-禁用 1-启用"`
	CreatedAt   *gtime.Time `json:"created_at"  description:""`
	UpdatedAt   *gtime.Time `json:"updated_at"  description:""`
}
