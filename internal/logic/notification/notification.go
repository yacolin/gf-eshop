package notification

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/notification/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sNotification struct{}

func init() {
	service.RegisterNotification(&sNotification{})
}

const channelInApp = 1
const categorySystem = 1

func (s *sNotification) getStaffId(ctx context.Context) int64 {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.StaffId
}

func (s *sNotification) checkAdmin(ctx context.Context) error {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return errcode.ErrUnauthorized
	}
	isAdmin, err := service.Roles().IsAdmin(ctx, claims.StaffId)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errcode.ErrInsufficientPermissions
	}
	return nil
}

func (s *sNotification) List(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error) {
	staffId := s.getStaffId(ctx)
	if staffId == 0 {
		return nil, errcode.ErrUnauthorized
	}
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	total, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId+" IN (0, ?)", staffId).
		Where(dao.Notifications.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.NotificationListRes{
			List:  make([]*v1.NotificationListItem, 0),
			Total: 0,
		}, nil
	}

	var notifications []*entity.Notifications
	err = dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId+" IN (0, ?)", staffId).
		Where(dao.Notifications.Columns().DeletedAt, nil).
		OrderAsc(dao.Notifications.Columns().Priority).
		OrderDesc(dao.Notifications.Columns().CreatedAt).
		Page(page, size).
		Scan(&notifications)
	if err != nil {
		return nil, err
	}

	readSet, err := s.getReadSet(ctx, staffId, notifications)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.NotificationListItem, len(notifications))
	for i, n := range notifications {
		list[i] = &v1.NotificationListItem{
			Id:          n.Id,
			UserId:      n.UserId,
			Title:       n.Title,
			Content:     n.Content,
			Channel:     n.Channel,
			Category:    n.Category,
			TargetType:  n.TargetType,
			TargetId:    n.TargetId,
			RedirectUrl: n.RedirectUrl,
			IconUrl:     n.IconUrl,
			Priority:    n.Priority,
			IsRead:      readSet[n.Id],
			CreatedBy:   n.CreatedBy,
			CreatedAt:   n.CreatedAt,
		}
	}
	return &v1.NotificationListRes{List: list, Total: total}, nil
}

func (s *sNotification) getReadSet(ctx context.Context, staffId int64, notifications []*entity.Notifications) (map[int64]bool, error) {
	if len(notifications) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(notifications))
	for i, n := range notifications {
		ids[i] = n.Id
	}
	var reads []*entity.NotificationReads
	err := dao.NotificationReads.Ctx(ctx).
		Where(dao.NotificationReads.Columns().UserId, staffId).
		WhereIn(dao.NotificationReads.Columns().NotificationId, ids).
		Scan(&reads)
	if err != nil {
		return nil, err
	}
	readSet := make(map[int64]bool, len(reads))
	for _, r := range reads {
		readSet[r.NotificationId] = true
	}
	return readSet, nil
}

func (s *sNotification) UnreadCount(ctx context.Context, req *v1.NotificationUnreadCountReq) (res *v1.NotificationUnreadCountRes, err error) {
	staffId := s.getStaffId(ctx)
	if staffId == 0 {
		return nil, errcode.ErrUnauthorized
	}
	count, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId+" IN (0, ?)", staffId).
		Where(dao.Notifications.Columns().DeletedAt, nil).
		Where("id NOT IN (SELECT notification_id FROM base_notification_reads WHERE user_id = ?)", staffId).
		Count()
	if err != nil {
		return nil, err
	}
	return &v1.NotificationUnreadCountRes{Count: int64(count)}, nil
}

func (s *sNotification) MarkAsRead(ctx context.Context, req *v1.NotificationMarkAsReadReq) (res *v1.NotificationMarkAsReadRes, err error) {
	staffId := s.getStaffId(ctx)
	if staffId == 0 {
		return nil, errcode.ErrUnauthorized
	}
	_, err = g.DB().Exec(ctx,
		"INSERT IGNORE INTO base_notification_reads (notification_id, user_id, read_at) VALUES (?, ?, ?)",
		req.Id, staffId, gtime.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &v1.NotificationMarkAsReadRes{}, nil
}

func (s *sNotification) MarkAllAsRead(ctx context.Context, req *v1.NotificationMarkAllAsReadReq) (res *v1.NotificationMarkAllAsReadRes, err error) {
	staffId := s.getStaffId(ctx)
	if staffId == 0 {
		return nil, errcode.ErrUnauthorized
	}
	_, err = g.DB().Exec(ctx, `
		INSERT IGNORE INTO base_notification_reads (notification_id, user_id, read_at)
		SELECT n.id, ?, ?
		FROM base_notifications n
		LEFT JOIN base_notification_reads nr ON nr.notification_id = n.id AND nr.user_id = ?
		WHERE n.user_id IN (0, ?)
		  AND nr.id IS NULL
		  AND n.deleted_at IS NULL
	`, staffId, gtime.Now(), staffId, staffId)
	if err != nil {
		return nil, err
	}
	return &v1.NotificationMarkAllAsReadRes{}, nil
}

func (s *sNotification) Delete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error) {
	staffId := s.getStaffId(ctx)
	if staffId == 0 {
		return nil, errcode.ErrUnauthorized
	}
	count, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().Id, req.Id).
		Where(dao.Notifications.Columns().UserId+" IN (0, ?)", staffId).
		Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrNotificationNotFound
	}
	_, err = dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().Id, req.Id).
		Delete()
	if err != nil {
		return nil, err
	}
	return &v1.NotificationDeleteRes{}, nil
}

func (s *sNotification) SendSystem(ctx context.Context, req *v1.NotificationSendSystemReq) (res *v1.NotificationSendSystemRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	title := req.Title
	content := req.Content
	if req.TemplateCode != "" {
		var tmpl *entity.NotificationTemplates
		tmpl, err = s.getTemplateByCode(ctx, req.TemplateCode)
		if err != nil {
			return nil, errcode.ErrNotificationTemplateNotFound
		}
		if title == "" {
			title = tmpl.TitleTemplate
		}
		if content == "" {
			content = tmpl.ContentTemplate
		}
	}
	if title == "" || content == "" {
		return nil, errcode.ErrInvalidParams
	}
	result, err := dao.Notifications.Ctx(ctx).Insert(do.Notifications{
		UserId:    req.UserId,
		Title:     title,
		Content:   content,
		Channel:   channelInApp,
		Category:  categorySystem,
		Priority:  1,
		CreatedBy: s.getStaffId(ctx),
		CreatedAt: gtime.Now(),
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()

	wsMsg := map[string]interface{}{
		"id":       id,
		"title":    title,
		"content":  content,
		"channel":  channelInApp,
		"category": categorySystem,
		"is_read":  false,
	}
	if req.UserId == 0 {
		data, _ := json.Marshal(wsMsg)
		service.WsHub().Broadcast(data)
	} else {
		service.WsHub().PushToUser(req.UserId, "notification", wsMsg)
	}
	return &v1.NotificationSendSystemRes{Id: id}, nil
}

func (s *sNotification) ListTemplates(ctx context.Context, req *v1.NotificationListTemplatesReq) (res *v1.NotificationListTemplatesRes, err error) {
	var list []*entity.NotificationTemplates
	err = dao.NotificationTemplates.Ctx(ctx).
		Where(dao.NotificationTemplates.Columns().Status, 1).
		OrderAsc(dao.NotificationTemplates.Columns().Priority).
		OrderAsc(dao.NotificationTemplates.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.NotificationTemplates, 0)
	}
	items := make([]*v1.NotificationListTemplatesItem, len(list))
	for i, t := range list {
		items[i] = &v1.NotificationListTemplatesItem{
			Id:              t.Id,
			TemplateCode:    t.TemplateCode,
			Channel:         t.Channel,
			TitleTemplate:   t.TitleTemplate,
			ContentTemplate: t.ContentTemplate,
			Category:        t.Category,
			Priority:        t.Priority,
			Status:          t.Status,
		}
	}
	return &v1.NotificationListTemplatesRes{List: items}, nil
}

func (s *sNotification) getTemplateByCode(ctx context.Context, code string) (*entity.NotificationTemplates, error) {
	var t *entity.NotificationTemplates
	err := dao.NotificationTemplates.Ctx(ctx).
		Where(dao.NotificationTemplates.Columns().TemplateCode, code).
		Scan(&t)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errcode.ErrNotificationTemplateNotFound
	}
	return t, nil
}
