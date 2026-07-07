// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsRulesDao is the data access object for table usr_points_rules.
type PointsRulesDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PointsRulesColumns // columns contains all the column names of Table for convenient usage.
}

// PointsRulesColumns defines and stores column names for table usr_points_rules.
type PointsRulesColumns struct {
	Id          string //
	Name        string // 规则名称
	RuleKey     string // 规则键名：earn_rate-消费返积分比例 expire_days-积分过期天数 signin_points-签到奖励积分 review_points-评价奖励积分
	RuleValue   string // 规则值
	Description string // 规则说明
	SortOrder   string // 排序
	Status      string // 状态：0-禁用 1-启用
	CreatedAt   string //
	UpdatedAt   string //
}

// pointsRulesColumns holds the columns for table usr_points_rules.
var pointsRulesColumns = PointsRulesColumns{
	Id:          "id",
	Name:        "name",
	RuleKey:     "rule_key",
	RuleValue:   "rule_value",
	Description: "description",
	SortOrder:   "sort_order",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewPointsRulesDao creates and returns a new DAO object for table data access.
func NewPointsRulesDao() *PointsRulesDao {
	return &PointsRulesDao{
		group:   "default",
		table:   "usr_points_rules",
		columns: pointsRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PointsRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PointsRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PointsRulesDao) Columns() PointsRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PointsRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PointsRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PointsRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
