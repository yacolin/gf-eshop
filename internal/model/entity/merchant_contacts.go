// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantContacts is the golang structure for table merchant_contacts.
type MerchantContacts struct {
	Id           int64       `json:"id"            description:"主键"`
	MerchantId   int64       `json:"merchant_id"   description:"商家ID"`
	ContactName  string      `json:"contact_name"  description:"联系人姓名"`
	ContactPhone string      `json:"contact_phone" description:"联系电话"`
	ContactRole  string      `json:"contact_role"  description:"联系人角色：finance-财务 legal-法人 operation-运营"`
	IsPrimary    int         `json:"is_primary"    description:"是否主要联系人 0-否 1-是"`
	CreatedAt    *gtime.Time `json:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updated_at"    description:""`
	DeletedAt    *gtime.Time `json:"deleted_at"    description:""`
}
