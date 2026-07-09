// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewAuditLogs is the golang structure for table review_audit_logs.
type ReviewAuditLogs struct {
	Id           int64       `json:"id"            description:"主键"`
	ReviewId     int64       `json:"review_id"     description:"评价ID"`
	Action       string      `json:"action"        description:"操作：submit/approve/reject/delete/shield"`
	OperatorId   int64       `json:"operator_id"   description:"操作人ID"`
	OperatorName string      `json:"operator_name" description:"操作人名称"`
	BeforeStatus int         `json:"before_status" description:"变更前状态"`
	AfterStatus  int         `json:"after_status"  description:"变更后状态"`
	Remark       string      `json:"remark"        description:"操作备注"`
	Snapshot     string      `json:"snapshot"      description:"评价快照（用于回溯）"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
}
