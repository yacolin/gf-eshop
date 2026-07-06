// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Carts is the golang structure for table carts.
type Carts struct {
	Id          int64       `json:"id"           description:"购物车ID"`
	UserId      int64       `json:"user_id"      description:"用户ID（已登录用户）"`
	SessionId   string      `json:"session_id"   description:"会话ID（未登录时的临时标识）"`
	ItemCount   int         `json:"item_count"   description:"商品种类数"`
	TotalAmount int64       `json:"total_amount" description:"总金额（分，聚合，减少查SKU次数）"`
	ExpiredAt   *gtime.Time `json:"expired_at"   description:"过期时间（session 型购物车自动清理）"`
	CreatedAt   *gtime.Time `json:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:""`
}
