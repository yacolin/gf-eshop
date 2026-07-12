// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRules is the golang structure of table usr_points_rules for DAO operations like Where/Data.
type PointsRules struct {
	g.Meta       `orm:"table:usr_points_rules, do:true"`
	Id           interface{} //
	Name         interface{} // 规则名称
	RuleKey      interface{} // 规则键名：earn_rate-消费返积分比例 expire_days-积分过期天数 signin_points-签到奖励积分 review_points-评价奖励积分
	ValueInt     interface{} // 整数值（如积分数量、天数）
	ValueDecimal interface{} // 小数值（如比例、倍数）
	ValueString  interface{} // 字符串值（如配置json、文本）
	Description  interface{} // 规则说明
	SortOrder    interface{} // 排序
	Status       interface{} // 状态：0-禁用 1-启用
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
