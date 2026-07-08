// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewReplies is the golang structure for table review_replies.
type ReviewReplies struct {
	Id         int64       `json:"id"          description:"回复ID"`
	ReviewId   int64       `json:"review_id"   description:"关联评价ID"`
	ParentId   int64       `json:"parent_id"   description:"父级回复ID（支持多级回复）"`
	ReplyType  int         `json:"reply_type"  description:"1-商家回复 2-用户追问 3-平台回复"`
	Content    string      `json:"content"     description:"回复内容"`
	OperatorId int64       `json:"operator_id" description:"操作人ID"`
	CreatedAt  *gtime.Time `json:"created_at"  description:""`
	DeletedAt  *gtime.Time `json:"deleted_at"  description:""`
}
