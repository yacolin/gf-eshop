// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Levels is the golang structure for table levels.
type Levels struct {
	Id               int64       `json:"id"                description:"等级ID"`
	Name             string      `json:"name"              description:"等级名称（如：青铜会员、白银会员、黄金会员、钻石会员）"`
	Icon             string      `json:"icon"              description:"等级图标URL"`
	Level            int         `json:"level"             description:"等级数值（1=青铜 2=白银 3=黄金 4=钻石）"`
	MinPoints        int64       `json:"min_points"        description:"该等级所需最低累计积分"`
	MaxPoints        int64       `json:"max_points"        description:"该等级所需最高累计积分（0表示无上限）"`
	DiscountRate     int64       `json:"discount_rate"     description:"折扣率（千分比，1000=无折扣，900=九折）"`
	FreeShipping     int         `json:"free_shipping"     description:"1-免运费"`
	PointsMultiplier float64     `json:"points_multiplier" description:"消费积分倍数（如 1.5 倍积分）"`
	Benefits         string      `json:"benefits"          description:"扩展权益JSON（如：{\"birthday_gift\": true, \"exclusive_coupon\": true}）"`
	Status           int         `json:"status"            description:"1-启用 0-禁用"`
	SortOrder        int         `json:"sort_order"        description:"排序"`
	CreatedAt        *gtime.Time `json:"created_at"        description:""`
	UpdatedAt        *gtime.Time `json:"updated_at"        description:""`
	DeletedAt        *gtime.Time `json:"deleted_at"        description:""`
}
