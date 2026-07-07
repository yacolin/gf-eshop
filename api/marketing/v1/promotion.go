package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type PromotionProductItem struct {
	Id          int64  `json:"id"`
	ProductType int    `json:"product_type"`
	ProductId   int64  `json:"product_id"`
	CategoryId  int64  `json:"category_id"`
	SpuName     string `json:"spu_name"`
	Subtitle    string `json:"subtitle"`
	MainImage   string `json:"main_image"`
	Unit        string `json:"unit"`
	MinPrice    int64  `json:"min_price"`
	MaxPrice    int64  `json:"max_price"`
	SalesCount  int    `json:"sales_count"`
	SpuStatus   int    `json:"spu_status"`
}

type PromotionListReq struct {
	g.Meta    `path:"/promotions" tags:"Marketing" method:"get" summary:"促销列表"`
	Page      int `json:"page"`       // 页码，默认1
	PageSize  int `json:"page_size"`  // 每页条数，默认10
	Status    *int `json:"status"`     // 按状态筛选
	PromoType int `json:"promo_type"` // 按促销类型筛选
}
type PromotionListRes struct {
	List  []*entity.Promotions `json:"list"`
	Total int                  `json:"total"`
}

type PromotionDetailReq struct {
	g.Meta `path:"/promotions/{id}" tags:"Marketing" method:"get" summary:"促销详情"`
	Id     int64 `json:"id"`
}
type PromotionDetailRes struct {
	*entity.Promotions
}

type PromotionFullDetailReq struct {
	g.Meta `path:"/promotions/{id}/detail" tags:"Marketing" method:"get" summary:"促销详细信息(含规则和商品)"`
	Id     int64 `json:"id"`
}
type PromotionFullDetailRes struct {
	Promotion *entity.Promotions       `json:"promotion"`
	Rule      *entity.PromotionRules   `json:"rule"`
	Products  []*PromotionProductItem  `json:"products"`
}

type PromotionCreateReq struct {
	g.Meta `path:"/promotions" tags:"Marketing" method:"post" summary:"新增促销"`

	PromoName     string `json:"promo_name"         v:"required|length:1,100"`
	PromoType     int    `json:"promo_type"          v:"required|between:1,6"`
	PromoCode     string `json:"promo_code"`
	StartTime     string `json:"start_time"          v:"required"`
	EndTime       string `json:"end_time"            v:"required"`
	TotalQuantity int    `json:"total_quantity"`
	PerUserLimit  int    `json:"per_user_limit"`

	RuleName       string `json:"rule_name"`
	ConditionType  int    `json:"condition_type"`   // 1-无门槛 2-满金额 3-满件数 4-指定用户等级
	ConditionValue int64  `json:"condition_value"`  // 分
	BenefitType    int    `json:"benefit_type"`      // 1-减固定金额 2-打折扣 3-赠品 4-免运费 5-送积分
	BenefitValue   int64  `json:"benefit_value"`     // 分/千分比
	IsStackable    int    `json:"is_stackable"`
	StackPriority  int    `json:"stack_priority"`

	ProductIDs []int64 `json:"product_ids"` // SPU ID 列表（product_type=3 时）
}
type PromotionCreateRes struct {
	Id int64 `json:"id"`
}

type PromotionUpdateReq struct {
	g.Meta `path:"/promotions/{id}" tags:"Marketing" method:"put" summary:"更新促销"`

	Id            int64  `json:"id"                v:"required"`
	PromoName     string `json:"promo_name"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	TotalQuantity int    `json:"total_quantity"`
	PerUserLimit  int    `json:"per_user_limit"`
	Status        int    `json:"status"` // 1-草稿 2-生效中 3-已结束 4-已作废

	RuleName       string `json:"rule_name"`
	ConditionType  int    `json:"condition_type"`
	ConditionValue int64  `json:"condition_value"`
	BenefitType    int    `json:"benefit_type"`
	BenefitValue   int64  `json:"benefit_value"`
	IsStackable    int    `json:"is_stackable"`
	StackPriority  int    `json:"stack_priority"`
}
type PromotionUpdateRes struct{}

type PromotionDeleteReq struct {
	g.Meta `path:"/promotions/{id}" tags:"Marketing" method:"delete" summary:"删除促销"`
	Id     int64 `json:"id"`
}
type PromotionDeleteRes struct{}
