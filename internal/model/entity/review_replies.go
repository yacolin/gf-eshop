// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewReplies is the golang structure for table review_replies.
type ReviewReplies struct {
	Id           int64       `json:"id"            description:"回复ID"`
	ReviewId     int64       `json:"review_id"     description:"关联评价ID"`
	RootReplyId  int64       `json:"root_reply_id" description:"根回复ID（一级回复为NULL）"`
	ParentId     int64       `json:"parent_id"     description:"父级回复ID（支持二级回复）"`
	ReplyType    int         `json:"reply_type"    description:"1-商家回复 2-用户追问 3-平台回复"`
	Content      string      `json:"content"       description:"回复内容"`
	OperatorId   int64       `json:"operator_id"   description:"操作人ID"`
	OperatorName string      `json:"operator_name" description:"操作人名称（冗余，避免JOIN用户表）"`
	Status       int         `json:"status"        description:"1-正常 2-隐藏 3-删除"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updated_at"    description:""`
}
