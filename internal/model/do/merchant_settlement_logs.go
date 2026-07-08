// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantSettlementLogs is the golang structure of table mch_merchant_settlement_logs for DAO operations like Where/Data.
type MerchantSettlementLogs struct {
	g.Meta           `orm:"table:mch_merchant_settlement_logs, do:true"`
	Id               interface{} // 主键
	MerchantId       interface{} // 商家ID
	SettlementNo     interface{} // 结算单号
	SettlementCycle  interface{} // 结算周期（如 2026-07-01~2026-07-15）
	TotalAmount      interface{} // 期内总金额（分）
	CommissionAmount interface{} // 平台佣金（分）
	SettlementAmount interface{} // 应结算金额（分）
	Status           interface{} // 0-待结算 1-已结算 2-已打款
	SettledAt        *gtime.Time // 结算时间
	PaidAt           *gtime.Time // 打款时间
	Remark           interface{} // 备注
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
