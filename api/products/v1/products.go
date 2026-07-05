package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// CreateFull 嵌套子结构体
type CreateSKUItem struct {
	SkuCode      string  `json:"sku_code"      description:"商家编码"`
	Barcode      string  `json:"barcode"       description:"条码"`
	Spec         string  `json:"spec"          description:"规格JSON"`
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
	AttributeId int64  `json:"attribute_id" description:"属性ID"`
	Value       string `json:"value"        description:"属性值"`
}

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

type ProductsListReq struct {
	g.Meta `path:"/products" tags:"Products" method:"get" summary:"商品列表(游标分页)"`

	Size       int    `json:"size"        description:"每页条数(默认10,最大100)"`
	Cursor     string `json:"cursor"      description:"游标(首次不传,后续使用返回的cursor)"`
	Name       string `json:"name"`
	CategoryId int64  `json:"category_id"`
	BrandId    int64  `json:"brand_id"`
	Status     int    `json:"status"`
	PriceMin   int64  `json:"price_min"`
	PriceMax   int64  `json:"price_max"`
}
type ProductsListRes struct {
	List    []*entity.Products `json:"list"`
	Cursor  string             `json:"cursor"`
	HasMore bool               `json:"has_more"`
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

type ProductsDetailRes struct {
	*entity.Products
	Attributes  []ProductAttrDetailResponse `json:"attributes"`
	Description *entity.ProductDescriptions `json:"description,omitempty"`
	SKUs        []*SkuDetailItem            `json:"skus"`
}

type ProductsDetailPureReq struct {
	g.Meta `path:"/products/pure/{id}" tags:"Products" method:"get" summary:"商品详情(纯实体,不含聚合)"`
	Id     int64 `json:"id"`
}
type ProductsDetailPureRes struct {
	*entity.Products
}

type ProductsCreateReq struct {
	g.Meta `path:"/products" tags:"Products" method:"post" summary:"新增商品"`

	Name       string `json:"name"        v:"required|length:1,200" description:"商品名称"`
	Subtitle   string `json:"subtitle"    description:"副标题"`
	CategoryId int64  `json:"category_id" v:"required"             description:"类目ID"`
	BrandId    int64  `json:"brand_id"    description:"品牌ID"`
	Unit       string `json:"unit"        description:"单位"`
	MainImage  string `json:"main_image"  v:"required"             description:"主图"`
	Images     string `json:"images"      description:"附图JSON"`
	VideoUrl   string `json:"video_url"   description:"视频URL"`
	SortOrder  int    `json:"sort_order"  description:"排序权重"`
	Status     int    `json:"status"      description:"状态"`
	CreatedBy  string `json:"created_by"  description:"创建人"`
}
type ProductsCreateRes struct {
	Id int64 `json:"id"`
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
