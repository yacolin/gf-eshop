// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUsers is the golang structure of table mch_merchant_users for DAO operations like Where/Data.
type MerchantUsers struct {
	g.Meta      `orm:"table:mch_merchant_users, do:true"`
	Id          interface{} // 主键
	MerchantId  interface{} // 商家ID
	StaffId     interface{} // 员工ID（关联 sys_staff.id）
	RoleId      interface{} // 店铺角色ID（关联 sys_roles.id）
	Status      interface{} // 1-正常 2-禁用
	InvitedAt   *gtime.Time // 邀请时间
	LastLoginAt *gtime.Time // 最后登录时间
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
