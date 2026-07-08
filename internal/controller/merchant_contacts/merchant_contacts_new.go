package merchantContacts

import "gf-eshop/api/merchant_contacts"

type ControllerV1 struct{}

func NewV1() merchantContacts.IMerchantContactsV1 {
	return &ControllerV1{}
}
