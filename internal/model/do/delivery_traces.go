// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeliveryTraces is the golang structure of table tx_delivery_traces for DAO operations like Where/Data.
type DeliveryTraces struct {
	g.Meta      `orm:"table:tx_delivery_traces, do:true"`
	Id          interface{} // 轨迹ID
	DeliveryId  interface{} // 关联 tx_deliveries.id
	TrackingNo  interface{} // 运单号（冗余，方便直接查）
	TraceTime   *gtime.Time // 轨迹发生时间
	Location    interface{} // 轨迹地点（如：深圳分拨中心）
	Status      interface{} // 轨迹节点（如：已揽收、已到达分拨中心、派送中）
	Description interface{} // 轨迹描述（如：快件已到达深圳分拨中心）
	CreatedAt   *gtime.Time //
}
