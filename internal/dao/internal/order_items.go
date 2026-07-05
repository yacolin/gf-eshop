// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderItemsDao is the data access object for table tx_order_items.
type OrderItemsDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns OrderItemsColumns // columns contains all the column names of Table for convenient usage.
}

// OrderItemsColumns defines and stores column names for table tx_order_items.
type OrderItemsColumns struct {
	Id           string // 订单项ID
	OrderId      string // 关联 tx_orders.id
	SubOrderId   string // 子订单ID
	MerchantId   string // 所属商家ID
	OrderNo      string // 订单号（冗余，方便按订单号查）
	SubOrderNo   string // 子订单号（冗余，方便按商家订单查）
	SkuId        string // 关联 skus.id(sp_skus)
	ProductId    string // 关联 products.id(sp_products)
	SkuCode      string // 商家编码（冗余快照）
	ProductName  string // 商品名（冗余快照）
	SkuSpec      string // 规格JSON快照
	Image        string // 商品图（冗余快照）
	Price        string // 单价（分，下单时价格）
	Quantity     string // 购买数量
	Subtotal     string // 小计（分 = price * quantity）
	RefundStatus string // 退款状态：none-无 refunding-退款中 refunded-已退款
	RefundAmount string // 已退款金额（分）
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
}

// orderItemsColumns holds the columns for table tx_order_items.
var orderItemsColumns = OrderItemsColumns{
	Id:           "id",
	OrderId:      "order_id",
	SubOrderId:   "sub_order_id",
	MerchantId:   "merchant_id",
	OrderNo:      "order_no",
	SubOrderNo:   "sub_order_no",
	SkuId:        "sku_id",
	ProductId:    "product_id",
	SkuCode:      "sku_code",
	ProductName:  "product_name",
	SkuSpec:      "sku_spec",
	Image:        "image",
	Price:        "price",
	Quantity:     "quantity",
	Subtotal:     "subtotal",
	RefundStatus: "refund_status",
	RefundAmount: "refund_amount",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewOrderItemsDao creates and returns a new DAO object for table data access.
func NewOrderItemsDao() *OrderItemsDao {
	return &OrderItemsDao{
		group:   "default",
		table:   "tx_order_items",
		columns: orderItemsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OrderItemsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OrderItemsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OrderItemsDao) Columns() OrderItemsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OrderItemsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OrderItemsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OrderItemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
