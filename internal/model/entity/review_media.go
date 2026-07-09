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
	FileSize  int         `json:"file_size"  description:"文件大小（字节）"`
	Width     int         `json:"width"      description:"宽度（图片/视频）"`
	Height    int         `json:"height"     description:"高度（图片/视频）"`
	Duration  int         `json:"duration"   description:"时长（视频，秒）"`
	SortOrder int         `json:"sort_order" description:"排序"`
	CreatedAt *gtime.Time `json:"created_at" description:""`
}
