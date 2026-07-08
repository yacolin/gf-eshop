// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewReplies is the golang structure of table rev_review_replies for DAO operations like Where/Data.
type ReviewReplies struct {
	g.Meta     `orm:"table:rev_review_replies, do:true"`
	Id         interface{} // 回复ID
	ReviewId   interface{} // 关联评价ID
	ParentId   interface{} // 父级回复ID（支持多级回复）
	ReplyType  interface{} // 1-商家回复 2-用户追问 3-平台回复
	Content    interface{} // 回复内容
	OperatorId interface{} // 操作人ID
	CreatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
