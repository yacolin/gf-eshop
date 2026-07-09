package categories

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	categoryIdsKey     = "category:ids"
	categoryIdsTTL     = 10 * time.Minute
	categoryEntityTTL  = 10 * time.Minute
	categoryScoreScale = 1_000_000_000_000
)

const categoryListLua = `
local zkey  = KEYS[1]
local pfx   = ARGV[1]   -- "category:"
local start = tonumber(ARGV[2])
local stop  = tonumber(ARGV[3])

local total = redis.call('ZCARD', zkey)
if total == 0 then
    return {0}
end

local ids = redis.call('ZRANGE', zkey, start, stop)
if #ids == 0 then
    return {total}
end

local keys = {}
for i, id in ipairs(ids) do
    keys[i] = pfx .. id
end
if redis.call('EXISTS', unpack(keys)) < #ids then
    return {-1}
end

local vals = redis.call('MGET', unpack(keys))
local result = {total}
for i = 1, #vals do
    result[#result + 1] = vals[i]
end
return result
`

func cacheKeyCategory(id int64) string { return fmt.Sprintf("category:%d", id) }

func encodeCategoryScore(sortOrder int, id int64) float64 {
	return float64(int64(sortOrder)*categoryScoreScale + (categoryScoreScale - id))
}

// --- Lua 分页读取 ---

func getCategoryPage(ctx context.Context, page, size int) (list []*entity.Categories, total int, err error) {
	start := (page - 1) * size
	stop := start + size - 1

	v, err := g.Redis().Do(ctx, "EVAL", categoryListLua, 1, categoryIdsKey, "category:", start, stop)
	if err != nil {
		return nil, 0, err
	}
	if v.IsNil() || len(v.Vars()) == 0 {
		return nil, 0, nil
	}
	elems := v.Vars()
	total = elems[0].Int()
	if total <= 0 {
		return nil, 0, nil
	}
	list = make([]*entity.Categories, 0, len(elems)-1)
	for _, e := range elems[1:] {
		var c entity.Categories
		if err := sonic.Unmarshal(e.Bytes(), &c); err != nil {
			return nil, 0, err
		}
		list = append(list, &c)
	}
	return list, total, nil
}

// --- Singleflight 重建 ---

var (
	catRebuildMu  sync.Mutex
	catRebuilding = map[string]chan struct{}{
		"categories": nil,
	}
)

func ensureCategoryCache(ctx context.Context) {
	card, _ := g.Redis().Do(ctx, "ZCARD", categoryIdsKey)
	if !card.IsNil() && card.Int() > 0 {
		return
	}

	catRebuildMu.Lock()
	ch, ok := catRebuilding["categories"]
	if !ok || ch == nil {
		ch = make(chan struct{})
		catRebuilding["categories"] = ch
		catRebuildMu.Unlock()

		_, _ = rebuildCategoryCache(context.Background())

		catRebuildMu.Lock()
		delete(catRebuilding, "categories")
		catRebuildMu.Unlock()
		close(ch)
		return
	}
	catRebuildMu.Unlock()
	<-ch
}

// --- 索引 ---

func addCategoryToIndex(ctx context.Context, id int64, sortOrder int) {
	g.Redis().Do(ctx, "ZADD", categoryIdsKey, encodeCategoryScore(sortOrder, id), id)
}

func removeCategoryFromIndex(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "ZREM", categoryIdsKey, id)
}

// --- 缓存重建 ---

func rebuildCategoryCache(ctx context.Context) (int, error) {
	var list []*entity.Categories
	if err := dao.Categories.Ctx(ctx).OrderAsc(dao.Categories.Columns().SortOrder).Scan(&list); err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, nil
	}
	g.Redis().Do(ctx, "DEL", categoryIdsKey)
	for _, c := range list {
		g.Redis().Do(ctx, "ZADD", categoryIdsKey, encodeCategoryScore(c.SortOrder, c.Id), c.Id)
		data, _ := sonic.Marshal(c)
		g.Redis().Do(ctx, "SETEX", cacheKeyCategory(c.Id), int(categoryEntityTTL.Seconds()), string(data))
	}
	g.Redis().Do(ctx, "EXPIRE", categoryIdsKey, int(categoryIdsTTL.Seconds()))
	return len(list), nil
}

func Warmup(ctx context.Context) (int, error) {
	n, err := rebuildCategoryCache(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "category cache warmup failed: %v", err)
		return 0, err
	}
	g.Log().Infof(ctx, "category cache warmed up")
	return n, nil
}

// --- 单条缓存 ---

func getCategoryEntityCache(ctx context.Context, id int64) (*entity.Categories, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyCategory(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var c entity.Categories
	if err := sonic.Unmarshal(v.Bytes(), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func setCategoryEntityCache(ctx context.Context, c *entity.Categories) error {
	data, err := sonic.Marshal(c)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyCategory(c.Id), int(categoryEntityTTL.Seconds()), string(data))
	return err
}

func delCategoryEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyCategory(id))
}
