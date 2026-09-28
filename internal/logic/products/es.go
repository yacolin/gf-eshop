package products

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/search"
)

const (
	// productEntity 对应索引别名 eshop_products。
	productEntity = "products"
	// productNGramMax 子串匹配的最大长度。商品名 varchar(200)，
	// 取 32 覆盖常见搜索词，同时控制 ngram 索引膨胀。
	productNGramMax = 32
	// productIndexSchema 索引结构版本。变更 mapping / analyzer / 字段语义时必须递增，
	// 否则启动对账发现不了差异（文档数不变），旧索引会被一直沿用。
	// v2：ngram 检索侧分析器由 keyword 改为与索引侧同一 tokenizer + operator=and，
	//     否则含空格的商品名（如「三星e 青春版」）整串搜不到。
	productIndexSchema = "v2"
	// productReindexBatch 重建时每批从 DB 读取的行数，
	// 避免一次性把全表（几十万行）读进内存。
	productReindexBatch = 2000

	// 以下两个字段不在 sp_products 中，由 SKU 聚合而来，仅用于价格区间筛选。
	fieldPriceMin = "price_min"
	fieldPriceMax = "price_max"
)

// productDoc 是写入 ES 的**检索投影**：只保留参与筛选/排序的字段。
//
// 为什么不存全量实体：列表响应统一由 MySQL 侧的
// listByIDs + listProductStats + buildListResponse 组装，
// 而 products 字段多且含 images/seo_* 等大字段，全量存进 ES 只会让索引显著膨胀。
// 检索投影让索引体积与「可检索字段数」成正比，而不是与实体宽度成正比。
type productDoc struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Subtitle      string  `json:"subtitle"`
	CategoryId    int64   `json:"category_id"`
	BrandId       int64   `json:"brand_id"`
	MerchantId    int64   `json:"merchant_id"`
	Status        int     `json:"status"`
	SortOrder     int     `json:"sort_order"`
	SalesCount    int     `json:"sales_count"`
	RatingAverage float64 `json:"rating_average"`
	PriceMin      int64   `json:"price_min"`
	PriceMax      int64   `json:"price_max"`
}

// productMapping 商品索引结构。
//
// dynamic=false：未知字段进 _source 但不建索引，避免 mapping 漂移。
// 注意 subtitle 用 IK 分词但不加 ngram 子字段：它最长 500 字符，
// 挂 ngram 会让索引急剧膨胀，而副标题只做词级检索已足够。
func productMapping() map[string]any {
	cols := dao.Products.Columns()
	return map[string]any{
		"dynamic": false,
		"properties": map[string]any{
			cols.Id:         map[string]any{"type": "long"},
			cols.Name:       search.TextWithNGram(),
			cols.Subtitle:   map[string]any{"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart"},
			cols.CategoryId: map[string]any{"type": "long"},
			cols.BrandId:    map[string]any{"type": "long"},
			cols.MerchantId: map[string]any{"type": "long"},
			cols.Status:     map[string]any{"type": "integer"},
			cols.SortOrder:  map[string]any{"type": "integer"},
			cols.SalesCount: map[string]any{"type": "integer"},
			// 评分保留两位小数，用 scaled_float 精确存储避免浮点误差
			cols.RatingAverage: map[string]any{"type": "scaled_float", "scaling_factor": 100},
			// 价格区间来自 sp_skus 聚合，见 buildProductDocs
			fieldPriceMin: map[string]any{"type": "long"},
			fieldPriceMax: map[string]any{"type": "long"},
		},
	}
}

// ── 文档构建 ────────────────────────────────────────────────────────────

// buildProductDocs 把商品实体转换为 ES 文档，价格区间复用 listProductStats
// 以保证与列表响应里的 price_min/price_max 同源。
func buildProductDocs(ctx context.Context, products []*entity.Products) []search.Doc {
	ids := make([]int64, 0, len(products))
	for _, p := range products {
		ids = append(ids, p.Id)
	}
	stats := listProductStats(ctx, ids)

	docs := make([]search.Doc, 0, len(products))
	for _, p := range products {
		st := stats[p.Id]
		docs = append(docs, search.Doc{
			ID: strconv.FormatInt(p.Id, 10),
			Source: productDoc{
				Id:            p.Id,
				Name:          p.Name,
				Subtitle:      p.Subtitle,
				CategoryId:    p.CategoryId,
				BrandId:       p.BrandId,
				MerchantId:    p.MerchantId,
				Status:        p.Status,
				SortOrder:     p.SortOrder,
				SalesCount:    p.SalesCount,
				RatingAverage: p.RatingAverage,
				PriceMin:      st.PriceMin,
				PriceMax:      st.PriceMax,
			},
		})
	}
	return docs
}

// ── 全量重建 ────────────────────────────────────────────────────────────

// ReindexProducts 全量重建商品索引。
// 与 brands 相同：新索引 → 分批 bulk → refresh → 原子切别名 → 删旧索引。
// 读取也分批（按 id 递增翻页），避免一次性加载全表。
func ReindexProducts(ctx context.Context) (int, error) {
	if !search.Available(ctx) {
		return 0, search.ErrUnavailable
	}
	alias := search.Alias(ctx, productEntity)
	index := search.VersionedIndex(ctx, productEntity, productIndexSchema+"_"+time.Now().Format("20060102150405"))
	cols := dao.Products.Columns()

	_ = search.DeleteIndex(ctx, index)
	if err := search.EnsureIndex(ctx, index, search.IndexSettings(productNGramMax), productMapping()); err != nil {
		return 0, err
	}

	total := 0
	lastId := int64(0)
	for {
		var batch []*entity.Products
		err := dao.Products.Ctx(ctx).
			WhereGT(cols.Id, lastId).
			OrderAsc(cols.Id).
			Limit(productReindexBatch).
			Scan(&batch)
		if err != nil {
			_ = search.DeleteIndex(ctx, index)
			return total, err
		}
		if len(batch) == 0 {
			break
		}
		docs := buildProductDocs(ctx, batch)
		n, err := search.BulkIndex(ctx, index, docs)
		total += n
		if err != nil {
			_ = search.DeleteIndex(ctx, index)
			return total, err
		}
		lastId = batch[len(batch)-1].Id
		if len(batch) < productReindexBatch {
			break
		}
	}

	if err := search.Refresh(ctx, index); err != nil {
		return total, err
	}
	removed, err := search.SwapAlias(ctx, alias, index)
	if err != nil {
		return total, err
	}
	for _, old := range removed {
		if err := search.DeleteIndex(ctx, old); err != nil {
			g.Log().Warningf(ctx, "删除旧商品索引 %s 失败: %v", old, err)
		}
	}
	g.Log().Infof(ctx, "商品索引重建完成：%s -> %s（%d 条，清理旧索引 %v）", alias, index, total, removed)
	return total, nil
}

// ── 启动对账自愈 ────────────────────────────────────────────────────────

// dbPriceAggregate 汇总「按商品聚合的 SKU 价格区间」的总量指纹。
// 与 ES 侧 CollectStats 的 sum(price_min)/sum(price_max) 对比即可发现价格漂移。
func dbPriceAggregate(ctx context.Context) (count, sumMin, sumMax int64, err error) {
	sql := fmt.Sprintf(`
		SELECT COUNT(*) AS cnt, IFNULL(SUM(mn),0) AS sum_min, IFNULL(SUM(mx),0) AS sum_max
		FROM (
			SELECT product_id, MIN(price) AS mn, MAX(price) AS mx
			FROM %s WHERE deleted_at IS NULL GROUP BY product_id
		) t`, dao.Skus.Table())

	var row struct {
		Cnt    int64 `json:"cnt"`
		SumMin int64 `json:"sum_min"`
		SumMax int64 `json:"sum_max"`
	}
	if err = g.DB().GetScan(ctx, &row, sql); err != nil {
		return 0, 0, 0, err
	}
	return row.Cnt, row.SumMin, row.SumMax, nil
}

// WarmupES 启动时保证商品索引与 DB 一致。
//
// 只比文档数是不够的：商品文档里的 price_min/price_max 来自 SKU 聚合，
// SKU 改价后文档数不变但价格已过期，价格区间筛选会出错。
// 因此这里额外比对两边价格合计，任一不一致即全量重建。
func WarmupES(ctx context.Context) (int, error) {
	if !search.Available(ctx) {
		return 0, search.ErrUnavailable
	}
	alias := search.Alias(ctx, productEntity)

	// 结构版本（mapping/analyzer）变更时计数与价格合计都发现不了，必须优先判断
	prefix := search.SchemaIndexPrefix(ctx, productEntity, productIndexSchema)
	if ok, err := search.AliasMatchesSchema(ctx, alias, prefix); err == nil && !ok {
		g.Log().Infof(ctx, "商品索引结构版本已变更（期望前缀 %s），执行全量重建", prefix)
		return ReindexProducts(ctx)
	}

	dbCount, err := dao.Products.Ctx(ctx).Count()
	if err != nil {
		return 0, err
	}
	dbAggCount, dbSumMin, dbSumMax, err := dbPriceAggregate(ctx)
	if err != nil {
		return 0, err
	}

	esCount, sums, err := search.CollectStats(ctx, alias, []string{fieldPriceMin, fieldPriceMax})
	switch {
	case err == nil && int64(dbCount) == esCount &&
		int64(sums[fieldPriceMin]) == dbSumMin && int64(sums[fieldPriceMax]) == dbSumMax &&
		esCount == dbAggCount:
		g.Log().Infof(ctx, "商品 ES 索引与 DB 一致（%d 条，价格合计 %d/%d），跳过重建",
			esCount, dbSumMin, dbSumMax)
		return 0, nil
	case err == nil:
		g.Log().Warningf(ctx,
			"商品 ES 索引与 DB 不一致，触发全量重建：文档数 es=%d db=%d，价格合计 es=%v/%v db=%d/%d",
			esCount, dbCount, sums[fieldPriceMin], sums[fieldPriceMax], dbSumMin, dbSumMax)
	case search.IsIndexNotFound(err):
		g.Log().Infof(ctx, "商品 ES 索引尚未建立，执行首次全量重建")
	default:
		g.Log().Warningf(ctx, "读取商品 ES 索引状态失败，执行全量重建: %v", err)
	}
	return ReindexProducts(ctx)
}

// ── 单条双写 ────────────────────────────────────────────────────────────

// syncProductDoc 把单个商品的最新状态（含最新价格区间）写入 ES。
// 失败只记日志，DB 仍是唯一真相源，由启动对账与 reindex 兜底。
func syncProductDoc(ctx context.Context, id int64) {
	if !search.Available(ctx) {
		return
	}
	p, err := getByID(ctx, id)
	if err != nil || p == nil {
		// 商品已被删除或查询失败：确保索引里不残留旧文档
		if err == nil && p == nil {
			removeProductDoc(ctx, id)
		}
		return
	}
	docs := buildProductDocs(ctx, []*entity.Products{p})
	if len(docs) == 0 {
		return
	}
	if err := search.IndexDoc(ctx, search.Alias(ctx, productEntity), docs[0].ID, docs[0].Source); err != nil {
		g.Log().Warningf(ctx, "商品 %d 同步到 ES 失败（等待启动对账修正）: %v", id, err)
	}
}

// removeProductDoc 从 ES 删除单个商品文档。
func removeProductDoc(ctx context.Context, id int64) {
	if !search.Available(ctx) {
		return
	}
	if err := search.DeleteDoc(ctx, search.Alias(ctx, productEntity), strconv.FormatInt(id, 10)); err != nil {
		g.Log().Warningf(ctx, "商品 %d 从 ES 删除失败（等待启动对账修正）: %v", id, err)
	}
}

// SyncSearchDoc 供其他模块（如 skus）在改动商品或其 SKU 价格后同步检索索引。
//
// 为什么必须暴露：索引里的 price_min/price_max 来自 sp_skus 聚合，
// 而 SKU 价格可以从 skus 模块独立修改。若不触发重新同步，
// 价格区间筛选就会用到过期数据（文档数不变，因此仅靠计数对账发现不了）。
func (s *sProducts) SyncSearchDoc(ctx context.Context, productId int64) {
	if productId <= 0 {
		return
	}
	syncProductDoc(ctx, productId)
}

// ── 检索 ────────────────────────────────────────────────────────────────

// descendantCategoryIDs 解析类目自身及其所有子类目 ID，
// SQL 语义与 listIDs/listIDsByFilter 中的子查询保持一致（按 path 前缀匹配）。
func descendantCategoryIDs(ctx context.Context, rootId int64) ([]int64, error) {
	sql := fmt.Sprintf(`
		SELECT id FROM %s
		WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM %s WHERE id = ?), ?, '/%%')`,
		dao.Categories.Table(), dao.Categories.Table())

	values, err := g.DB().GetArray(ctx, sql, rootId, rootId, strconv.FormatInt(rootId, 10))
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(values))
	for _, v := range values {
		ids = append(ids, v.Int64())
	}
	return ids, nil
}

// searchProductIDsES 用 ES 产出有序商品 ID（id 倒序），语义对齐 listIDsByFilter：
//
//   - 游标分页：cursor 就是上一页末位的商品 ID，等价于 WHERE id < cursor
//     （排序键本身就是 id，因此用 range 过滤即可，不需要 search_after）
//   - 多取一条（size+1）供调用方判断 hasMore
//   - 价格区间：DB 侧是 EXISTS(price >= min) AND EXISTS(price <= max)，
//     等价于 price_max >= min AND price_min <= max，故用反范式的两个字段做 range
//   - 名称：name.ngram 子串匹配（等价 LIKE '%x%'），并同时命中 subtitle（IK 词级）
//
// 返回的 ID 仍交由 listByIDs/listProductStats 组装，ES 只负责「选出哪些 ID、什么顺序」。
func searchProductIDsES(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) ([]int64, bool, error) {
	cols := dao.Products.Columns()

	filters := make([]any, 0, 6)
	if req.CategoryId > 0 {
		ids, err := descendantCategoryIDs(ctx, req.CategoryId)
		if err != nil {
			return nil, false, err
		}
		if len(ids) == 0 {
			return nil, false, nil
		}
		filters = append(filters, map[string]any{"terms": map[string]any{cols.CategoryId: ids}})
	}
	if req.BrandId > 0 {
		filters = append(filters, map[string]any{"term": map[string]any{cols.BrandId: req.BrandId}})
	}
	if req.Status != nil {
		filters = append(filters, map[string]any{"term": map[string]any{cols.Status: *req.Status}})
	}
	if req.PriceMin > 0 {
		filters = append(filters, map[string]any{
			"range": map[string]any{fieldPriceMax: map[string]any{"gte": req.PriceMin}},
		})
	}
	if req.PriceMax > 0 {
		filters = append(filters, map[string]any{
			"range": map[string]any{fieldPriceMin: map[string]any{"lte": req.PriceMax}},
		})
	}
	if cursorId > 0 {
		filters = append(filters, map[string]any{
			"range": map[string]any{cols.Id: map[string]any{"lt": cursorId}},
		})
	}

	must := make([]any, 0, 1)
	if req.Name != "" {
		// 只按 name 匹配，与 DB 的 name LIKE '%x%' 保持同样的字段范围。
		//
		// 刻意**不**把 subtitle 放进 should：subtitle 是 IK 词级检索，
		// 与 name 的 ngram 子串语义不同，混在 should 里会让
		// 「按商品名搜」返回一堆只命中了副标题无关词的商品。
		// 多字段/相关度排序需要配合按 _score 排序，属于后续工作
		// （当前 listByIDs 会按 id 倒序重排，会覆盖相关性顺序）。
		must = append(must, search.NGramMatch(cols.Name, req.Name))
	}

	body := map[string]any{
		"size": size + 1,
		// 只需要 id 参与组装，不回传整个 _source
		"_source": []string{cols.Id},
		"sort":    []any{map[string]any{cols.Id: "desc"}},
		"query":   map[string]any{"bool": map[string]any{"filter": filters, "must": must}},
	}

	result, err := search.Search(ctx, search.Alias(ctx, productEntity), body)
	if err != nil {
		return nil, false, err
	}
	ids := make([]int64, 0, len(result.Hits))
	for _, raw := range result.Hits {
		var d struct {
			Id int64 `json:"id"`
		}
		if err := sonic.Unmarshal(raw, &d); err != nil {
			return nil, false, err
		}
		ids = append(ids, d.Id)
	}

	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}
	return ids, hasMore, nil
}
