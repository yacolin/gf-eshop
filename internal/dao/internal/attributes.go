// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AttributesDao is the data access object for table sp_attributes.
type AttributesDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns AttributesColumns // columns contains all the column names of Table for convenient usage.
}

// AttributesColumns defines and stores column names for table sp_attributes.
type AttributesColumns struct {
	Id         string //
	Name       string // 属性名称（如：处理器、屏幕尺寸）
	CategoryId string // 所属类目ID（该属性只出现在这个类目下）
	ValueType  string // 1-文本 2-数值 3-颜色
	Filterable string // 是否参与前台筛选
	Unit       string // 单位（如：英寸、GB）
	Required   string // 1-必填（该属性在该类目下创建商品时必须填写）
	Searchable string // 1-作为前台筛选条件（列表页筛选项来源）
	IsSkuSpec  string // 1-是SKU规格（如颜色、内存） 0-仅SPU属性（如上市时间）
	SortOrder  string //
	Status     string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// attributesColumns holds the columns for table sp_attributes.
var attributesColumns = AttributesColumns{
	Id:         "id",
	Name:       "name",
	CategoryId: "category_id",
	ValueType:  "value_type",
	Filterable: "filterable",
	Unit:       "unit",
	Required:   "required",
	Searchable: "searchable",
	IsSkuSpec:  "is_sku_spec",
	SortOrder:  "sort_order",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewAttributesDao creates and returns a new DAO object for table data access.
func NewAttributesDao() *AttributesDao {
	return &AttributesDao{
		group:   "default",
		table:   "sp_attributes",
		columns: attributesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AttributesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AttributesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AttributesDao) Columns() AttributesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AttributesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AttributesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AttributesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
