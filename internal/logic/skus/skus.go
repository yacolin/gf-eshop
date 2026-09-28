package skus

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/skus/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sSkus struct{}

func init() {
	service.RegisterSkus(&sSkus{})
}

func (s *sSkus) List(ctx context.Context, req *v1.SkusListReq) (res *v1.SkusListRes, err error) {
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

	m := dao.Skus.Ctx(ctx)
	if req.ProductId > 0 {
		m = m.Where(dao.Skus.Columns().ProductId, req.ProductId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.SkusListRes{
			List:  make([]*entity.Skus, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Skus
	err = m.Page(page, size).OrderAsc(dao.Skus.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.SkusListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sSkus) Detail(ctx context.Context, req *v1.SkusDetailReq) (res *v1.SkusDetailRes, err error) {
	var sku *entity.Skus
	err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).Scan(&sku)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errcode.ErrSKUNotFound
	}
	available, status := loadSKUInventory(ctx, sku.Id)
	return &v1.SkusDetailRes{
		Skus:              sku,
		AvailableQuantity: available,
		InventoryStatus:   status,
	}, nil
}

func (s *sSkus) GetByCode(ctx context.Context, req *v1.SkusGetByCodeReq) (res *v1.SkusGetByCodeRes, err error) {
	var sku *entity.Skus
	err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().SkuCode, req.SkuCode).Scan(&sku)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errcode.ErrSKUNotFound
	}
	available, status := loadSKUInventory(ctx, sku.Id)
	return &v1.SkusGetByCodeRes{
		Skus:              sku,
		AvailableQuantity: available,
		InventoryStatus:   status,
	}, nil
}

func (s *sSkus) Create(ctx context.Context, req *v1.SkusCreateReq) (res *v1.SkusCreateRes, err error) {
	result, err := dao.Skus.Ctx(ctx).Insert(do.Skus{
		ProductId:    req.ProductId,
		SkuCode:      req.SkuCode,
		Barcode:      req.Barcode,
		SpecSummary:  req.SpecSummary,
		Price:        req.Price,
		MarketPrice:  req.MarketPrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
		Volume:       req.Volume,
		Length:       req.Length,
		Width:        req.Width,
		Height:       req.Height,
		MinPurchaseQty: req.MinPurchaseQty,
		MaxPurchaseQty: req.MaxPurchaseQty,
		Image:        req.Image,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	// 新增 SKU 会改变所属商品的价格区间，同步检索索引
	syncProductSearchIndex(ctx, req.ProductId)
	return &v1.SkusCreateRes{Id: id}, nil
}

func (s *sSkus) Update(ctx context.Context, req *v1.SkusUpdateReq) (res *v1.SkusUpdateRes, err error) {
	// 先检查是否存在
	sku, err := dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).One()
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errcode.ErrSKUNotFound
	}

	// 仅组装非 nil 字段，实现部分更新
	data := g.Map{}
	if req.Price != nil {
		data[dao.Skus.Columns().Price] = *req.Price
	}
	if req.MarketPrice != nil {
		data[dao.Skus.Columns().MarketPrice] = *req.MarketPrice
	}
	if req.CostPrice != nil {
		data[dao.Skus.Columns().CostPrice] = *req.CostPrice
	}
	if req.Status != nil {
		data[dao.Skus.Columns().Status] = *req.Status
	}
	if req.Image != nil {
		data[dao.Skus.Columns().Image] = *req.Image
	}
	if req.Barcode != nil {
		data[dao.Skus.Columns().Barcode] = *req.Barcode
	}
	if req.Weight != nil {
		data[dao.Skus.Columns().Weight] = *req.Weight
	}
	if req.Volume != nil {
		data[dao.Skus.Columns().Volume] = *req.Volume
	}
	if req.Length != nil {
		data[dao.Skus.Columns().Length] = *req.Length
	}
	if req.Width != nil {
		data[dao.Skus.Columns().Width] = *req.Width
	}
	if req.Height != nil {
		data[dao.Skus.Columns().Height] = *req.Height
	}
	if req.MinPurchaseQty != nil {
		data[dao.Skus.Columns().MinPurchaseQty] = *req.MinPurchaseQty
	}
	if req.MaxPurchaseQty != nil {
		data[dao.Skus.Columns().MaxPurchaseQty] = *req.MaxPurchaseQty
	}

	if len(data) == 0 {
		return &v1.SkusUpdateRes{}, nil
	}

	_, err = dao.Skus.Ctx(ctx).Data(data).Where(dao.Skus.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	// 价格可能已变，同步所属商品的检索索引（sku 为上面已加载的记录）
	syncProductSearchIndex(ctx, sku[dao.Skus.Columns().ProductId].Int64())
	return &v1.SkusUpdateRes{}, nil
}

func (s *sSkus) Delete(ctx context.Context, req *v1.SkusDeleteReq) (res *v1.SkusDeleteRes, err error) {
	// 必须在软删除之前取出所属商品 ID（删除后该行会被 dao 过滤掉）
	pid := productIDOfSku(ctx, req.Id)

	_, err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	// 删除 SKU 会改变商品价格区间，同步检索索引
	syncProductSearchIndex(ctx, pid)
	return &v1.SkusDeleteRes{}, nil
}

// productIDOfSku 查询 SKU 所属商品 ID；查不到返回 0。
func productIDOfSku(ctx context.Context, skuId int64) int64 {
	v, err := dao.Skus.Ctx(ctx).
		Fields(dao.Skus.Columns().ProductId).
		Where(dao.Skus.Columns().Id, skuId).
		Value()
	if err != nil || v == nil || v.IsNil() {
		return 0
	}
	return v.Int64()
}

// syncProductSearchIndex 同步商品检索索引（best-effort）。
//
// 索引中的 price_min/price_max 由 sp_skus 聚合而来，而 SKU 价格可以从本模块
// 独立修改。不触发同步的话，价格区间筛选会用到过期数据 —— 这类漂移不会改变
// 文档数，因此仅靠启动时的计数对账发现不了。
func syncProductSearchIndex(ctx context.Context, productId int64) {
	if productId <= 0 {
		return
	}
	service.Products().SyncSearchDoc(ctx, productId)
}

// loadSKUInventory 从 sp_inventories 查询 SKU 的可用库存
func loadSKUInventory(ctx context.Context, skuId int64) (available int64, status string) {
	type invRow struct {
		Quantity int64 `orm:"quantity"`
		Reserved int64 `orm:"reserved"`
	}
	var row invRow
	err := dao.Inventories.Ctx(ctx).
		Fields("SUM(quantity) AS quantity", "SUM(reserved) AS reserved").
		Where(dao.Inventories.Columns().SkuId, skuId).
		Where("deleted_at IS NULL").
		Group(dao.Inventories.Columns().SkuId).
		Scan(&row)
	if err != nil || row.Quantity <= 0 {
		return 0, "缺货"
	}
	available = row.Quantity - row.Reserved
	if available <= 0 {
		return 0, "缺货"
	}
	return available, "充足"
}
