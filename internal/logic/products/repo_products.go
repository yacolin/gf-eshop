package products

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
)

// 本文件是 sp_products 主表的数据访问层：读取（多级缓存 / 单条 / 批量 / 条件筛选）与写入（新增 / 更新 / 软删除）。
// 上层 products.go 只做编排，不直接触碰 dao 与 g.DB()。

// ── 事务 ────────────────────────────────────────────────────────────────

// withProductsTx 在商品库事务中执行 fn。
func withProductsTx(ctx context.Context, fn func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Products.Transaction(ctx, fn)
}

// ── 读取 ────────────────────────────────────────────────────────────────

// getByID 按主键查询商品实体；不存在时返回 (nil, nil)。
func getByID(ctx context.Context, id int64) (*entity.Products, error) {
	var p *entity.Products
	if err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Scan(&p); err != nil {
		return nil, err
	}
	return p, nil
}

// getByIDCached 多级缓存读取：L1 本地 → Bloom Filter → L2 Redis → DB 兜底并回填各级缓存。
func getByIDCached(ctx context.Context, id int64) (*entity.Products, error) {
	// 1. L1 本地缓存（最快）
	if p, ok := productLocalCache.get(id); ok {
		return p, nil
	}

	// 2. Bloom Filter 快速拦截（未预热时 count==0 放行）
	if !productBloom.mayExist(id) {
		return nil, errcode.ErrProductNotFound
	}

	// 3. L2 Redis
	if p, err := getProductEntityCache(ctx, id); err == nil && p != nil {
		productLocalCache.set(p.Id, p)
		return p, nil
	}

	// 4. DB 兜底
	p, err := getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errcode.ErrProductNotFound
	}
	// 回填所有缓存层级
	_ = setProductEntityCache(context.Background(), p)
	return p, nil
}

// existsByID 判断商品是否存在。
func existsByID(ctx context.Context, id int64) (bool, error) {
	count, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// listIDs 按类目/品牌/状态查询商品 ID（ID 倒序），用于重建列表 ZSET。
// categoryId 命中自身及其所有子类目。
func listIDs(ctx context.Context, categoryId, brandId int64, status int) ([]int64, error) {
	m := dao.Products.Ctx(ctx).Fields(dao.Products.Columns().Id)
	if categoryId > 0 {
		m = m.Where(dao.Products.Columns().CategoryId+" IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))",
			categoryId, categoryId, strconv.FormatInt(categoryId, 10))
	}
	if brandId > 0 {
		m = m.Where(dao.Products.Columns().BrandId, brandId)
	}
	if status > 0 {
		m = m.Where(dao.Products.Columns().Status, status)
	}

	values, err := m.OrderDesc(dao.Products.Columns().Id).Array()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(values))
	for i, v := range values {
		ids[i] = v.Int64()
	}
	return ids, nil
}

// listIDsByFilter keyset 游标分页查询商品 ID（支持名称/价格区间筛选）。
// 多取一条（size+1）供调用方判断是否还有下一页。
func listIDsByFilter(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) ([]int64, error) {
	m := dao.Products.Ctx(ctx)
	if req.Name != "" {
		m = m.WhereLike(dao.Products.Columns().Name, "%"+req.Name+"%")
	}
	if req.CategoryId > 0 {
		m = m.Where(dao.Products.Columns().CategoryId+" IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))",
			req.CategoryId, req.CategoryId, strconv.FormatInt(req.CategoryId, 10))
	}
	if req.BrandId > 0 {
		m = m.Where(dao.Products.Columns().BrandId, req.BrandId)
	}
	if req.Status != nil {
		m = m.Where(dao.Products.Columns().Status, *req.Status)
	}
	if req.PriceMin > 0 {
		m = m.Where(dao.Products.Columns().Id+" IN (SELECT DISTINCT product_id FROM sp_skus WHERE price >= ? AND deleted_at IS NULL)", req.PriceMin)
	}
	if req.PriceMax > 0 {
		m = m.Where(dao.Products.Columns().Id+" IN (SELECT DISTINCT product_id FROM sp_skus WHERE price <= ? AND deleted_at IS NULL)", req.PriceMax)
	}
	if cursorId > 0 {
		m = m.WhereLT(dao.Products.Columns().Id, cursorId)
	}

	values, err := m.Fields(dao.Products.Columns().Id).OrderDesc(dao.Products.Columns().Id).Limit(size + 1).Array()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(values))
	for i, v := range values {
		ids[i] = v.Int64()
	}
	return ids, nil
}

// listByIDs 按 ID 批量查询商品实体（ID 倒序）；空结果返回空切片而非 nil。
func listByIDs(ctx context.Context, ids []int64) ([]*entity.Products, error) {
	var products []*entity.Products
	err := dao.Products.Ctx(ctx).
		Where(dao.Products.Columns().Id+" IN (?)", ids).
		OrderDesc(dao.Products.Columns().Id).
		Scan(&products)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = make([]*entity.Products, 0)
	}
	return products, nil
}

// ── 写入 ────────────────────────────────────────────────────────────────

// insertProduct 在事务中新增商品主表记录，返回自增主键。
func insertProduct(ctx context.Context, tx gdb.TX, data do.Products) (int64, error) {
	result, err := tx.Model("sp_products").Insert(data)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// updateProduct 按主键更新商品主表（无事务场景）。
func updateProduct(ctx context.Context, id int64, data do.Products) error {
	_, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Update(data)
	return err
}

// updateProductTx 在事务中按主键更新商品主表。
func updateProductTx(ctx context.Context, tx gdb.TX, id int64, data do.Products) error {
	_, err := tx.Model("sp_products").Where("id", id).Update(data)
	return err
}

// setHasDescription 写入 has_description 标记（0/1）。
func setHasDescription(ctx context.Context, tx gdb.TX, id int64, hasDesc int) error {
	_, err := tx.Model("sp_products").Where("id", id).Update(g.Map{
		"has_description": hasDesc,
	})
	return err
}

// deleteProduct 软删除商品（表含 deleted_at，DELETE 自动转 UPDATE）。
func deleteProduct(ctx context.Context, id int64) error {
	_, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Delete()
	return err
}
