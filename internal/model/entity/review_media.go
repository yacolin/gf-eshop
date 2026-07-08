// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewMedia is the golang structure for table review_media.
type ReviewMedia struct {
	Id        int64       `json:"id"         description:"媒体ID"`
	ReviewId  int64       `json:"review_id"  description:"关联评价ID"`
	MediaType int         `json:"media_type" description:"1-图片 2-视频"`
	MediaUrl  string      `json:"media_url"  description:"媒体文件URL"`
	SortOrder int         `json:"sort_order" description:"排序"`
	CreatedAt *gtime.Time `json:"created_at" description:""`
}
