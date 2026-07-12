// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LevelRulesDao is the data access object for table usr_level_rules.
type LevelRulesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns LevelRulesColumns // columns contains all the column names of Table for convenient usage.
}

// LevelRulesColumns defines and stores column names for table usr_level_rules.
type LevelRulesColumns struct {
	Id             string //
	Name           string // 规则名称
	RuleType       string // 规则类型：upgrade-自动升级 downgrade-自动降级
	FromLevelId    string // 源等级ID（0=任意等级）
	ToLevelId      string // 目标等级ID
	ConditionType  string // 条件类型：points-累计积分 order_count-订单数 order_amount-消费金额
	ConditionValue string // 条件阈值
	Description    string // 规则说明
	SortOrder      string // 排序
	Status         string // 状态：0-禁用 1-启用
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// levelRulesColumns holds the columns for table usr_level_rules.
var levelRulesColumns = LevelRulesColumns{
	Id:             "id",
	Name:           "name",
	RuleType:       "rule_type",
	FromLevelId:    "from_level_id",
	ToLevelId:      "to_level_id",
	ConditionType:  "condition_type",
	ConditionValue: "condition_value",
	Description:    "description",
	SortOrder:      "sort_order",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewLevelRulesDao creates and returns a new DAO object for table data access.
func NewLevelRulesDao() *LevelRulesDao {
	return &LevelRulesDao{
		group:   "default",
		table:   "usr_level_rules",
		columns: levelRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *LevelRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *LevelRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *LevelRulesDao) Columns() LevelRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *LevelRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *LevelRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *LevelRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
