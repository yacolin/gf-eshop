package merchants

import (
	"context"

	"gf-eshop/api/merchant_qualifications/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantQualifications struct{}

func init() {
	service.RegisterMerchantQualifications(&sMerchantQualifications{})
}

func (s *sMerchantQualifications) List(ctx context.Context, req *v1.MerchantQualificationsListReq) (res *v1.MerchantQualificationsListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.MerchantQualifications.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.MerchantQualifications.Columns().MerchantId, req.MerchantId)
	}
	if req.Status > 0 {
		m = m.Where(dao.MerchantQualifications.Columns().Status, req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.MerchantQualificationsListRes{
			List:  make([]*entity.MerchantQualifications, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.MerchantQualifications
	err = m.Page(page, size).OrderDesc(dao.MerchantQualifications.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantQualificationsListRes{List: list, Total: total}, nil
}

func (s *sMerchantQualifications) Detail(ctx context.Context, req *v1.MerchantQualificationsDetailReq) (res *v1.MerchantQualificationsDetailRes, err error) {
	var entity *entity.MerchantQualifications
	err = dao.MerchantQualifications.Ctx(ctx).Where(dao.MerchantQualifications.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrMerchantQualificationNotFound
	}
	return &v1.MerchantQualificationsDetailRes{MerchantQualifications: entity}, nil
}

func (s *sMerchantQualifications) Create(ctx context.Context, req *v1.MerchantQualificationsCreateReq) (res *v1.MerchantQualificationsCreateRes, err error) {
	result, err := dao.MerchantQualifications.Ctx(ctx).Insert(do.MerchantQualifications{
		MerchantId:        req.MerchantId,
		QualificationType: req.QualificationType,
		QualificationName: req.QualificationName,
		FileUrl:           req.FileUrl,
		ExpireAt:          req.ExpireAt,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MerchantQualificationsCreateRes{Id: id}, nil
}

func (s *sMerchantQualifications) Update(ctx context.Context, req *v1.MerchantQualificationsUpdateReq) (res *v1.MerchantQualificationsUpdateRes, err error) {
	count, err := dao.MerchantQualifications.Ctx(ctx).Where(dao.MerchantQualifications.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantQualificationNotFound
	}
	_, err = dao.MerchantQualifications.Ctx(ctx).Data(do.MerchantQualifications{
		QualificationName: req.QualificationName,
		FileUrl:           req.FileUrl,
		ExpireAt:          req.ExpireAt,
	}).Where(dao.MerchantQualifications.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantQualificationsUpdateRes{}, nil
}

func (s *sMerchantQualifications) Delete(ctx context.Context, req *v1.MerchantQualificationsDeleteReq) (res *v1.MerchantQualificationsDeleteRes, err error) {
	_, err = dao.MerchantQualifications.Ctx(ctx).Where(dao.MerchantQualifications.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantQualificationsDeleteRes{}, nil
}

func (s *sMerchantQualifications) Audit(ctx context.Context, req *v1.MerchantQualificationsAuditReq) (res *v1.MerchantQualificationsAuditRes, err error) {
	count, err := dao.MerchantQualifications.Ctx(ctx).Where(dao.MerchantQualifications.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantQualificationNotFound
	}
	_, err = dao.MerchantQualifications.Ctx(ctx).Data(do.MerchantQualifications{
		Status:      req.Status,
		AuditRemark: req.AuditRemark,
	}).Where(dao.MerchantQualifications.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantQualificationsAuditRes{}, nil
}
