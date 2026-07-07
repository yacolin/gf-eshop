// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LevelsDao is the data access object for table usr_levels.
type LevelsDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns LevelsColumns // columns contains all the column names of Table for convenient usage.
}

// LevelsColumns defines and stores column names for table usr_levels.
type LevelsColumns struct {
	Id               string // 等级ID
	Name             string // 等级名称（如：青铜会员、白银会员、黄金会员、钻石会员）
	Level            string // 等级数值（1=青铜 2=白银 3=黄金 4=钻石）
	MinPoints        string // 该等级所需最低累计积分
	MaxPoints        string // 该等级所需最高累计积分（0表示无上限）
	DiscountRate     string // 折扣率（千分比，1000=无折扣，900=九折）
	FreeShipping     string // 1-免运费
	PointsMultiplier string // 消费积分倍数（如 1.5 倍积分）
	Benefits         string // 扩展权益JSON（如：{"birthday_gift": true, "exclusive_coupon": true}）
	Status           string // 1-启用 0-禁用
	SortOrder        string // 排序
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// levelsColumns holds the columns for table usr_levels.
var levelsColumns = LevelsColumns{
	Id:               "id",
	Name:             "name",
	Level:            "level",
	MinPoints:        "min_points",
	MaxPoints:        "max_points",
	DiscountRate:     "discount_rate",
	FreeShipping:     "free_shipping",
	PointsMultiplier: "points_multiplier",
	Benefits:         "benefits",
	Status:           "status",
	SortOrder:        "sort_order",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewLevelsDao creates and returns a new DAO object for table data access.
func NewLevelsDao() *LevelsDao {
	return &LevelsDao{
		group:   "default",
		table:   "usr_levels",
		columns: levelsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *LevelsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *LevelsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *LevelsDao) Columns() LevelsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *LevelsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *LevelsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *LevelsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
