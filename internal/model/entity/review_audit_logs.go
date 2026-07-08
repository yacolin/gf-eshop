// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewAuditLogs is the golang structure for table review_audit_logs.
type ReviewAuditLogs struct {
	Id         int64       `json:"id"          description:"主键"`
	ReviewId   int64       `json:"review_id"   description:"评价ID"`
	Action     string      `json:"action"      description:"操作：submit/approve/reject/delete"`
	OperatorId int64       `json:"operator_id" description:"操作人ID"`
	Remark     string      `json:"remark"      description:"操作备注"`
	CreatedAt  *gtime.Time `json:"created_at"  description:""`
}
