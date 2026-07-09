// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeliveryTraces is the golang structure for table delivery_traces.
type DeliveryTraces struct {
	Id          int64       `json:"id"          description:"轨迹ID"`
	DeliveryId  int64       `json:"delivery_id" description:"关联 tx_deliveries.id"`
	TrackingNo  string      `json:"tracking_no" description:"运单号（冗余，方便直接查）"`
	TraceTime   *gtime.Time `json:"trace_time"  description:"轨迹发生时间"`
	Location    string      `json:"location"    description:"轨迹地点（如：深圳分拨中心）"`
	Status      string      `json:"status"      description:"轨迹节点（如：已揽收、已到达分拨中心、派送中）"`
	Description string      `json:"description" description:"轨迹描述（如：快件已到达深圳分拨中心）"`
	CreatedAt   *gtime.Time `json:"created_at"  description:""`
}
