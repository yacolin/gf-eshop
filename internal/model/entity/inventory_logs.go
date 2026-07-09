// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLogs is the golang structure for table inventory_logs.
type InventoryLogs struct {
	Id             int64       `json:"id"              description:"日志ID"`
	SkuId          int64       `json:"sku_id"          description:"关联 skus.id"`
	MerchantId     int64       `json:"merchant_id"     description:"所属商家ID"`
	WarehouseId    int64       `json:"warehouse_id"    description:"仓库ID"`
	ChangeType     string      `json:"change_type"     description:"变更类型：order_lock-下单预占 order_unlock-取消释放 order_deduct-支付扣减 inbound-入库 outbound-出库 return-退货入库 adjust-盘盈亏修正"`
	BeforeQuantity int64       `json:"before_quantity" description:"变更前物理库存"`
	AfterQuantity  int64       `json:"after_quantity"  description:"变更后物理库存"`
	BeforeReserved int64       `json:"before_reserved" description:"变更前预占库存"`
	AfterReserved  int64       `json:"after_reserved"  description:"变更后预占库存"`
	ChangeAmount   int64       `json:"change_amount"   description:"变更数量（正=增加，负=减少）"`
	ReferenceId    string      `json:"reference_id"    description:"关联单据ID（如订单号、入库单号）"`
	Operator       string      `json:"operator"        description:"操作人（系统操作填 system）"`
	Note           string      `json:"note"            description:"备注"`
	CreatedAt      *gtime.Time `json:"created_at"      description:""`
}
