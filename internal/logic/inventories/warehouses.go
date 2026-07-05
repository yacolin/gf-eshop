package inventories

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/inventories/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sWarehouses struct{}

func init() {
	service.RegisterWarehouses(&sWarehouses{})
}

func (s *sWarehouses) List(ctx context.Context, req *v1.WarehousesListReq) (res *v1.WarehousesListRes, err error) {
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

	total, err := dao.Warehouses.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.WarehousesListRes{
			List:  make([]*entity.Warehouses, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Warehouses
	err = dao.Warehouses.Ctx(ctx).Page(page, size).OrderAsc(dao.Warehouses.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.WarehousesListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sWarehouses) Detail(ctx context.Context, req *v1.WarehousesDetailReq) (res *v1.WarehousesDetailRes, err error) {
	var entity *entity.Warehouses
	err = dao.Warehouses.Ctx(ctx).Where(dao.Warehouses.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "仓库不存在")
	}
	return &v1.WarehousesDetailRes{Warehouses: entity}, nil
}

func (s *sWarehouses) Create(ctx context.Context, req *v1.WarehousesCreateReq) (res *v1.WarehousesCreateRes, err error) {
	result, err := dao.Warehouses.Ctx(ctx).Insert(do.Warehouses{
		WarehouseName: req.WarehouseName,
		WarehouseType: req.WarehouseType,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		Address:       req.Address,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.WarehousesCreateRes{Id: id}, nil
}

func (s *sWarehouses) Update(ctx context.Context, req *v1.WarehousesUpdateReq) (res *v1.WarehousesUpdateRes, err error) {
	count, err := dao.Warehouses.Ctx(ctx).Where(dao.Warehouses.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "仓库不存在")
	}

	_, err = dao.Warehouses.Ctx(ctx).Data(do.Warehouses{
		WarehouseName: req.WarehouseName,
		WarehouseType: req.WarehouseType,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		Address:       req.Address,
		Status:        req.Status,
	}).Where(dao.Warehouses.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.WarehousesUpdateRes{}, nil
}

func (s *sWarehouses) Delete(ctx context.Context, req *v1.WarehousesDeleteReq) (res *v1.WarehousesDeleteRes, err error) {
	_, err = dao.Warehouses.Ctx(ctx).Where(dao.Warehouses.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.WarehousesDeleteRes{}, nil
}
