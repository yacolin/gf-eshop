// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductVersions is the golang structure for table product_versions.
type ProductVersions struct {
	Id            int64       `json:"id"             description:"主键"`
	ProductId     int64       `json:"product_id"     description:"关联 sp_products.id"`
	Version       int         `json:"version"        description:"版本号（从1递增）"`
	Diff          string      `json:"diff"           description:"变更JSON（{\"before\": {...}, \"after\": {...}}）"`
	ChangedFields string      `json:"changed_fields" description:"变更字段列表（如：[\"name\", \"price\", \"status\"]）"`
	Operator      string      `json:"operator"       description:"操作人"`
	OperatorId    int64       `json:"operator_id"    description:"操作人ID"`
	Reason        string      `json:"reason"         description:"变更原因"`
	CreatedAt     *gtime.Time `json:"created_at"     description:""`
}
