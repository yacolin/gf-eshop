// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeliveryTracesDao is the data access object for table tx_delivery_traces.
type DeliveryTracesDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns DeliveryTracesColumns // columns contains all the column names of Table for convenient usage.
}

// DeliveryTracesColumns defines and stores column names for table tx_delivery_traces.
type DeliveryTracesColumns struct {
	Id          string // 轨迹ID
	DeliveryId  string // 关联 tx_deliveries.id
	TrackingNo  string // 运单号（冗余，方便直接查）
	TraceTime   string // 轨迹发生时间
	Location    string // 轨迹地点（如：深圳分拨中心）
	Status      string // 轨迹节点（如：已揽收、已到达分拨中心、派送中）
	Description string // 轨迹描述（如：快件已到达深圳分拨中心）
	CreatedAt   string //
}

// deliveryTracesColumns holds the columns for table tx_delivery_traces.
var deliveryTracesColumns = DeliveryTracesColumns{
	Id:          "id",
	DeliveryId:  "delivery_id",
	TrackingNo:  "tracking_no",
	TraceTime:   "trace_time",
	Location:    "location",
	Status:      "status",
	Description: "description",
	CreatedAt:   "created_at",
}

// NewDeliveryTracesDao creates and returns a new DAO object for table data access.
func NewDeliveryTracesDao() *DeliveryTracesDao {
	return &DeliveryTracesDao{
		group:   "default",
		table:   "tx_delivery_traces",
		columns: deliveryTracesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DeliveryTracesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DeliveryTracesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DeliveryTracesDao) Columns() DeliveryTracesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DeliveryTracesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DeliveryTracesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DeliveryTracesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
