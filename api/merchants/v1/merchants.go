package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/model/entity"
)


// ---------- List ----------
type MerchantsListReq struct {
	g.Meta `path:"/merchants" tags:"Merchants" method:"get" summary:"商家列表"`

	Page       int    `json:"page"`        // 页码，默认1
	PageSize   int    `json:"page_size"`    // 每页条数，默认20
}

type MerchantsListRes struct {
	List []*entity.Merchants 	`json:"list"`
	Total int 				 	`json:"total"`
}


// ---------- Detail ----------
type MerchantsDetailReq struct {
	g.Meta `path:"/merchants/{id}" tags:"Merchants" method:"get" summary:"商家详情"`
	Id     int64 `json:"id"`
}
type MerchantsDetailRes struct {
	*entity.Merchants
}

// ---------- Create ----------
type MerchantsCreateReq struct {
	g.Meta `path:"/merchants" tags:"Merchants" method:"post" summary:"新增商家"`

	MerchantName    string  `json:"merchant_name"     v:"required|length:1,100" description:"商家名称"`
	MerchantType    int     `json:"merchant_type"     description:"1-个人商家 2-企业商家 3-品牌直营"`
	MerchantLevel   int     `json:"merchant_level"    description:"商家等级 1-普通 2-银牌 3-金牌 4-钻石"`
	BusinessScope   string  `json:"business_scope"    description:"经营范围"`
	BusinessYears   int     `json:"business_years"    description:"经营年限"`
	ContactPerson   string  `json:"contact_person"    description:"主要联系人"`
	ContactPhone    string  `json:"contact_phone"     description:"联系电话"`
	ContactEmail    string  `json:"contact_email"     description:"联系邮箱"`
	LogoUrl         string  `json:"logo_url"          description:"店铺Logo"`
	BannerUrl       string  `json:"banner_url"        description:"店铺Banner图"`
	ShopDescription string  `json:"shop_description"  description:"店铺简介"`
	Status          int     `json:"status"            description:"0-待审核 1-正常 2-冻结 3-已注销"`
	CommissionRate  int64   `json:"commission_rate"   description:"平台抽佣比例（千分比）"`
	SettlementCycle int     `json:"settlement_cycle"  description:"结算周期 1-T+1 2-T+7 3-月结"`
	SettledAt       *gtime.Time `json:"settled_at"    description:"入驻时间"`
	ExpireAt        *gtime.Time `json:"expire_at"     description:"合同到期时间"`
}

type MerchantsCreateRes struct {
	Id int64 `json:"id"`
}


// ---------- Update ----------
type MerchantsUpdateReq struct {
	g.Meta `path:"/merchants/{id}" tags:"Merchants" method:"put" summary:"更新商家"`

	Id              int64       `json:"id"               v:"required"`
	MerchantName    string      `json:"merchant_name"    description:"商家名称"`
	MerchantType    int         `json:"merchant_type"    description:"1-个人商家 2-企业商家 3-品牌直营"`
	MerchantLevel   int         `json:"merchant_level"   description:"商家等级 1-普通 2-银牌 3-金牌 4-钻石"`
	BusinessScope   string      `json:"business_scope"   description:"经营范围"`
	BusinessYears   int         `json:"business_years"   description:"经营年限"`
	ContactPerson   string      `json:"contact_person"   description:"主要联系人"`
	ContactPhone    string      `json:"contact_phone"    description:"联系电话"`
	ContactEmail    string      `json:"contact_email"    description:"联系邮箱"`
	LogoUrl         string      `json:"logo_url"         description:"店铺Logo"`
	BannerUrl       string      `json:"banner_url"       description:"店铺Banner图"`
	ShopDescription string      `json:"shop_description" description:"店铺简介"`
	Status          int         `json:"status"           description:"0-待审核 1-正常 2-冻结 3-已注销"`
	AuditStatus     int         `json:"audit_status"     description:"0-待审核 1-审核通过 2-审核拒绝"`
	AuditReason     string      `json:"audit_reason"     description:"审核拒绝原因"`
	AuditedAt       *gtime.Time `json:"audited_at"       description:"审核时间"`
	FrozenReason    string      `json:"frozen_reason"    description:"冻结原因"`
	CommissionRate  int64       `json:"commission_rate"  description:"平台抽佣比例（千分比）"`
	SettlementCycle int         `json:"settlement_cycle" description:"结算周期 1-T+1 2-T+7 3-月结"`
	ExpireAt        *gtime.Time `json:"expire_at"        description:"合同到期时间"`
}
type MerchantsUpdateRes struct{}



// ---------- Delete ----------
type MerchantsDeleteReq struct {
	g.Meta `path:"/merchants/{id}" tags:"Merchants" method:"delete" summary:"删除商家"`
	Id     int64 `json:"id"`
}
type MerchantsDeleteRes struct{}