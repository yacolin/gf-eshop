package notification

import (
	"context"

	"gf-eshop/api/notification/v1"
)

type INotificationV1 interface {
	List(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error)
	UnreadCount(ctx context.Context, req *v1.NotificationUnreadCountReq) (res *v1.NotificationUnreadCountRes, err error)
	MarkAsRead(ctx context.Context, req *v1.NotificationMarkAsReadReq) (res *v1.NotificationMarkAsReadRes, err error)
	MarkAllAsRead(ctx context.Context, req *v1.NotificationMarkAllAsReadReq) (res *v1.NotificationMarkAllAsReadRes, err error)
	Delete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error)
	SendSystem(ctx context.Context, req *v1.NotificationSendSystemReq) (res *v1.NotificationSendSystemRes, err error)
	ListTemplates(ctx context.Context, req *v1.NotificationListTemplatesReq) (res *v1.NotificationListTemplatesRes, err error)
}
