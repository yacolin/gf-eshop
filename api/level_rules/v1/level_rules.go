package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type LevelRulesListReq struct {
	g.Meta      `path:"/level-rules" tags:"LevelRules" method:"get" summary:"升降级规则列表"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	RuleType    string `json:"rule_type"    description:"规则类型筛选：upgrade-升级 downgrade-降级"`
	Status      *int   `json:"status"       description:"状态筛选"`
}

type LevelRulesListRes struct {
	List  []*entity.LevelRules `json:"list"`
	Total int                  `json:"total"`
}

type LevelRulesDetailReq struct {
	g.Meta `path:"/level-rules/{id}" tags:"LevelRules" method:"get" summary:"规则详情"`
	Id     int `json:"id" v:"required#规则ID不能为空"`
}

type LevelRulesDetailRes struct {
	*entity.LevelRules
}

type LevelRulesCreateReq struct {
	g.Meta         `path:"/level-rules" tags:"LevelRules" method:"post" summary:"新增升降级规则"`
	Name           string `json:"name"            v:"required#规则名称不能为空" description:"规则名称"`
	RuleType       string `json:"rule_type"       v:"required#规则类型不能为空" description:"规则类型：upgrade-升级 downgrade-降级"`
	FromLevelId    int64  `json:"from_level_id"   description:"源等级ID（0=任意等级）"`
	ToLevelId      int64  `json:"to_level_id"     description:"目标等级ID"`
	ConditionType  string `json:"condition_type"  v:"required#条件类型不能为空" description:"条件类型：points-累计积分 order_count-订单数 order_amount-消费金额 inactive_days-无消费天数"`
	ConditionValue int64  `json:"condition_value" v:"required#条件阈值不能为空" description:"条件阈值"`
	Description    string `json:"description"     description:"规则说明"`
	SortOrder      int    `json:"sort_order"      description:"排序"`
	Status         int    `json:"status"          description:"状态"`
}

type LevelRulesCreateRes struct {
	Id int `json:"id"`
}

type LevelRulesUpdateReq struct {
	g.Meta         `path:"/level-rules/{id}" tags:"LevelRules" method:"put" summary:"更新升降级规则"`
	Id             int    `json:"id"              v:"required#规则ID不能为空"`
	Name           string `json:"name"            description:"规则名称"`
	RuleType       string `json:"rule_type"       description:"规则类型"`
	FromLevelId    int64  `json:"from_level_id"   description:"源等级ID"`
	ToLevelId      int64  `json:"to_level_id"     description:"目标等级ID"`
	ConditionType  string `json:"condition_type"  description:"条件类型"`
	ConditionValue int64  `json:"condition_value" description:"条件阈值"`
	Description    string `json:"description"     description:"规则说明"`
	SortOrder      *int   `json:"sort_order"      description:"排序"`
	Status         *int   `json:"status"          description:"状态"`
}

type LevelRulesUpdateRes struct{}

type LevelRulesDeleteReq struct {
	g.Meta `path:"/level-rules/{id}" tags:"LevelRules" method:"delete" summary:"删除升降级规则"`
	Id     int `json:"id" v:"required#规则ID不能为空"`
}

type LevelRulesDeleteRes struct{}
