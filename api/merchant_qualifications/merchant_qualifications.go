package merchantQualifications

import (
	"context"

	"gf-eshop/api/merchant_qualifications/v1"
)

type IMerchantQualificationsV1 interface {
	List(ctx context.Context, req *v1.MerchantQualificationsListReq) (res *v1.MerchantQualificationsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantQualificationsDetailReq) (res *v1.MerchantQualificationsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantQualificationsCreateReq) (res *v1.MerchantQualificationsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantQualificationsUpdateReq) (res *v1.MerchantQualificationsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantQualificationsDeleteReq) (res *v1.MerchantQualificationsDeleteRes, err error)
	Audit(ctx context.Context, req *v1.MerchantQualificationsAuditReq) (res *v1.MerchantQualificationsAuditRes, err error)
}
