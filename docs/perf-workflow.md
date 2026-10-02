# 缓存接入流程：ZSET + Lua 原子分页

> 为新业务模块接入高性能列表缓存的标准流程，参考 `brands` 模块的 ZSET+Lua+Singleflight 方案。
> 整个流程只需手写一个 `cache.go` + 修改三处业务代码。

---

## 全景图

```
请求链路：

List 请求 ──→ 无筛选? ──→ EVAL Lua 脚本 ──→ ZCARD              ← 总数
                    │                     └→ ZRANGE            ← 分页 ID
                    │                      └→ EXISTS + MGET   ← 批量实体
                    │                       └→ sonic.Unmarshal ← 反序列化
                    │                          一次 Redis 往返，零次 SQL
                    │
                    └─→ 有筛选 ──→ MySQL 查询（COUNT + SELECT）
```

```
存储模型：

module:ids    (ZSET)    score=id                        member=id    TTL=10min
module:1      (String)  json 序列化的实体                            TTL=10min
module:2      (String)  json 序列化的实体                            TTL=10min
...
```

```
写入维护：

Create ──→ DB Insert ──→ ZADD module:ids (异步, background ctx)
Update ──→ DB Update ──→ ZADD module:ids + DEL module:{id}
Delete ──→ DB Delete ──→ ZREM module:ids + DEL module:{id}
```

---

## 一、创建 `internal/logic/<模块>/cache.go`

这是唯一需要完整手写的文件，**复制以下模板，替换 `product` 为实际模块名**：

```go
package products

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/bytedance/sonic"
    "github.com/gogf/gf/v2/frame/g"

    "gf-eshop/internal/dao"
    "gf-eshop/internal/model/entity"
)

// ── 常量 ──────────────────────────────────────────────────────────

const (
    idsKey     = "product:ids"       // 替换 product → 模块名
    idsTTL     = 10 * time.Minute
    entityTTL  = 10 * time.Minute
    scoreScale = 1_000_000_000_000   // 10^12
)

// ── Lua 脚本（无需修改） ───────────────────────────────────────────

const listLua = `
local zkey  = KEYS[1]
local pfx   = ARGV[1]
local start = tonumber(ARGV[2])
local stop  = tonumber(ARGV[3])

local total = redis.call('ZCARD', zkey)
if total == 0 then return {0} end

local ids = redis.call('ZRANGE', zkey, start, stop)
if #ids == 0 then return {total} end

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

// ── Key 与 Score ──────────────────────────────────────────────────

func cacheKey(id int64) string { return fmt.Sprintf("product:%d", id) }

// score 直接使用 id，保证 ZRANGE 结果按 id ASC 排列。
func encodeScore(sortOrder int, id int64) float64 {
    return float64(int64(sortOrder)*scoreScale + (scoreScale - id))
}

// ── Lua 分页读取 ──────────────────────────────────────────────────

func getModulePage(ctx context.Context, page, size int) (list []*entity.Products, total int, err error) {
    start := (page - 1) * size
    stop := start + size - 1
    v, err := g.Redis().Do(ctx, "EVAL", listLua, 1, idsKey, "product:", start, stop)
    if err != nil || v.IsNil() || len(v.Vars()) == 0 {
        return nil, 0, err
    }
    elems := v.Vars()
    total = elems[0].Int()
    if total <= 0 {
        return nil, 0, nil // 0=空列表, -1=缓存不完整
    }
    list = make([]*entity.Products, 0, len(elems)-1)
    for _, e := range elems[1:] {
        var p entity.Products
        if err := sonic.Unmarshal(e.Bytes(), &p); err != nil {
            return nil, 0, err
        }
        list = append(list, &p)
    }
    return list, total, nil
}

// ── Singleflight 重建 ─────────────────────────────────────────────

var (
    mu    sync.Mutex
    chMap = map[string]chan struct{}{"products": nil}
)

// ensureCache 确保缓存可用；若不存在则触发 singleflight 重建，
// 200 并发下仅 1 个 goroutine 执行重建，其余等待。
func ensureCache(ctx context.Context) {
    card, _ := g.Redis().Do(ctx, "ZCARD", idsKey)
    if !card.IsNil() && card.Int() > 0 {
        return
    }
    mu.Lock()
    ch, ok := chMap["products"]
    if !ok || ch == nil {
        ch = make(chan struct{})
        chMap["products"] = ch
        mu.Unlock()

        _ = rebuildCache(context.Background())

        mu.Lock()
        delete(chMap, "products")
        mu.Unlock()
        close(ch)
        return
    }
    mu.Unlock()
    <-ch
}

// ── 缓存重建 ──────────────────────────────────────────────────────

// rebuildCache 从 DB 重建完整缓存（ZSET + 单条），在预热和缓存 miss 时调用。
func rebuildCache(ctx context.Context) error {
    var list []*entity.Products
    if err := dao.Products.Ctx(ctx).OrderAsc(dao.Products.Columns().SortOrder).Scan(&list); err != nil {
        return err
    }
    if len(list) == 0 {
        return nil
    }
    g.Redis().Do(ctx, "DEL", idsKey)
    for _, p := range list {
        g.Redis().Do(ctx, "ZADD", idsKey, encodeScore(p.SortOrder, p.Id), p.Id)
        data, _ := sonic.Marshal(p)
        g.Redis().Do(ctx, "SETEX", cacheKey(p.Id), int(entityTTL.Seconds()), string(data))
    }
    // ZSET 必须与实体缓存同步过期，否则实体过期后所有请求都触发重建
    g.Redis().Do(ctx, "EXPIRE", idsKey, int(idsTTL.Seconds()))
    return nil
}

// Warmup 启动预热（在 cmd.go 中调用）
func Warmup(ctx context.Context) {
    if err := rebuildCache(ctx); err != nil {
        g.Log().Warningf(ctx, "product cache warmup failed: %v", err)
        return
    }
    g.Log().Infof(ctx, "product cache warmed up")
}

// ── 索引维护（写入时调用） ────────────────────────────────────────

func addToIndex(ctx context.Context, id int64, sortOrder int) {
    g.Redis().Do(ctx, "ZADD", idsKey, encodeScore(sortOrder, id), id)
}

func removeFromIndex(ctx context.Context, id int64) {
    g.Redis().Do(ctx, "ZREM", idsKey, id)
}

// ── 单条缓存（Detail 复用） ──────────────────────────────────────

func getEntityCache(ctx context.Context, id int64) (*entity.Products, error) {
    v, err := g.Redis().Do(ctx, "GET", cacheKey(id))
    if err != nil || v.IsNil() {
        return nil, err
    }
    var p entity.Products
    if err := sonic.Unmarshal(v.Bytes(), &p); err != nil {
        return nil, err
    }
    return &p, nil
}

func setEntityCache(ctx context.Context, p *entity.Products) error {
    data, err := sonic.Marshal(p)
    if err != nil {
        return err
    }
    _, err = g.Redis().Do(ctx, "SETEX", cacheKey(p.Id), int(entityTTL.Seconds()), string(data))
    return err
}

func delEntityCache(ctx context.Context, id int64) {
    g.Redis().Do(ctx, "DEL", cacheKey(id))
}
```

### 需要替换的占位符

| 占位符 | 替换值 | 次数 |
|--------|--------|------|
| `product` / `Products` | 模块名（如 `brand`、`category`） | 全文 |
| `entity.Products` | 对应的 entity 类型 | 4 处 |
| `dao.Products` | 对应的 DAO | 2 处 |
| `products`（chMap key） | 模块名，用于 singleflight | 2 处 |

---

## 二、修改业务逻辑（增/删/改/查四方法）

### List — 替换为三级兜底

```go
func (s *sProducts) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
    page, size := normalizePage(req.Page, req.PageSize)

    // 无筛选：走 ZSET+Lua（热路径）
    if noFilters(req) {
        list, total, err := getModulePage(ctx, page, size)
        if err == nil && total > 0 {
            return &v1.ListRes{List: list, Total: total}, nil
        }
        if ctx.Err() != nil {
            return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
        }
        ensureCache(ctx)                                             // 自动重建
        list, total, err = getModulePage(ctx, page, size)            // 重试一次
        if err == nil && total > 0 {
            return &v1.ListRes{List: list, Total: total}, nil
        }
        // 最终兜底：直接查 DB
        return queryAllAndPaginate(ctx, page, size), nil
    }

    // 有筛选：直接查 DB
    return queryDBWithFilters(ctx, req, page, size)
}
```

### Detail — 增加缓存回写（与现有代码兼容）

```go
func (s *sProducts) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
    cached, err := getEntityCache(ctx, req.Id)
    if err == nil && cached != nil {
        return &v1.DetailRes{Products: cached}, nil
    }
    if ctx.Err() != nil {
        return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
    }

    var entity *entity.Products
    err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Scan(&entity)
    if err != nil || entity == nil {
        return nil, gerror.NewCode(gcode.CodeNotFound, "记录不存在")
    }
    // 回写缓存使用 background context，不绑定请求生命周期
    setEntityCache(context.Background(), entity)
    return &v1.DetailRes{Products: entity}, nil
}
```

### Create — 加入 ZSET

```go
func (s *sProducts) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
    result, err := dao.Products.Ctx(ctx).Insert(do.Products{...})
    if err != nil {
        return nil, err
    }
    id, _ := result.LastInsertId()
    addToIndex(context.Background(), id, req.SortOrder)  // 异步加入 ZSET
    return &v1.CreateRes{Id: id}, nil
}
```

### Update — 更新 ZSET score + 失效详情缓存

```go
func (s *sProducts) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
    // ... 检查存在 ...
    _, err = dao.Products.Ctx(ctx).Data(do.Products{...}).Where(...).Update()
    if err != nil {
        return nil, err
    }
    // 更新 ZSET（score=id，幂等）+ 使详情缓存失效
    g.Redis().Do(context.Background(), "ZADD", idsKey, encodeScore(req.SortOrder, req.Id), req.Id)
    delEntityCache(context.Background(), req.Id)
    return &v1.UpdateRes{}, nil
}
```

### Delete — 从 ZSET 移除

```go
func (s *sProducts) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
    _, err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Delete()
    if err != nil {
        return nil, err
    }
    removeFromIndex(context.Background(), req.Id)
    delEntityCache(context.Background(), req.Id)
    return &v1.DeleteRes{}, nil
}
```

---

## 三、启动预热（cmd.go）

在 `internal/cmd/cmd.go` 的 `Func` 中添加一行预热调用：

```go
import productsLogic "gf-eshop/internal/logic/products"

func Func(ctx context.Context, parser *gcmd.Parser) error {
    brandsLogic.Warmup(ctx)      // 已有
    categoriesLogic.Warmup(ctx)  // 已有
    productsLogic.Warmup(ctx)    // ★ 新增
    // ...
}
```

---

## 四、配置检查

确保 `manifest/config/config.yaml` 的 Redis 连接池足够承载目标并发：

```yaml
redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    maxActive: 300       # 略大于峰值并发
    maxIdle: 200          # 接近 maxActive，避免频繁建连
    minIdle: 50           # 保底空闲连接
    idleTimeout: "300s"   # 长空闲超时，减少回收
    waitTimeout: "3s"
    readTimeout: "2s"
    writeTimeout: "2s"
    maxConnLifetime: "30m"
```

> **关键**：`maxIdle` 必须接近 `maxActive`，否则高并发下大部分请求需要新建 TCP 连接，建连开销会拖垮性能。

---

## 五、接入清单

### 新建文件

| 文件 | 行数 | 说明 |
|------|------|------|
| `internal/logic/<模块>/cache.go` | ~200 | 模板代码，替换 module/entity/dao 名 |

### 修改文件

| 文件 | 改动 | 行数 |
|------|------|------|
| `internal/logic/<模块>/<模块>.go` | **List** 三级兜底 + **Detail** 缓存 + **Create** ZADD + **Update** ZADD+DEL + **Delete** ZREM+DEL | ~20 行 |
| `internal/cmd/cmd.go` | 加一行 `Warmup(ctx)` | +1 行 |

---

## 六、FAQ

### Q: 这个方案适用于多少条数据的模块？

任意规模。ZSET 的 `ZRANGE` 和 `ZCARD` 均为 O(log N)，MGET O(1)，瓶颈不在数据量而在序列化速度。

### Q: 实体有自定义排序字段怎么办？

当前列表 ZSET 的 `score` 就是 `id`，即 id 升序。如需其他排序逻辑，改 score 写入处即可。

### Q: ZSET 和实体缓存的 TTL 为什么要一样？

如果 ZSET 不过期而实体过期，Lua 脚本中的 `EXISTS` 检测会失败，每次都触发重建。**两个 key 必须设相同的 TTL**。

### Q: 为什么写入操作（ZADD/ZREM）要传 `context.Background()`？

请求的 HTTP context 在响应返回后即被取消。若 Redis 操作绑定请求 context，可能被中途取消导致缓存不一致。

### Q: 怎么验证缓存正常工作？

```bash
# 1. 看启动日志
./main 2>&1 | grep "cache warmed up"

# 2. 看 Redis key
redis-cli --raw ZCARD brand:ids
redis-cli GET "brand:1"

# 3. 压测验证零 SQL
wrk -t4 -c200 -d30s --latency http://localhost:8000/api/v1/brands
# 确认日志中没有 SQL 查询打印
```

### Q: 如果模块不需要 Detail 接口怎么办？

只保留 `getModulePage`、`ensureCache`、`Warmup` 和索引维护方法，去掉单条缓存的 `getEntityCache` / `setEntityCache` / `delEntityCache`。ZSET 方案仍然适用于纯列表场景。

### Q: 有筛选条件的查询为什么不走缓存？

两个原因：

1. **命中率低** — 筛选条件排列组合无限（`name=a&page=1`、`name=b&page=2`...），缓存几秒钟就过期，收益微乎其微
2. **失效复杂** — 任一条数据变更都可能影响所有筛选结果，维护 N × N 的缓存失效关系不现实

如果后续业务出现**高频复用的筛选条件**（如首页"已上架"筛选占总流量 30% 以上），可以对热条件单独加缓存：

```go
if req.Status == 1 && req.Name == "" && req.FirstLetter == "" {
    cacheKey := fmt.Sprintf("product:list:status:enabled:%d:%d", page, size)
    // 走缓存，写入时额外 DEL 这个 key
}
```

### Q: `hack/config.yaml` 的 `jsonOmitEmpty` 不起作用？

gf v2.6.1 不支持 `jsonOmitEmpty`。需要用 `sed` 手动给 `entity/*.go` 和 `do/*.go` 的 string 和 `*gtime.Time` 字段加 `omitempty`：

```bash
for f in internal/model/entity/*.go internal/model/do/*.go; do
  sed -i '' 's/\(string.*json:"[^"]*\)"/\1,omitempty"/g' "$f"
  sed -i '' 's/\(\*gtime\.Time.*json:"[^"]*\)"/\1,omitempty"/g' "$f"
done
```

注意：int 字段不要加 `omitempty`，否则 `status: 0` 会被忽略。
