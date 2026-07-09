// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSales is the golang structure of table tx_after_sales for DAO operations like Where/Data.
type AfterSales struct {
	g.Meta             `orm:"table:tx_after_sales, do:true"`
	Id                 interface{} // 主键
	AfterSaleNo        interface{} // 售后单号
	OrderId            interface{} // 订单ID
	OrderItemId        interface{} // 订单明细ID
	MerchantId         interface{} // 商家ID
	UserId             interface{} // 用户ID
	AfterSaleType      interface{} // 1-退款 2-退货 3-换货
	ApplyQuantity      interface{} // 申请售后数量
	Reason             interface{} // 申请原因
	Amount             interface{} // 退款金额（分）
	RefundId           interface{} // 关联退款单ID
	RefundNo           interface{} // 关联退款单号
	Status             interface{} // 0-待审核 1-审核通过(待退货) 2-退货中 3-待退款 4-已完成 5-已拒绝 6-已取消
	ReturnCarrier      interface{} // 退货物流商
	ReturnTrackingNo   interface{} // 退货运单号
	ReturnShippedAt    *gtime.Time // 买家退货发出时间
	MerchantReceivedAt *gtime.Time // 商家收货时间
	ApplyAt            *gtime.Time // 申请时间
	AuditedAt          *gtime.Time // 审核时间
	CompletedAt        *gtime.Time // 完成时间
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
	DeletedAt          *gtime.Time //
}
