package productDescriptions

import (
	"context"

	"gf-eshop/api/product_descriptions/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sProductDescriptions struct{}

func init() {
	service.RegisterProductDescriptions(&sProductDescriptions{})
}

func (s *sProductDescriptions) Detail(ctx context.Context, req *v1.ProductDescriptionsDetailReq) (res *v1.ProductDescriptionsDetailRes, err error) {
	var entity *entity.ProductDescriptions
	err = dao.ProductDescriptions.Ctx(ctx).Where(dao.ProductDescriptions.Columns().ProductId, req.ProductId).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrNotFound
	}
	return &v1.ProductDescriptionsDetailRes{ProductDescriptions: entity}, nil
}

func (s *sProductDescriptions) Save(ctx context.Context, req *v1.ProductDescriptionsSaveReq) (res *v1.ProductDescriptionsSaveRes, err error) {
	count, err := dao.ProductDescriptions.Ctx(ctx).Where(dao.ProductDescriptions.Columns().ProductId, req.ProductId).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		_, err = dao.ProductDescriptions.Ctx(ctx).Data(do.ProductDescriptions{
			Description:       req.Description,
			MobileDescription: req.MobileDescription,
		}).Where(dao.ProductDescriptions.Columns().ProductId, req.ProductId).Update()
	} else {
		_, err = dao.ProductDescriptions.Ctx(ctx).Insert(do.ProductDescriptions{
			ProductId:         req.ProductId,
			Description:       req.Description,
			MobileDescription: req.MobileDescription,
		})
	}
	if err != nil {
		return nil, err
	}
	return &v1.ProductDescriptionsSaveRes{}, nil
}
