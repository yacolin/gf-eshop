// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionRulesDao is the data access object for table mkt_promotion_rules.
type PromotionRulesDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns PromotionRulesColumns // columns contains all the column names of Table for convenient usage.
}

// PromotionRulesColumns defines and stores column names for table mkt_promotion_rules.
type PromotionRulesColumns struct {
	Id             string // 规则ID
	PromotionId    string // 所属促销ID
	MerchantId     string // 所属商家ID（0表示平台级规则）
	RuleName       string // 规则名称（便于理解）
	ConditionType  string // 1-无门槛 2-满金额 3-满件数 4-指定用户等级
	ConditionValue string // 门槛值（分，满20000则存20000）
	BenefitType    string // 1-减固定金额 2-打折扣 3-赠品 4-免运费 5-送积分
	BenefitValue   string // 优惠值（减固定金额填分如3000；打折扣填千分比如800=8折）
	IsStackable    string // 是否可与其他促销叠加 0-否 1-是
	StackPriority  string // 叠加优先级（数字越小越优先计算）
	CreatedBy      string // 创建人
	UpdatedBy      string // 更新人
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 软删除时间
}

// promotionRulesColumns holds the columns for table mkt_promotion_rules.
var promotionRulesColumns = PromotionRulesColumns{
	Id:             "id",
	PromotionId:    "promotion_id",
	MerchantId:     "merchant_id",
	RuleName:       "rule_name",
	ConditionType:  "condition_type",
	ConditionValue: "condition_value",
	BenefitType:    "benefit_type",
	BenefitValue:   "benefit_value",
	IsStackable:    "is_stackable",
	StackPriority:  "stack_priority",
	CreatedBy:      "created_by",
	UpdatedBy:      "updated_by",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewPromotionRulesDao creates and returns a new DAO object for table data access.
func NewPromotionRulesDao() *PromotionRulesDao {
	return &PromotionRulesDao{
		group:   "default",
		table:   "mkt_promotion_rules",
		columns: promotionRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PromotionRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PromotionRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PromotionRulesDao) Columns() PromotionRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PromotionRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PromotionRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PromotionRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
