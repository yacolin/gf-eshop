package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type MerchantQualificationsListReq struct {
	g.Meta `path:"/merchant-qualifications" tags:"MerchantQualifications" method:"get" summary:"资质列表"`

	Page       int `json:"page"`         // 页码，默认1
	PageSize   int `json:"page_size"`    // 每页条数，默认20
	MerchantId int `json:"merchant_id" dc:"商家ID（可选）"`
	Status     int `json:"status"     dc:"状态（可选）"`
}

type MerchantQualificationsListRes struct {
	List  []*entity.MerchantQualifications `json:"list"`
	Total int                              `json:"total"`
}

// ---------- Detail ----------
type MerchantQualificationsDetailReq struct {
	g.Meta `path:"/merchant-qualifications/{id}" tags:"MerchantQualifications" method:"get" summary:"资质详情"`
	Id     int64 `json:"id"`
}
type MerchantQualificationsDetailRes struct {
	*entity.MerchantQualifications
}

// ---------- Create ----------
type MerchantQualificationsCreateReq struct {
	g.Meta `path:"/merchant-qualifications" tags:"MerchantQualifications" method:"post" summary:"新增资质"`

	MerchantId        int64       `json:"merchant_id"        v:"required" description:"商家ID"`
	QualificationType string      `json:"qualification_type" v:"required" description:"资质类型"`
	QualificationName string      `json:"qualification_name" v:"required" description:"资质名称"`
	FileUrl           string      `json:"file_url"           v:"required" description:"资质文件URL"`
	ExpireAt          *gtime.Time `json:"expire_at"                       description:"有效期"`
}

type MerchantQualificationsCreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type MerchantQualificationsUpdateReq struct {
	g.Meta `path:"/merchant-qualifications/{id}" tags:"MerchantQualifications" method:"put" summary:"更新资质"`

	Id                int64       `json:"id"                 v:"required"`
	QualificationName string      `json:"qualification_name" description:"资质名称"`
	FileUrl           string      `json:"file_url"           description:"资质文件URL"`
	ExpireAt          *gtime.Time `json:"expire_at"          description:"有效期"`
}
type MerchantQualificationsUpdateRes struct{}

// ---------- Delete ----------
type MerchantQualificationsDeleteReq struct {
	g.Meta `path:"/merchant-qualifications/{id}" tags:"MerchantQualifications" method:"delete" summary:"删除资质"`
	Id     int64 `json:"id"`
}
type MerchantQualificationsDeleteRes struct{}

// ---------- Audit ----------
type MerchantQualificationsAuditReq struct {
	g.Meta `path:"/merchant-qualifications/{id}/audit" tags:"MerchantQualifications" method:"put" summary:"审核资质"`

	Id          int64  `json:"id"     v:"required"`
	Status      int    `json:"status" v:"required|in:1,3" description:"1-审核通过 3-审核拒绝"`
	AuditRemark string `json:"audit_remark"               description:"审核备注"`
}
type MerchantQualificationsAuditRes struct{}
