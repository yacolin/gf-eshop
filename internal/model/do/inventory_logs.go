// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLogs is the golang structure of table sp_inventory_logs for DAO operations like Where/Data.
type InventoryLogs struct {
	g.Meta         `orm:"table:sp_inventory_logs, do:true"`
	Id             interface{} // 日志ID
	SkuId          interface{} // 关联 skus.id
	MerchantId     interface{} // 所属商家ID
	WarehouseId    interface{} // 仓库ID
	ChangeType     interface{} // 变更类型：order_lock-下单预占 order_unlock-取消释放 order_deduct-支付扣减 inbound-入库 outbound-出库 return-退货入库 adjust-盘盈亏修正
	BeforeQuantity interface{} // 变更前物理库存
	AfterQuantity  interface{} // 变更后物理库存
	BeforeReserved interface{} // 变更前预占库存
	AfterReserved  interface{} // 变更后预占库存
	ChangeAmount   interface{} // 变更数量（正=增加，负=减少）
	ReferenceId    interface{} // 关联单据ID（如订单号、入库单号）
	Operator       interface{} // 操作人（系统操作填 system）
	Note           interface{} // 备注
	CreatedAt      *gtime.Time //
}
