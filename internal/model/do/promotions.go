// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Promotions is the golang structure of table mkt_promotions for DAO operations like Where/Data.
type Promotions struct {
	g.Meta        `orm:"table:mkt_promotions, do:true"`
	Id            interface{} // 促销ID
	MerchantId    interface{} // 所属商家ID（0表示平台级活动）
	PromoName     interface{} // 活动名称
	PromoType     interface{} // 1-满减券 2-折扣券 3-秒杀 4-满额减 5-满件折 6-会员价
	PromoCode     interface{} // 优惠码（优惠券专用）
	StartTime     *gtime.Time // 开始时间
	EndTime       *gtime.Time // 结束时间
	TotalQuantity interface{} // 发行总量（0表示不限）
	PerUserLimit  interface{} // 每人限领/限购数量
	UsedQuantity  interface{} // 已使用/已售数量
	RuleId        interface{} // 关联规则表（mkt_promotion_rules）
	Status        interface{} // 1-草稿 2-生效中 3-已结束 4-已作废
	CreatedBy     interface{} // 创建人
	UpdatedBy     interface{} // 更新人
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间
}
