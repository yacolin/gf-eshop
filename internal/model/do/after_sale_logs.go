// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleLogs is the golang structure of table tx_after_sale_logs for DAO operations like Where/Data.
type AfterSaleLogs struct {
	g.Meta       `orm:"table:tx_after_sale_logs, do:true"`
	Id           interface{} // 主键
	AfterSaleId  interface{} // 售后单ID
	OperatorId   interface{} // 操作人ID
	OperatorType interface{} // operator/user/merchant/admin
	Action       interface{} // 动作
	Remark       interface{} // 备注
	CreatedAt    *gtime.Time //
}
