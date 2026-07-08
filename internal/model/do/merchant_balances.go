// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBalances is the golang structure of table mch_merchant_balances for DAO operations like Where/Data.
type MerchantBalances struct {
	g.Meta           `orm:"table:mch_merchant_balances, do:true"`
	Id               interface{} // 主键
	MerchantId       interface{} // 商家ID
	AvailableBalance interface{} // 可提现余额（分）
	FreezeBalance    interface{} // 冻结余额（分）
	Currency         interface{} // 币种
	Version          interface{} // 版本号（并发控制）
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
