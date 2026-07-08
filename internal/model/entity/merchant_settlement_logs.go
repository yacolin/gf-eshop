// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantSettlementLogs is the golang structure for table merchant_settlement_logs.
type MerchantSettlementLogs struct {
	Id               int64       `json:"id"                description:"主键"`
	MerchantId       int64       `json:"merchant_id"       description:"商家ID"`
	SettlementNo     string      `json:"settlement_no"     description:"结算单号"`
	SettlementCycle  string      `json:"settlement_cycle"  description:"结算周期（如 2026-07-01~2026-07-15）"`
	TotalAmount      int64       `json:"total_amount"      description:"期内总金额（分）"`
	CommissionAmount int64       `json:"commission_amount" description:"平台佣金（分）"`
	SettlementAmount int64       `json:"settlement_amount" description:"应结算金额（分）"`
	Status           int         `json:"status"            description:"0-待结算 1-已结算 2-已打款"`
	SettledAt        *gtime.Time `json:"settled_at"        description:"结算时间"`
	PaidAt           *gtime.Time `json:"paid_at"           description:"打款时间"`
	Remark           string      `json:"remark"            description:"备注"`
	CreatedAt        *gtime.Time `json:"created_at"        description:""`
	UpdatedAt        *gtime.Time `json:"updated_at"        description:""`
	DeletedAt        *gtime.Time `json:"deleted_at"        description:""`
}
