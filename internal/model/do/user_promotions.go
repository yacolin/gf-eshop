// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPromotions is the golang structure of table mkt_user_promotions for DAO operations like Where/Data.
type UserPromotions struct {
	g.Meta          `orm:"table:mkt_user_promotions, do:true"`
	Id              interface{} // 主键ID
	UserPromotionNo interface{} // 用户促销资产编号
	UserId          interface{} // 用户ID
	PromotionId     interface{} // 促销ID
	MerchantId      interface{} // 所属商家ID
	AcquireTime     *gtime.Time // 领取时间
	ExpireTime      *gtime.Time // 过期时间
	Status          interface{} // 1-未使用 2-锁定中(下单未付) 3-已使用 4-已过期 5-已作废
	LockOrderId     interface{} // 锁定的订单ID（用于回滚）
	UsedTime        *gtime.Time // 使用时间
	OrderId         interface{} // 最终使用的订单ID
	QueueToken      interface{} // 秒杀排队令牌
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
