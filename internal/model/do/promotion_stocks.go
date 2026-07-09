// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionStocks is the golang structure of table mkt_promotion_stocks for DAO operations like Where/Data.
type PromotionStocks struct {
	g.Meta         `orm:"table:mkt_promotion_stocks, do:true"`
	Id             interface{} //
	PromotionId    interface{} // 促销ID
	SkuId          interface{} // SKU ID（秒杀专用，通用活动可为空）
	TotalStock     interface{} // 总库存
	AvailableStock interface{} // 可用库存
	LockedStock    interface{} // 锁定库存（下单未付）
	Version        interface{} // 乐观锁版本号
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
