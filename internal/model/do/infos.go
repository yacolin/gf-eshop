// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Infos is the golang structure of table usr_infos for DAO operations like Where/Data.
type Infos struct {
	g.Meta      `orm:"table:usr_infos, do:true"`
	Id          interface{} // 主键
	UserId      interface{} // 用户ID
	Gender      interface{} // 性别：0-未知 1-男 2-女
	Birthday    *gtime.Time // 生日
	Bio         interface{} // 个人简介
	Country     interface{} // 国家
	Province    interface{} // 省
	City        interface{} // 市
	ZipCode     interface{} // 邮编
	Language    interface{} // 语言
	Timezone    interface{} // 时区
	LevelId     interface{} // 当前等级ID（关联 usr_levels.id）
	TotalPoints interface{} // 累计积分
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
