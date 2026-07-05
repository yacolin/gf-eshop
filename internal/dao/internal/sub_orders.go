// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubOrdersDao is the data access object for table tx_sub_orders.
type SubOrdersDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns SubOrdersColumns // columns contains all the column names of Table for convenient usage.
}

// SubOrdersColumns defines and stores column names for table tx_sub_orders.
type SubOrdersColumns struct {
	Id             string // 子订单ID
	SubOrderNo     string // 子订单号（按商家拆单后的业务唯一键）
	ParentOrderId  string // 父订单ID
	ParentOrderNo  string // 父订单号（冗余，便于查询）
	UserId         string // 用户ID
	MerchantId     string // 所属商家ID
	TotalAmount    string // 商品总金额（分）
	DiscountAmount string // 优惠金额（分）
	ShippingFee    string // 运费（分）
	PayAmount      string // 子订单实付金额（分）
	Status         string // 子订单状态：pending-待支付 paid-已支付 shipped-已发货 delivered-已签收 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款
	RefundStatus   string // 退款状态：none-无 partial_refunded-部分退款 refunded-已退款
	SellerRemark   string // 卖家备注
	PaidAt         string // 支付时间
	ShippedAt      string // 发货时间
	DeliveredAt    string // 签收时间
	CompletedAt    string // 完成时间
	ClosedAt       string // 关闭时间
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// subOrdersColumns holds the columns for table tx_sub_orders.
var subOrdersColumns = SubOrdersColumns{
	Id:             "id",
	SubOrderNo:     "sub_order_no",
	ParentOrderId:  "parent_order_id",
	ParentOrderNo:  "parent_order_no",
	UserId:         "user_id",
	MerchantId:     "merchant_id",
	TotalAmount:    "total_amount",
	DiscountAmount: "discount_amount",
	ShippingFee:    "shipping_fee",
	PayAmount:      "pay_amount",
	Status:         "status",
	RefundStatus:   "refund_status",
	SellerRemark:   "seller_remark",
	PaidAt:         "paid_at",
	ShippedAt:      "shipped_at",
	DeliveredAt:    "delivered_at",
	CompletedAt:    "completed_at",
	ClosedAt:       "closed_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewSubOrdersDao creates and returns a new DAO object for table data access.
func NewSubOrdersDao() *SubOrdersDao {
	return &SubOrdersDao{
		group:   "default",
		table:   "tx_sub_orders",
		columns: subOrdersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubOrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubOrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubOrdersDao) Columns() SubOrdersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubOrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubOrdersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubOrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
