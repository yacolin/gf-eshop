package merchantBankAccounts

import (
	"context"

	"gf-eshop/api/merchant_bank_accounts/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantBankAccountsListReq) (res *v1.MerchantBankAccountsListRes, err error) {
	return service.MerchantBankAccounts().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.MerchantBankAccountsDetailReq) (res *v1.MerchantBankAccountsDetailRes, err error) {
	return service.MerchantBankAccounts().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantBankAccountsCreateReq) (res *v1.MerchantBankAccountsCreateRes, err error) {
	return service.MerchantBankAccounts().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantBankAccountsUpdateReq) (res *v1.MerchantBankAccountsUpdateRes, err error) {
	return service.MerchantBankAccounts().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantBankAccountsDeleteReq) (res *v1.MerchantBankAccountsDeleteRes, err error) {
	return service.MerchantBankAccounts().Delete(ctx, req)
}
