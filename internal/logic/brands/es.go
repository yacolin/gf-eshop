package brands

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/brands/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/search"
)

const (
	// brandEntity 对应索引别名 eshop_brands（物理索引 eshop_brands_<时间戳>）。
	brandEntity = "brands"
	// brandNGramMax 子串匹配的最大长度，覆盖 varchar(100) 的常见品牌名。
	brandNGramMax = 32
	// brandIndexSchema 索引结构版本。变更 mapping / analyzer / 字段语义时必须递增，
	// 否则启动对账发现不了差异（文档数不变），旧索引会被一直沿用。
	// v2：ngram 检索侧分析器由 keyword 改为与索引侧同一 tokenizer + operator=and。
	brandIndexSchema = "v2"
)

// brandMapping 品牌索引结构。
//
// dynamic=false：实体以后新增字段只会存进 _source 供接口原样回显，不会触发
// 动态 mapping 漂移；需要参与检索时再显式登记字段即可。
//
// _source 里存的是完整 entity.Brands，因此 ES 返回的记录与 DB 查询结果字段完全一致。
func brandMapping() map[string]any {
	cols := dao.Brands.Columns()
	return map[string]any{
		"dynamic": false,
		"properties": map[string]any{
			cols.Id:          map[string]any{"type": "long"},
			cols.Name:        search.TextWithNGram(),
			cols.EnglishName: search.TextWithNGram(),
			cols.FirstLetter: map[string]any{"type": "keyword"},
			cols.SortOrder:   map[string]any{"type": "integer"},
			cols.Status:      map[string]any{"type": "integer"},
			cols.LogoUrl:     map[string]any{"type": "keyword", "index": false},
			cols.Description: search.TextWithNGramNotIndexed(),
			cols.CreatedAt:   search.StoredOnly(),
			cols.UpdatedAt:   search.StoredOnly(),
			cols.DeletedAt:   search.StoredOnly(),
		},
	}
}

// ReindexBrands 全量重建品牌索引。
//
// 流程：新建带时间戳的物理索引 → 全量 bulk → refresh → 原子切换别名 → 删除旧索引。
// 由于「先 add 别名后 remove 旧索引」，检索侧任意时刻都有可用索引，不出现空窗。
func ReindexBrands(ctx context.Context) (int, error) {
	if !search.Available(ctx) {
		return 0, search.ErrUnavailable
	}
	var list []*entity.Brands
	if err := dao.Brands.Ctx(ctx).OrderAsc(dao.Brands.Columns().SortOrder).Scan(&list); err != nil {
		return 0, err
	}

	alias := search.Alias(ctx, brandEntity)
	index := search.VersionedIndex(ctx, brandEntity, brandIndexSchema+"_"+time.Now().Format("20060102150405"))

	// 同一秒内重复重建会撞名，先清理以保证幂等
	_ = search.DeleteIndex(ctx, index)
	if err := search.EnsureIndex(ctx, index, search.IndexSettings(brandNGramMax), brandMapping()); err != nil {
		return 0, err
	}

	docs := make([]search.Doc, 0, len(list))
	for _, b := range list {
		docs = append(docs, search.Doc{ID: strconv.FormatInt(b.Id, 10), Source: b})
	}
	if _, err := search.BulkIndex(ctx, index, docs); err != nil {
		// 重建失败时清理半成品索引，避免留下垃圾
		_ = search.DeleteIndex(ctx, index)
		return 0, err
	}
	if err := search.Refresh(ctx, index); err != nil {
		return 0, err
	}

	removed, err := search.SwapAlias(ctx, alias, index)
	if err != nil {
		return 0, err
	}
	for _, old := range removed {
		if err := search.DeleteIndex(ctx, old); err != nil {
			g.Log().Warningf(ctx, "删除旧品牌索引 %s 失败: %v", old, err)
		}
	}
	g.Log().Infof(ctx, "品牌索引重建完成：%s -> %s（%d 条，清理旧索引 %v）", alias, index, len(docs), removed)
	return len(docs), nil
}

// WarmupES 启动时保证 ES 侧可用：别名缺失或文档数与 DB 不一致时全量重建。
// 这是「双写丢数据 / ES 被清空 / mapping 升级」的统一自愈入口。
func WarmupES(ctx context.Context) (int, error) {
	if !search.Available(ctx) {
		return 0, search.ErrUnavailable
	}
	alias := search.Alias(ctx, brandEntity)

	// 结构版本（mapping/analyzer）变更时计数对账发现不了，必须优先判断
	prefix := search.SchemaIndexPrefix(ctx, brandEntity, brandIndexSchema)
	if ok, err := search.AliasMatchesSchema(ctx, alias, prefix); err == nil && !ok {
		g.Log().Infof(ctx, "品牌索引结构版本已变更（期望前缀 %s），执行全量重建", prefix)
		return ReindexBrands(ctx)
	}

	dbCount, err := dao.Brands.Ctx(ctx).Count()
	if err != nil {
		return 0, err
	}
	esCount, err := search.Count(ctx, alias)
	switch {
	case err == nil && int64(dbCount) == esCount:
		g.Log().Infof(ctx, "品牌 ES 索引与 DB 一致（%d 条），跳过重建", esCount)
		return 0, nil
	case err == nil:
		g.Log().Warningf(ctx, "品牌 ES 文档数 %d 与 DB %d 不一致，触发全量重建", esCount, dbCount)
	case errors.Is(err, search.ErrIndexNotFound):
		g.Log().Infof(ctx, "品牌 ES 索引尚未建立，执行首次全量重建")
	default:
		g.Log().Warningf(ctx, "读取品牌 ES 索引状态失败，执行全量重建: %v", err)
	}
	return ReindexBrands(ctx)
}

// syncBrandDoc 双写：把单个品牌的最新状态写入 ES。
// 失败只记日志（熔断 + 启动自愈兜底），不影响主流程 —— DB 才是唯一真相源。
func syncBrandDoc(ctx context.Context, id int64) {
	if !search.Available(ctx) {
		return
	}
	var b *entity.Brands
	if err := dao.Brands.Ctx(ctx).Where(dao.Brands.Columns().Id, id).Scan(&b); err != nil || b == nil {
		return
	}
	err := search.IndexDoc(
		ctx,
		search.Alias(ctx, brandEntity),
		strconv.FormatInt(id, 10),
		b,
	)
	if err != nil {
		g.Log().Warningf(ctx, "品牌 %d 同步到 ES 失败（等待启动自愈修正）: %v", id, err)
	}
}

// removeBrandDoc 双写：从 ES 删除单个品牌文档。
func removeBrandDoc(ctx context.Context, id int64) {
	if !search.Available(ctx) {
		return
	}
	if err := search.DeleteDoc(ctx, search.Alias(ctx, brandEntity), strconv.FormatInt(id, 10)); err != nil {
		g.Log().Warningf(ctx, "品牌 %d 从 ES 删除失败（等待启动自愈修正）: %v", id, err)
	}
}

// searchBrandsES 用 ES 完成「筛选 + 排序 + 翻页」，返回结果与 DB 分支结构一致。
//
// 查询语义：
//   - first_letter / status → term filter（精确筛选，不参与算分）
//   - name → ngram 子串匹配，并同时命中 english_name（多字段检索，替代 LIKE '%x%'）
//   - 排序 → sort_order asc, id desc，与改造前的 DB 分支完全一致
func searchBrandsES(ctx context.Context, req *v1.BrandsListReq, page, size int) (*v1.BrandsListRes, error) {
	cfg := search.Cfg(ctx)
	from := (page - 1) * size
	if from+size > cfg.MaxResultWindow {
		return nil, search.ErrDeepPaging
	}

	cols := dao.Brands.Columns()
	filters := make([]any, 0, 2)
	if req.FirstLetter != "" {
		filters = append(filters, map[string]any{"term": map[string]any{cols.FirstLetter: req.FirstLetter}})
	}
	if req.Status != nil {
		filters = append(filters, map[string]any{"term": map[string]any{cols.Status: *req.Status}})
	}

	must := make([]any, 0, 1)
	if req.Name != "" {
		must = append(must, map[string]any{
			"bool": map[string]any{
				"should": []any{
					search.NGramMatch(cols.Name, req.Name),
					search.NGramMatch(cols.EnglishName, req.Name),
				},
				"minimum_should_match": 1,
			},
		})
	}

	body := map[string]any{
		"from": from,
		"size": size,
		// 关闭 10000 条封顶的相对估算，返回精确 total
		"track_total_hits": true,
		"sort": []any{
			map[string]any{cols.SortOrder: "asc"},
			map[string]any{cols.Id: "desc"},
		},
		"query": map[string]any{
			"bool": map[string]any{"filter": filters, "must": must},
		},
	}

	result, err := search.Search(ctx, search.Alias(ctx, brandEntity), body)
	if err != nil {
		return nil, err
	}

	list := make([]*entity.Brands, 0, len(result.Hits))
	for _, raw := range result.Hits {
		var b entity.Brands
		if err := sonic.Unmarshal(raw, &b); err != nil {
			return nil, err
		}
		list = append(list, &b)
	}
	return &v1.BrandsListRes{List: list, Total: int(result.Total)}, nil
}
