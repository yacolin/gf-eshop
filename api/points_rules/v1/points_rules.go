package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type PointsRulesListReq struct {
	g.Meta   `path:"/points-rules" tags:"PointsRules" method:"get" summary:"积分规则列表"`
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	RuleKey  string `json:"rule_key"`  // 规则键名筛选
	Status   *int   `json:"status"`    // 状态筛选
}

type PointsRulesListRes struct {
	List  []*entity.PointsRules `json:"list"`
	Total int                   `json:"total"`
}

type PointsRulesDetailReq struct {
	g.Meta `path:"/points-rules/{id}" tags:"PointsRules" method:"get" summary:"积分规则详情"`
	Id     int `json:"id" v:"required#规则ID不能为空"`
}

type PointsRulesDetailRes struct {
	*entity.PointsRules
}

type PointsRulesCreateReq struct {
	g.Meta      `path:"/points-rules" tags:"PointsRules" method:"post" summary:"新增积分规则"`
	Name        string `json:"name"        v:"required#规则名称不能为空"` // 规则名称
	RuleKey     string `json:"rule_key"     v:"required#规则键名不能为空"` // 规则键名
	RuleValue   string `json:"rule_value"   v:"required#规则值不能为空"` // 规则值
	Description string `json:"description"`                          // 规则说明
	SortOrder   int    `json:"sort_order"`                           // 排序
	Status      int    `json:"status"`                               // 状态
}

type PointsRulesCreateRes struct {
	Id int `json:"id"`
}

type PointsRulesUpdateReq struct {
	g.Meta      `path:"/points-rules/{id}" tags:"PointsRules" method:"put" summary:"更新积分规则"`
	Id          int    `json:"id"          v:"required#规则ID不能为空"`
	Name        string `json:"name"`        // 规则名称
	RuleKey     string `json:"rule_key"`     // 规则键名
	RuleValue   string `json:"rule_value"`   // 规则值
	Description string `json:"description"`  // 规则说明
	SortOrder   *int   `json:"sort_order"`   // 排序
	Status      *int   `json:"status"`       // 状态
}

type PointsRulesUpdateRes struct{}

type PointsRulesDeleteReq struct {
	g.Meta `path:"/points-rules/{id}" tags:"PointsRules" method:"delete" summary:"删除积分规则"`
	Id     int `json:"id" v:"required#规则ID不能为空"`
}

type PointsRulesDeleteRes struct{}
