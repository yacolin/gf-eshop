// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrdersDao is the data access object for table tx_orders.
type OrdersDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns OrdersColumns // columns contains all the column names of Table for convenient usage.
}

// OrdersColumns defines and stores column names for table tx_orders.
type OrdersColumns struct {
	Id             string // 自增主键
	OrderNo        string // 父订单号（业务唯一键，如 202612010001）
	UserId         string // 用户ID
	TotalAmount    string // 商品总金额（分）
	DiscountAmount string // 优惠金额（分，含优惠券/满减）
	ShippingFee    string // 运费（分）
	PayAmount      string // 实付金额（分 = total - discount + shipping）
	Status         string // 父订单状态：pending-待支付 paid-已支付 partial_shipped-部分发货 completed-已完成 cancelled-已取消 closed-已关闭 refunding-退款中 refunded-已退款
	PaymentStatus  string // 支付状态：unpaid-未支付 paying-支付中 paid-已支付 refunding-退款中 refunded-已退款
	PaymentMethod  string // 支付方式：wechat-微信 alipay-支付宝 wallet-余额
	Consignee      string // 收货人
	Phone          string // 联系电话
	Province       string // 省
	City           string // 市
	District       string // 区
	DetailAddr     string // 详细地址
	ZipCode        string // 邮编
	CouponId       string // 使用的优惠券ID
	CouponSnapshot string // 优惠券快照（名称/面值等，便于售后追溯）
	BuyerRemark    string // 买家备注
	SellerRemark   string // 卖家备注
	Source         string // 订单来源：pc-电脑端 app-APP miniapp-小程序 h5-H5
	PaidAt         string // 支付时间
	ShippedAt      string // 发货时间
	DeliveredAt    string // 签收时间
	CompletedAt    string // 完成时间
	ClosedAt       string // 关闭时间
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// ordersColumns holds the columns for table tx_orders.
var ordersColumns = OrdersColumns{
	Id:             "id",
	OrderNo:        "order_no",
	UserId:         "user_id",
	TotalAmount:    "total_amount",
	DiscountAmount: "discount_amount",
	ShippingFee:    "shipping_fee",
	PayAmount:      "pay_amount",
	Status:         "status",
	PaymentStatus:  "payment_status",
	PaymentMethod:  "payment_method",
	Consignee:      "consignee",
	Phone:          "phone",
	Province:       "province",
	City:           "city",
	District:       "district",
	DetailAddr:     "detail_addr",
	ZipCode:        "zip_code",
	CouponId:       "coupon_id",
	CouponSnapshot: "coupon_snapshot",
	BuyerRemark:    "buyer_remark",
	SellerRemark:   "seller_remark",
	Source:         "source",
	PaidAt:         "paid_at",
	ShippedAt:      "shipped_at",
	DeliveredAt:    "delivered_at",
	CompletedAt:    "completed_at",
	ClosedAt:       "closed_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewOrdersDao creates and returns a new DAO object for table data access.
func NewOrdersDao() *OrdersDao {
	return &OrdersDao{
		group:   "default",
		table:   "tx_orders",
		columns: ordersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OrdersDao) Columns() OrdersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OrdersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
