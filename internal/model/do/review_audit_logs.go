// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewAuditLogs is the golang structure of table rev_review_audit_logs for DAO operations like Where/Data.
type ReviewAuditLogs struct {
	g.Meta     `orm:"table:rev_review_audit_logs, do:true"`
	Id         interface{} // 主键
	ReviewId   interface{} // 评价ID
	Action     interface{} // 操作：submit/approve/reject/delete
	OperatorId interface{} // 操作人ID
	Remark     interface{} // 操作备注
	CreatedAt  *gtime.Time //
}
