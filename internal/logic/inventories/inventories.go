package inventories

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

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
			return gerror.NewCode(gcode.CodeValidationFailed, "库存不足")
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
			return gerror.NewCode(gcode.CodeValidationFailed, "预占库存不足，无法释放")
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
			return gerror.NewCode(gcode.CodeValidationFailed, "库存不足")
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
		return nil, gerror.NewCode(gcode.CodeNotFound, "库存记录不存在")
	}
	return &v1.InventoriesGetStockRes{Inventories: inv}, nil
}

// ── 内部辅助方法 ──────────────────────────────────────────

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
