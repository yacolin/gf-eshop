// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartItemsDao is the data access object for table tx_cart_items.
type CartItemsDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns CartItemsColumns // columns contains all the column names of Table for convenient usage.
}

// CartItemsColumns defines and stores column names for table tx_cart_items.
type CartItemsColumns struct {
	Id          string // 购物车项ID
	CartId      string // 关联 tx_carts.id
	SkuId       string // 关联 skus.id(sp_skus)
	ProductId   string // 关联 products.id(sp_products)，冗余用于展示
	ProductName string // 商品名（冗余快照）
	SkuSpec     string // 规格JSON快照（如{"颜色":"红色","内存":"256G"}）
	Image       string // 商品图（冗余快照）
	Price       string // 加入时的价格（分，防止下单时价格变动导致纠纷）
	Quantity    string // 数量
	CreatedAt   string //
	UpdatedAt   string //
}

// cartItemsColumns holds the columns for table tx_cart_items.
var cartItemsColumns = CartItemsColumns{
	Id:          "id",
	CartId:      "cart_id",
	SkuId:       "sku_id",
	ProductId:   "product_id",
	ProductName: "product_name",
	SkuSpec:     "sku_spec",
	Image:       "image",
	Price:       "price",
	Quantity:    "quantity",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewCartItemsDao creates and returns a new DAO object for table data access.
func NewCartItemsDao() *CartItemsDao {
	return &CartItemsDao{
		group:   "default",
		table:   "tx_cart_items",
		columns: cartItemsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CartItemsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CartItemsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CartItemsDao) Columns() CartItemsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CartItemsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CartItemsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CartItemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
