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

type sInventories struct{}

func init() {
	service.RegisterInventories(&sInventories{})
}

func (s *sInventories) List(ctx context.Context, req *v1.InventoriesListReq) (res *v1.InventoriesListRes, err error) {
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

	m := dao.Inventories.Ctx(ctx)
	if req.SkuId > 0 {
		m = m.Where(dao.Inventories.Columns().SkuId, req.SkuId)
	}
	if req.WarehouseId > 0 {
		m = m.Where(dao.Inventories.Columns().WarehouseId, req.WarehouseId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.InventoriesListRes{
			List:  make([]*entity.Inventories, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Inventories
	err = m.Page(page, size).OrderDesc(dao.Inventories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sInventories) Detail(ctx context.Context, req *v1.InventoriesDetailReq) (res *v1.InventoriesDetailRes, err error) {
	var entity *entity.Inventories
	err = dao.Inventories.Ctx(ctx).Where(dao.Inventories.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "库存记录不存在")
	}
	return &v1.InventoriesDetailRes{Inventories: entity}, nil
}

func (s *sInventories) Create(ctx context.Context, req *v1.InventoriesCreateReq) (res *v1.InventoriesCreateRes, err error) {
	result, err := dao.Inventories.Ctx(ctx).Insert(do.Inventories{
		SkuId:       req.SkuId,
		WarehouseId: req.WarehouseId,
		Quantity:    req.Quantity,
		Threshold:   req.Threshold,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.InventoriesCreateRes{Id: id}, nil
}

func (s *sInventories) Update(ctx context.Context, req *v1.InventoriesUpdateReq) (res *v1.InventoriesUpdateRes, err error) {
	count, err := dao.Inventories.Ctx(ctx).Where(dao.Inventories.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "库存记录不存在")
	}

	_, err = dao.Inventories.Ctx(ctx).Data(do.Inventories{
		Quantity:    req.Quantity,
		Reserved:    req.Reserved,
		Threshold:   req.Threshold,
		MaxThreshold: req.MaxThreshold,
	}).Where(dao.Inventories.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesUpdateRes{}, nil
}

func (s *sInventories) Delete(ctx context.Context, req *v1.InventoriesDeleteReq) (res *v1.InventoriesDeleteRes, err error) {
	_, err = dao.Inventories.Ctx(ctx).Where(dao.Inventories.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesDeleteRes{}, nil
}
