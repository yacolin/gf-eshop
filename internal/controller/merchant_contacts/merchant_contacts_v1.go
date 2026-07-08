package merchantContacts

import (
	"context"

	"gf-eshop/api/merchant_contacts/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantContactsListReq) (res *v1.MerchantContactsListRes, err error) {
	return service.MerchantContacts().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.MerchantContactsDetailReq) (res *v1.MerchantContactsDetailRes, err error) {
	return service.MerchantContacts().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantContactsCreateReq) (res *v1.MerchantContactsCreateRes, err error) {
	return service.MerchantContacts().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantContactsUpdateReq) (res *v1.MerchantContactsUpdateRes, err error) {
	return service.MerchantContacts().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantContactsDeleteReq) (res *v1.MerchantContactsDeleteRes, err error) {
	return service.MerchantContacts().Delete(ctx, req)
}
