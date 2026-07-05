// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SkusDao is the data access object for table sp_skus.
type SkusDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns SkusColumns // columns contains all the column names of Table for convenient usage.
}

// SkusColumns defines and stores column names for table sp_skus.
type SkusColumns struct {
	Id             string // SKU ID
	ProductId      string // 关联 products.id
	MerchantId     string // 所属商家ID
	SkuCode        string // 商家编码（唯一，用于ERP/WMS对接）
	Barcode        string // 条码/EAN/UPC（仓库扫描用，NULL表示无条码）
	Spec           string // 规格JSON（如{"颜色":"红色","内存":"256G"}）
	SpecSignature  string // 规格MD5签名（用于快速匹配，由应用层计算）
	Price          string // 销售价（分）
	MarketPrice    string // 划线价/市场价（分）
	CostPrice      string // 成本价（分，仅后台可见）
	Weight         string // 重量（克）
	Volume         string // 体积（立方厘米）
	Length         string // 长（厘米）
	Width          string // 宽（厘米）
	Height         string // 高（厘米）
	MinPurchaseQty string // 最少购买数量
	MaxPurchaseQty string // 最大购买数量（0=不限）
	Image          string // SKU专属图（如不同颜色展示不同图片）
	Status         string // 1-正常 0-禁用（如某规格暂时缺货下架）
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// skusColumns holds the columns for table sp_skus.
var skusColumns = SkusColumns{
	Id:             "id",
	ProductId:      "product_id",
	MerchantId:     "merchant_id",
	SkuCode:        "sku_code",
	Barcode:        "barcode",
	Spec:           "spec",
	SpecSignature:  "spec_signature",
	Price:          "price",
	MarketPrice:    "market_price",
	CostPrice:      "cost_price",
	Weight:         "weight",
	Volume:         "volume",
	Length:         "length",
	Width:          "width",
	Height:         "height",
	MinPurchaseQty: "min_purchase_qty",
	MaxPurchaseQty: "max_purchase_qty",
	Image:          "image",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewSkusDao creates and returns a new DAO object for table data access.
func NewSkusDao() *SkusDao {
	return &SkusDao{
		group:   "default",
		table:   "sp_skus",
		columns: skusColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SkusDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SkusDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SkusDao) Columns() SkusColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SkusDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SkusDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SkusDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
