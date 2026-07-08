package merchantContacts

import (
	"context"

	"gf-eshop/api/merchant_contacts/v1"
)

type IMerchantContactsV1 interface {
	List(ctx context.Context, req *v1.MerchantContactsListReq) (res *v1.MerchantContactsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantContactsDetailReq) (res *v1.MerchantContactsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantContactsCreateReq) (res *v1.MerchantContactsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantContactsUpdateReq) (res *v1.MerchantContactsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantContactsDeleteReq) (res *v1.MerchantContactsDeleteRes, err error)
}
