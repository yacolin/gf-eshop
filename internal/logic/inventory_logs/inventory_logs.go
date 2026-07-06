package inventoryLogs

import (
	"context"

	"gf-eshop/api/inventory_logs/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sInventoryLogs struct{}

func init() {
	service.RegisterInventoryLogs(&sInventoryLogs{})
}

func (s *sInventoryLogs) List(ctx context.Context, req *v1.InventoryLogsListReq) (res *v1.InventoryLogsListRes, err error) {
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

	m := dao.InventoryLogs.Ctx(ctx)
	if req.SkuId > 0 {
		m = m.Where(dao.InventoryLogs.Columns().SkuId, req.SkuId)
	}
	if req.ChangeType != "" {
		m = m.Where(dao.InventoryLogs.Columns().ChangeType, req.ChangeType)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.InventoryLogsListRes{
			List:  make([]*entity.InventoryLogs, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.InventoryLogs
	err = m.Page(page, size).OrderDesc(dao.InventoryLogs.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.InventoryLogsListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sInventoryLogs) Detail(ctx context.Context, req *v1.InventoryLogsDetailReq) (res *v1.InventoryLogsDetailRes, err error) {
	var entity *entity.InventoryLogs
	err = dao.InventoryLogs.Ctx(ctx).Where(dao.InventoryLogs.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrNotFound
	}
	return &v1.InventoryLogsDetailRes{InventoryLogs: entity}, nil
}
