// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Levels is the golang structure of table usr_levels for DAO operations like Where/Data.
type Levels struct {
	g.Meta           `orm:"table:usr_levels, do:true"`
	Id               interface{} // 等级ID
	Name             interface{} // 等级名称（如：青铜会员、白银会员、黄金会员、钻石会员）
	Level            interface{} // 等级数值（1=青铜 2=白银 3=黄金 4=钻石）
	Icon             interface{} // 等级图标URL
	MinPoints        interface{} // 该等级所需最低累计积分
	DiscountRate     interface{} // 折扣率（千分比，1000=无折扣，900=九折）
	FreeShipping     interface{} // 1-免运费
	PointsMultiplier interface{} // 消费积分倍数（如 1.5 倍积分）
	Benefits         interface{} // 扩展权益JSON（如：{"birthday_gift": true, "exclusive_coupon": true}）
	Status           interface{} // 1-启用 0-禁用
	SortOrder        interface{} // 排序
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
