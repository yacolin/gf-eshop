// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Infos is the golang structure for table infos.
type Infos struct {
	Id          int64       `json:"id"           description:"主键"`
	UserId      int64       `json:"user_id"      description:"用户ID"`
	Gender      int         `json:"gender"       description:"性别：0-未知 1-男 2-女"`
	Birthday    *gtime.Time `json:"birthday"     description:"生日"`
	Bio         string      `json:"bio"          description:"个人简介"`
	Country     string      `json:"country"      description:"国家"`
	Province    string      `json:"province"     description:"省"`
	City        string      `json:"city"         description:"市"`
	ZipCode     string      `json:"zip_code"     description:"邮编"`
	Language    string      `json:"language"     description:"语言"`
	Timezone    string      `json:"timezone"     description:"时区"`
	LevelId     int64       `json:"level_id"     description:"当前等级ID（关联 usr_levels.id）"`
	TotalPoints int64       `json:"total_points" description:"累计积分"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:"删除时间"`
}
