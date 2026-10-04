package products

import (
	"encoding/base64"
	"strconv"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
)

// 本文件负责「数据 → 响应结构」的组装，不访问数据库。

// ── 列表 ────────────────────────────────────────────────────────────────

// 每页条数默认值与上限，与 orders 列表保持一致。
const (
	defaultListSize = 20
	maxListSize     = 100
)

// normalizeListSize 归一化每页条数：<=0 取默认值，超上限则截断。
func normalizeListSize(size int) int {
	if size <= 0 {
		return defaultListSize
	}
	if size > maxListSize {
		return maxListSize
	}
	return size
}

// buildListItems 组装列表项：商品实体 + 价格区间/总库存。
func buildListItems(products []*entity.Products, stats map[int64]productStats) []*v1.ProductsListItem {
	list := make([]*v1.ProductsListItem, len(products))
	for i, p := range products {
		s := stats[p.Id]
		list[i] = &v1.ProductsListItem{
			Products:   p,
			PriceMin:   s.PriceMin,
			PriceMax:   s.PriceMax,
			TotalStock: s.TotalStock,
		}
	}
	return list
}

// buildListResponse 组装列表响应：末位条目 ID 作为下一页游标。
// 只有 hasMore 时才给出 next_cursor —— 保证「next_cursor 为空 ⟺ 没有更多」这一
// 两个模块共用的语义（否则末页会回一个指不到东西的游标）。
// 刻意不返回 total：游标分页不做 COUNT（与 orders 一致）。
func buildListResponse(products []*entity.Products, stats map[int64]productStats, hasMore bool) *v1.ProductsListRes {
	nextCursor := ""
	if hasMore && len(products) > 0 {
		nextCursor = encodeCursor(products[len(products)-1].Id)
	}
	return &v1.ProductsListRes{
		List:       buildListItems(products, stats),
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

// emptyListResponse 空结果响应：空切片而非 nil，字段值与「查到了但没有下一页」一致。
func emptyListResponse() *v1.ProductsListRes {
	return &v1.ProductsListRes{
		List: make([]*v1.ProductsListItem, 0),
	}
}

// encodeCursor 生成游标（base64 编码的末位商品 ID）。
func encodeCursor(id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

// decodeCursor 解析游标；空字符串代表首页（返回 0）。
// 非法游标**返回错误**而不是静默当成首页 —— 静默降级会让调用方以为翻到了下一页，
// 实际又拿到第一页数据（此处与 orders 行为一致）。
func decodeCursor(cursor string) (int64, error) {
	if cursor == "" {
		return 0, nil
	}
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, errcode.Newf(errcode.CodeInvalidParams, "游标格式非法（应为 base64(id)）")
	}
	id, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil || id <= 0 {
		return 0, errcode.Newf(errcode.CodeInvalidParams, "游标内容非法: %q", string(b))
	}
	return id, nil
}

// ── 详情 ────────────────────────────────────────────────────────────────

// buildSkuItems SKU 实体 → 详情项（库存字段由 fillSkuInventory 填充）。
func buildSkuItems(skus []*entity.Skus) []*v1.SkuDetailItem {
	items := make([]*v1.SkuDetailItem, len(skus))
	for i, sku := range skus {
		items[i] = &v1.SkuDetailItem{Skus: sku}
	}
	return items
}

// skuItemIDs 提取详情项的 SKU ID 列表，用于批量加载库存。
func skuItemIDs(items []*v1.SkuDetailItem) []int64 {
	ids := make([]int64, len(items))
	for i, item := range items {
		ids[i] = item.Id
	}
	return ids
}

// fillSkuInventory 用库存聚合结果填充可用数量与库存状态文案；
// 无库存记录按 0 / 缺货处理，可用量负数归零。
func fillSkuInventory(items []*v1.SkuDetailItem, inv map[int64]skuInventory) {
	for _, item := range items {
		var avail int64
		if v, ok := inv[item.Id]; ok {
			avail = v.Quantity - v.Reserved
			if avail < 0 {
				avail = 0
			}
		}
		item.AvailableQuantity = avail
		if avail > 0 {
			item.InventoryStatus = "充足"
		} else {
			item.InventoryStatus = "缺货"
		}
	}
}

// buildSpecResponse 组装可/不可选规格：
// 补充 attribute_id / sort_order，并按属性名去重（已出现在可选规格中的不再重复展示）。
func buildSpecResponse(specAttrs, productAttrs []v1.ProductAttrDetailResponse, idByName map[string]int64, orderByName map[string]int) *v1.ProductSpecResponse {
	// 富化可选的规格（attribute_id / sort_order）
	for i := range specAttrs {
		specAttrs[i].AttributeId = idByName[specAttrs[i].AttributeName]
		specAttrs[i].SortOrder = orderByName[specAttrs[i].AttributeName]
	}
	if specAttrs == nil {
		specAttrs = make([]v1.ProductAttrDetailResponse, 0)
	}

	// 不可选的规格：去重（排除已在可选规格中出现的）
	seen := make(map[string]bool, len(specAttrs))
	for _, s := range specAttrs {
		seen[s.AttributeName] = true
	}
	nonSelectable := make([]v1.ProductAttrDetailResponse, 0, len(productAttrs))
	for _, a := range productAttrs {
		if seen[a.AttributeName] {
			continue
		}
		seen[a.AttributeName] = true
		a.AttributeId = idByName[a.AttributeName]
		a.SortOrder = orderByName[a.AttributeName]
		nonSelectable = append(nonSelectable, a)
	}
	if nonSelectable == nil {
		nonSelectable = make([]v1.ProductAttrDetailResponse, 0)
	}

	return &v1.ProductSpecResponse{
		Selectable:    specAttrs,
		NonSelectable: nonSelectable,
	}
}

// buildDetailResponse 聚合 SKU 价格区间与可售总库存，组装商品详情响应。
func buildDetailResponse(product *entity.Products, description *entity.ProductDescriptions, skuItems []*v1.SkuDetailItem, specs *v1.ProductSpecResponse) *v1.ProductsDetailRes {
	var priceMin, priceMax, totalStock int64
	for _, item := range skuItems {
		if item.Price > 0 {
			if priceMin == 0 || item.Price < priceMin {
				priceMin = item.Price
			}
			if item.Price > priceMax {
				priceMax = item.Price
			}
		}
		totalStock += item.AvailableQuantity
	}

	return &v1.ProductsDetailRes{
		Products:    product,
		PriceMin:    priceMin,
		PriceMax:    priceMax,
		TotalStock:  totalStock,
		Description: description,
		SKUs:        skuItems,
		Specs:       specs,
	}
}
