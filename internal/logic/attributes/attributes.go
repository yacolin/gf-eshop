package attributes

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/attributes/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sAttributes struct{}

func init() {
	service.RegisterAttributes(&sAttributes{})
}

func (s *sAttributes) List(ctx context.Context, req *v1.AttributesListReq) (res *v1.AttributesListRes, err error) {
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

	m := dao.Attributes.Ctx(ctx)
	if req.CategoryId > 0 {
		m = m.Where(dao.Attributes.Columns().CategoryId, req.CategoryId)
	}
	if req.Searchable > 0 {
		m = m.Where(dao.Attributes.Columns().Searchable, req.Searchable)
	}
	if req.IsSkuSpec > 0 {
		m = m.Where(dao.Attributes.Columns().IsSkuSpec, req.IsSkuSpec)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.AttributesListRes{
			List:  make([]*entity.Attributes, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Attributes
	err = m.Page(page, size).OrderAsc(dao.Attributes.Columns().SortOrder).OrderAsc(dao.Attributes.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.AttributesListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sAttributes) Detail(ctx context.Context, req *v1.AttributesDetailReq) (res *v1.AttributesDetailRes, err error) {
	var entity *entity.Attributes
	err = dao.Attributes.Ctx(ctx).Where(dao.Attributes.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "属性不存在")
	}
	return &v1.AttributesDetailRes{Attributes: entity}, nil
}

func (s *sAttributes) Create(ctx context.Context, req *v1.AttributesCreateReq) (res *v1.AttributesCreateRes, err error) {
	result, err := dao.Attributes.Ctx(ctx).Insert(do.Attributes{
		Name:       req.Name,
		CategoryId: req.CategoryId,
		InputType:  req.InputType,
		Values:     req.Values,
		Unit:       req.Unit,
		Required:   req.Required,
		Searchable: req.Searchable,
		IsSkuSpec:  req.IsSkuSpec,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.AttributesCreateRes{Id: id}, nil
}

func (s *sAttributes) Update(ctx context.Context, req *v1.AttributesUpdateReq) (res *v1.AttributesUpdateRes, err error) {
	count, err := dao.Attributes.Ctx(ctx).Where(dao.Attributes.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "属性不存在")
	}

	_, err = dao.Attributes.Ctx(ctx).Data(do.Attributes{
		Name:       req.Name,
		CategoryId: req.CategoryId,
		InputType:  req.InputType,
		Values:     req.Values,
		Unit:       req.Unit,
		Required:   req.Required,
		Searchable: req.Searchable,
		IsSkuSpec:  req.IsSkuSpec,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
	}).Where(dao.Attributes.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.AttributesUpdateRes{}, nil
}

func (s *sAttributes) Delete(ctx context.Context, req *v1.AttributesDeleteReq) (res *v1.AttributesDeleteRes, err error) {
	_, err = dao.Attributes.Ctx(ctx).Where(dao.Attributes.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.AttributesDeleteRes{}, nil
}
