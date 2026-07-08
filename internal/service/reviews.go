package service

import (
	"context"

	"gf-eshop/api/reviews/v1"
)

type IReviews interface {
	List(ctx context.Context, req *v1.ReviewsListReq) (res *v1.ReviewsListRes, err error)
	Detail(ctx context.Context, req *v1.ReviewsDetailReq) (res *v1.ReviewsDetailRes, err error)
	Create(ctx context.Context, req *v1.ReviewsCreateReq) (res *v1.ReviewsCreateRes, err error)
	Update(ctx context.Context, req *v1.ReviewsUpdateReq) (res *v1.ReviewsUpdateRes, err error)
	Delete(ctx context.Context, req *v1.ReviewsDeleteReq) (res *v1.ReviewsDeleteRes, err error)
	Audit(ctx context.Context, req *v1.ReviewsAuditReq) (res *v1.ReviewsAuditRes, err error)
	ListReplies(ctx context.Context, req *v1.ReviewsListRepliesReq) (res *v1.ReviewsListRepliesRes, err error)
	CreateReply(ctx context.Context, req *v1.ReviewsCreateReplyReq) (res *v1.ReviewsCreateReplyRes, err error)
	DeleteReply(ctx context.Context, req *v1.ReviewsDeleteReplyReq) (res *v1.ReviewsDeleteReplyRes, err error)
	ListAuditLogs(ctx context.Context, req *v1.ReviewsListAuditLogsReq) (res *v1.ReviewsListAuditLogsRes, err error)
}

var localReviews IReviews

func Reviews() IReviews {
	if localReviews == nil {
		panic("implement not found for interface IReviews, forgot register?")
	}
	return localReviews
}

func RegisterReviews(i IReviews) {
	localReviews = i
}
