// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantQualifications is the golang structure of table mch_merchant_qualifications for DAO operations like Where/Data.
type MerchantQualifications struct {
	g.Meta            `orm:"table:mch_merchant_qualifications, do:true"`
	Id                interface{} // 主键
	MerchantId        interface{} // 商家ID
	QualificationType interface{} // 资质类型：business_license-营业执照 food-食品经营许可 brand_authorization-品牌授权
	QualificationName interface{} // 资质名称
	FileUrl           interface{} // 资质文件URL
	ExpireAt          *gtime.Time // 有效期
	Status            interface{} // 0-待审核 1-审核通过 2-已过期 3-审核拒绝
	AuditRemark       interface{} // 审核备注
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
	DeletedAt         *gtime.Time //
}
