// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SettlementDetails is the golang structure for table settlement_details.
type SettlementDetails struct {
	Id               int64       `json:"id"                description:"主键"`
	MerchantId       int64       `json:"merchant_id"       description:"商家ID"`
	OrderId          int64       `json:"order_id"          description:"订单ID"`
	SettlementLogId  int64       `json:"settlement_log_id" description:"结算流水ID"`
	OrderAmount      int64       `json:"order_amount"      description:"订单实付金额（分）"`
	CommissionAmount int64       `json:"commission_amount" description:"平台佣金（分）"`
	SettlementAmount int64       `json:"settlement_amount" description:"应结算金额（分）"`
	RefundAmount     int64       `json:"refund_amount"     description:"退款冲减金额（分）"`
	Status           int         `json:"status"            description:"0-待结算 1-已结算 2-已冲减"`
	CreatedAt        *gtime.Time `json:"created_at"        description:""`
	UpdatedAt        *gtime.Time `json:"updated_at"        description:""`
	DeletedAt        *gtime.Time `json:"deleted_at"        description:""`
}
