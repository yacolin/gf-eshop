// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBalances is the golang structure for table merchant_balances.
type MerchantBalances struct {
	Id               int64       `json:"id"                description:"主键"`
	MerchantId       int64       `json:"merchant_id"       description:"商家ID"`
	AvailableBalance int64       `json:"available_balance" description:"可提现余额（分）"`
	FreezeBalance    int64       `json:"freeze_balance"    description:"冻结余额（分）"`
	Currency         string      `json:"currency"          description:"币种"`
	Version          int64       `json:"version"           description:"版本号（并发控制）"`
	CreatedAt        *gtime.Time `json:"created_at"        description:""`
	UpdatedAt        *gtime.Time `json:"updated_at"        description:""`
	DeletedAt        *gtime.Time `json:"deleted_at"        description:""`
}
