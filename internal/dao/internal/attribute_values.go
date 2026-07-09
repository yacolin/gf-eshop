// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AttributeValuesDao is the data access object for table sp_attribute_values.
type AttributeValuesDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns AttributeValuesColumns // columns contains all the column names of Table for convenient usage.
}

// AttributeValuesColumns defines and stores column names for table sp_attribute_values.
type AttributeValuesColumns struct {
	Id           string //
	AttributeId  string // 关联属性ID
	Value        string // 属性值（如：256G、红色）
	Alias        string // 别名列表，如["深空灰","黑灰"]，用于搜索纠错、模糊匹配
	SearchWeight string // 搜索权重（值越大匹配优先级越高）
	NumericValue string // 数值型值（用于区间筛选）
	ColorHex     string // 颜色色值（#FF0000）
	SortOrder    string // 排序权重
	Status       string // 1-启用 0-禁用
	CreatedAt    string //
}

// attributeValuesColumns holds the columns for table sp_attribute_values.
var attributeValuesColumns = AttributeValuesColumns{
	Id:           "id",
	AttributeId:  "attribute_id",
	Value:        "value",
	Alias:        "alias",
	SearchWeight: "search_weight",
	NumericValue: "numeric_value",
	ColorHex:     "color_hex",
	SortOrder:    "sort_order",
	Status:       "status",
	CreatedAt:    "created_at",
}

// NewAttributeValuesDao creates and returns a new DAO object for table data access.
func NewAttributeValuesDao() *AttributeValuesDao {
	return &AttributeValuesDao{
		group:   "default",
		table:   "sp_attribute_values",
		columns: attributeValuesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AttributeValuesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AttributeValuesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AttributeValuesDao) Columns() AttributeValuesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AttributeValuesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AttributeValuesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AttributeValuesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
