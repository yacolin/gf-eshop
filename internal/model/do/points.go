// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Points is the golang structure of table usr_points for DAO operations like Where/Data.
type Points struct {
	g.Meta       `orm:"table:usr_points, do:true"`
	Id           interface{} // 流水ID
	UserId       interface{} // 用户ID
	Points       interface{} // 积分变动（正=增加，负=扣减）
	BalanceAfter interface{} // 变动后积分余额
	Source       interface{} // 积分来源：order-下单消费 review-评价 signin-签到 admin-管理员调整 refund-退款扣减 expire-过期清零
	SourceId     interface{} // 来源ID（如订单号、评价ID）
	ExpireAt     *gtime.Time // 过期时间（NULL=永不过期）
	Status       interface{} // 0-待确认 1-已确认 2-已过期 3-已作废
	Remark       interface{} // 备注
	CreatedAt    *gtime.Time //
}
