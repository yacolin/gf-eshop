// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleLogs is the golang structure for table after_sale_logs.
type AfterSaleLogs struct {
	Id           int64       `json:"id"            description:"主键"`
	AfterSaleId  int64       `json:"after_sale_id" description:"售后单ID"`
	OperatorId   int64       `json:"operator_id"   description:"操作人ID"`
	OperatorType string      `json:"operator_type" description:"operator/user/merchant/admin"`
	Action       string      `json:"action"        description:"动作"`
	Remark       string      `json:"remark"        description:"备注"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
}
