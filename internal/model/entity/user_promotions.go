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
	ExpireTime      *gtime.Time `json:"expire_time"       description:"过期时间"`
	Status          int         `json:"status"            description:"1-未使用 2-锁定中(下单未付) 3-已使用 4-已过期 5-已作废"`
	LockOrderId     int64       `json:"lock_order_id"     description:"锁定的订单ID（用于回滚）"`
	UsedTime        *gtime.Time `json:"used_time"         description:"使用时间"`
	OrderId         int64       `json:"order_id"          description:"最终使用的订单ID"`
	QueueToken      string      `json:"queue_token"       description:"秒杀排队令牌"`
	CreatedAt       *gtime.Time `json:"created_at"        description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"        description:""`
	DeletedAt       *gtime.Time `json:"deleted_at"        description:""`
}
