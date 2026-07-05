package products

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sProducts struct{}

func init() {
	service.RegisterProducts(&sProducts{})
}

func (s *sProducts) List(ctx context.Context, req *v1.ProductsListReq) (res *v1.ProductsListRes, err error) {
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

	m := dao.Products.Ctx(ctx)
	if req.Name != "" {
		m = m.WhereLike(dao.Products.Columns().Name, "%"+req.Name+"%")
	}
	if req.CategoryId > 0 {
		m = m.Where(dao.Products.Columns().CategoryId, req.CategoryId)
	}
	if req.BrandId > 0 {
		m = m.Where(dao.Products.Columns().BrandId, req.BrandId)
	}
	if req.Status > 0 {
		m = m.Where(dao.Products.Columns().Status, req.Status)
	}
	if req.PriceMin > 0 {
		m = m.WhereGTE(dao.Products.Columns().MinPrice, req.PriceMin)
	}
	if req.PriceMax > 0 {
		m = m.WhereLTE(dao.Products.Columns().MaxPrice, req.PriceMax)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ProductsListRes{
			List:  make([]*entity.Products, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Products
	err = m.Page(page, size).OrderDesc(dao.Products.Columns().SortOrder).OrderDesc(dao.Products.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ProductsListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sProducts) Detail(ctx context.Context, req *v1.ProductsDetailReq) (res *v1.ProductsDetailRes, err error) {
	cached, err := getProductEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.ProductsDetailRes{Products: cached}, nil
	}
	if ctx.Err() != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
	}

	var entity *entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "商品不存在")
	}
	if err := setProductEntityCache(context.Background(), entity); err != nil {
		g.Log().Warning(ctx, "setProductEntityCache failed: %v", err)
	}
	return &v1.ProductsDetailRes{Products: entity}, nil
}

func (s *sProducts) Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error) {
	result, err := dao.Products.Ctx(ctx).Insert(do.Products{
		Name:       req.Name,
		Subtitle:   req.Subtitle,
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
		Unit:       req.Unit,
		MainImage:  req.MainImage,
		Images:     req.Images,
		VideoUrl:   req.VideoUrl,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
		CreatedBy:  req.CreatedBy,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.ProductsCreateRes{Id: id}, nil
}

func (s *sProducts) Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error) {
	count, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "商品不存在")
	}

	_, err = dao.Products.Ctx(ctx).Data(do.Products{
		Name:       req.Name,
		Subtitle:   req.Subtitle,
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
		Unit:       req.Unit,
		MainImage:  req.MainImage,
		Images:     req.Images,
		VideoUrl:   req.VideoUrl,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
		UpdatedBy:  req.UpdatedBy,
	}).Where(dao.Products.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	delProductEntityCache(context.Background(), req.Id)
	return &v1.ProductsUpdateRes{}, nil
}

func (s *sProducts) Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error) {
	_, err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	delProductEntityCache(context.Background(), req.Id)
	return &v1.ProductsDeleteRes{}, nil
}
