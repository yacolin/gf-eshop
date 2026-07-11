package inventories

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"strconv"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/inventories/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
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
	if req.Status != nil {
		m = m.Where(dao.Inventories.Columns().Status, *req.Status)
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
		return nil, errcode.ErrInventoryNotFound
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
		return nil, errcode.ErrInventoryNotFound
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

// ── 库存业务方法 ──────────────────────────────────────────

// Lock 下单预占库存（行锁防超卖）
func (s *sInventories) Lock(ctx context.Context, req *v1.InventoriesLockReq) (res *v1.InventoriesLockRes, err error) {
	var inv *entity.Inventories
	err = dao.Inventories.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		inv, err = findOrCreateWithTx(ctx, tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		available := inv.Quantity - inv.Reserved
		if available < req.Quantity {
			return errcode.ErrInsufficientStock
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Reserved += req.Quantity

		if err := updateInvFields(ctx, tx, inv); err != nil {
			return err
		}
		return createLog(ctx, tx, &do.InventoryLogs{
			SkuId:          req.SkuID,
			ChangeType:     "order_lock",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   req.Quantity,
			ReferenceId:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesLockRes{Inventories: inv}, nil
}

// Unlock 取消释放预占
func (s *sInventories) Unlock(ctx context.Context, req *v1.InventoriesUnlockReq) (res *v1.InventoriesUnlockRes, err error) {
	var inv *entity.Inventories
	err = dao.Inventories.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		inv, err = findOrCreateWithTx(ctx, tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		if inv.Reserved < req.Quantity {
			return errcode.ErrInvalidStockChange
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Reserved -= req.Quantity

		if err := updateInvFields(ctx, tx, inv); err != nil {
			return err
		}
		return createLog(ctx, tx, &do.InventoryLogs{
			SkuId:          req.SkuID,
			ChangeType:     "order_unlock",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   -req.Quantity,
			ReferenceId:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesUnlockRes{Inventories: inv}, nil
}

// Deduct 支付扣减库存
func (s *sInventories) Deduct(ctx context.Context, req *v1.InventoriesDeductReq) (res *v1.InventoriesDeductRes, err error) {
	var inv *entity.Inventories
	err = dao.Inventories.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		inv, err = findOrCreateWithTx(ctx, tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		if inv.Quantity < req.Quantity || inv.Reserved < req.Quantity {
			return errcode.ErrInsufficientStock
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Quantity -= req.Quantity
		inv.Reserved -= req.Quantity

		if err := updateInvFields(ctx, tx, inv); err != nil {
			return err
		}
		return createLog(ctx, tx, &do.InventoryLogs{
			SkuId:          req.SkuID,
			ChangeType:     "order_deduct",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   -req.Quantity,
			ReferenceId:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesDeductRes{Inventories: inv}, nil
}

// Restock 入库/补货
func (s *sInventories) Restock(ctx context.Context, req *v1.InventoriesRestockReq) (res *v1.InventoriesRestockRes, err error) {
	var inv *entity.Inventories
	err = dao.Inventories.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		inv, err = findOrCreateWithTx(ctx, tx, req.SkuID, req.WarehouseID)
		if err != nil {
			return err
		}
		beforeQty := inv.Quantity
		inv.Quantity += req.Quantity
		if inv.MaxThreshold > 0 && inv.Quantity > inv.MaxThreshold {
			inv.Quantity = inv.MaxThreshold
		}

		if err := updateInvFields(ctx, tx, inv); err != nil {
			return err
		}
		return createLog(ctx, tx, &do.InventoryLogs{
			SkuId:          req.SkuID,
			WarehouseId:    req.WarehouseID,
			ChangeType:     "inbound",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			ChangeAmount:   req.Quantity,
			ReferenceId:    req.ReferenceID,
			Operator:       req.Operator,
			Note:           req.Note,
		})
	})
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesRestockRes{Inventories: inv}, nil
}

// GetStock 按 SKU 查询库存
func (s *sInventories) GetStock(ctx context.Context, req *v1.InventoriesGetStockReq) (res *v1.InventoriesGetStockRes, err error) {
	var inv *entity.Inventories
	err = dao.Inventories.Ctx(ctx).
		Where(dao.Inventories.Columns().SkuId, req.SkuID).
		Scan(&inv)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, errcode.ErrInventoryNotFound
	}
	return &v1.InventoriesGetStockRes{Inventories: inv}, nil
}

// Alerts 库存预警列表（低于安全阈值）
func (s *sInventories) Alerts(ctx context.Context, req *v1.InventoriesAlertsReq) (res *v1.InventoriesAlertsRes, err error) {
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

	m := dao.Inventories.Ctx(ctx).
		Where(dao.Inventories.Columns().Quantity+" <= "+dao.Inventories.Columns().Threshold)
	if req.Status != nil {
		m = m.Where(dao.Inventories.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.InventoriesAlertsRes{
			List:  make([]*entity.Inventories, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Inventories
	err = m.Page(page, size).OrderDesc(dao.Inventories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesAlertsRes{
		List:  list,
		Total: total,
	}, nil
}

// AlertResolve 确认预警/忽略
func (s *sInventories) AlertResolve(ctx context.Context, req *v1.InventoriesAlertResolveReq) (res *v1.InventoriesAlertResolveRes, err error) {
	count, err := dao.Inventories.Ctx(ctx).Where(dao.Inventories.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrInventoryNotFound
	}
	return &v1.InventoriesAlertResolveRes{}, nil
}

// Export 导出库存列表 CSV
func (s *sInventories) Export(ctx context.Context, req *v1.InventoriesExportReq) (res *v1.InventoriesExportRes, err error) {
	m := dao.Inventories.Ctx(ctx)
	if req.SkuId > 0 {
		m = m.Where(dao.Inventories.Columns().SkuId, req.SkuId)
	}
	if req.WarehouseId > 0 {
		m = m.Where(dao.Inventories.Columns().WarehouseId, req.WarehouseId)
	}
	if req.Status != nil {
		m = m.Where(dao.Inventories.Columns().Status, *req.Status)
	}

	var list []*entity.Inventories
	err = m.OrderDesc(dao.Inventories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}

	buf, err := generateInventoryCSV(list)
	if err != nil {
		return nil, err
	}
	return &v1.InventoriesExportRes{Data: buf.Bytes()}, nil
}

// ── 内部辅助方法 ──────────────────────────────────────────

// generateInventoryCSV 生成库存列表 CSV
func generateInventoryCSV(list []*entity.Inventories) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	defer writer.Flush()

	headers := []string{"ID", "SKU ID", "仓库ID", "物理库存", "预占库存", "可用库存", "安全阈值", "上限", "状态", "最后盘点时间", "最后盘点人"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, inv := range list {
		available := inv.Quantity - inv.Reserved
		lastCounted := ""
		if inv.LastCountedAt != nil {
			lastCounted = inv.LastCountedAt.String()
		}
		row := []string{
			strconv.FormatInt(inv.Id, 10),
			strconv.FormatInt(inv.SkuId, 10),
			strconv.FormatInt(inv.WarehouseId, 10),
			strconv.FormatInt(inv.Quantity, 10),
			strconv.FormatInt(inv.Reserved, 10),
			strconv.FormatInt(available, 10),
			strconv.FormatInt(inv.Threshold, 10),
			strconv.FormatInt(inv.MaxThreshold, 10),
			strconv.Itoa(inv.Status),
			lastCounted,
			inv.LastCountedBy,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

// findOrCreateWithTx 在事务中查找库存记录（FOR UPDATE），不存在则创建
func findOrCreateWithTx(ctx context.Context, tx gdb.TX, skuID, warehouseID int64) (*entity.Inventories, error) {
	var inv entity.Inventories
	err := tx.Model("sp_inventories").
		Where("sku_id", skuID).
		Where("warehouse_id", warehouseID).
		LockUpdate().
		Scan(&inv)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if inv.Id > 0 {
		return &inv, nil
	}

	// 不存在则创建
	inv = entity.Inventories{
		SkuId:       skuID,
		WarehouseId: warehouseID,
		Status:      1, // 充足
	}
	_, err = tx.Model("sp_inventories").Insert(inv)
	if err != nil {
		return nil, err
	}
	// 重新读取以获得完整记录
	err = tx.Model("sp_inventories").
		Where("sku_id", skuID).
		Where("warehouse_id", warehouseID).
		Scan(&inv)
	return &inv, err
}

// updateInvFields 更新库存数量/预占/状态字段
func updateInvFields(ctx context.Context, tx gdb.TX, inv *entity.Inventories) error {
	status := int64(1)
	if inv.Quantity == 0 {
		status = 3 // 无货
	} else if inv.Quantity-inv.Reserved <= 0 {
		status = 2 // 缺货
	}
	_, err := tx.Model("sp_inventories").
		Where("id", inv.Id).
		Data(g.Map{
			"quantity": inv.Quantity,
			"reserved": inv.Reserved,
			"status":   status,
		}).
		Update()
	return err
}

// createLog 创建库存变更流水
func createLog(ctx context.Context, tx gdb.TX, log *do.InventoryLogs) error {
	_, err := tx.Model("sp_inventory_logs").Insert(log)
	return err
}
