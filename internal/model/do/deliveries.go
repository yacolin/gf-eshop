// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Deliveries is the golang structure of table tx_deliveries for DAO operations like Where/Data.
type Deliveries struct {
	g.Meta       `orm:"table:tx_deliveries, do:true"`
	Id           interface{} // 物流单ID
	DeliveryNo   interface{} // 物流单号（业务唯一）
	OrderId      interface{} // 关联 tx_orders.id
	OrderNo      interface{} // 订单号（冗余）
	MerchantId   interface{} // 所属商家ID
	Carrier      interface{} // 物流商：sf-顺丰 yto-圆通 zto-中通 yunda-韵达 jd-京东物流 other-其他
	TrackingNo   interface{} // 运单号（物流商单号）
	WarehouseId  interface{} // 发货仓库ID
	PackageCount interface{} // 包裹数量
	Consignee    interface{} // 收货人
	Phone        interface{} // 联系电话
	Province     interface{} // 省
	City         interface{} // 市
	District     interface{} // 区
	DetailAddr   interface{} // 详细地址
	ShippingFee  interface{} // 运费（分）
	Status       interface{} // 物流状态：pending-待发货 picked-已拣货 shipped-已发货 delivering-配送中 delivered-已签收 returned-已退回
	ShippedAt    *gtime.Time // 发货时间
	DeliveredAt  *gtime.Time // 签收时间
	CreatedBy    interface{} // 操作人ID
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
}
