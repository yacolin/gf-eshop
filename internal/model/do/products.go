// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Products is the golang structure of table sp_products for DAO operations like Where/Data.
type Products struct {
	g.Meta         `orm:"table:sp_products, do:true"`
	Id             interface{} // SPU ID
	MerchantId     interface{} // 所属商家ID
	Name           interface{} // 商品名称（用于搜索和展示）
	Subtitle       interface{} // 商品副标题（卖点文案，如"2026新款"）
	CategoryId     interface{} // 前台主类目ID（叶子节点）
	BrandId        interface{} // 品牌ID
	Unit           interface{} // 单位（件/箱/台/套）
	MainImage      interface{} // 主图URL（CDN地址）
	Images         interface{} // 附图JSON数组（最多10张）
	VideoUrl       interface{} // 主图视频URL
	SalesCount     interface{} // 总销量（从订单明细聚合，每日更新）
	RatingAverage  interface{} // 平均评分（1-5）
	RatingCount    interface{} // 评价总数
	Status         interface{} // 0-草稿 1-待审 2-已上架 3-已下架 4-违规封禁
	SortOrder      interface{} // 排序权重（越大越靠前，运营手动调整）
	HasDescription interface{} // 1-有图文详情（存于 sp_product_descriptions 表）
	SeoTitle       interface{} // SEO标题（自定义title，留空则使用name）
	SeoKeywords    interface{} // SEO关键词（逗号分隔）
	SeoDescription interface{} // SEO描述（页面meta description）
	CreatedBy      interface{} // 创建人（运营工号）
	UpdatedBy      interface{} // 最后更新人
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 软删除时间（NULL表示未删除）
}
