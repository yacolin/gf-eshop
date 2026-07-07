// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsDao is the data access object for table usr_points.
type PointsDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns PointsColumns // columns contains all the column names of Table for convenient usage.
}

// PointsColumns defines and stores column names for table usr_points.
type PointsColumns struct {
	Id           string // 流水ID
	UserId       string // 用户ID
	Points       string // 积分变动（正=增加，负=扣减）
	BalanceAfter string // 变动后积分余额
	Source       string // 积分来源：order-下单消费 review-评价 signin-签到 admin-管理员调整 refund-退款扣减 expire-过期清零
	SourceId     string // 来源ID（如订单号、评价ID）
	ExpireAt     string // 过期时间（NULL=永不过期）
	Status       string // 0-待确认 1-已确认 2-已过期 3-已作废
	Remark       string // 备注
	CreatedAt    string //
}

// pointsColumns holds the columns for table usr_points.
var pointsColumns = PointsColumns{
	Id:           "id",
	UserId:       "user_id",
	Points:       "points",
	BalanceAfter: "balance_after",
	Source:       "source",
	SourceId:     "source_id",
	ExpireAt:     "expire_at",
	Status:       "status",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewPointsDao creates and returns a new DAO object for table data access.
func NewPointsDao() *PointsDao {
	return &PointsDao{
		group:   "default",
		table:   "usr_points",
		columns: pointsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PointsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PointsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PointsDao) Columns() PointsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PointsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PointsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PointsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
