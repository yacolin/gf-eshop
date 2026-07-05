// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductDescriptions is the golang structure for table product_descriptions.
type ProductDescriptions struct {
	Id                int64       `json:"id"                 description:""`
	ProductId         int64       `json:"product_id"         description:"关联 products.id"`
	Description       string      `json:"description"        description:"商品详情（富文本HTML）"`
	MobileDescription string      `json:"mobile_description" description:"移动端详情（可选，适配手机展示）"`
	CreatedAt         *gtime.Time `json:"created_at"         description:""`
	UpdatedAt         *gtime.Time `json:"updated_at"         description:""`
}
