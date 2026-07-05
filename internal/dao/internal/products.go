// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductsDao is the data access object for table sp_products.
type ProductsDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns ProductsColumns // columns contains all the column names of Table for convenient usage.
}

// ProductsColumns defines and stores column names for table sp_products.
type ProductsColumns struct {
	Id             string // SPU ID
	MerchantId     string // 所属商家ID
	Name           string // 商品名称（用于搜索和展示）
	Subtitle       string // 商品副标题（卖点文案，如"2026新款"）
	CategoryId     string // 前台主类目ID（叶子节点）
	BrandId        string // 品牌ID
	Unit           string // 单位（件/箱/台/套）
	MainImage      string // 主图URL（CDN地址）
	Images         string // 附图JSON数组（最多10张）
	VideoUrl       string // 主图视频URL
	MinPrice       string // SKU最低价（分）
	MaxPrice       string // SKU最高价（分）
	TotalStock     string // 可售库存总和（SUM(quantity - reserved)）
	SalesCount     string // 总销量（从订单明细聚合，每日更新）
	RatingAverage  string // 平均评分（1-5）
	RatingCount    string // 评价总数
	Status         string // 0-草稿 1-待审 2-已上架 3-已下架 4-违规封禁
	SortOrder      string // 排序权重（越大越靠前，运营手动调整）
	HasDescription string // 1-有图文详情（存于 sp_product_descriptions 表）
	SeoTitle       string // SEO标题（自定义title，留空则使用name）
	SeoKeywords    string // SEO关键词（逗号分隔）
	SeoDescription string // SEO描述（页面meta description）
	CreatedBy      string // 创建人（运营工号）
	UpdatedBy      string // 最后更新人
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 软删除时间（NULL表示未删除）
}

// productsColumns holds the columns for table sp_products.
var productsColumns = ProductsColumns{
	Id:             "id",
	MerchantId:     "merchant_id",
	Name:           "name",
	Subtitle:       "subtitle",
	CategoryId:     "category_id",
	BrandId:        "brand_id",
	Unit:           "unit",
	MainImage:      "main_image",
	Images:         "images",
	VideoUrl:       "video_url",
	MinPrice:       "min_price",
	MaxPrice:       "max_price",
	TotalStock:     "total_stock",
	SalesCount:     "sales_count",
	RatingAverage:  "rating_average",
	RatingCount:    "rating_count",
	Status:         "status",
	SortOrder:      "sort_order",
	HasDescription: "has_description",
	SeoTitle:       "seo_title",
	SeoKeywords:    "seo_keywords",
	SeoDescription: "seo_description",
	CreatedBy:      "created_by",
	UpdatedBy:      "updated_by",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewProductsDao creates and returns a new DAO object for table data access.
func NewProductsDao() *ProductsDao {
	return &ProductsDao{
		group:   "default",
		table:   "sp_products",
		columns: productsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductsDao) Columns() ProductsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
