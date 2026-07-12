// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Products is the golang structure for table products.
type Products struct {
	Id             int64       `json:"id"              description:"SPU ID"`
	MerchantId     int64       `json:"merchant_id"     description:"所属商家ID"`
	Name           string      `json:"name"            description:"商品名称（用于搜索和展示）"`
	Subtitle       string      `json:"subtitle"        description:"商品副标题（卖点文案，如\"2026新款\"）"`
	CategoryId     int64       `json:"category_id"     description:"前台主类目ID（叶子节点）"`
	BrandId        int64       `json:"brand_id"        description:"品牌ID"`
	Unit           string      `json:"unit"            description:"单位（件/箱/台/套）"`
	MainImage      string      `json:"main_image"      description:"主图URL（CDN地址）"`
	Images         string      `json:"images"          description:"附图JSON数组（最多10张）"`
	VideoUrl       string      `json:"video_url"       description:"主图视频URL"`
	SalesCount     int         `json:"sales_count"     description:"总销量（从订单明细聚合，每日更新）"`
	RatingAverage  float64     `json:"rating_average"  description:"平均评分（1-5）"`
	RatingCount    int         `json:"rating_count"    description:"评价总数"`
	Status         int         `json:"status"          description:"0-草稿 1-待审 2-已上架 3-已下架 4-违规封禁"`
	SortOrder      int         `json:"sort_order"      description:"排序权重（越大越靠前，运营手动调整）"`
	HasDescription int         `json:"has_description" description:"1-有图文详情（存于 sp_product_descriptions 表）"`
	SeoTitle       string      `json:"seo_title"       description:"SEO标题（自定义title，留空则使用name）"`
	SeoKeywords    string      `json:"seo_keywords"    description:"SEO关键词（逗号分隔）"`
	SeoDescription string      `json:"seo_description" description:"SEO描述（页面meta description）"`
	CreatedBy      int64       `json:"created_by"      description:"创建人ID"`
	UpdatedBy      int64       `json:"updated_by"      description:"最后更新人ID"`
	CreatedAt      *gtime.Time `json:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updated_at"      description:"更新时间"`
	DeletedAt      *gtime.Time `json:"deleted_at"      description:"软删除时间（NULL表示未删除）"`
}
