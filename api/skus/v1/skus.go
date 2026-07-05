package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type SkusListReq struct {
	g.Meta `path:"/skus" tags:"Skus" method:"get" summary:"SKU列表"`

	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	ProductId int64 `json:"product_id"`
}
type SkusListRes struct {
	List  []*entity.Skus `json:"list"`
	Total int            `json:"total"`
}

type SkusDetailReq struct {
	g.Meta `path:"/skus/{id}" tags:"Skus" method:"get" summary:"SKU详情"`
	Id     int64 `json:"id"`
}
type SkusDetailRes struct {
	*entity.Skus
	AvailableQuantity int64  `json:"available_quantity"`
	InventoryStatus   string `json:"inventory_status,omitempty"`
}

type SkusGetByCodeReq struct {
	g.Meta `path:"/skus/code/{sku_code}" tags:"Skus" method:"get" summary:"根据编码查询SKU"`
	SkuCode string `json:"sku_code"`
}
type SkusGetByCodeRes struct {
	*entity.Skus
	AvailableQuantity int64  `json:"available_quantity"`
	InventoryStatus   string `json:"inventory_status,omitempty"`
}

type SkusCreateReq struct {
	g.Meta `path:"/skus" tags:"Skus" method:"post" summary:"新增SKU"`

	ProductId    int64   `json:"product_id"    v:"required"            description:"商品ID"`
	SkuCode      string  `json:"sku_code"      v:"required|length:1,100" description:"商家编码"`
	Barcode      string  `json:"barcode"       description:"条码"`
	Spec         string  `json:"spec"          v:"required"            description:"规格JSON"`
	Price        int64   `json:"price"         v:"required|min:1"      description:"销售价(分)"`
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
type SkusCreateRes struct {
	Id int64 `json:"id"`
}

type SkusUpdateReq struct {
	g.Meta `path:"/skus/{id}" tags:"Skus" method:"put" summary:"更新SKU"`

	Id             int64    `json:"id"              v:"required"`
	Price          *int64   `json:"price"           description:"销售价(分)"`
	MarketPrice    *int64   `json:"market_price"    description:"划线价(分)"`
	CostPrice      *int64   `json:"cost_price"      description:"成本价(分)"`
	Status         *int     `json:"status"          description:"1-正常 0-禁用"`
	Image          *string  `json:"image"           description:"SKU图"`
	Barcode        *string  `json:"barcode"         description:"条码"`
	Weight         *float64 `json:"weight"          description:"重量(克)"`
	Volume         *float64 `json:"volume"          description:"体积(cm³)"`
	Length         *float64 `json:"length"          description:"长(cm)"`
	Width          *float64 `json:"width"           description:"宽(cm)"`
	Height         *float64 `json:"height"          description:"高(cm)"`
	MinPurchaseQty *int     `json:"min_purchase_qty" description:"最少购买数量"`
	MaxPurchaseQty *int     `json:"max_purchase_qty" description:"最大购买数量"`
}
type SkusUpdateRes struct{}

type SkusDeleteReq struct {
	g.Meta `path:"/skus/{id}" tags:"Skus" method:"delete" summary:"删除SKU"`
	Id     int64 `json:"id"`
}
type SkusDeleteRes struct{}
