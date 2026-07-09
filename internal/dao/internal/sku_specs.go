// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SkuSpecsDao is the data access object for table sp_sku_specs.
type SkuSpecsDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns SkuSpecsColumns // columns contains all the column names of Table for convenient usage.
}

// SkuSpecsColumns defines and stores column names for table sp_sku_specs.
type SkuSpecsColumns struct {
	Id               string //
	SkuId            string // 关联SKU ID
	AttributeId      string // 关联属性ID
	AttributeValueId string // 关联属性值ID
	SortOrder        string // 展示顺序
}

// skuSpecsColumns holds the columns for table sp_sku_specs.
var skuSpecsColumns = SkuSpecsColumns{
	Id:               "id",
	SkuId:            "sku_id",
	AttributeId:      "attribute_id",
	AttributeValueId: "attribute_value_id",
	SortOrder:        "sort_order",
}

// NewSkuSpecsDao creates and returns a new DAO object for table data access.
func NewSkuSpecsDao() *SkuSpecsDao {
	return &SkuSpecsDao{
		group:   "default",
		table:   "sp_sku_specs",
		columns: skuSpecsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SkuSpecsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SkuSpecsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SkuSpecsDao) Columns() SkuSpecsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SkuSpecsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SkuSpecsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SkuSpecsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
