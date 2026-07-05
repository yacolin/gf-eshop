// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderLogs is the golang structure of table tx_order_logs for DAO operations like Where/Data.
type OrderLogs struct {
	g.Meta       `orm:"table:tx_order_logs, do:true"`
	Id           interface{} // 日志ID
	OrderId      interface{} // 关联 tx_orders.id
	OrderNo      interface{} // 订单号（冗余，便于按号查日志）
	FromStatus   interface{} // 变更前状态
	ToStatus     interface{} // 变更后状态
	Operator     interface{} // 操作人
	OperatorType interface{} // 操作人类型：system-系统 user-用户 admin-管理员
	Note         interface{} // 备注（如：支付成功、超时取消）
	CreatedAt    *gtime.Time //
}
