// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Skus is the golang structure of table sp_skus for DAO operations like Where/Data.
type Skus struct {
	g.Meta         `orm:"table:sp_skus, do:true"`
	Id             interface{} // SKU ID
	ProductId      interface{} // 关联 products.id
	MerchantId     interface{} // 所属商家ID
	SkuCode        interface{} // 商家编码（唯一，用于ERP/WMS对接）
	Barcode        interface{} // 条码/EAN/UPC（仓库扫描用，NULL表示无条码）
	Spec           interface{} // 规格JSON（如{"颜色":"红色","内存":"256G"}）
	SpecSignature  interface{} // 规格MD5签名（用于快速匹配，由应用层计算）
	Price          interface{} // 销售价（分）
	MarketPrice    interface{} // 划线价/市场价（分）
	CostPrice      interface{} // 成本价（分，仅后台可见）
	Weight         interface{} // 重量（克）
	Volume         interface{} // 体积（立方厘米）
	Length         interface{} // 长（厘米）
	Width          interface{} // 宽（厘米）
	Height         interface{} // 高（厘米）
	MinPurchaseQty interface{} // 最少购买数量
	MaxPurchaseQty interface{} // 最大购买数量（0=不限）
	Image          interface{} // SKU专属图（如不同颜色展示不同图片）
	Status         interface{} // 1-正常 0-禁用（如某规格暂时缺货下架）
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
