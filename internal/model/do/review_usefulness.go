// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewUsefulness is the golang structure of table rev_review_usefulness for DAO operations like Where/Data.
type ReviewUsefulness struct {
	g.Meta    `orm:"table:rev_review_usefulness, do:true"`
	Id        interface{} // 主键
	ReviewId  interface{} // 评价ID
	UserId    interface{} // 用户ID
	CreatedAt *gtime.Time //
}
