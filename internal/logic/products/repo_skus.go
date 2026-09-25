package products

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
)

// 本文件是 SKU / 库存 / 规格（sp_skus、sp_sku_specs、sp_inventories、sp_warehouses）的数据访问层。
// 只返回原始数据，响应文案与结构组装交给 assembler.go。

// ── SKU ────────────────────────────────────────────────────────────────

// listSkus 查询商品下未删除的 SKU（ID 正序）；空结果返回空切片而非 nil。
func listSkus(ctx context.Context, productId int64) ([]*entity.Skus, error) {
	var list []*entity.Skus
	err := dao.Skus.Ctx(ctx).
		Where(dao.Skus.Columns().ProductId, productId).
		OrderAsc(dao.Skus.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.Skus, 0)
	}
	return list, nil
}

// countSkusIncludingDeleted 统计商品下 SKU 数量（含软删除）。
// 用于批量生成 SKU 时推算编码序号，避免重新创建时唯一键冲突。
func countSkusIncludingDeleted(ctx context.Context, productId int64) int {
	count, err := dao.Skus.Ctx(ctx).
		Unscoped().
		Where(dao.Skus.Columns().ProductId, productId).
		Count()
	if err != nil {
		return 0
	}
	return count
}

// insertSku 在事务中新增单个 SKU，返回自增主键。
func insertSku(ctx context.Context, tx gdb.TX, data do.Skus) (int64, error) {
	result, err := tx.Model("sp_skus").Insert(data)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// updateSku 在事务中按主键更新单个 SKU。
func updateSku(ctx context.Context, tx gdb.TX, id int64, data do.Skus) error {
	_, err := tx.Model("sp_skus").Where("id", id).Update(data)
	return err
}

// nullableString 空字符串转 NULL：条码等唯一/可空列存 NULL 而非空串。
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// ── 仓库与库存 ──────────────────────────────────────────────────────────

// listWarehouses 查询全部仓库，用于新 SKU 初始化库存行。
func listWarehouses(ctx context.Context) ([]*entity.Warehouses, error) {
	var list []*entity.Warehouses
	if err := dao.Warehouses.Ctx(ctx).Scan(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// insertInventory 在事务中新增一条仓库库存记录。
func insertInventory(ctx context.Context, tx gdb.TX, data do.Inventories) error {
	_, err := tx.Model("sp_inventories").Insert(data)
	return err
}

// skuInventory 单 SKU 的库存原始聚合（多仓库求和）。
type skuInventory struct {
	Quantity int64
	Reserved int64
}

// loadSkuInventory 批量加载 SKU 库存（各仓库求和），键为 sku_id；查询失败返回空映射。
func loadSkuInventory(ctx context.Context, skuIDs []int64) map[int64]skuInventory {
	result := make(map[int64]skuInventory, len(skuIDs))
	if len(skuIDs) == 0 {
		return result
	}

	type invRow struct {
		SkuId    int64 `orm:"sku_id"`
		Quantity int64
		Reserved int64
	}
	var rows []invRow
	err := dao.Inventories.Ctx(ctx).
		Fields("sku_id", "SUM(quantity) AS quantity", "SUM(reserved) AS reserved").
		Where("sku_id IN (?)", skuIDs).
		Where("deleted_at IS NULL").
		Group("sku_id").
		Scan(&rows)
	if err != nil {
		return result
	}

	for _, r := range rows {
		result[r.SkuId] = skuInventory{Quantity: r.Quantity, Reserved: r.Reserved}
	}
	return result
}

// productStats 商品级价格区间与可售总库存。
type productStats struct {
	PriceMin   int64
	PriceMax   int64
	TotalStock int64
}

// listProductStats 批量查询商品级价格区间与总库存。
// 返回的映射保证包含全部入参 ID（无数据时为零值），查询失败时保留零值不报错。
func listProductStats(ctx context.Context, productIDs []int64) map[int64]productStats {
	result := make(map[int64]productStats, len(productIDs))
	if len(productIDs) == 0 {
		return result
	}
	for _, id := range productIDs {
		result[id] = productStats{}
	}

	// 价格区间：从 sp_skus 聚合
	type priceRow struct {
		ProductId int64
		PriceMin  int64
		PriceMax  int64
	}
	var prices []priceRow
	err := dao.Skus.Ctx(ctx).
		Fields("product_id", "MIN(price) AS price_min", "MAX(price) AS price_max").
		Where("product_id IN (?)", productIDs).
		Where("deleted_at IS NULL").
		Group("product_id").
		Scan(&prices)
	if err == nil {
		for _, r := range prices {
			s := result[r.ProductId]
			s.PriceMin = r.PriceMin
			s.PriceMax = r.PriceMax
			result[r.ProductId] = s
		}
	}

	// 总库存：SUM(quantity - reserved) per product
	type stockRow struct {
		ProductId  int64
		TotalStock int64
	}
	var stocks []stockRow
	err = g.DB().Model("sp_skus s").
		Fields("s.product_id", "SUM(COALESCE(i.quantity,0) - COALESCE(i.reserved,0)) AS total_stock").
		LeftJoin("sp_inventories i", "i.sku_id = s.id AND i.deleted_at IS NULL").
		Where("s.product_id IN (?)", productIDs).
		Where("s.deleted_at IS NULL").
		Group("s.product_id").
		Scan(&stocks)
	if err == nil {
		for _, r := range stocks {
			s := result[r.ProductId]
			if r.TotalStock > 0 {
				s.TotalStock = r.TotalStock
			}
			result[r.ProductId] = s
		}
	}

	return result
}

// ── 规格（EAV） ─────────────────────────────────────────────────────────

// listSkuSpecAttrs 查询 SKU 关联的规格属性 ID → 属性名（去重），用于 UpdateFull 的旧规格组校验。
func listSkuSpecAttrs(ctx context.Context, skuIDs []int64) (map[int64]string, error) {
	result := make(map[int64]string)
	if len(skuIDs) == 0 {
		return result, nil
	}

	type specRow struct {
		AttributeId   int64  `orm:"attribute_id"`
		AttributeName string `orm:"attribute_name"`
	}
	var rows []specRow
	err := g.DB().Model("sp_sku_specs ss").
		Fields("DISTINCT ss.attribute_id", "a.name AS attribute_name").
		LeftJoin("sp_attributes a", "a.id = ss.attribute_id").
		Where("ss.sku_id IN (?)", skuIDs).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.AttributeId] = r.AttributeName
	}
	return result, nil
}

// aggregateSkuSpecs 从 sp_sku_specs EAV 表聚合规格维度（如 颜色→[红,蓝]）。
// 查询失败或无数据返回空切片。
func aggregateSkuSpecs(ctx context.Context, skus []*entity.Skus) []v1.ProductAttrDetailResponse {
	if len(skus) == 0 {
		return make([]v1.ProductAttrDetailResponse, 0)
	}

	skuIDs := make([]int64, 0, len(skus))
	for _, sku := range skus {
		skuIDs = append(skuIDs, sku.Id)
	}

	type specRow struct {
		AttributeName    string `orm:"attribute_name"`
		AttributeValueID int64  `orm:"attribute_value_id"`
		AttributeValue   string `orm:"attribute_value"`
	}
	var rows []specRow
	err := g.DB().Model("sp_sku_specs ss").
		Fields("a.name AS attribute_name", "ss.attribute_value_id", "av.value AS attribute_value").
		LeftJoin("sp_attributes a", "a.id = ss.attribute_id").
		LeftJoin("sp_attribute_values av", "av.id = ss.attribute_value_id").
		Where("ss.sku_id IN (?)", skuIDs).
		Order("ss.sort_order ASC").
		Scan(&rows)
	if err != nil || len(rows) == 0 {
		return make([]v1.ProductAttrDetailResponse, 0)
	}

	type attrValues struct {
		set   map[int64]struct{}
		order []*v1.ProductAttrValueResponse
	}
	attrs := make(map[string]*attrValues)
	keyOrder := make([]string, 0)

	for _, r := range rows {
		av, ok := attrs[r.AttributeName]
		if !ok {
			av = &attrValues{set: make(map[int64]struct{})}
			attrs[r.AttributeName] = av
			keyOrder = append(keyOrder, r.AttributeName)
		}
		if _, seen := av.set[r.AttributeValueID]; !seen {
			av.set[r.AttributeValueID] = struct{}{}
			av.order = append(av.order, &v1.ProductAttrValueResponse{
				Id:    r.AttributeValueID,
				Value: r.AttributeValue,
			})
		}
	}

	result := make([]v1.ProductAttrDetailResponse, len(keyOrder))
	for i, name := range keyOrder {
		result[i] = v1.ProductAttrDetailResponse{
			AttributeName: name,
			Values:        attrs[name].order,
		}
	}
	return result
}

// fallbackSkuSpecAttrs 当 sp_sku_specs 无数据时，从 sp_product_attributes 兜底：
// 筛选 is_sku_spec=1 的销售属性呈现为可选规格项。
func fallbackSkuSpecAttrs(ctx context.Context, productId int64) []v1.ProductAttrDetailResponse {
	type attrRow struct {
		AttributeId      int64  `orm:"attribute_id"`
		AttributeName    string `orm:"attribute_name"`
		AttributeValueId int64  `orm:"attribute_value_id"`
		AttributeValue   string `orm:"attribute_value"`
	}
	var rows []attrRow
	err := g.DB().Model("sp_product_attributes pa").
		Fields("pa.attribute_id", "a.name AS attribute_name", "pa.attribute_value_id", "COALESCE(av.value, pa.value) AS attribute_value").
		LeftJoin("sp_attributes a", "a.id = pa.attribute_id").
		LeftJoin("sp_attribute_values av", "av.id = pa.attribute_value_id").
		Where("pa.product_id", productId).
		Where("a.is_sku_spec", 1).
		Where("pa.deleted_at IS NULL").
		Order("a.sort_order ASC, pa.id ASC").
		Scan(&rows)
	if err != nil || len(rows) == 0 {
		return make([]v1.ProductAttrDetailResponse, 0)
	}

	seen := make(map[int64]int)
	result := make([]v1.ProductAttrDetailResponse, 0)
	for _, r := range rows {
		var vals []*v1.ProductAttrValueResponse
		if r.AttributeValueId > 0 || r.AttributeValue != "" {
			vals = []*v1.ProductAttrValueResponse{{Id: r.AttributeValueId, Value: r.AttributeValue}}
		}
		if idx, ok := seen[r.AttributeId]; ok {
			result[idx].Values = append(result[idx].Values, vals...)
		} else {
			seen[r.AttributeId] = len(result)
			result = append(result, v1.ProductAttrDetailResponse{
				AttributeId:   r.AttributeId,
				AttributeName: r.AttributeName,
				Values:        vals,
			})
		}
	}
	return result
}

// insertSkuSpec 在事务中新增一条 SKU 规格维度记录。
func insertSkuSpec(ctx context.Context, tx gdb.TX, data do.SkuSpecs) error {
	_, err := tx.Model("sp_sku_specs").Insert(data)
	return err
}
