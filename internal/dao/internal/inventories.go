// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoriesDao is the data access object for table sp_inventories.
type InventoriesDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns InventoriesColumns // columns contains all the column names of Table for convenient usage.
}

// InventoriesColumns defines and stores column names for table sp_inventories.
type InventoriesColumns struct {
	Id            string // 库存记录ID
	SkuId         string // 关联 skus.id
	MerchantId    string // 所属商家ID
	WarehouseId   string // 仓库ID（关联 sp_warehouses.id）
	Quantity      string // 物理库存总量（含预占）
	Reserved      string // 预占库存（下单未支付）
	Available     string // 可售库存（虚拟列，无需持久化）
	InTransit     string // 在途库存（采购中/调拨中）
	Threshold     string // 安全库存预警阈值（低于此值触发告警）
	MaxThreshold  string // 最大库存上限（入库不能超过此值）
	Status        string // 1-充足 2-缺货 3-无货
	LastCountedAt string // 最后盘点时间
	LastCountedBy string // 最后盘点人
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// inventoriesColumns holds the columns for table sp_inventories.
var inventoriesColumns = InventoriesColumns{
	Id:            "id",
	SkuId:         "sku_id",
	MerchantId:    "merchant_id",
	WarehouseId:   "warehouse_id",
	Quantity:      "quantity",
	Reserved:      "reserved",
	Available:     "available",
	InTransit:     "in_transit",
	Threshold:     "threshold",
	MaxThreshold:  "max_threshold",
	Status:        "status",
	LastCountedAt: "last_counted_at",
	LastCountedBy: "last_counted_by",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewInventoriesDao creates and returns a new DAO object for table data access.
func NewInventoriesDao() *InventoriesDao {
	return &InventoriesDao{
		group:   "default",
		table:   "sp_inventories",
		columns: inventoriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InventoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InventoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InventoriesDao) Columns() InventoriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InventoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InventoriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InventoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
