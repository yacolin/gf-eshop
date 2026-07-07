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
	ProductId   int64       `json:"product_id"   description:"产品ID（product_type=3/4时使用）"`
	CategoryId  int64       `json:"category_id"  description:"分类ID（product_type=2时使用）"`
	CreatedAt   *gtime.Time `json:"created_at"   description:""`
	DeletedAt   *gtime.Time `json:"deleted_at"   description:""`
}
