package merchantQualifications

import (
	"context"

	"gf-eshop/api/merchant_qualifications/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.MerchantQualificationsListReq) (res *v1.MerchantQualificationsListRes, err error) {
	return service.MerchantQualifications().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.MerchantQualificationsDetailReq) (res *v1.MerchantQualificationsDetailRes, err error) {
	return service.MerchantQualifications().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.MerchantQualificationsCreateReq) (res *v1.MerchantQualificationsCreateRes, err error) {
	return service.MerchantQualifications().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.MerchantQualificationsUpdateReq) (res *v1.MerchantQualificationsUpdateRes, err error) {
	return service.MerchantQualifications().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.MerchantQualificationsDeleteReq) (res *v1.MerchantQualificationsDeleteRes, err error) {
	return service.MerchantQualifications().Delete(ctx, req)
}

func (c *ControllerV1) Audit(ctx context.Context, req *v1.MerchantQualificationsAuditReq) (res *v1.MerchantQualificationsAuditRes, err error) {
	return service.MerchantQualifications().Audit(ctx, req)
}
