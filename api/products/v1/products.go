package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ──────────────────────────────────────────────
// 公共嵌套结构体
// ──────────────────────────────────────────────

type CreateSKUItem struct {
	SkuCode      string  `json:"sku_code"      description:"商家编码"`
	Barcode      string  `json:"barcode"       description:"条码"`
	SpecSummary  string  `json:"spec_summary"  description:"规格文本快照（如：红色 / 256G）"`
	Price        int64   `json:"price"         description:"销售价(分)"`
	MarketPrice  int64   `json:"market_price"  description:"划线价(分)"`
	CostPrice    int64   `json:"cost_price"    description:"成本价(分)"`
	Weight       float64 `json:"weight"        description:"重量(克)"`
	Volume       float64 `json:"volume"        description:"体积(cm³)"`
	Length       float64 `json:"length"        description:"长(cm)"`
	Width        float64 `json:"width"         description:"宽(cm)"`
	Height       float64 `json:"height"        description:"高(cm)"`
	MinPurchaseQty int   `json:"min_purchase_qty" description:"最少购买数量"`
	MaxPurchaseQty int   `json:"max_purchase_qty" description:"最大购买数量"`
	Image        string  `json:"image"         description:"SKU图"`
}

type CreateProductAttrItem struct {
	AttributeId      int64  `json:"attribute_id"       description:"属性ID"`
	AttributeValueId int64  `json:"attribute_value_id"  description:"引用属性值字典ID（可选），优先使用"`
	Value            string `json:"value"              description:"属性值"`
}

// ──────────────────────────────────────────────
// 批量创建 SKU（Phase 2 — 笛卡尔积）
// ──────────────────────────────────────────────

type SpecGroupItem struct {
	AttributeId int64   `json:"attribute_id" v:"required" description:"销售属性ID（is_sku_spec=1）"`
	ValueIds    []int64 `json:"value_ids"    v:"required" description:"选中的属性值ID列表"`
}

type BatchCreateSkuItem struct {
	Id          int64  `json:"id"           description:"SKU ID"`
	SkuCode     string `json:"sku_code"     description:"商家编码"`
	SpecSummary string `json:"spec_summary" description:"规格文本快照"`
	Price       int64  `json:"price"        description:"销售价(分)"`
}

type ProductsBatchCreateSKUsReq struct {
	g.Meta `path:"/products/{product_id}/skus/batch" tags:"Products" method:"post" summary:"批量生成SKU（笛卡尔积）"`

	ProductId      int64           `json:"product_id"       v:"required" description:"商品ID"`
	SpecGroups     []SpecGroupItem `json:"spec_groups"      v:"required" description:"规格分组，系统自动计算笛卡尔积"`
	BasePrice      int64           `json:"base_price"       description:"基准价(分)，各SKU统一使用此价格"`
	SkuCodePrefix  string          `json:"sku_code_prefix"  description:"SKU编码前缀，为空则使用SKU-{product_id}-"`
}
type ProductsBatchCreateSKUsRes struct {
	Total int                   `json:"total" description:"生成SKU数量"`
	SKUs  []*BatchCreateSkuItem `json:"skus"`
}

// ──────────────────────────────────────────────
// 创建商品（Phase 1 — 商品概况 + 非销售属性）
// ──────────────────────────────────────────────

type ProductsCreateReq struct {
	g.Meta `path:"/products" tags:"Products" method:"post" summary:"新增商品（含非销售属性绑定）"`

	Name       string                 `json:"name"        v:"required|length:1,200" description:"商品名称"`
	Subtitle   string                 `json:"subtitle"    description:"副标题"`
	CategoryId int64                  `json:"category_id" v:"required"             description:"类目ID"`
	BrandId    int64                  `json:"brand_id"    description:"品牌ID"`
	Unit       string                 `json:"unit"        description:"单位"`
	MainImage  string                 `json:"main_image"  v:"required"             description:"主图"`
	Images     string                 `json:"images"      description:"附图JSON"`
	VideoUrl   string                 `json:"video_url"   description:"视频URL"`
	SortOrder  int                    `json:"sort_order"  description:"排序权重"`
	Status     int                    `json:"status"      description:"状态"`
	CreatedBy  string                 `json:"created_by"  description:"创建人"`
	Attributes []CreateProductAttrItem `json:"attributes" description:"非销售属性值列表（仅 is_sku_spec=0 的属性）"`
}
type ProductsCreateRes struct {
	Id int64 `json:"id"`
}

// ──────────────────────────────────────────────
// 创建商品（Phase 1+2 全量 — 兼容遗留）
// ──────────────────────────────────────────────

type ProductsCreateFullReq struct {
	g.Meta `path:"/products/full" tags:"Products" method:"post" summary:"创建商品（含SKU/属性/描述）"`

	Name        string                 `json:"name"        v:"required|length:1,200" description:"商品名称"`
	Subtitle    string                 `json:"subtitle"    description:"副标题"`
	CategoryId  int64                  `json:"category_id" v:"required"              description:"类目ID"`
	BrandId     int64                  `json:"brand_id"    description:"品牌ID"`
	Unit        string                 `json:"unit"        description:"单位"`
	MainImage   string                 `json:"main_image"  v:"required"              description:"主图"`
	Images      string                 `json:"images"      description:"附图JSON"`
	VideoUrl    string                 `json:"video_url"   description:"视频URL"`
	SortOrder   int                    `json:"sort_order"  description:"排序权重"`
	CreatedBy   string                 `json:"created_by"  description:"创建人"`
	Description string                 `json:"description" description:"商品详情HTML"`
	MobileDesc  string                 `json:"mobile_description" description:"移动端详情"`
	SKUs        []CreateSKUItem        `json:"skus"        description:"SKU列表"`
	Attributes  []CreateProductAttrItem `json:"attributes" description:"属性值列表"`
}
type ProductsCreateFullRes struct {
	Id int64 `json:"id"`
}

// ──────────────────────────────────────────────
// 列表、详情、更新、删除
// ──────────────────────────────────────────────

type ProductsListReq struct {
	g.Meta `path:"/products" tags:"Products" method:"get" summary:"商品列表(游标分页)"`

	Size       int    `json:"size"        description:"每页条数(默认10,最大100)"`
	Cursor     string `json:"cursor"      description:"游标(首次不传,后续使用返回的cursor)"`
	Name       string `json:"name"`
	CategoryId int64  `json:"category_id"`
	BrandId    int64  `json:"brand_id"`
	Status     *int   `json:"status"`
	PriceMin   int64  `json:"price_min"`
	PriceMax   int64  `json:"price_max"`
}
type ProductsListItem struct {
	*entity.Products
	PriceMin   int64 `json:"price_min"   description:"最低销售价(分)"`
	PriceMax   int64 `json:"price_max"   description:"最高销售价(分)"`
	TotalStock int64 `json:"total_stock" description:"总库存（可售）"`
}

type ProductsListRes struct {
	List    []*ProductsListItem `json:"list"`
	Cursor  string              `json:"cursor"`
	HasMore bool                `json:"has_more"`
}

type ProductsDetailReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"get" summary:"商品详情"`
	Id     int64 `json:"id"`
}
type ProductAttrDetailResponse struct {
	AttributeId   int64    `json:"attribute_id"`
	AttributeName string   `json:"attribute_name"`
	Values        []string `json:"values"`
	SortOrder     int      `json:"sort_order"`
}

type SkuDetailItem struct {
	*entity.Skus
	AvailableQuantity int64  `json:"available_quantity"`
	InventoryStatus   string `json:"inventory_status,omitempty"`
}

type ProductSpecResponse struct {
	Selectable    []ProductAttrDetailResponse `json:"selectable"     description:"可选的规格（驱动SKU选择器，如颜色、内存）"`
	NonSelectable []ProductAttrDetailResponse `json:"non_selectable" description:"不可选的规格（仅展示，如处理器、屏幕尺寸）"`
}

type ProductsDetailRes struct {
	*entity.Products
	PriceMin    int64                        `json:"price_min"   description:"最低销售价(分)"`
	PriceMax    int64                        `json:"price_max"   description:"最高销售价(分)"`
	TotalStock  int64                        `json:"total_stock" description:"总库存（可售）"`
	Description *entity.ProductDescriptions  `json:"description,omitempty"`
	SKUs        []*SkuDetailItem             `json:"skus"`
	Specs       *ProductSpecResponse         `json:"specs"`
}

// ──────────────────────────────────────────────
// 产品属性管理（CRUD）
// ──────────────────────────────────────────────

type ProductsGetAttributesReq struct {
	g.Meta `path:"/products/{id}/attributes" tags:"Products" method:"get" summary:"获取商品绑定属性"`

	Id int64 `json:"id" v:"required"`
}
type ProductsAttributeItem struct {
	Id               int64  `json:"id"                 description:"关联ID"`
	AttributeId      int64  `json:"attribute_id"       description:"属性ID"`
	AttributeName    string `json:"attribute_name"     description:"属性名称"`
	AttributeValueId int64  `json:"attribute_value_id" description:"引用属性值字典ID"`
	Value            string `json:"value"              description:"属性值"`
}
type ProductsGetAttributesRes struct {
	List []*ProductsAttributeItem `json:"list"`
}

type ProductsUpdateAttributesReq struct {
	g.Meta `path:"/products/{id}/attributes" tags:"Products" method:"put" summary:"更新商品属性（全量替换）"`

	Id         int64                  `json:"id"         v:"required"`
	Attributes []CreateProductAttrItem `json:"attributes" v:"required" description:"新的属性值列表，全量替换"`
}
type ProductsUpdateAttributesRes struct{}

type ProductsEnrichedDetailReq struct {
	g.Meta `path:"/products/enriched/{id}" tags:"Products" method:"get" summary:"商品富化详情（含SKU/规格/库存/描述）"`
	Id     int64 `json:"id"`
}
type ProductsEnrichedDetailRes struct {
	*ProductsDetailRes
}

type ProductsDetailPureReq struct {
	g.Meta `path:"/products/pure/{id}" tags:"Products" method:"get" summary:"商品详情(纯实体,不含聚合)"`
	Id     int64 `json:"id"`
}
type ProductsDetailPureRes struct {
	*entity.Products
}

type ProductsUpdateReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"put" summary:"更新商品"`

	Id         int64  `json:"id"          v:"required"`
	Name       string `json:"name"        v:"length:1,200" description:"商品名称"`
	Subtitle   string `json:"subtitle"    description:"副标题"`
	CategoryId int64  `json:"category_id" description:"类目ID"`
	BrandId    int64  `json:"brand_id"    description:"品牌ID"`
	Unit       string `json:"unit"        description:"单位"`
	MainImage  string `json:"main_image"  description:"主图"`
	Images     string `json:"images"      description:"附图JSON"`
	VideoUrl   string `json:"video_url"   description:"视频URL"`
	SortOrder  int    `json:"sort_order"  description:"排序权重"`
	Status     int    `json:"status"      description:"状态"`
	UpdatedBy  string `json:"updated_by"  description:"更新人"`
}
type ProductsUpdateRes struct{}

type ProductsDeleteReq struct {
	g.Meta `path:"/products/{id}" tags:"Products" method:"delete" summary:"删除商品"`
	Id     int64 `json:"id"`
}
type ProductsDeleteRes struct{}
