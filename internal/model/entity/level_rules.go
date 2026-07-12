// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LevelRules is the golang structure for table level_rules.
type LevelRules struct {
	Id             int         `json:"id"              description:""`
	Name           string      `json:"name"            description:"规则名称"`
	RuleType       string      `json:"rule_type"       description:"规则类型：upgrade-自动升级 downgrade-自动降级"`
	FromLevelId    int64       `json:"from_level_id"   description:"源等级ID（0=任意等级）"`
	ToLevelId      int64       `json:"to_level_id"     description:"目标等级ID"`
	ConditionType  string      `json:"condition_type"  description:"条件类型：points-累计积分 order_count-订单数 order_amount-消费金额"`
	ConditionValue int64       `json:"condition_value" description:"条件阈值"`
	Description    string      `json:"description"     description:"规则说明"`
	SortOrder      int         `json:"sort_order"      description:"排序"`
	Status         int         `json:"status"          description:"状态：0-禁用 1-启用"`
	CreatedAt      *gtime.Time `json:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:"更新时间"`
}
