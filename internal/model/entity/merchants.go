// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Merchants is the golang structure for table merchants.
type Merchants struct {
	Id              int64       `json:"id"               description:"商家ID"`
	MerchantName    string      `json:"merchant_name"    description:"商家名称（店铺名）"`
	MerchantCode    string      `json:"merchant_code"    description:"商家编码（系统生成，唯一）"`
	MerchantType    int         `json:"merchant_type"    description:"1-个人商家 2-企业商家 3-品牌直营"`
	MerchantLevel   int         `json:"merchant_level"   description:"商家等级 1-普通 2-银牌 3-金牌 4-钻石（影响佣金率/权限）"`
	BusinessScope   string      `json:"business_scope"   description:"经营范围"`
	BusinessYears   int         `json:"business_years"   description:"经营年限（入驻年限）"`
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
	CommissionRate  int64       `json:"commission_rate"  description:"平台抽佣比例（千分比，如 50 表示5%）"`
	SettlementCycle int         `json:"settlement_cycle" description:"结算周期 1-T+1 2-T+7 3-月结"`
	TotalOrders     int         `json:"total_orders"     description:"历史总订单数"`
	TotalSales      int64       `json:"total_sales"      description:"历史总销售额（分）"`
	AvgRating       float64     `json:"avg_rating"       description:"店铺平均评分"`
	ProductCount    int         `json:"product_count"    description:"在售商品数量"`
	SettledAt       *gtime.Time `json:"settled_at"       description:"入驻时间"`
	ExpireAt        *gtime.Time `json:"expire_at"        description:"合同到期时间"`
	CreatedBy       int64       `json:"created_by"       description:"创建人"`
	UpdatedBy       int64       `json:"updated_by"       description:"更新人"`
	CreatedAt       *gtime.Time `json:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updated_at"       description:""`
	DeletedAt       *gtime.Time `json:"deleted_at"       description:""`
}
