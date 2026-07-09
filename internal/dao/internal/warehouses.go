// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WarehousesDao is the data access object for table sp_warehouses.
type WarehousesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns WarehousesColumns // columns contains all the column names of Table for convenient usage.
}

// WarehousesColumns defines and stores column names for table sp_warehouses.
type WarehousesColumns struct {
	Id            string // 仓库ID
	MerchantId    string // 所属商家ID（0表示平台仓）
	WarehouseName string // 仓库名称
	WarehouseType string // 1-平台仓 2-商家仓 3-第三方仓
	WarehouseCode string // 仓库编码
	Status        string // 1-启用 2-禁用
	Province      string // 省
	City          string // 市
	District      string // 区
	Address       string // 详细地址
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// warehousesColumns holds the columns for table sp_warehouses.
var warehousesColumns = WarehousesColumns{
	Id:            "id",
	MerchantId:    "merchant_id",
	WarehouseName: "warehouse_name",
	WarehouseType: "warehouse_type",
	WarehouseCode: "warehouse_code",
	Status:        "status",
	Province:      "province",
	City:          "city",
	District:      "district",
	Address:       "address",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewWarehousesDao creates and returns a new DAO object for table data access.
func NewWarehousesDao() *WarehousesDao {
	return &WarehousesDao{
		group:   "default",
		table:   "sp_warehouses",
		columns: warehousesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *WarehousesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *WarehousesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *WarehousesDao) Columns() WarehousesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *WarehousesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *WarehousesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *WarehousesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
