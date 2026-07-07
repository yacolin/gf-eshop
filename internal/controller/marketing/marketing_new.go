package marketing

import (
	"gf-eshop/api/marketing"
)

type ControllerV1 struct{}

func NewV1() marketing.IMarketingV1 {
	return &ControllerV1{}
}
