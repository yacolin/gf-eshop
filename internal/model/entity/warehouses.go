// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Warehouses is the golang structure for table warehouses.
type Warehouses struct {
	Id            int64       `json:"id"             description:"仓库ID"`
	MerchantId    int64       `json:"merchant_id"    description:"所属商家ID（0表示平台仓）"`
	WarehouseName string      `json:"warehouse_name" description:"仓库名称"`
	WarehouseType int         `json:"warehouse_type" description:"1-平台仓 2-商家仓 3-第三方仓"`
	Status        int         `json:"status"         description:"1-启用 2-禁用"`
	Province      string      `json:"province"       description:"省"`
	City          string      `json:"city"           description:"市"`
	District      string      `json:"district"       description:"区"`
	Address       string      `json:"address"        description:"详细地址"`
	CreatedAt     *gtime.Time `json:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updated_at"     description:""`
	DeletedAt     *gtime.Time `json:"deleted_at"     description:""`
}
