package productAttributes

import (
	"context"

	"gf-eshop/api/product_attributes/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sProductAttributes struct{}

func init() {
	service.RegisterProductAttributes(&sProductAttributes{})
}

func (s *sProductAttributes) List(ctx context.Context, req *v1.ProductAttributesListReq) (res *v1.ProductAttributesListRes, err error) {
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

	m := dao.ProductAttributes.Ctx(ctx)
	if req.ProductId > 0 {
		m = m.Where(dao.ProductAttributes.Columns().ProductId, req.ProductId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ProductAttributesListRes{
			List:  make([]*entity.ProductAttributes, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.ProductAttributes
	err = m.Page(page, size).OrderAsc(dao.ProductAttributes.Columns().SortOrder).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ProductAttributesListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sProductAttributes) Create(ctx context.Context, req *v1.ProductAttributesCreateReq) (res *v1.ProductAttributesCreateRes, err error) {
	result, err := dao.ProductAttributes.Ctx(ctx).Insert(do.ProductAttributes{
		ProductId:   req.ProductId,
		AttributeId: req.AttributeId,
		Value:       req.Value,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.ProductAttributesCreateRes{Id: id}, nil
}

func (s *sProductAttributes) Delete(ctx context.Context, req *v1.ProductAttributesDeleteReq) (res *v1.ProductAttributesDeleteRes, err error) {
	count, err := dao.ProductAttributes.Ctx(ctx).Where(dao.ProductAttributes.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrAttributeNotFound
	}

	_, err = dao.ProductAttributes.Ctx(ctx).Where(dao.ProductAttributes.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.ProductAttributesDeleteRes{}, nil
}
