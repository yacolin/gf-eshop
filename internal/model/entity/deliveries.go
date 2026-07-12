// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Deliveries is the golang structure for table deliveries.
type Deliveries struct {
	Id           int64       `json:"id"            description:"物流单ID"`
	DeliveryNo   string      `json:"delivery_no"   description:"物流单号（业务唯一）"`
	OrderId      int64       `json:"order_id"      description:"关联 tx_orders.id"`
	OrderNo      string      `json:"order_no"      description:"订单号（冗余）"`
	MerchantId   int64       `json:"merchant_id"   description:"所属商家ID"`
	Carrier      string      `json:"carrier"       description:"物流商：sf-顺丰 yto-圆通 zto-中通 yunda-韵达 jd-京东物流 other-其他"`
	TrackingNo   string      `json:"tracking_no"   description:"运单号（物流商单号）"`
	WarehouseId  int64       `json:"warehouse_id"  description:"发货仓库ID"`
	PackageCount int         `json:"package_count" description:"包裹数量"`
	Consignee    string      `json:"consignee"     description:"收货人"`
	Phone        string      `json:"phone"         description:"联系电话"`
	Province     string      `json:"province"      description:"省"`
	City         string      `json:"city"          description:"市"`
	District     string      `json:"district"      description:"区"`
	DetailAddr   string      `json:"detail_addr"   description:"详细地址"`
	ShippingFee  int64       `json:"shipping_fee"  description:"运费（分）"`
	Status       string      `json:"status"        description:"物流状态：pending-待发货 picked-已拣货 shipped-已发货 delivering-配送中 delivered-已签收 returned-已退回"`
	ShippedAt    *gtime.Time `json:"shipped_at"    description:"发货时间"`
	DeliveredAt  *gtime.Time `json:"delivered_at"  description:"签收时间"`
	CreatedBy    int64       `json:"created_by"    description:"操作人ID"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updated_at"    description:""`
	DeletedAt    *gtime.Time `json:"deleted_at"    description:""`
}
