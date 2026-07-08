package merchantBalances

import (
	"context"

	"gf-eshop/api/merchant_balances/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantBalances struct{}

func init() {
	service.RegisterMerchantBalances(&sMerchantBalances{})
}

func (s *sMerchantBalances) GetByMerchant(ctx context.Context, req *v1.GetByMerchantReq) (res *v1.GetByMerchantRes, err error) {
	var bal *entity.MerchantBalances
	err = dao.MerchantBalances.Ctx(ctx).Where(dao.MerchantBalances.Columns().MerchantId, req.MerchantId).Scan(&bal)
	if err != nil {
		return nil, err
	}
	if bal == nil {
		return &v1.GetByMerchantRes{MerchantBalances: &entity.MerchantBalances{
			MerchantId:       req.MerchantId,
			AvailableBalance: 0,
			FreezeBalance:    0,
			Currency:         "CNY",
		}}, nil
	}
	return &v1.GetByMerchantRes{MerchantBalances: bal}, nil
}
