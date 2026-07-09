// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductAttributesDao is the data access object for table sp_product_attributes.
type ProductAttributesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns ProductAttributesColumns // columns contains all the column names of Table for convenient usage.
}

// ProductAttributesColumns defines and stores column names for table sp_product_attributes.
type ProductAttributesColumns struct {
	Id               string //
	ProductId        string // 关联 products.id
	AttributeId      string // 关联 attributes.id
	AttributeValueId string // 引用属性值字典ID（可选，关联 attribute_values.id）
	Value            string // 属性值（如：A16）。有字典值时冗余存储便于展示，无字典值时存自由文本
	SortOrder        string // 排序权重（越小越靠前，用于控制前台展示顺序）
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string // 软删除
}

// productAttributesColumns holds the columns for table sp_product_attributes.
var productAttributesColumns = ProductAttributesColumns{
	Id:               "id",
	ProductId:        "product_id",
	AttributeId:      "attribute_id",
	AttributeValueId: "attribute_value_id",
	Value:            "value",
	SortOrder:        "sort_order",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewProductAttributesDao creates and returns a new DAO object for table data access.
func NewProductAttributesDao() *ProductAttributesDao {
	return &ProductAttributesDao{
		group:   "default",
		table:   "sp_product_attributes",
		columns: productAttributesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductAttributesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductAttributesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductAttributesDao) Columns() ProductAttributesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductAttributesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductAttributesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductAttributesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
