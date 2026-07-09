// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionUsageLogsDao is the data access object for table mkt_promotion_usage_logs.
type PromotionUsageLogsDao struct {
	table   string                    // table is the underlying table name of the DAO.
	group   string                    // group is the database configuration group name of current DAO.
	columns PromotionUsageLogsColumns // columns contains all the column names of Table for convenient usage.
}

// PromotionUsageLogsColumns defines and stores column names for table mkt_promotion_usage_logs.
type PromotionUsageLogsColumns struct {
	Id                string // 主键ID
	PromotionId       string // 促销ID
	UserPromotionId   string // 用户促销资产ID
	UserId            string // 用户ID
	OrderId           string // 使用的订单ID
	DiscountAmount    string // 优惠金额（分）
	PromotionSnapshot string // 优惠快照（名称、规则等）
	CreatedAt         string //
}

// promotionUsageLogsColumns holds the columns for table mkt_promotion_usage_logs.
var promotionUsageLogsColumns = PromotionUsageLogsColumns{
	Id:                "id",
	PromotionId:       "promotion_id",
	UserPromotionId:   "user_promotion_id",
	UserId:            "user_id",
	OrderId:           "order_id",
	DiscountAmount:    "discount_amount",
	PromotionSnapshot: "promotion_snapshot",
	CreatedAt:         "created_at",
}

// NewPromotionUsageLogsDao creates and returns a new DAO object for table data access.
func NewPromotionUsageLogsDao() *PromotionUsageLogsDao {
	return &PromotionUsageLogsDao{
		group:   "default",
		table:   "mkt_promotion_usage_logs",
		columns: promotionUsageLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PromotionUsageLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PromotionUsageLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PromotionUsageLogsDao) Columns() PromotionUsageLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PromotionUsageLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PromotionUsageLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PromotionUsageLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
