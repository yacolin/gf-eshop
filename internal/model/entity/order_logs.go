// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderLogs is the golang structure for table order_logs.
type OrderLogs struct {
	Id           int64       `json:"id"            description:"日志ID"`
	OrderId      int64       `json:"order_id"      description:"关联 tx_orders.id"`
	OrderNo      string      `json:"order_no"      description:"订单号（冗余，便于按号查日志）"`
	FromStatus   string      `json:"from_status"   description:"变更前状态"`
	ToStatus     string      `json:"to_status"     description:"变更后状态"`
	Operator     string      `json:"operator"      description:"操作人"`
	OperatorType string      `json:"operator_type" description:"操作人类型：system-系统 user-用户 admin-管理员"`
	Note         string      `json:"note"          description:"备注（如：支付成功、超时取消）"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
}
