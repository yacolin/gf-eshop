// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewUsefulness is the golang structure for table review_usefulness.
type ReviewUsefulness struct {
	Id        int64       `json:"id"         description:"主键"`
	ReviewId  int64       `json:"review_id"  description:"评价ID"`
	UserId    int64       `json:"user_id"    description:"用户ID"`
	CreatedAt *gtime.Time `json:"created_at" description:""`
}
