package notification

import (
	"context"

	"gf-eshop/api/notification/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error) {
	return service.Notification().List(ctx, req)
}

func (c *ControllerV1) UnreadCount(ctx context.Context, req *v1.NotificationUnreadCountReq) (res *v1.NotificationUnreadCountRes, err error) {
	return service.Notification().UnreadCount(ctx, req)
}

func (c *ControllerV1) MarkAsRead(ctx context.Context, req *v1.NotificationMarkAsReadReq) (res *v1.NotificationMarkAsReadRes, err error) {
	return service.Notification().MarkAsRead(ctx, req)
}

func (c *ControllerV1) MarkAllAsRead(ctx context.Context, req *v1.NotificationMarkAllAsReadReq) (res *v1.NotificationMarkAllAsReadRes, err error) {
	return service.Notification().MarkAllAsRead(ctx, req)
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error) {
	return service.Notification().Delete(ctx, req)
}

func (c *ControllerV1) SendSystem(ctx context.Context, req *v1.NotificationSendSystemReq) (res *v1.NotificationSendSystemRes, err error) {
	return service.Notification().SendSystem(ctx, req)
}

func (c *ControllerV1) ListTemplates(ctx context.Context, req *v1.NotificationListTemplatesReq) (res *v1.NotificationListTemplatesRes, err error) {
	return service.Notification().ListTemplates(ctx, req)
}
