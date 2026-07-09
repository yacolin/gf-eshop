// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AfterSalesDao is the data access object for table tx_after_sales.
type AfterSalesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns AfterSalesColumns // columns contains all the column names of Table for convenient usage.
}

// AfterSalesColumns defines and stores column names for table tx_after_sales.
type AfterSalesColumns struct {
	Id                 string // 主键
	AfterSaleNo        string // 售后单号
	OrderId            string // 订单ID
	OrderItemId        string // 订单明细ID
	MerchantId         string // 商家ID
	UserId             string // 用户ID
	AfterSaleType      string // 1-退款 2-退货 3-换货
	ApplyQuantity      string // 申请售后数量
	Reason             string // 申请原因
	Amount             string // 退款金额（分）
	RefundId           string // 关联退款单ID
	RefundNo           string // 关联退款单号
	Status             string // 0-待审核 1-审核通过(待退货) 2-退货中 3-待退款 4-已完成 5-已拒绝 6-已取消
	ReturnCarrier      string // 退货物流商
	ReturnTrackingNo   string // 退货运单号
	ReturnShippedAt    string // 买家退货发出时间
	MerchantReceivedAt string // 商家收货时间
	ApplyAt            string // 申请时间
	AuditedAt          string // 审核时间
	CompletedAt        string // 完成时间
	CreatedAt          string //
	UpdatedAt          string //
	DeletedAt          string //
}

// afterSalesColumns holds the columns for table tx_after_sales.
var afterSalesColumns = AfterSalesColumns{
	Id:                 "id",
	AfterSaleNo:        "after_sale_no",
	OrderId:            "order_id",
	OrderItemId:        "order_item_id",
	MerchantId:         "merchant_id",
	UserId:             "user_id",
	AfterSaleType:      "after_sale_type",
	ApplyQuantity:      "apply_quantity",
	Reason:             "reason",
	Amount:             "amount",
	RefundId:           "refund_id",
	RefundNo:           "refund_no",
	Status:             "status",
	ReturnCarrier:      "return_carrier",
	ReturnTrackingNo:   "return_tracking_no",
	ReturnShippedAt:    "return_shipped_at",
	MerchantReceivedAt: "merchant_received_at",
	ApplyAt:            "apply_at",
	AuditedAt:          "audited_at",
	CompletedAt:        "completed_at",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewAfterSalesDao creates and returns a new DAO object for table data access.
func NewAfterSalesDao() *AfterSalesDao {
	return &AfterSalesDao{
		group:   "default",
		table:   "tx_after_sales",
		columns: afterSalesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AfterSalesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AfterSalesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AfterSalesDao) Columns() AfterSalesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AfterSalesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AfterSalesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AfterSalesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
