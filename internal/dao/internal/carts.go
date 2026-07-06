// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartsDao is the data access object for table tx_carts.
type CartsDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns CartsColumns // columns contains all the column names of Table for convenient usage.
}

// CartsColumns defines and stores column names for table tx_carts.
type CartsColumns struct {
	Id          string // 购物车ID
	UserId      string // 用户ID（已登录用户）
	SessionId   string // 会话ID（未登录时的临时标识）
	ItemCount   string // 商品种类数
	TotalAmount string // 总金额（分，聚合，减少查SKU次数）
	ExpiredAt   string // 过期时间（session 型购物车自动清理）
	CreatedAt   string //
	UpdatedAt   string //
}

// cartsColumns holds the columns for table tx_carts.
var cartsColumns = CartsColumns{
	Id:          "id",
	UserId:      "user_id",
	SessionId:   "session_id",
	ItemCount:   "item_count",
	TotalAmount: "total_amount",
	ExpiredAt:   "expired_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewCartsDao creates and returns a new DAO object for table data access.
func NewCartsDao() *CartsDao {
	return &CartsDao{
		group:   "default",
		table:   "tx_carts",
		columns: cartsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CartsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CartsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CartsDao) Columns() CartsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CartsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CartsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CartsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
