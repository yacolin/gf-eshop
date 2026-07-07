// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Points is the golang structure for table points.
type Points struct {
	Id           int64       `json:"id"            description:"流水ID"`
	UserId       int64       `json:"user_id"       description:"用户ID"`
	Points       int64       `json:"points"        description:"积分变动（正=增加，负=扣减）"`
	BalanceAfter int64       `json:"balance_after" description:"变动后积分余额"`
	Source       string      `json:"source"        description:"积分来源：order-下单消费 review-评价 signin-签到 admin-管理员调整 refund-退款扣减 expire-过期清零"`
	SourceId     string      `json:"source_id"     description:"来源ID（如订单号、评价ID）"`
	ExpireAt     *gtime.Time `json:"expire_at"     description:"过期时间（NULL=永不过期）"`
	Status       int         `json:"status"        description:"0-待确认 1-已确认 2-已过期 3-已作废"`
	Remark       string      `json:"remark"        description:"备注"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
}
