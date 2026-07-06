// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AddressesDao is the data access object for table usr_addresses.
type AddressesDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns AddressesColumns // columns contains all the column names of Table for convenient usage.
}

// AddressesColumns defines and stores column names for table usr_addresses.
type AddressesColumns struct {
	Id        string // 主键
	UserId    string // 用户ID
	Consignee string // 收货人姓名
	Phone     string // 联系电话
	Country   string // 国家
	Province  string // 省
	City      string // 市
	District  string // 区/县
	Detail    string // 详细地址
	ZipCode   string // 邮编
	Tag       string // 地址标签：home/office/company/other
	IsDefault string // 是否默认地址
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	DeletedAt string // 删除时间
}

// addressesColumns holds the columns for table usr_addresses.
var addressesColumns = AddressesColumns{
	Id:        "id",
	UserId:    "user_id",
	Consignee: "consignee",
	Phone:     "phone",
	Country:   "country",
	Province:  "province",
	City:      "city",
	District:  "district",
	Detail:    "detail",
	ZipCode:   "zip_code",
	Tag:       "tag",
	IsDefault: "is_default",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewAddressesDao creates and returns a new DAO object for table data access.
func NewAddressesDao() *AddressesDao {
	return &AddressesDao{
		group:   "default",
		table:   "usr_addresses",
		columns: addressesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AddressesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AddressesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AddressesDao) Columns() AddressesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AddressesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AddressesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AddressesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
