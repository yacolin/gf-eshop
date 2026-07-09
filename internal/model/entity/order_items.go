// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItems is the golang structure for table order_items.
type OrderItems struct {
	Id             int64       `json:"id"               description:"订单项ID"`
	OrderId        int64       `json:"order_id"         description:"关联 tx_orders.id"`
	SubOrderId     int64       `json:"sub_order_id"     description:"子订单ID"`
	MerchantId     int64       `json:"merchant_id"      description:"所属商家ID"`
	OrderNo        string      `json:"order_no"         description:"订单号（冗余，方便按订单号查）"`
	SubOrderNo     string      `json:"sub_order_no"     description:"子订单号（冗余，方便按商家订单查）"`
	SkuId          int64       `json:"sku_id"           description:"关联 skus.id(sp_skus)"`
	ProductId      int64       `json:"product_id"       description:"关联 products.id(sp_products)"`
	SkuCode        string      `json:"sku_code"         description:"商家编码（冗余快照）"`
	ProductName    string      `json:"product_name"     description:"商品名（冗余快照）"`
	SkuSpecSummary string      `json:"sku_spec_summary" description:"规格摘要（如：红色 / 256G）"`
	SkuSpec        string      `json:"sku_spec"         description:"规格JSON快照"`
	Image          string      `json:"image"            description:"商品图（冗余快照）"`
	Price          int64       `json:"price"            description:"单价（分，下单时价格）"`
	Quantity       int         `json:"quantity"         description:"购买数量"`
	Subtotal       int64       `json:"subtotal"         description:"小计（分 = price * quantity）"`
	RefundStatus   string      `json:"refund_status"    description:"退款状态：none-无 refunding-退款中 refunded-已退款"`
	RefundAmount   int64       `json:"refund_amount"    description:"已退款金额（分）"`
	CreatedAt      *gtime.Time `json:"created_at"       description:""`
	UpdatedAt      *gtime.Time `json:"updated_at"       description:""`
	DeletedAt      *gtime.Time `json:"deleted_at"       description:""`
}
