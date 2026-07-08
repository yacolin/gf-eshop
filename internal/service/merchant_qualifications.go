package service

import (
	"context"

	"gf-eshop/api/merchant_qualifications/v1"
)

type IMerchantQualifications interface {
	List(ctx context.Context, req *v1.MerchantQualificationsListReq) (res *v1.MerchantQualificationsListRes, err error)
	Detail(ctx context.Context, req *v1.MerchantQualificationsDetailReq) (res *v1.MerchantQualificationsDetailRes, err error)
	Create(ctx context.Context, req *v1.MerchantQualificationsCreateReq) (res *v1.MerchantQualificationsCreateRes, err error)
	Update(ctx context.Context, req *v1.MerchantQualificationsUpdateReq) (res *v1.MerchantQualificationsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.MerchantQualificationsDeleteReq) (res *v1.MerchantQualificationsDeleteRes, err error)
	Audit(ctx context.Context, req *v1.MerchantQualificationsAuditReq) (res *v1.MerchantQualificationsAuditRes, err error)
}

var localMerchantQualifications IMerchantQualifications

func MerchantQualifications() IMerchantQualifications {
	if localMerchantQualifications == nil {
		panic("implement not found for interface IMerchantQualifications, forgot register?")
	}
	return localMerchantQualifications
}

func RegisterMerchantQualifications(i IMerchantQualifications) {
	localMerchantQualifications = i
}
