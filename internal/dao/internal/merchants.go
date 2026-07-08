// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantsDao is the data access object for table mch_merchants.
type MerchantsDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns MerchantsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantsColumns defines and stores column names for table mch_merchants.
type MerchantsColumns struct {
	Id              string // 商家ID
	MerchantName    string // 商家名称（店铺名）
	MerchantCode    string // 商家编码（系统生成，唯一）
	MerchantType    string // 1-个人商家 2-企业商家 3-品牌直营
	MerchantLevel   string // 商家等级 1-普通 2-银牌 3-金牌 4-钻石（影响佣金率/权限）
	BusinessScope   string // 经营范围
	BusinessYears   string // 经营年限（入驻年限）
	ContactPerson   string // 主要联系人
	ContactPhone    string // 联系电话
	ContactEmail    string // 联系邮箱
	LogoUrl         string // 店铺Logo
	BannerUrl       string // 店铺Banner图
	ShopDescription string // 店铺简介
	Status          string // 0-待审核 1-正常 2-冻结 3-已注销
	AuditStatus     string // 0-待审核 1-审核通过 2-审核拒绝
	AuditReason     string // 审核拒绝原因
	AuditedAt       string // 审核时间
	FrozenReason    string // 冻结原因
	CommissionRate  string // 平台抽佣比例（千分比，如 50 表示5%）
	SettlementCycle string // 结算周期 1-T+1 2-T+7 3-月结
	TotalOrders     string // 历史总订单数
	TotalSales      string // 历史总销售额（分）
	AvgRating       string // 店铺平均评分
	ProductCount    string // 在售商品数量
	SettledAt       string // 入驻时间
	ExpireAt        string // 合同到期时间
	CreatedBy       string // 创建人
	UpdatedBy       string // 更新人
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// merchantsColumns holds the columns for table mch_merchants.
var merchantsColumns = MerchantsColumns{
	Id:              "id",
	MerchantName:    "merchant_name",
	MerchantCode:    "merchant_code",
	MerchantType:    "merchant_type",
	MerchantLevel:   "merchant_level",
	BusinessScope:   "business_scope",
	BusinessYears:   "business_years",
	ContactPerson:   "contact_person",
	ContactPhone:    "contact_phone",
	ContactEmail:    "contact_email",
	LogoUrl:         "logo_url",
	BannerUrl:       "banner_url",
	ShopDescription: "shop_description",
	Status:          "status",
	AuditStatus:     "audit_status",
	AuditReason:     "audit_reason",
	AuditedAt:       "audited_at",
	FrozenReason:    "frozen_reason",
	CommissionRate:  "commission_rate",
	SettlementCycle: "settlement_cycle",
	TotalOrders:     "total_orders",
	TotalSales:      "total_sales",
	AvgRating:       "avg_rating",
	ProductCount:    "product_count",
	SettledAt:       "settled_at",
	ExpireAt:        "expire_at",
	CreatedBy:       "created_by",
	UpdatedBy:       "updated_by",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewMerchantsDao creates and returns a new DAO object for table data access.
func NewMerchantsDao() *MerchantsDao {
	return &MerchantsDao{
		group:   "default",
		table:   "mch_merchants",
		columns: merchantsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantsDao) Columns() MerchantsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
