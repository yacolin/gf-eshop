// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Promotions is the golang structure for table promotions.
type Promotions struct {
	Id            int64       `json:"id"             description:"促销ID"`
	MerchantId    int64       `json:"merchant_id"    description:"所属商家ID（0表示平台级活动）"`
	PromoName     string      `json:"promo_name"     description:"活动名称"`
	PromoType     int         `json:"promo_type"     description:"1-满减券 2-折扣券 3-秒杀 4-满额减 5-满件折 6-会员价"`
	PromoCode     string      `json:"promo_code"     description:"优惠码（优惠券专用）"`
	StartTime     *gtime.Time `json:"start_time"     description:"开始时间"`
	EndTime       *gtime.Time `json:"end_time"       description:"结束时间"`
	TotalQuantity int         `json:"total_quantity" description:"发行总量（0表示不限）"`
	PerUserLimit  int         `json:"per_user_limit" description:"每人限领/限购数量"`
	UsedQuantity  int         `json:"used_quantity"  description:"已使用/已售数量"`
	RuleId        int64       `json:"rule_id"        description:"关联规则表（mkt_promotion_rules）"`
	Status        int         `json:"status"         description:"1-草稿 2-生效中 3-已结束 4-已作废"`
	CreatedBy     int64       `json:"created_by"     description:"创建人"`
	UpdatedBy     int64       `json:"updated_by"     description:"更新人"`
	CreatedAt     *gtime.Time `json:"created_at"     description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updated_at"     description:"更新时间"`
	DeletedAt     *gtime.Time `json:"deleted_at"     description:"软删除时间"`
}
