package category_attributes

import (
	"context"

	"gf-eshop/api/category_attributes/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sCategoryAttributes struct{}

func init() {
	service.RegisterCategoryAttributes(&sCategoryAttributes{})
}

func (s *sCategoryAttributes) List(ctx context.Context, req *v1.CategoryAttributesListReq) (res *v1.CategoryAttributesListRes, err error) {
	m := dao.CategoryAttributes.Ctx(ctx).Where(dao.CategoryAttributes.Columns().CategoryId, req.CategoryId)
	var list []*entity.CategoryAttributes
	err = m.OrderAsc(dao.CategoryAttributes.Columns().SortOrder).Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.CategoryAttributes, 0)
	}
	return &v1.CategoryAttributesListRes{List: list}, nil
}

func (s *sCategoryAttributes) Create(ctx context.Context, req *v1.CategoryAttributesCreateReq) (res *v1.CategoryAttributesCreateRes, err error) {
	result, err := dao.CategoryAttributes.Ctx(ctx).Insert(do.CategoryAttributes{
		CategoryId:      req.CategoryId,
		AttributeId:     req.AttributeId,
		Required:        req.Required,
		IsDefaultFilter: req.IsDefaultFilter,
		SortOrder:       req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.CategoryAttributesCreateRes{Id: id}, nil
}

func (s *sCategoryAttributes) BatchCreate(ctx context.Context, req *v1.CategoryAttributesBatchCreateReq) (res *v1.CategoryAttributesBatchCreateRes, err error) {
	for _, item := range req.Items {
		_, err = dao.CategoryAttributes.Ctx(ctx).Insert(do.CategoryAttributes{
			CategoryId:      req.CategoryId,
			AttributeId:     item.AttributeId,
			Required:        item.Required,
			IsDefaultFilter: item.IsDefaultFilter,
			SortOrder:       item.SortOrder,
		})
		if err != nil {
			return nil, err
		}
	}
	return &v1.CategoryAttributesBatchCreateRes{}, nil
}

func (s *sCategoryAttributes) Delete(ctx context.Context, req *v1.CategoryAttributesDeleteReq) (res *v1.CategoryAttributesDeleteRes, err error) {
	_, err = dao.CategoryAttributes.Ctx(ctx).Where(dao.CategoryAttributes.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.CategoryAttributesDeleteRes{}, nil
}
