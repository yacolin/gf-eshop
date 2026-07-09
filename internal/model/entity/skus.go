// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Skus is the golang structure for table skus.
type Skus struct {
	Id             int64       `json:"id"               description:"SKU ID"`
	ProductId      int64       `json:"product_id"       description:"关联 products.id"`
	MerchantId     int64       `json:"merchant_id"      description:"所属商家ID"`
	SkuCode        string      `json:"sku_code"         description:"商家编码（唯一，用于ERP/WMS对接）"`
	Barcode        string      `json:"barcode"          description:"条码/EAN/UPC（仓库扫描用，NULL表示无条码）"`
	SpecSummary    string      `json:"spec_summary"     description:"规格文本快照（如：红色 / 256G）"`
	Price          int64       `json:"price"            description:"销售价（分）"`
	MarketPrice    int64       `json:"market_price"     description:"划线价/市场价（分）"`
	CostPrice      int64       `json:"cost_price"       description:"成本价（分，仅后台可见）"`
	Weight         float64     `json:"weight"           description:"重量（克）"`
	Volume         float64     `json:"volume"           description:"体积（立方厘米）"`
	Length         float64     `json:"length"           description:"长（厘米）"`
	Width          float64     `json:"width"            description:"宽（厘米）"`
	Height         float64     `json:"height"           description:"高（厘米）"`
	MinPurchaseQty int         `json:"min_purchase_qty" description:"最少购买数量"`
	MaxPurchaseQty int         `json:"max_purchase_qty" description:"最大购买数量（0=不限）"`
	Image          string      `json:"image"            description:"SKU专属图（如不同颜色展示不同图片）"`
	Status         int         `json:"status"           description:"1-正常 0-禁用（如某规格暂时缺货下架）"`
	CreatedAt      *gtime.Time `json:"created_at"       description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"       description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"       description:""`
}
