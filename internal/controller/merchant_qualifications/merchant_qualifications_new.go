package merchantQualifications

import "gf-eshop/api/merchant_qualifications"

type ControllerV1 struct{}

func NewV1() merchantQualifications.IMerchantQualificationsV1 {
	return &ControllerV1{}
}
