// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeliveriesDao is the data access object for table tx_deliveries.
type DeliveriesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns DeliveriesColumns // columns contains all the column names of Table for convenient usage.
}

// DeliveriesColumns defines and stores column names for table tx_deliveries.
type DeliveriesColumns struct {
	Id           string // 物流单ID
	DeliveryNo   string // 物流单号（业务唯一）
	OrderId      string // 关联 tx_orders.id
	OrderNo      string // 订单号（冗余）
	MerchantId   string // 所属商家ID
	Carrier      string // 物流商：sf-顺丰 yto-圆通 zto-中通 yunda-韵达 jd-京东物流 other-其他
	TrackingNo   string // 运单号（物流商单号）
	WarehouseId  string // 发货仓库ID
	PackageCount string // 包裹数量
	Consignee    string // 收货人
	Phone        string // 联系电话
	Province     string // 省
	City         string // 市
	District     string // 区
	DetailAddr   string // 详细地址
	ShippingFee  string // 运费（分）
	Status       string // 物流状态：pending-待发货 picked-已拣货 shipped-已发货 delivering-配送中 delivered-已签收 returned-已退回
	ShippedAt    string // 发货时间
	DeliveredAt  string // 签收时间
	CreatedBy    string // 操作人
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
}

// deliveriesColumns holds the columns for table tx_deliveries.
var deliveriesColumns = DeliveriesColumns{
	Id:           "id",
	DeliveryNo:   "delivery_no",
	OrderId:      "order_id",
	OrderNo:      "order_no",
	MerchantId:   "merchant_id",
	Carrier:      "carrier",
	TrackingNo:   "tracking_no",
	WarehouseId:  "warehouse_id",
	PackageCount: "package_count",
	Consignee:    "consignee",
	Phone:        "phone",
	Province:     "province",
	City:         "city",
	District:     "district",
	DetailAddr:   "detail_addr",
	ShippingFee:  "shipping_fee",
	Status:       "status",
	ShippedAt:    "shipped_at",
	DeliveredAt:  "delivered_at",
	CreatedBy:    "created_by",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDeliveriesDao creates and returns a new DAO object for table data access.
func NewDeliveriesDao() *DeliveriesDao {
	return &DeliveriesDao{
		group:   "default",
		table:   "tx_deliveries",
		columns: deliveriesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DeliveriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DeliveriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DeliveriesDao) Columns() DeliveriesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DeliveriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DeliveriesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DeliveriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
