// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPromotions is the golang structure for table user_promotions.
type UserPromotions struct {
	Id              int64       `json:"id"                description:"主键ID"`
	UserPromotionNo string      `json:"user_promotion_no" description:"用户促销资产编号"`
	UserId          int64       `json:"user_id"           description:"用户ID"`
	PromotionId     int64       `json:"promotion_id"      description:"促销ID"`
	MerchantId      int64       `json:"merchant_id"       description:"所属商家ID"`
	AcquireTime     *gtime.Time `json:"acquire_time"      description:"领取时间"`
	ExpireTime      *gtime.Time `json:"expire_time"       description:"过期时间（优惠券必填）"`
	Status          int         `json:"status"            description:"1-未使用 2-已使用 3-已过期 4-已作废"`
	UsedTime        *gtime.Time `json:"used_time"         description:"使用时间"`
	OrderId         int64       `json:"order_id"          description:"使用的订单ID"`
	QueueToken      string      `json:"queue_token"       description:"秒杀排队令牌"`
	CreatedBy       int64       `json:"created_by"        description:"创建人"`
	UpdatedBy       int64       `json:"updated_by"        description:"更新人"`
	CreatedAt       *gtime.Time `json:"created_at"        description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updated_at"        description:"更新时间"`
	DeletedAt       *gtime.Time `json:"deleted_at"        description:"软删除时间"`
}
