// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SettlementDetails is the golang structure of table mch_settlement_details for DAO operations like Where/Data.
type SettlementDetails struct {
	g.Meta           `orm:"table:mch_settlement_details, do:true"`
	Id               interface{} // 主键
	MerchantId       interface{} // 商家ID
	OrderId          interface{} // 订单ID
	SettlementLogId  interface{} // 结算流水ID
	OrderAmount      interface{} // 订单实付金额（分）
	CommissionAmount interface{} // 平台佣金（分）
	SettlementAmount interface{} // 应结算金额（分）
	RefundAmount     interface{} // 退款冲减金额（分）
	Status           interface{} // 0-待结算 1-已结算 2-已冲减
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
