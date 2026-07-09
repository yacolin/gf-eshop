// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionsDao is the data access object for table mkt_promotions.
type PromotionsDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns PromotionsColumns // columns contains all the column names of Table for convenient usage.
}

// PromotionsColumns defines and stores column names for table mkt_promotions.
type PromotionsColumns struct {
	Id            string // 促销ID
	PromotionNo   string // 促销业务编号
	MerchantId    string // 所属商家ID（0表示平台级活动）
	PromoName     string // 活动名称
	PromoType     string // 1-满减券 2-折扣券 3-秒杀 4-满额减 5-满件折 6-会员价
	PromoCode     string // 优惠码（优惠券专用）
	StartTime     string // 开始时间
	EndTime       string // 结束时间
	TotalQuantity string // 发行总量（0表示不限）
	PerUserLimit  string // 每人限领/限购数量
	UsedQuantity  string // 已使用/已售数量（异步统计，非实时）
	RuleId        string // 关联规则表ID
	Status        string // 1-草稿 2-待生效 3-生效中 4-已暂停 5-已结束 6-已作废
	Priority      string // 优先级（数字越大越优先，同类型互斥）
	CreatedBy     string // 创建人
	UpdatedBy     string // 更新人
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// promotionsColumns holds the columns for table mkt_promotions.
var promotionsColumns = PromotionsColumns{
	Id:            "id",
	PromotionNo:   "promotion_no",
	MerchantId:    "merchant_id",
	PromoName:     "promo_name",
	PromoType:     "promo_type",
	PromoCode:     "promo_code",
	StartTime:     "start_time",
	EndTime:       "end_time",
	TotalQuantity: "total_quantity",
	PerUserLimit:  "per_user_limit",
	UsedQuantity:  "used_quantity",
	RuleId:        "rule_id",
	Status:        "status",
	Priority:      "priority",
	CreatedBy:     "created_by",
	UpdatedBy:     "updated_by",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewPromotionsDao creates and returns a new DAO object for table data access.
func NewPromotionsDao() *PromotionsDao {
	return &PromotionsDao{
		group:   "default",
		table:   "mkt_promotions",
		columns: promotionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PromotionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PromotionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PromotionsDao) Columns() PromotionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PromotionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PromotionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PromotionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
