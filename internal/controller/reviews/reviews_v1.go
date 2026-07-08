package reviews

import (
	"context"

	"gf-eshop/api/reviews/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ReviewsListReq) (res *v1.ReviewsListRes, err error) {
	return service.Reviews().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.ReviewsDetailReq) (res *v1.ReviewsDetailRes, err error) {
	return service.Reviews().Detail(ctx, req)
}

func (c *ControllerV1) Create(ctx context.Context, req *v1.ReviewsCreateReq) (res *v1.ReviewsCreateRes, err error) {
	return service.Reviews().Create(ctx, req)
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.ReviewsUpdateReq) (res *v1.ReviewsUpdateRes, err error) {
	return service.Reviews().Update(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.ReviewsDeleteReq) (res *v1.ReviewsDeleteRes, err error) {
	return service.Reviews().Delete(ctx, req)
}

func (c *ControllerV1) Audit(ctx context.Context, req *v1.ReviewsAuditReq) (res *v1.ReviewsAuditRes, err error) {
	return service.Reviews().Audit(ctx, req)
}

func (c *ControllerV1) ListReplies(ctx context.Context, req *v1.ReviewsListRepliesReq) (res *v1.ReviewsListRepliesRes, err error) {
	return service.Reviews().ListReplies(ctx, req)
}

func (c *ControllerV1) CreateReply(ctx context.Context, req *v1.ReviewsCreateReplyReq) (res *v1.ReviewsCreateReplyRes, err error) {
	return service.Reviews().CreateReply(ctx, req)
}

func (c *ControllerV1) DeleteReply(ctx context.Context, req *v1.ReviewsDeleteReplyReq) (res *v1.ReviewsDeleteReplyRes, err error) {
	return service.Reviews().DeleteReply(ctx, req)
}

func (c *ControllerV1) ListAuditLogs(ctx context.Context, req *v1.ReviewsListAuditLogsReq) (res *v1.ReviewsListAuditLogsRes, err error) {
	return service.Reviews().ListAuditLogs(ctx, req)
}
