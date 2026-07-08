// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantContacts is the golang structure of table mch_merchant_contacts for DAO operations like Where/Data.
type MerchantContacts struct {
	g.Meta       `orm:"table:mch_merchant_contacts, do:true"`
	Id           interface{} // 主键
	MerchantId   interface{} // 商家ID
	ContactName  interface{} // 联系人姓名
	ContactPhone interface{} // 联系电话
	ContactRole  interface{} // 联系人角色：finance-财务 legal-法人 operation-运营
	IsPrimary    interface{} // 是否主要联系人 0-否 1-是
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
}
