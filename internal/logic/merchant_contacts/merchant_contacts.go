package merchantContacts

import (
	"context"

	"gf-eshop/api/merchant_contacts/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantContacts struct{}

func init() {
	service.RegisterMerchantContacts(&sMerchantContacts{})
}

func (s *sMerchantContacts) List(ctx context.Context, req *v1.MerchantContactsListReq) (res *v1.MerchantContactsListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.MerchantContacts.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.MerchantContacts.Columns().MerchantId, req.MerchantId)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.MerchantContactsListRes{
			List:  make([]*entity.MerchantContacts, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.MerchantContacts
	err = m.Page(page, size).OrderDesc(dao.MerchantContacts.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantContactsListRes{List: list, Total: total}, nil
}

func (s *sMerchantContacts) Detail(ctx context.Context, req *v1.MerchantContactsDetailReq) (res *v1.MerchantContactsDetailRes, err error) {
	var entity *entity.MerchantContacts
	err = dao.MerchantContacts.Ctx(ctx).Where(dao.MerchantContacts.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrMerchantContactNotFound
	}
	return &v1.MerchantContactsDetailRes{MerchantContacts: entity}, nil
}

func (s *sMerchantContacts) Create(ctx context.Context, req *v1.MerchantContactsCreateReq) (res *v1.MerchantContactsCreateRes, err error) {
	result, err := dao.MerchantContacts.Ctx(ctx).Insert(do.MerchantContacts{
		MerchantId:   req.MerchantId,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactRole:  req.ContactRole,
		IsPrimary:    req.IsPrimary,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MerchantContactsCreateRes{Id: id}, nil
}

func (s *sMerchantContacts) Update(ctx context.Context, req *v1.MerchantContactsUpdateReq) (res *v1.MerchantContactsUpdateRes, err error) {
	count, err := dao.MerchantContacts.Ctx(ctx).Where(dao.MerchantContacts.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantContactNotFound
	}
	_, err = dao.MerchantContacts.Ctx(ctx).Data(do.MerchantContacts{
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactRole:  req.ContactRole,
		IsPrimary:    req.IsPrimary,
	}).Where(dao.MerchantContacts.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantContactsUpdateRes{}, nil
}

func (s *sMerchantContacts) Delete(ctx context.Context, req *v1.MerchantContactsDeleteReq) (res *v1.MerchantContactsDeleteRes, err error) {
	_, err = dao.MerchantContacts.Ctx(ctx).Where(dao.MerchantContacts.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantContactsDeleteRes{}, nil
}
