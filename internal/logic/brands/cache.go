package brands

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
	brandIdsKey     = "brand:ids"
	brandIdsTTL     = 10 * time.Minute
	brandEntityTTL  = 10 * time.Minute
	brandScoreScale = 1_000_000_000_000
)

// brandListLua 在一个原子操作中完成 ZCARD + ZRANGE + EXISTS + MGET，
// 只需一次 Redis 往返。
const brandListLua = `
local zkey  = KEYS[1]
local pfx   = ARGV[1]   -- "brand:"
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

-- 先检查所有 key 是否存在，避免部分 MGET miss
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

func cacheKeyBrand(id int64) string { return fmt.Sprintf("brand:%d", id) }

func encodeBrandScore(sortOrder int, id int64) float64 {
	return float64(int64(sortOrder)*brandScoreScale + (brandScoreScale - id))
}

// --- Lua 分页读取 ---

func getBrandPage(ctx context.Context, page, size int) (list []*entity.Brands, total int, err error) {
	start := (page - 1) * size
	stop := start + size - 1

	v, err := g.Redis().Do(ctx, "EVAL", brandListLua, 1, brandIdsKey, "brand:", start, stop)
	if err != nil {
		return nil, 0, err
	}
	if v.IsNil() || len(v.Vars()) == 0 {
		return nil, 0, nil
	}
	elems := v.Vars()
	total = elems[0].Int()
	if total <= 0 {
		return nil, 0, nil // 缓存不完整（-1）或为空（0）
	}
	list = make([]*entity.Brands, 0, len(elems)-1)
	for _, e := range elems[1:] {
		var b entity.Brands
		if err := sonic.Unmarshal(e.Bytes(), &b); err != nil {
			return nil, 0, err
		}
		list = append(list, &b)
	}
	return list, total, nil
}

// --- Singleflight 重建 ---

var (
	rebuildMu  sync.Mutex
	rebuilding = map[string]chan struct{}{
		"brands":    nil,
	}
)

// ensureBrandCache 确保缓存可用：若缓存不存在则触发 singleflight 重建
func ensureBrandCache(ctx context.Context) {
	card, _ := g.Redis().Do(ctx, "ZCARD", brandIdsKey)
	if !card.IsNil() && card.Int() > 0 {
		return
	}

	rebuildMu.Lock()
	ch, ok := rebuilding["brands"]
	if !ok || ch == nil {
		ch = make(chan struct{})
		rebuilding["brands"] = ch
		rebuildMu.Unlock()

		_, _ = rebuildBrandCache(context.Background())

		rebuildMu.Lock()
		delete(rebuilding, "brands")
		rebuildMu.Unlock()
		close(ch)
		return
	}
	rebuildMu.Unlock()
	<-ch
}

// --- 索引与批量 ---

func addBrandToIndex(ctx context.Context, id int64, sortOrder int) {
	g.Redis().Do(ctx, "ZADD", brandIdsKey, encodeBrandScore(sortOrder, id), id)
}

func removeBrandFromIndex(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "ZREM", brandIdsKey, id)
}

// --- 缓存重建 ---

func rebuildBrandCache(ctx context.Context) (int, error) {
	var list []*entity.Brands
	if err := dao.Brands.Ctx(ctx).OrderAsc(dao.Brands.Columns().SortOrder).Scan(&list); err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, nil
	}
	g.Redis().Do(ctx, "DEL", brandIdsKey)
	for _, b := range list {
		g.Redis().Do(ctx, "ZADD", brandIdsKey, encodeBrandScore(b.SortOrder, b.Id), b.Id)
		data, _ := sonic.Marshal(b)
		g.Redis().Do(ctx, "SETEX", cacheKeyBrand(b.Id), int(brandEntityTTL.Seconds()), string(data))
	}
	// ZSET 与实体缓存同步过期
	g.Redis().Do(ctx, "EXPIRE", brandIdsKey, int(brandIdsTTL.Seconds()))
	return len(list), nil
}

// Warmup 启动预热
func Warmup(ctx context.Context) (int, error) {
	n, err := rebuildBrandCache(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "brand cache warmup failed: %v", err)
		return 0, err
	}
	g.Log().Infof(ctx, "brand cache warmed up")
	return n, nil
}

// --- 单条缓存 ---

func getBrandEntityCache(ctx context.Context, id int64) (*entity.Brands, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyBrand(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var b entity.Brands
	if err := sonic.Unmarshal(v.Bytes(), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func setBrandEntityCache(ctx context.Context, b *entity.Brands) error {
	data, err := sonic.Marshal(b)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyBrand(b.Id), int(brandEntityTTL.Seconds()), string(data))
	return err
}

func delBrandEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyBrand(id))
}
