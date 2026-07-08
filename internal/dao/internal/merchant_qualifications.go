// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantQualificationsDao is the data access object for table mch_merchant_qualifications.
type MerchantQualificationsDao struct {
	table   string                        // table is the underlying table name of the DAO.
	group   string                        // group is the database configuration group name of current DAO.
	columns MerchantQualificationsColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantQualificationsColumns defines and stores column names for table mch_merchant_qualifications.
type MerchantQualificationsColumns struct {
	Id                string // 主键
	MerchantId        string // 商家ID
	QualificationType string // 资质类型：business_license-营业执照 food-食品经营许可 brand_authorization-品牌授权
	QualificationName string // 资质名称
	FileUrl           string // 资质文件URL
	ExpireAt          string // 有效期
	Status            string // 0-待审核 1-审核通过 2-已过期 3-审核拒绝
	AuditRemark       string // 审核备注
	CreatedAt         string //
	UpdatedAt         string //
	DeletedAt         string //
}

// merchantQualificationsColumns holds the columns for table mch_merchant_qualifications.
var merchantQualificationsColumns = MerchantQualificationsColumns{
	Id:                "id",
	MerchantId:        "merchant_id",
	QualificationType: "qualification_type",
	QualificationName: "qualification_name",
	FileUrl:           "file_url",
	ExpireAt:          "expire_at",
	Status:            "status",
	AuditRemark:       "audit_remark",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
}

// NewMerchantQualificationsDao creates and returns a new DAO object for table data access.
func NewMerchantQualificationsDao() *MerchantQualificationsDao {
	return &MerchantQualificationsDao{
		group:   "default",
		table:   "mch_merchant_qualifications",
		columns: merchantQualificationsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantQualificationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantQualificationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantQualificationsDao) Columns() MerchantQualificationsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantQualificationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantQualificationsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantQualificationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
