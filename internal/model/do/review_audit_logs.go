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
	g.Meta       `orm:"table:rev_review_audit_logs, do:true"`
	Id           interface{} // 主键
	ReviewId     interface{} // 评价ID
	Action       interface{} // 操作：submit/approve/reject/delete/shield
	OperatorId   interface{} // 操作人ID
	OperatorName interface{} // 操作人名称
	BeforeStatus interface{} // 变更前状态
	AfterStatus  interface{} // 变更后状态
	Remark       interface{} // 操作备注
	Snapshot     interface{} // 评价快照（用于回溯）
	CreatedAt    *gtime.Time //
}
