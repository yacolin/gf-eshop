// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AfterSaleEvidencesDao is the data access object for table tx_after_sale_evidences.
type AfterSaleEvidencesDao struct {
	table   string                    // table is the underlying table name of the DAO.
	group   string                    // group is the database configuration group name of current DAO.
	columns AfterSaleEvidencesColumns // columns contains all the column names of Table for convenient usage.
}

// AfterSaleEvidencesColumns defines and stores column names for table tx_after_sale_evidences.
type AfterSaleEvidencesColumns struct {
	Id          string // 主键
	AfterSaleId string // 售后单ID
	MediaType   string // 1-图片 2-视频
	MediaUrl    string // 凭证URL
	SortOrder   string // 排序
	CreatedAt   string //
}

// afterSaleEvidencesColumns holds the columns for table tx_after_sale_evidences.
var afterSaleEvidencesColumns = AfterSaleEvidencesColumns{
	Id:          "id",
	AfterSaleId: "after_sale_id",
	MediaType:   "media_type",
	MediaUrl:    "media_url",
	SortOrder:   "sort_order",
	CreatedAt:   "created_at",
}

// NewAfterSaleEvidencesDao creates and returns a new DAO object for table data access.
func NewAfterSaleEvidencesDao() *AfterSaleEvidencesDao {
	return &AfterSaleEvidencesDao{
		group:   "default",
		table:   "tx_after_sale_evidences",
		columns: afterSaleEvidencesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AfterSaleEvidencesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AfterSaleEvidencesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AfterSaleEvidencesDao) Columns() AfterSaleEvidencesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AfterSaleEvidencesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AfterSaleEvidencesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AfterSaleEvidencesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
