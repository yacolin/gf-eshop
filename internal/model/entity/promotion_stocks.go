// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionStocks is the golang structure for table promotion_stocks.
type PromotionStocks struct {
	Id             int64       `json:"id"              description:""`
	PromotionId    int64       `json:"promotion_id"    description:"促销ID"`
	SkuId          int64       `json:"sku_id"          description:"SKU ID（秒杀专用，通用活动可为空）"`
	TotalStock     int         `json:"total_stock"     description:"总库存"`
	AvailableStock int         `json:"available_stock" description:"可用库存"`
	LockedStock    int         `json:"locked_stock"    description:"锁定库存（下单未付）"`
	Version        int         `json:"version"         description:"乐观锁版本号"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:""`
}
