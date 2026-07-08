package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type MerchantContactsListReq struct {
	g.Meta `path:"/merchant-contacts" tags:"MerchantContacts" method:"get" summary:"联系人列表"`

	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	MerchantId int64 `json:"merchant_id" description:"商家ID，可选筛选"`
}

type MerchantContactsListRes struct {
	List  []*entity.MerchantContacts `json:"list"`
	Total int                        `json:"total"`
}

// ---------- Detail ----------
type MerchantContactsDetailReq struct {
	g.Meta `path:"/merchant-contacts/{id}" tags:"MerchantContacts" method:"get" summary:"联系人详情"`
	Id     int64 `json:"id"`
}
type MerchantContactsDetailRes struct {
	*entity.MerchantContacts
}

// ---------- Create ----------
type MerchantContactsCreateReq struct {
	g.Meta       `path:"/merchant-contacts" tags:"MerchantContacts" method:"post" summary:"新增联系人"`

	MerchantId   int64  `json:"merchant_id"   v:"required" description:"商家ID"`
	ContactName  string `json:"contact_name"  v:"required" description:"联系人姓名"`
	ContactPhone string `json:"contact_phone" v:"required" description:"联系电话"`
	ContactRole  string `json:"contact_role"  description:"联系人角色：finance-财务 legal-法人 operation-运营"`
	IsPrimary    int    `json:"is_primary"    description:"是否主要联系人 0-否 1-是"`
}

type MerchantContactsCreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type MerchantContactsUpdateReq struct {
	g.Meta       `path:"/merchant-contacts/{id}" tags:"MerchantContacts" method:"put" summary:"更新联系人"`

	Id           int64  `json:"id"            v:"required"`
	ContactName  string `json:"contact_name"  description:"联系人姓名"`
	ContactPhone string `json:"contact_phone" description:"联系电话"`
	ContactRole  string `json:"contact_role"  description:"联系人角色：finance-财务 legal-法人 operation-运营"`
	IsPrimary    int    `json:"is_primary"    description:"是否主要联系人 0-否 1-是"`
}

type MerchantContactsUpdateRes struct{}

// ---------- Delete ----------
type MerchantContactsDeleteReq struct {
	g.Meta `path:"/merchant-contacts/{id}" tags:"MerchantContacts" method:"delete" summary:"删除联系人"`
	Id     int64 `json:"id"`
}
type MerchantContactsDeleteRes struct{}
