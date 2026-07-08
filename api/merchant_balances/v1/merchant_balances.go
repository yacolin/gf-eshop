package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- GetByMerchant ----------
type GetByMerchantReq struct {
	g.Meta    `path:"/merchant-balances/{merchant_id}" tags:"MerchantBalances" method:"get" summary:"商家余额"`
	MerchantId int64 `json:"merchant_id" v:"required"`
}

type GetByMerchantRes struct {
	*entity.MerchantBalances
}
