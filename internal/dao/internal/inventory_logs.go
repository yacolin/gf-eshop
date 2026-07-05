// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryLogsDao is the data access object for table sp_inventory_logs.
type InventoryLogsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns InventoryLogsColumns // columns contains all the column names of Table for convenient usage.
}

// InventoryLogsColumns defines and stores column names for table sp_inventory_logs.
type InventoryLogsColumns struct {
	Id             string // 日志ID
	SkuId          string // 关联 skus.id
	MerchantId     string // 所属商家ID
	WarehouseId    string // 仓库ID
	ChangeType     string // 变更类型：order_lock-下单预占 order_unlock-取消释放 order_dedut-支付扣减 inbound-入库 outbound-出库 return-退货入库 adjust-盘盈亏修正
	BeforeQuantity string // 变更前物理库存
	AfterQuantity  string // 变更后物理库存
	BeforeReserved string // 变更前预占库存
	AfterReserved  string // 变更后预占库存
	ChangeAmount   string // 变更数量（正=增加，负=减少）
	ReferenceId    string // 关联单据ID（如订单号、入库单号）
	Operator       string // 操作人（系统操作填 system）
	Note           string // 备注
	CreatedAt      string //
}

// inventoryLogsColumns holds the columns for table sp_inventory_logs.
var inventoryLogsColumns = InventoryLogsColumns{
	Id:             "id",
	SkuId:          "sku_id",
	MerchantId:     "merchant_id",
	WarehouseId:    "warehouse_id",
	ChangeType:     "change_type",
	BeforeQuantity: "before_quantity",
	AfterQuantity:  "after_quantity",
	BeforeReserved: "before_reserved",
	AfterReserved:  "after_reserved",
	ChangeAmount:   "change_amount",
	ReferenceId:    "reference_id",
	Operator:       "operator",
	Note:           "note",
	CreatedAt:      "created_at",
}

// NewInventoryLogsDao creates and returns a new DAO object for table data access.
func NewInventoryLogsDao() *InventoryLogsDao {
	return &InventoryLogsDao{
		group:   "default",
		table:   "sp_inventory_logs",
		columns: inventoryLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InventoryLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InventoryLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InventoryLogsDao) Columns() InventoryLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InventoryLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InventoryLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InventoryLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
