// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionProducts is the golang structure for table promotion_products.
type PromotionProducts struct {
	Id          int64       `json:"id"           description:"主键ID"`
	PromotionId int64       `json:"promotion_id" description:"促销ID"`
	MerchantId  int64       `json:"merchant_id"  description:"所属商家ID"`
	ProductType int         `json:"product_type" description:"1-全站 2-指定分类 3-指定SPU 4-指定SKU"`
	TargetId    int64       `json:"target_id"    description:"目标ID（SPU_ID或SKU_ID或Category_ID）"`
	CreatedAt   *gtime.Time `json:"created_at"   description:""`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:""`
}
