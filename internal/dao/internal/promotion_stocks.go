// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionStocksDao is the data access object for table mkt_promotion_stocks.
type PromotionStocksDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns PromotionStocksColumns // columns contains all the column names of Table for convenient usage.
}

// PromotionStocksColumns defines and stores column names for table mkt_promotion_stocks.
type PromotionStocksColumns struct {
	Id             string //
	PromotionId    string // 促销ID
	SkuId          string // SKU ID（秒杀专用，通用活动可为空）
	TotalStock     string // 总库存
	AvailableStock string // 可用库存
	LockedStock    string // 锁定库存（下单未付）
	Version        string // 乐观锁版本号
	CreatedAt      string //
	UpdatedAt      string //
}

// promotionStocksColumns holds the columns for table mkt_promotion_stocks.
var promotionStocksColumns = PromotionStocksColumns{
	Id:             "id",
	PromotionId:    "promotion_id",
	SkuId:          "sku_id",
	TotalStock:     "total_stock",
	AvailableStock: "available_stock",
	LockedStock:    "locked_stock",
	Version:        "version",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewPromotionStocksDao creates and returns a new DAO object for table data access.
func NewPromotionStocksDao() *PromotionStocksDao {
	return &PromotionStocksDao{
		group:   "default",
		table:   "mkt_promotion_stocks",
		columns: promotionStocksColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PromotionStocksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PromotionStocksDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PromotionStocksDao) Columns() PromotionStocksColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PromotionStocksDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PromotionStocksDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PromotionStocksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
