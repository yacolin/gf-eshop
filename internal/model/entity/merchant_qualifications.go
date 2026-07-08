// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantQualifications is the golang structure for table merchant_qualifications.
type MerchantQualifications struct {
	Id                int64       `json:"id"                 description:"主键"`
	MerchantId        int64       `json:"merchant_id"        description:"商家ID"`
	QualificationType string      `json:"qualification_type" description:"资质类型：business_license-营业执照 food-食品经营许可 brand_authorization-品牌授权"`
	QualificationName string      `json:"qualification_name" description:"资质名称"`
	FileUrl           string      `json:"file_url"           description:"资质文件URL"`
	ExpireAt          *gtime.Time `json:"expire_at"          description:"有效期"`
	Status            int         `json:"status"             description:"0-待审核 1-审核通过 2-已过期 3-审核拒绝"`
	AuditRemark       string      `json:"audit_remark"       description:"审核备注"`
	CreatedAt         *gtime.Time `json:"created_at"         description:""`
	UpdatedAt         *gtime.Time `json:"updated_at"         description:""`
	DeletedAt         *gtime.Time `json:"deleted_at"         description:""`
}
