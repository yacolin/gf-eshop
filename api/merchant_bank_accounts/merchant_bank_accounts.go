package merchantBankAccounts

import (
	"context"

	"gf-eshop/api/merchant_bank_accounts/v1"
)

type IMerchantBankAccountsV1 interface {
	List(ctx context.Context, req *v1.MerchantBankAccountsListReq) (res *v1.MerchantBankAccountsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantBankAccountsDetailReq) (res *v1.MerchantBankAccountsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantBankAccountsCreateReq) (res *v1.MerchantBankAccountsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantBankAccountsUpdateReq) (res *v1.MerchantBankAccountsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantBankAccountsDeleteReq) (res *v1.MerchantBankAccountsDeleteRes, err error)
}
