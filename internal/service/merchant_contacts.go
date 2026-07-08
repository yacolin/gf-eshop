package service

import (
	"context"

	"gf-eshop/api/merchant_contacts/v1"
)

type IMerchantContacts interface {
	List(ctx context.Context, req *v1.MerchantContactsListReq) (res *v1.MerchantContactsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantContactsDetailReq) (res *v1.MerchantContactsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantContactsCreateReq) (res *v1.MerchantContactsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantContactsUpdateReq) (res *v1.MerchantContactsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantContactsDeleteReq) (res *v1.MerchantContactsDeleteRes, err error)
}

var localMerchantContacts IMerchantContacts

func MerchantContacts() IMerchantContacts {
	if localMerchantContacts == nil {
		panic("implement not found for interface IMerchantContacts, forgot register?")
	}
	return localMerchantContacts
}

func RegisterMerchantContacts(i IMerchantContacts) {
	localMerchantContacts = i
}
