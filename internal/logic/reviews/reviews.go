package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/reviews/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

func generateReviewNo() string {
	return fmt.Sprintf("REV%d%04d", time.Now().UnixMilli(), rand.Intn(10000))
}

type sReviews struct{}

func init() {
	service.RegisterReviews(&sReviews{})
}

func (s *sReviews) List(ctx context.Context, req *v1.ReviewsListReq) (res *v1.ReviewsListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.Reviews.Ctx(ctx)
	if req.Id > 0 {
		m = m.Where(dao.Reviews.Columns().Id, req.Id)
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.Reviews.Columns().MerchantId, req.MerchantId)
	}
	if req.SpuId > 0 {
		m = m.Where(dao.Reviews.Columns().SpuId, req.SpuId)
	}
	if req.Status != nil {
		m = m.Where(dao.Reviews.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ReviewsListRes{
			List:  make([]*entity.Reviews, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Reviews
	err = m.Page(page, size).OrderDesc(dao.Reviews.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ReviewsListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sReviews) Detail(ctx context.Context, req *v1.ReviewsDetailReq) (res *v1.ReviewsDetailRes, err error) {
	var entity *entity.Reviews
	err = dao.Reviews.Ctx(ctx).Where(dao.Reviews.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrReviewNotFound
	}
	return &v1.ReviewsDetailRes{Reviews: entity}, nil
}

func (s *sReviews) Create(ctx context.Context, req *v1.ReviewsCreateReq) (res *v1.ReviewsCreateRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").Int64()
	contentLen := len([]rune(req.Content))
	if contentLen > 32767 {
		contentLen = 32767
	}

	result, err := dao.Reviews.Ctx(ctx).Insert(do.Reviews{
		UserId:          userID,
		ReviewNo:        generateReviewNo(),
		OrderId:         req.OrderId,
		OrderItemId:     req.OrderItemId,
		SpuId:           req.SpuId,
		SkuId:           req.SkuId,
		OverallRating:   req.OverallRating,
		QualityRating:   req.QualityRating,
		LogisticsRating: req.LogisticsRating,
		ServiceRating:   req.ServiceRating,
		Content:         req.Content,
		ContentLength:   contentLen,
		IsAnonymous:     req.IsAnonymous,
		Status:          0,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()

	created := entity.Reviews{Id: id, Status: 0}
	snapshot, _ := json.Marshal(created)

	_, err = dao.ReviewAuditLogs.Ctx(ctx).Insert(do.ReviewAuditLogs{
		ReviewId:     id,
		Action:       "submit",
		OperatorId:   userID,
		AfterStatus:  0,
		Snapshot:     string(snapshot),
	})
	if err != nil {
		return nil, err
	}

	return &v1.ReviewsCreateRes{Id: id}, nil
}

func (s *sReviews) Update(ctx context.Context, req *v1.ReviewsUpdateReq) (res *v1.ReviewsUpdateRes, err error) {
	count, err := dao.Reviews.Ctx(ctx).Where(dao.Reviews.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrReviewNotFound
	}

	_, err = dao.Reviews.Ctx(ctx).Data(do.Reviews{
		OverallRating:   req.OverallRating,
		QualityRating:   req.QualityRating,
		LogisticsRating: req.LogisticsRating,
		ServiceRating:   req.ServiceRating,
		Content:         req.Content,
		IsAnonymous:     req.IsAnonymous,
	}).Where(dao.Reviews.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.ReviewsUpdateRes{}, nil
}

func (s *sReviews) Delete(ctx context.Context, req *v1.ReviewsDeleteReq) (res *v1.ReviewsDeleteRes, err error) {
	var review entity.Reviews
	err = dao.Reviews.Ctx(ctx).Where(dao.Reviews.Columns().Id, req.Id).Scan(&review)
	if err != nil {
		return nil, err
	}
	if review.Id == 0 {
		return nil, errcode.ErrReviewNotFound
	}

	claims := utility.GetUserClaims(ctx)
	operatorID := int64(0)
	operatorName := ""
	if claims != nil {
		operatorID = claims.UserId
		operatorName = claims.Username
	}

	_, err = dao.Reviews.Ctx(ctx).Data(do.Reviews{
		Status: 3,
	}).Where(dao.Reviews.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}

	snapshot, _ := json.Marshal(review)

	_, err = dao.ReviewAuditLogs.Ctx(ctx).Insert(do.ReviewAuditLogs{
		ReviewId:     req.Id,
		Action:       "delete",
		OperatorId:   operatorID,
		OperatorName: operatorName,
		BeforeStatus: review.Status,
		AfterStatus:  3,
		Snapshot:     string(snapshot),
	})
	if err != nil {
		return nil, err
	}

	return &v1.ReviewsDeleteRes{}, nil
}

func (s *sReviews) Audit(ctx context.Context, req *v1.ReviewsAuditReq) (res *v1.ReviewsAuditRes, err error) {
	var review entity.Reviews
	err = dao.Reviews.Ctx(ctx).Where(dao.Reviews.Columns().Id, req.Id).Scan(&review)
	if err != nil {
		return nil, err
	}
	if review.Id == 0 {
		return nil, errcode.ErrReviewNotFound
	}

	claims := utility.GetUserClaims(ctx)
	operatorID := int64(0)
	operatorName := ""
	if claims != nil {
		operatorID = claims.UserId
		operatorName = claims.Username
	}

	_, err = dao.Reviews.Ctx(ctx).Data(do.Reviews{
		Status:       req.Status,
		RejectReason: req.RejectReason,
		AuditedBy:    operatorID,
		AuditedAt:    gtime.Now(),
	}).Where(dao.Reviews.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}

	action := "approve"
	if req.Status == 2 {
		action = "reject"
	}

	snapshot, _ := json.Marshal(review)

	_, err = dao.ReviewAuditLogs.Ctx(ctx).Insert(do.ReviewAuditLogs{
		ReviewId:     req.Id,
		Action:       action,
		OperatorId:   operatorID,
		OperatorName: operatorName,
		BeforeStatus: review.Status,
		AfterStatus:  req.Status,
		Remark:       req.RejectReason,
		Snapshot:     string(snapshot),
	})
	if err != nil {
		return nil, err
	}

	return &v1.ReviewsAuditRes{}, nil
}

func (s *sReviews) ListReplies(ctx context.Context, req *v1.ReviewsListRepliesReq) (res *v1.ReviewsListRepliesRes, err error) {
	var list []*entity.ReviewReplies
	err = dao.ReviewReplies.Ctx(ctx).
		Where(dao.ReviewReplies.Columns().ReviewId, req.Id).
		OrderAsc(dao.ReviewReplies.Columns().CreatedAt).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.ReviewReplies, 0)
	}
	return &v1.ReviewsListRepliesRes{List: list}, nil
}

func (s *sReviews) CreateReply(ctx context.Context, req *v1.ReviewsCreateReplyReq) (res *v1.ReviewsCreateReplyRes, err error) {
	claims := utility.GetUserClaims(ctx)
	operatorID := int64(0)
	operatorName := ""
	if claims != nil {
		operatorID = claims.UserId
		operatorName = claims.Username
	}

	rootReplyId := int64(0)
	if req.ParentId > 0 {
		var parent entity.ReviewReplies
		err := dao.ReviewReplies.Ctx(ctx).Where(dao.ReviewReplies.Columns().Id, req.ParentId).Scan(&parent)
		if err == nil && parent.Id > 0 {
			if parent.RootReplyId > 0 {
				rootReplyId = parent.RootReplyId
			} else {
				rootReplyId = parent.Id
			}
		}
	}

	result, err := dao.ReviewReplies.Ctx(ctx).Insert(do.ReviewReplies{
		ReviewId:     req.Id,
		RootReplyId:  rootReplyId,
		ParentId:     req.ParentId,
		ReplyType:    req.ReplyType,
		Content:      req.Content,
		OperatorId:   operatorID,
		OperatorName: operatorName,
		Status:       1,
	})
	if err != nil {
		return nil, err
	}
	replyId, _ := result.LastInsertId()

	_, err = dao.Reviews.Ctx(ctx).
		Data(g.Map{
			dao.Reviews.Columns().ReplyCount: gdb.Raw("reply_count + 1"),
		}).
		Where(dao.Reviews.Columns().Id, req.Id).
		Update()
	if err != nil {
		return nil, err
	}

	return &v1.ReviewsCreateReplyRes{Id: replyId}, nil
}

func (s *sReviews) DeleteReply(ctx context.Context, req *v1.ReviewsDeleteReplyReq) (res *v1.ReviewsDeleteReplyRes, err error) {
	_, err = dao.ReviewReplies.Ctx(ctx).
		Where(dao.ReviewReplies.Columns().Id, req.ReplyId).
		Where(dao.ReviewReplies.Columns().ReviewId, req.Id).
		Delete()
	if err != nil {
		return nil, err
	}
	return &v1.ReviewsDeleteReplyRes{}, nil
}

func (s *sReviews) ListAuditLogs(ctx context.Context, req *v1.ReviewsListAuditLogsReq) (res *v1.ReviewsListAuditLogsRes, err error) {
	var list []*entity.ReviewAuditLogs
	err = dao.ReviewAuditLogs.Ctx(ctx).
		Where(dao.ReviewAuditLogs.Columns().ReviewId, req.Id).
		OrderDesc(dao.ReviewAuditLogs.Columns().CreatedAt).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.ReviewAuditLogs, 0)
	}
	return &v1.ReviewsListAuditLogsRes{List: list}, nil
}
