// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleEvidences is the golang structure for table after_sale_evidences.
type AfterSaleEvidences struct {
	Id          int64       `json:"id"            description:"主键"`
	AfterSaleId int64       `json:"after_sale_id" description:"售后单ID"`
	MediaType   int         `json:"media_type"    description:"1-图片 2-视频"`
	MediaUrl    string      `json:"media_url"     description:"凭证URL"`
	SortOrder   int         `json:"sort_order"    description:"排序"`
	CreatedAt   *gtime.Time `json:"created_at"    description:""`
}
