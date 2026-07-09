// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Carts is the golang structure of table tx_carts for DAO operations like Where/Data.
type Carts struct {
	g.Meta    `orm:"table:tx_carts, do:true"`
	Id        interface{} // 购物车ID
	UserId    interface{} // 用户ID（已登录用户）
	SessionId interface{} // 会话ID（未登录时的临时标识）
	ExpiredAt *gtime.Time // 过期时间（session 型购物车自动清理）
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
