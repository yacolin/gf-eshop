// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeliveryItemsDao is the data access object for table tx_delivery_items.
type DeliveryItemsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns DeliveryItemsColumns // columns contains all the column names of Table for convenient usage.
}

// DeliveryItemsColumns defines and stores column names for table tx_delivery_items.
type DeliveryItemsColumns struct {
	Id          string // 主键
	DeliveryId  string // 关联 tx_deliveries.id
	OrderItemId string // 关联 tx_order_items.id
	SkuId       string // SKU ID（冗余）
	ProductName string // 商品名（冗余快照）
	Quantity    string // 本次发货数量
	CreatedAt   string //
}

// deliveryItemsColumns holds the columns for table tx_delivery_items.
var deliveryItemsColumns = DeliveryItemsColumns{
	Id:          "id",
	DeliveryId:  "delivery_id",
	OrderItemId: "order_item_id",
	SkuId:       "sku_id",
	ProductName: "product_name",
	Quantity:    "quantity",
	CreatedAt:   "created_at",
}

// NewDeliveryItemsDao creates and returns a new DAO object for table data access.
func NewDeliveryItemsDao() *DeliveryItemsDao {
	return &DeliveryItemsDao{
		group:   "default",
		table:   "tx_delivery_items",
		columns: deliveryItemsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DeliveryItemsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DeliveryItemsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DeliveryItemsDao) Columns() DeliveryItemsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DeliveryItemsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DeliveryItemsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DeliveryItemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
