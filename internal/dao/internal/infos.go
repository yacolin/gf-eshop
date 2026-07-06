// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InfosDao is the data access object for table usr_infos.
type InfosDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns InfosColumns // columns contains all the column names of Table for convenient usage.
}

// InfosColumns defines and stores column names for table usr_infos.
type InfosColumns struct {
	Id          string // 主键
	UserId      string // 用户ID
	Gender      string // 性别：0-未知 1-男 2-女
	Birthday    string // 生日
	Bio         string // 个人简介
	Country     string // 国家
	Province    string // 省
	City        string // 市
	ZipCode     string // 邮编
	Language    string // 语言
	Timezone    string // 时区
	LevelId     string // 当前等级ID（关联 usr_levels.id）
	TotalPoints string // 累计积分
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// infosColumns holds the columns for table usr_infos.
var infosColumns = InfosColumns{
	Id:          "id",
	UserId:      "user_id",
	Gender:      "gender",
	Birthday:    "birthday",
	Bio:         "bio",
	Country:     "country",
	Province:    "province",
	City:        "city",
	ZipCode:     "zip_code",
	Language:    "language",
	Timezone:    "timezone",
	LevelId:     "level_id",
	TotalPoints: "total_points",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewInfosDao creates and returns a new DAO object for table data access.
func NewInfosDao() *InfosDao {
	return &InfosDao{
		group:   "default",
		table:   "usr_infos",
		columns: infosColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InfosDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InfosDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InfosDao) Columns() InfosColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InfosDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InfosDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InfosDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
