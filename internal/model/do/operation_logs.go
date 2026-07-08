// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLogs is the golang structure of table sys_operation_logs for DAO operations like Where/Data.
type OperationLogs struct {
	g.Meta        `orm:"table:sys_operation_logs, do:true"`
	Id            interface{} // 主键
	StaffId       interface{} // 操作员工ID（关联 sys_staff.id）
	StaffName     interface{} // 操作员工姓名（冗余，便于查询）
	Operation     interface{} // 操作类型：update_price/disable_user/create_coupon/...
	Resource      interface{} // 操作资源：product/order/user/coupon/...
	ResourceId    interface{} // 资源ID
	Detail        interface{} // 操作详情JSON（记录变更前后快照）
	Result        interface{} // 1-成功 0-失败
	FailureReason interface{} // 失败原因
	Ip            interface{} // 操作IP
	CreatedAt     *gtime.Time // 创建时间
}
