package products

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/singleflight"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	productEntityTTL = 10 * time.Minute

	// L1 local cache
	spuLocalCacheSize = 8192
	spuLocalCacheTTL  = 60 * time.Second
	spuLocalJitter    = 0.2

	// Bloom Filter
	bloomN = 100000
	bloomP = 0.01
)

func cacheKeyProduct(id int64) string { return fmt.Sprintf("product:%d", id) }

// ── L1 Local Cache ──────────────────────────────────────────────────────

type l1Entry struct {
	item      *entity.Products
	expiresAt time.Time
}

type spuLocalCache struct {
	mu      sync.RWMutex
	entries map[int64]*l1Entry
	maxSize int
	ttl     time.Duration
}

func newSPULocalCache() *spuLocalCache {
	return &spuLocalCache{
		entries: make(map[int64]*l1Entry, spuLocalCacheSize),
		maxSize: spuLocalCacheSize,
		ttl:     spuLocalCacheTTL,
	}
}

func (c *spuLocalCache) get(id int64) (*entity.Products, bool) {
	c.mu.RLock()
	entry, ok := c.entries[id]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if entry.expiresAt.Before(time.Now()) {
		return nil, false
	}
	return entry.item, true
}

func (c *spuLocalCache) set(id int64, p *entity.Products) {
	c.mu.Lock()
	if len(c.entries) >= c.maxSize {
		c.evictLocked()
	}
	c.entries[id] = &l1Entry{
		item:      p,
		expiresAt: time.Now().Add(jitteredTTL(c.ttl, spuLocalJitter)),
	}
	c.mu.Unlock()
}

func (c *spuLocalCache) warmupSingle(id int64, p *entity.Products) {
	c.mu.Lock()
	if len(c.entries) >= c.maxSize {
		c.evictLocked()
	}
	c.entries[id] = &l1Entry{
		item:      p,
		expiresAt: time.Now().Add(jitteredTTL(c.ttl, spuLocalJitter)),
	}
	c.mu.Unlock()
}

func (c *spuLocalCache) remove(id int64) {
	c.mu.Lock()
	delete(c.entries, id)
	c.mu.Unlock()
}

func (c *spuLocalCache) evictLocked() {
	now := time.Now()
	for id, entry := range c.entries {
		if entry.expiresAt.Before(now) {
			delete(c.entries, id)
		}
	}
	if len(c.entries) >= c.maxSize {
		for id := range c.entries {
			delete(c.entries, id)
			break
		}
	}
}

// ── Bloom Filter ────────────────────────────────────────────────────────

type spuBloomFilter struct {
	mu     sync.RWMutex
	filter *bloom.BloomFilter
	count  int64
}

func newSPUBloomFilter() *spuBloomFilter {
	return &spuBloomFilter{
		filter: bloom.NewWithEstimates(bloomN, bloomP),
	}
}

func (b *spuBloomFilter) add(id int64) {
	b.mu.Lock()
	b.filter.AddString(strconv.FormatInt(id, 10))
	b.count++
	b.mu.Unlock()
}

func (b *spuBloomFilter) addAll(ids []int64) {
	b.mu.Lock()
	for _, id := range ids {
		b.filter.AddString(strconv.FormatInt(id, 10))
	}
	b.count += int64(len(ids))
	b.mu.Unlock()
}

// mayExist 返回 true 表示可能存在（或未预热）。返回 false 表示一定不存在。
func (b *spuBloomFilter) mayExist(id int64) bool {
	b.mu.RLock()
	c := b.count
	if c == 0 {
		b.mu.RUnlock()
		return true // 未预热，放行
	}
	ok := b.filter.TestString(strconv.FormatInt(id, 10))
	b.mu.RUnlock()
	return ok
}

// ── 全局实例 ────────────────────────────────────────────────────────────

var (
	productLocalCache = newSPULocalCache()
	productBloom      = newSPUBloomFilter()
)

// ── L2 Redis (单条) ─────────────────────────────────────────────────────

func getProductEntityCache(ctx context.Context, id int64) (*entity.Products, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyProduct(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var p entity.Products
	if err := sonic.Unmarshal(v.Bytes(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func setProductEntityCache(ctx context.Context, p *entity.Products) error {
	productBloom.add(p.Id)
	productLocalCache.set(p.Id, p)
	data, err := sonic.Marshal(p)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyProduct(p.Id), int(productEntityTTL.Seconds()), string(data))
	return err
}

func delProductEntityCache(ctx context.Context, id int64) {
	productLocalCache.remove(id)
	g.Redis().Do(ctx, "DEL", cacheKeyProduct(id))
}

// ── 启动预热 ─────────────────────────────────────────────────────────────

var rebuildOnce sync.Once

func Warmup(ctx context.Context) (int, error) {
	var count int
	var errVal error
	rebuildOnce.Do(func() {
		var list []*entity.Products
		if err := dao.Products.Ctx(ctx).OrderDesc(dao.Products.Columns().Id).Scan(&list); err != nil {
			g.Log().Warningf(ctx, "product cache warmup failed: %v", err)
			errVal = err
			return
		}
		if len(list) == 0 {
			return
		}
		count = len(list)

		// Bloom Filter
		ids := make([]int64, count)
		for i, p := range list {
			ids[i] = p.Id
		}
		productBloom.addAll(ids)

		// L1 + L2
		for _, p := range list {
			productLocalCache.warmupSingle(p.Id, p)
			data, _ := sonic.Marshal(p)
			g.Redis().Do(ctx, "SETEX", cacheKeyProduct(p.Id), int(productEntityTTL.Seconds()), string(data))
		}

		g.Log().Infof(ctx, "product cache warmed up: %d items (L1 + Bloom + L2)", count)
	})
	return count, errVal
}

// ── TTL Jitter ──────────────────────────────────────────────────────────

func jitteredTTL(base time.Duration, jitter float64) time.Duration {
	delta := time.Duration(float64(base) * jitter)
	return base + time.Duration(rand.Int63n(int64(delta*2+1))) - delta
}

// ── ZSET List Cache ─────────────────────────────────────────────────────

const productListCacheTTL = 600 // 10 minutes

// productListZSETSf 保证同一筛选组合只有一个 goroutine 重建 ZSET
var productListZSETSf singleflight.Group

// cacheKeyProductListIDs 构造 ZSET 缓存键
func cacheKeyProductListIDs(categoryId, brandId int64, status int) string {
	return fmt.Sprintf("product:list:ids:cat=%d:brand=%d:status=%d", categoryId, brandId, status)
}

// ensureProductListZSET 确保列表 ZSET 存在：缺失时经 singleflight 调用 build 重建并写入。
func ensureProductListZSET(ctx context.Context, key string, build func() ([]int64, error)) error {
	_, err, _ := productListZSETSf.Do(key, func() (interface{}, error) {
		exists, err := g.Redis().Do(ctx, "EXISTS", key)
		if err != nil || exists.Int() == 0 {
			ids, err := build()
			if err != nil {
				return nil, err
			}
			_ = setProductListZSET(ctx, key, ids)
		}
		return nil, nil
	})
	return err
}

// productListFetchScript ZREVRANK + ZREVRANGE 一次 EVAL 往返（倒序）
// 返回: [count, status, id1, id2, ...]
//   count>0: status=0, 后面跟 ID 列表（从高到低）
//   count=0: status=-1 → 游标不在缓存; status=0 → 无更多数据
const productListFetchScript = `
local cursor_id = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local start = 0
if cursor_id > 0 then
    local rank = redis.call("ZREVRANK", KEYS[1], string.format("%020d", cursor_id))
    if not rank then
        return {0, -1}
    end
    start = rank + 1
end
local members = redis.call("ZREVRANGE", KEYS[1], start, start + limit - 1)
local n = #members
if n == 0 then
    return {0, 0}
end
local result = {n, 0}
for i = 1, n do
    result[#result + 1] = tonumber(members[i])
end
return result
`

// setProductListZSET 将 ID 列表写入 ZSET（覆盖重建）
func setProductListZSET(ctx context.Context, key string, ids []int64) error {
	g.Redis().Do(ctx, "DEL", key)
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		g.Redis().Do(ctx, "ZADD", key, float64(id), fmt.Sprintf("%020d", id))
	}
	g.Redis().Do(ctx, "EXPIRE", key, productListCacheTTL)
	return nil
}

// fetchProductListIDs 从 ZSET 游标分页获取 ID 列表
func fetchProductListIDs(ctx context.Context, key string, cursorID int64, limit int) ([]int64, error) {
	v, err := g.Redis().Do(ctx, "EVAL", productListFetchScript, 1, key, cursorID, limit)
	if err != nil || v.IsNil() {
		return nil, err
	}
	arr := v.Interfaces()
	if len(arr) < 2 {
		return nil, nil
	}
	n, _ := strconv.ParseInt(fmt.Sprint(arr[0]), 10, 64)
	st, _ := strconv.ParseInt(fmt.Sprint(arr[1]), 10, 64)
	if n == 0 {
		if st == -1 {
			return nil, fmt.Errorf("cursor not in cache")
		}
		return []int64{}, nil
	}
	ids := make([]int64, 0, n)
	for i := 2; i < len(arr); i++ {
		id, _ := strconv.ParseInt(fmt.Sprint(arr[i]), 10, 64)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// delAllProductListCaches 删除所有商品列表 ZSET 缓存
func delAllProductListCaches(ctx context.Context) {
	v, err := g.Redis().Do(ctx, "KEYS", "product:list:ids:*")
	if err != nil || v.IsNil() {
		return
	}
	for _, key := range v.Interfaces() {
		if ks, ok := key.(string); ok {
			g.Redis().Do(ctx, "DEL", ks)
		}
	}
}
