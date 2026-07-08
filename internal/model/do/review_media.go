// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewMedia is the golang structure of table rev_review_media for DAO operations like Where/Data.
type ReviewMedia struct {
	g.Meta    `orm:"table:rev_review_media, do:true"`
	Id        interface{} // 媒体ID
	ReviewId  interface{} // 关联评价ID
	MediaType interface{} // 1-图片 2-视频
	MediaUrl  interface{} // 媒体文件URL
	SortOrder interface{} // 排序
	CreatedAt *gtime.Time //
}
