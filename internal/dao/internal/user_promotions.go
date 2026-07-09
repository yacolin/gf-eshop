// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserPromotionsDao is the data access object for table mkt_user_promotions.
type UserPromotionsDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns UserPromotionsColumns // columns contains all the column names of Table for convenient usage.
}

// UserPromotionsColumns defines and stores column names for table mkt_user_promotions.
type UserPromotionsColumns struct {
	Id              string // 主键ID
	UserPromotionNo string // 用户促销资产编号
	UserId          string // 用户ID
	PromotionId     string // 促销ID
	MerchantId      string // 所属商家ID
	AcquireTime     string // 领取时间
	ExpireTime      string // 过期时间
	Status          string // 1-未使用 2-锁定中(下单未付) 3-已使用 4-已过期 5-已作废
	LockOrderId     string // 锁定的订单ID（用于回滚）
	UsedTime        string // 使用时间
	OrderId         string // 最终使用的订单ID
	QueueToken      string // 秒杀排队令牌
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// userPromotionsColumns holds the columns for table mkt_user_promotions.
var userPromotionsColumns = UserPromotionsColumns{
	Id:              "id",
	UserPromotionNo: "user_promotion_no",
	UserId:          "user_id",
	PromotionId:     "promotion_id",
	MerchantId:      "merchant_id",
	AcquireTime:     "acquire_time",
	ExpireTime:      "expire_time",
	Status:          "status",
	LockOrderId:     "lock_order_id",
	UsedTime:        "used_time",
	OrderId:         "order_id",
	QueueToken:      "queue_token",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewUserPromotionsDao creates and returns a new DAO object for table data access.
func NewUserPromotionsDao() *UserPromotionsDao {
	return &UserPromotionsDao{
		group:   "default",
		table:   "mkt_user_promotions",
		columns: userPromotionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserPromotionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserPromotionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserPromotionsDao) Columns() UserPromotionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserPromotionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserPromotionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserPromotionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
