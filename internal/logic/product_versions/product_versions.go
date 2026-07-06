package productVersions

import (
	"context"

	"gf-eshop/api/product_versions/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sProductVersions struct{}

func init() {
	service.RegisterProductVersions(&sProductVersions{})
}

func (s *sProductVersions) List(ctx context.Context, req *v1.ProductVersionsListReq) (res *v1.ProductVersionsListRes, err error) {
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

	m := dao.ProductVersions.Ctx(ctx)
	if req.ProductId > 0 {
		m = m.Where(dao.ProductVersions.Columns().ProductId, req.ProductId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ProductVersionsListRes{
			List:  make([]*entity.ProductVersions, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.ProductVersions
	err = m.Page(page, size).OrderDesc(dao.ProductVersions.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ProductVersionsListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sProductVersions) Detail(ctx context.Context, req *v1.ProductVersionsDetailReq) (res *v1.ProductVersionsDetailRes, err error) {
	var entity *entity.ProductVersions
	err = dao.ProductVersions.Ctx(ctx).Where(dao.ProductVersions.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrNotFound
	}
	return &v1.ProductVersionsDetailRes{ProductVersions: entity}, nil
}
