// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLogs is the golang structure for table operation_logs.
type OperationLogs struct {
	Id            int64       `json:"id"             description:"主键"`
	StaffId       int64       `json:"staff_id"       description:"操作员工ID（关联 sys_staff.id）"`
	StaffName     string      `json:"staff_name"     description:"操作员工姓名（冗余，便于查询）"`
	Operation     string      `json:"operation"      description:"操作类型：update_price/disable_user/create_coupon/..."`
	Resource      string      `json:"resource"       description:"操作资源：product/order/user/coupon/..."`
	ResourceId    string      `json:"resource_id"    description:"资源ID"`
	Detail        string      `json:"detail"         description:"操作详情JSON（记录变更前后快照）"`
	Result        int         `json:"result"         description:"1-成功 0-失败"`
	FailureReason string      `json:"failure_reason" description:"失败原因"`
	Ip            string      `json:"ip"             description:"操作IP"`
	CreatedAt     *gtime.Time `json:"created_at"     description:"创建时间"`
}
