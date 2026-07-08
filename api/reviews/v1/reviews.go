package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type ReviewsListReq struct {
	g.Meta `path:"/reviews" tags:"Reviews" method:"get" summary:"评价列表"`

	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Id         int   `json:"id"`
	MerchantId int   `json:"merchant_id"`
	SpuId      int   `json:"spu_id"`
	Status     *int  `json:"status"`
}
type ReviewsListRes struct {
	List  []*entity.Reviews `json:"list"`
	Total int               `json:"total"`
}

// ---------- Detail ----------
type ReviewsDetailReq struct {
	g.Meta `path:"/reviews/{id}" tags:"Reviews" method:"get" summary:"评价详情"`
	Id     int64 `json:"id"`
}
type ReviewsDetailRes struct {
	*entity.Reviews
}

// ---------- Create ----------
type ReviewsCreateReq struct {
	g.Meta `path:"/reviews" tags:"Reviews" method:"post" summary:"新增评价"`

	OrderId         int64  `json:"order_id"          v:"required"`
	OrderItemId     int64  `json:"order_item_id"`
	SpuId           int64  `json:"spu_id"            v:"required"`
	SkuId           int64  `json:"sku_id"`
	OverallRating   int    `json:"overall_rating"    v:"required|min:1|max:5"`
	QualityRating   int    `json:"quality_rating"    v:"min:1|max:5"`
	LogisticsRating int    `json:"logistics_rating"  v:"min:1|max:5"`
	ServiceRating   int    `json:"service_rating"    v:"min:1|max:5"`
	Content         string `json:"content"`
	IsAnonymous     int    `json:"is_anonymous"`
}
type ReviewsCreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type ReviewsUpdateReq struct {
	g.Meta `path:"/reviews/{id}" tags:"Reviews" method:"put" summary:"更新评价"`

	Id              int64  `json:"id"               v:"required"`
	OverallRating   int    `json:"overall_rating"   v:"min:1|max:5"`
	QualityRating   int    `json:"quality_rating"   v:"min:1|max:5"`
	LogisticsRating int    `json:"logistics_rating" v:"min:1|max:5"`
	ServiceRating   int    `json:"service_rating"   v:"min:1|max:5"`
	Content         string `json:"content"`
	IsAnonymous     int    `json:"is_anonymous"`
}
type ReviewsUpdateRes struct{}

// ---------- Delete ----------
type ReviewsDeleteReq struct {
	g.Meta `path:"/reviews/{id}" tags:"Reviews" method:"delete" summary:"删除评价"`
	Id     int64 `json:"id"`
}
type ReviewsDeleteRes struct{}

// ---------- Audit ----------
type ReviewsAuditReq struct {
	g.Meta      `path:"/reviews/{id}/audit" tags:"Reviews" method:"put" summary:"审核评价"`
	Id          int64  `json:"id"            v:"required"`
	Status      int    `json:"status"        v:"required|in:1,2"`
	RejectReason string `json:"reject_reason"`
}
type ReviewsAuditRes struct{}

// ---------- ListReplies ----------
type ReviewsListRepliesReq struct {
	g.Meta `path:"/reviews/{id}/replies" tags:"Reviews" method:"get" summary:"评价回复列表"`
	Id     int64 `json:"id"`
}
type ReviewsListRepliesRes struct {
	List []*entity.ReviewReplies `json:"list"`
}

// ---------- CreateReply ----------
type ReviewsCreateReplyReq struct {
	g.Meta    `path:"/reviews/{id}/replies" tags:"Reviews" method:"post" summary:"新增回复"`
	Id        int64  `json:"id"         v:"required"`
	Content   string `json:"content"    v:"required"`
	ReplyType int    `json:"reply_type"`
}
type ReviewsCreateReplyRes struct {
	Id int64 `json:"id"`
}

// ---------- DeleteReply ----------
type ReviewsDeleteReplyReq struct {
	g.Meta  `path:"/reviews/{id}/replies/{replyId}" tags:"Reviews" method:"delete" summary:"删除回复"`
	Id      int64 `json:"id"`
	ReplyId int64 `json:"replyId"`
}
type ReviewsDeleteReplyRes struct{}

// ---------- ListAuditLogs ----------
type ReviewsListAuditLogsReq struct {
	g.Meta `path:"/reviews/{id}/audit-logs" tags:"Reviews" method:"get" summary:"审核日志列表"`
	Id     int64 `json:"id"`
}
type ReviewsListAuditLogsRes struct {
	List []*entity.ReviewAuditLogs `json:"list"`
}
