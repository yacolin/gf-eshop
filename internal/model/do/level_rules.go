// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LevelRules is the golang structure of table usr_level_rules for DAO operations like Where/Data.
type LevelRules struct {
	g.Meta         `orm:"table:usr_level_rules, do:true"`
	Id             interface{} //
	Name           interface{} // 规则名称
	RuleType       interface{} // 规则类型：upgrade-自动升级 downgrade-自动降级
	FromLevelId    interface{} // 源等级ID（0=任意等级）
	ToLevelId      interface{} // 目标等级ID
	ConditionType  interface{} // 条件类型：points-累计积分 order_count-订单数 order_amount-消费金额
	ConditionValue interface{} // 条件阈值
	Description    interface{} // 规则说明
	SortOrder      interface{} // 排序
	Status         interface{} // 状态：0-禁用 1-启用
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
