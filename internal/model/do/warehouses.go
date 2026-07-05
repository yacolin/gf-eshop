// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Warehouses is the golang structure of table sp_warehouses for DAO operations like Where/Data.
type Warehouses struct {
	g.Meta        `orm:"table:sp_warehouses, do:true"`
	Id            interface{} // 仓库ID
	MerchantId    interface{} // 所属商家ID（0表示平台仓）
	WarehouseName interface{} // 仓库名称
	WarehouseType interface{} // 1-平台仓 2-商家仓 3-第三方仓
	Status        interface{} // 1-启用 2-禁用
	Province      interface{} // 省
	City          interface{} // 市
	District      interface{} // 区
	Address       interface{} // 详细地址
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
