package level_rules

import "gf-eshop/api/level_rules"

type ControllerV1 struct{}

func NewV1() level_rules.ILevelRulesV1 {
	return &ControllerV1{}
}
