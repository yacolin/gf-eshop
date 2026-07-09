package attribute_values

import (
	"context"
	"encoding/json"

	"gf-eshop/api/attribute_values/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sAttributeValues struct{}

func init() {
	service.RegisterAttributeValues(&sAttributeValues{})
}

func (s *sAttributeValues) List(ctx context.Context, req *v1.AttributeValuesListReq) (res *v1.AttributeValuesListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.AttributeValues.Ctx(ctx)
	if req.AttributeId > 0 {
		m = m.Where(dao.AttributeValues.Columns().AttributeId, req.AttributeId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.AttributeValuesListRes{
			List:  make([]*entity.AttributeValues, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.AttributeValues
	err = m.Page(page, size).OrderDesc(dao.AttributeValues.Columns().SearchWeight).OrderAsc(dao.AttributeValues.Columns().SortOrder).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.AttributeValuesListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sAttributeValues) Detail(ctx context.Context, req *v1.AttributeValuesDetailReq) (res *v1.AttributeValuesDetailRes, err error) {
	var ent *entity.AttributeValues
	err = dao.AttributeValues.Ctx(ctx).Where(dao.AttributeValues.Columns().Id, req.Id).Scan(&ent)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, errcode.ErrAttributeValueNotFound
	}
	return &v1.AttributeValuesDetailRes{AttributeValues: ent}, nil
}

func (s *sAttributeValues) Create(ctx context.Context, req *v1.AttributeValuesCreateReq) (res *v1.AttributeValuesCreateRes, err error) {
	status := req.Status
	if status <= 0 {
		status = 1
	}
	data := do.AttributeValues{
		AttributeId:  req.AttributeId,
		Value:        req.Value,
		SearchWeight: req.SearchWeight,
		NumericValue: req.NumericValue,
		ColorHex:     req.ColorHex,
		SortOrder:    req.SortOrder,
		Status:       status,
	}
	if len(req.Alias) > 0 {
		aliasJSON, _ := json.Marshal(req.Alias)
		data.Alias = string(aliasJSON)
	}

	result, err := dao.AttributeValues.Ctx(ctx).Insert(data)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.AttributeValuesCreateRes{Id: id}, nil
}

func (s *sAttributeValues) Update(ctx context.Context, req *v1.AttributeValuesUpdateReq) (res *v1.AttributeValuesUpdateRes, err error) {
	count, err := dao.AttributeValues.Ctx(ctx).Where(dao.AttributeValues.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrAttributeValueNotFound
	}

	data := do.AttributeValues{
		Value:        req.Value,
		SearchWeight: req.SearchWeight,
		NumericValue: req.NumericValue,
		ColorHex:     req.ColorHex,
		SortOrder:    req.SortOrder,
	}
	if req.Status > 0 {
		data.Status = req.Status
	}
	if len(req.Alias) > 0 {
		aliasJSON, _ := json.Marshal(req.Alias)
		data.Alias = string(aliasJSON)
	}

	_, err = dao.AttributeValues.Ctx(ctx).Data(data).Where(dao.AttributeValues.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.AttributeValuesUpdateRes{}, nil
}

func (s *sAttributeValues) Delete(ctx context.Context, req *v1.AttributeValuesDeleteReq) (res *v1.AttributeValuesDeleteRes, err error) {
	_, err = dao.AttributeValues.Ctx(ctx).Where(dao.AttributeValues.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.AttributeValuesDeleteRes{}, nil
}
