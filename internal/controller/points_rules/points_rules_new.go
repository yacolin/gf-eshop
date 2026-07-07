package points_rules

import "gf-eshop/api/points_rules"

type ControllerV1 struct{}

func NewV1() points_rules.IPointsRulesV1 {
	return &ControllerV1{}
}
