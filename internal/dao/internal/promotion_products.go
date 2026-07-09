// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionProductsDao is the data access object for table mkt_promotion_products.
type PromotionProductsDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns PromotionProductsColumns // columns contains all the column names of Table for convenient usage.
}

// PromotionProductsColumns defines and stores column names for table mkt_promotion_products.
type PromotionProductsColumns struct {
	Id          string // 主键ID
	PromotionId string // 促销ID
	MerchantId  string // 所属商家ID
	ProductType string // 1-全站 2-指定分类 3-指定SPU 4-指定SKU
	TargetId    string // 目标ID（SPU_ID或SKU_ID或Category_ID）
	CreatedAt   string //
	DeletedAt   string //
}

// promotionProductsColumns holds the columns for table mkt_promotion_products.
var promotionProductsColumns = PromotionProductsColumns{
	Id:          "id",
	PromotionId: "promotion_id",
	MerchantId:  "merchant_id",
	ProductType: "product_type",
	TargetId:    "target_id",
	CreatedAt:   "created_at",
	DeletedAt:   "deleted_at",
}

// NewPromotionProductsDao creates and returns a new DAO object for table data access.
func NewPromotionProductsDao() *PromotionProductsDao {
	return &PromotionProductsDao{
		group:   "default",
		table:   "mkt_promotion_products",
		columns: promotionProductsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PromotionProductsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PromotionProductsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PromotionProductsDao) Columns() PromotionProductsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PromotionProductsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PromotionProductsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PromotionProductsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
