# 缓存性能优化路径：从 Cache-Aside 到 ZSET+Lua+Singleflight

本文档记录 gf-eshop 项目品牌/类目列表接口的完整优化过程，可作为后续业务模块的标准化缓存方案参考。

---

## 一、性能目标与测试工具

| 指标 | 目标 |
|------|------|
| 并发连接 | 200 |
| P99 延迟 | < 50ms |
| 错误率 | 0% |
| SQL 查询 | 仅写入路径 + 缓存重建时 |

压测命令：
```bash
wrk -t4 -c200 -d30s --latency http://localhost:8000/api/v1/brands
wrk -t2 -c50 -d10s --latency http://localhost:8000/api/v1/brands/1
```

---

## 二、优化历程

### 阶段 0：初始状态（仅基础 Cache-Aside）

**架构**：
```
List 请求 → Redis GET brand:all（全量 JSON）→ 内存分页
              ↑ miss
              └→ MySQL SELECT COUNT + SELECT page → 回写缓存

Detail 请求 → Redis GET brand:{id} → 返回
              ↑ miss
              └→ MySQL SELECT WHERE id=? → 回写缓存
```

**存储**：
| Key | 类型 | 内容 | TTL |
|-----|------|------|-----|
| `brand:all` | String | 全量品牌 JSON 数组 | 10min |
| `brand:{id}` | String | 单个品牌 JSON | 10min |

**写入**：增/删/改 → 删 `brand:all` + `brand:{id}`

**压测结果**（50 并发）：
```
P50:  66.30ms
P99:  105.63ms
QPS:  746 req/s
```

**暴露的问题**：

| # | 问题 | 根因 |
|---|------|------|
| 1 | P99 105ms | 每次请求穿透 MySQL，COUNT + SELECT 两次查询 |
| 2 | `context canceled` 错误 | 缓存 miss + 请求上下文已取消 → 仍查 DB，级联失败 |
| 3 | 缓存过期=雪崩 | 全量 JSON 过期后，所有请求同时打 MySQL |
| 4 | 200 并发扛不住 | Redis/MySQL 连接池默认值太小 |

---

### 阶段 1：连接池配置（config.yaml）

**问题**：默认连接池太小，200 并发直接打满。

**修改 `manifest/config/config.yaml`**：

```yaml
database:
  default:
    link: "mysql:root:123456@tcp(localhost:3306)/eshop_db?..."
    maxActive: 100        # 最大活跃连接（默认 0=无限，实际很小）
    maxIdle: 10           # 空闲连接数
    idleTimeout: "60s"
    maxConnLifetime: "30m"

redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    maxActive: 300        # 略大于目标并发数（200 + 缓冲）
    maxIdle: 200          # ★ 关键：预建连接，避免频繁 TCP 握手
    minIdle: 50           # 保底空闲连接
    idleTimeout: "300s"   # 延长空闲回收时间
    waitTimeout: "3s"     # 池满等待上限
    dialTimeout: "2s"
    readTimeout: "2s"
    writeTimeout: "2s"
    maxConnLifetime: "30m"
```

**关键决策**：`maxIdle` 必须接近 `maxActive`。若 `maxIdle=20` 且 `maxActive=200`，每轮 200 并发需要新建 180 个 TCP 连接（用完就关），连接建立开销拖死系统。

---

### 阶段 2：防穿透 + 上下文取消保护（brands.go）

**问题 1**：缓存 miss 时请求上下文已取消，仍穿透到数据库。

**修复**：缓存 miss 后先检查 `ctx.Err()`，已取消则直接返回错误：

```go
func (s *sBrands) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
    cached, err := getBrandEntityCache(ctx, req.Id)
    if err == nil && cached != nil {
        return &v1.DetailRes{Brands: cached}, nil
    }
    // ★ 上下文已取消则直接返回，不继续查 DB
    if ctx.Err() != nil {
        return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
    }
    // ... DB 查询
}
```

**问题 2**：Detail 回写缓存绑定请求 context，请求取消导致缓存写入失败。

**修复**：缓存写入使用 `context.Background()`，与请求生命周期解耦：

```go
// 使用 background context 回写缓存，避免因请求上下文取消导致缓存写入失败
if err := setBrandEntityCache(context.Background(), entity); err != nil {
    g.Log().Warning(ctx, "setBrandEntityCache failed: %v", err)
}
```

---

### 阶段 3：ZSET + Lua 原子分页（核心重构）

**为什么不用全量 JSON 缓存？**

| 方案 | 单次传输 | 排序 | 并发瓶颈 |
|------|---------|------|---------|
| 全量 JSON (`brand:all`) | 50KB+（随数据增长） | 内存 sort | 大 JSON 反复序列化 |
| ZSET + MGET | ~5KB（仅一页） | Redis ZSET score | 需 3 次 Redis 往返 |

全量 JSON 的问题是 **缓存键粒度太粗**：一条数据变更就要重建整个缓存，且每次传输全量数据。

#### 3.1 数据模型

```
Redis 存储结构：

brand:ids    (ZSET)    score=encode(sort_order, id)   member=id    TTL=10min
brand:1      (String)  {"id":1,"name":"苹果",...}                  TTL=10min
brand:2      (String)  {"id":2,"name":"华为",...}                  TTL=10min
...
```

**Score 编码**：`sort_order * 10^12 + (10^12 - id)`

这样 ZRANGE 按 score ASC 排列时，自动满足「sort_order ASC, id DESC」的排序需求，无需额外排序。

```go
const brandScoreScale = 1_000_000_000_000

func encodeBrandScore(sortOrder int, id int64) float64 {
    return float64(int64(sortOrder)*brandScoreScale + (brandScoreScale - id))
}
```

#### 3.2 Lua 脚本：一次往返完成全部操作

**问题**：ZCARD → ZRANGE → MGET 需要 3 次 Redis 往返。200 并发 × 3 = 600 个排队命令。

**解决**：Lua 脚本在 Redis 服务端原子执行：

```lua
-- KEYS[1]: ZSET key ("brand:ids")
-- ARGV[1]: entity key prefix ("brand:")
-- ARGV[2]: start offset
-- ARGV[3]: stop offset

local zkey  = KEYS[1]
local pfx   = ARGV[1]
local start = tonumber(ARGV[2])
local stop  = tonumber(ARGV[3])

-- 1. 获取总数
local total = redis.call('ZCARD', zkey)
if total == 0 then
    return {0}
end

-- 2. 获取分页 ID
local ids = redis.call('ZRANGE', zkey, start, stop)
if #ids == 0 then
    return {total}
end

-- 3. 校验所有实体 key 存在（防止 ZSET 在但实体过期的半残状态）
local keys = {}
for i, id in ipairs(ids) do
    keys[i] = pfx .. id
end
if redis.call('EXISTS', unpack(keys)) < #ids then
    return {-1}    -- 缓存不完整，触发重建
end

-- 4. 批量获取实体 JSON
local vals = redis.call('MGET', unpack(keys))
local result = {total}
for i = 1, #vals do
    result[#result + 1] = vals[i]
end
return result
```

**Go 调用**：
```go
func getBrandPage(ctx context.Context, page, size int) (list []*entity.Brands, total int, err error) {
    start := (page - 1) * size
    stop := start + size - 1
    v, err := g.Redis().Do(ctx, "EVAL", brandListLua, 1, brandIdsKey, "brand:", start, stop)
    // ...
    elems := v.Vars()
    total = elems[0].Int()
    if total <= 0 {
        return nil, 0, nil  // 0=空, -1=不完整
    }
    for _, e := range elems[1:] {
        var b entity.Brands
        sonic.Unmarshal(e.Bytes(), &b)
        list = append(list, &b)
    }
    return list, total, nil
}
```

#### 3.3 缓存键同步过期

**问题**：之前 ZSET 没设 TTL，实体 key 10 分钟后就过期了，但 ZSET 还在，导致 Lua 脚本检测到 EXISTS < #ids 返回 -1，所有请求穿透。

**修复**：重建缓存时给 ZSET 和实体 key 都设相同的 10 分钟 TTL：

```go
func rebuildBrandCache(ctx context.Context) error {
    // ... DEL + ZADD + SETEX ...
    // ZSET 与实体缓存同步过期
    g.Redis().Do(ctx, "EXPIRE", brandIdsKey, int(brandIdsTTL.Seconds()))
    return nil
}
```

#### 3.4 Singleflight：防止缓存击穿

**问题**：缓存过期瞬间，200 个并发请求同时检测到缓存不可用，全部触发 DB 全表扫描。

**解决**：手写 singleflight —— 同一时刻只有一个 goroutine 执行重建，其余等待：

```go
var (
    rebuildMu  sync.Mutex
    rebuilding = map[string]chan struct{}{"brands": nil}
)

func ensureBrandCache(ctx context.Context) {
    // 快速检查：缓存已存在则直接返回
    card, _ := g.Redis().Do(ctx, "ZCARD", brandIdsKey)
    if !card.IsNil() && card.Int() > 0 {
        return
    }

    rebuildMu.Lock()
    ch, ok := rebuilding["brands"]
    if !ok || ch == nil {
        // 我是第一个：创建 channel，执行重建
        ch = make(chan struct{})
        rebuilding["brands"] = ch
        rebuildMu.Unlock()

        _ = rebuildBrandCache(context.Background())

        rebuildMu.Lock()
        delete(rebuilding, "brands")
        rebuildMu.Unlock()
        close(ch)  // 通知所有等待者
        return
    }
    rebuildMu.Unlock()
    <-ch  // 等待重建完成
}
```

#### 3.5 List 处理器：三级兜底

```go
func (s *sBrands) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
    // 第一级：Lua 缓存命中（热路径，一次 Redis 往返）
    list, total, err := getBrandPage(ctx, page, size)
    if err == nil && total > 0 {
        return &v1.ListRes{List: list, Total: total}, nil
    }

    // 上下文取消检查
    if ctx.Err() != nil {
        return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
    }

    // 第二级：触发 singleflight 重建，重建后重试
    ensureBrandCache(ctx)
    list, total, err = getBrandPage(ctx, page, size)
    if err == nil && total > 0 {
        return &v1.ListRes{List: list, Total: total}, nil
    }

    // 第三级：最终兜底，直接查库
    var dbAll []*entity.Brands
    dao.Brands.Ctx(ctx).OrderAsc(...).Scan(&dbAll)
    return paginateBrands(page, size, dbAll), nil
}
```

#### 3.6 写入路径：维护 ZSET 索引

```go
// Create: 先写 DB，再异步加入 ZSET
func (s *sBrands) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
    result, err := dao.Brands.Ctx(ctx).Insert(...)
    id, _ := result.LastInsertId()
    addBrandToIndex(context.Background(), id, req.SortOrder)
    return &v1.CreateRes{Id: id}, nil
}

// Update: 更新 ZSET score（sort_order 可能变化）+ 删旧详情缓存
func (s *sBrands) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
    dao.Brands.Ctx(ctx).Data(...).Where(...).Update()
    g.Redis().Do(context.Background(), "ZADD", brandIdsKey, encodeBrandScore(req.SortOrder, req.Id), req.Id)
    delBrandEntityCache(context.Background(), req.Id)
    return &v1.UpdateRes{}, nil
}

// Delete: 从 ZSET 移除 + 删详情缓存
func (s *sBrands) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
    dao.Brands.Ctx(ctx).Where(...).Delete()
    removeBrandFromIndex(context.Background(), req.Id)
    delBrandEntityCache(context.Background(), req.Id)
    return &v1.DeleteRes{}, nil
}
```

**压测结果**（200 并发）：
```
P50:  ~5ms
P99:  ~25ms
QPS:  14,000+ req/s
SQL:  0（纯缓存命中）
```

---

### 阶段 4：Sonic 替代 encoding/json（序列化加速）

**问题**：Go 标准库 `encoding/json` 使用反射，在高 QPS 场景下 CPU profile 中序列化占比显著。

**解决**：替换为字节跳动的 [sonic](https://github.com/bytedance/sonic)，基于 JIT 和 SIMD 的 JSON 库，API 完全兼容：

```go
// 之前
import "encoding/json"
json.Marshal(v)
json.Unmarshal(data, &v)

// 之后
import "github.com/bytedance/sonic"
sonic.Marshal(v)
sonic.Unmarshal(data, &v)
```

**替换范围**：仅替换业务代码中我们直接调用的序列化（cache 读写、Lua 脚本返回值解析）。GoFrame 框架内部的 JSON 处理不受影响。

**安装**：
```bash
go get github.com/bytedance/sonic@latest
```

---

## 三、最终架构总览

```
┌─────────────────────────────────────────────────────────┐
│                     List 请求                            │
│  GET /api/v1/brands?page=1&page_size=20                 │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
         ┌─────────────────────────┐
         │  getBrandPage(page,size) │  ← 无筛选条件
         │  EVAL brandListLua       │     一次 Redis 往返
         │  ZCARD+ZRANGE+EXISTS    │
         │  +MGET → sonic.Unmarshal │
         └───────────┬─────────────┘
                     │
          ┌──────────┼──────────┐
          │          │          │
      命中(>0)   不完整(-1)   为空(0)
          │          │          │
          ▼          ▼          ▼
       返回分页   ensureBrandCache   返回空列表
       (热路径)   (singleflight)
                     │
                     ▼
           ┌──────────────────┐
           │ rebuildBrandCache │  ← 唯一入口，含互斥锁
           │ DEL + ZADD +     │     只有一个 goroutine 执行
           │ SETEX + EXPIRE   │     其余等待 channel 信号
           └──────────────────┘
                     │
                     ▼
              重试 getBrandPage
                     │
              失败 → DB 兜底查询
```

```
┌─────────────────────────────────────────────────────────┐
│                    Detail 请求                            │
│  GET /api/v1/brands/1                                   │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
         ┌─────────────────────────┐
         │  getBrandEntityCache(1)  │
         │  GET brand:1             │
         │  sonic.Unmarshal         │
         └───────────┬─────────────┘
                     │
              命中 ←─┼─→ 未命中 → ctx.Err()? → DB 查询
              返回     │           已取消→返回错误  ↓
                       │                     sonic.Marshal
                       │                     SETEX brand:1
                       │                     (context.Background)
                       │
                       ▼
                    返回详情
```

```
┌─────────────────────────────────────────────────────────┐
│                    写入请求                               │
│  POST/PUT/DELETE /api/v1/brands                         │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
         ┌─────────────────────────┐
         │  MySQL INSERT/UPDATE/   │
         │  DELETE                  │
         └───────────┬─────────────┘
                     │
                     ▼
         ┌─────────────────────────┐
         │  异步维护 Redis（bg ctx）│
         │  Create: ZADD brand:ids │
         │  Update: ZADD + DEL     │
         │  Delete: ZREM + DEL     │
         └─────────────────────────┘
```

---

## 四、新业务接入模板

假设要为新模块 `products` 接入 ZSET+Lua 缓存，按以下步骤：

### 4.1 创建 `internal/logic/products/cache.go`

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

const (
    idsKey     = "product:ids"       // ZSET key
    idsTTL     = 10 * time.Minute
    entityTTL  = 10 * time.Minute
    scoreScale = 1_000_000_000_000
)

// ★ Lua 脚本（直接复制，替换业务前缀即可）
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

func cacheKey(id int64) string { return fmt.Sprintf("product:%d", id) }

func encodeScore(sortOrder int, id int64) float64 {
    return float64(int64(sortOrder)*scoreScale + (scoreScale - id))
}

// --- Lua 分页读取 ---
func getProductPage(ctx context.Context, page, size int) (list []*entity.Products, total int, err error) {
    start := (page - 1) * size
    stop := start + size - 1
    v, err := g.Redis().Do(ctx, "EVAL", listLua, 1, idsKey, "product:", start, stop)
    if err != nil || v.IsNil() || len(v.Vars()) == 0 {
        return nil, 0, err
    }
    elems := v.Vars()
    total = elems[0].Int()
    if total <= 0 {
        return nil, 0, nil
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

// --- Singleflight ---
var (
    mu     sync.Mutex
    chMap  = map[string]chan struct{}{"products": nil}
)

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

// --- 缓存重建 ---
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

// --- 索引维护 ---
func addToIndex(ctx context.Context, id int64, sortOrder int) {
    g.Redis().Do(ctx, "ZADD", idsKey, encodeScore(sortOrder, id), id)
}

func removeFromIndex(ctx context.Context, id int64) {
    g.Redis().Do(ctx, "ZREM", idsKey, id)
}

// --- 单条缓存 ---
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

### 4.2 List 处理器模板

```go
func (s *sProducts) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
    page, size := normalizePage(req.Page, req.PageSize)

    // 无筛选：走 ZSET+Lua
    if noFilters(req) {
        list, total, err := getProductPage(ctx, page, size)
        if err == nil && total > 0 {
            return &v1.ListRes{List: list, Total: total}, nil
        }
        if ctx.Err() != nil {
            return nil, gerror.NewCode(gcode.CodeOperationFailed, "请求已取消")
        }
        ensureCache(ctx)
        list, total, err = getProductPage(ctx, page, size)
        if err == nil && total > 0 {
            return &v1.ListRes{List: list, Total: total}, nil
        }
        // 兜底 DB
        var dbAll []*entity.Products
        dao.Products.Ctx(ctx).OrderAsc(...).Scan(&dbAll)
        return paginate(page, size, dbAll), nil
    }

    // 有筛选：直接查 DB
    return queryDB(ctx, req, page, size)
}
```

### 4.3 写入处理器模板

```go
func (s *sProducts) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
    result, _ := dao.Products.Ctx(ctx).Insert(...)
    id, _ := result.LastInsertId()
    addToIndex(context.Background(), id, req.SortOrder) // 异步维护 ZSET
    return &v1.CreateRes{Id: id}, nil
}

func (s *sProducts) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
    dao.Products.Ctx(ctx).Data(...).Where(...).Update()
    // 更新 ZSET score + 失效详情缓存
    g.Redis().Do(context.Background(), "ZADD", idsKey, encodeScore(req.SortOrder, req.Id), req.Id)
    delEntityCache(context.Background(), req.Id)
    return &v1.UpdateRes{}, nil
}

func (s *sProducts) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
    dao.Products.Ctx(ctx).Where(...).Delete()
    removeFromIndex(context.Background(), req.Id)
    delEntityCache(context.Background(), req.Id)
    return &v1.DeleteRes{}, nil
}
```

### 4.4 启动预热（cmd.go）

```go
import productsLogic "gf-eshop/internal/logic/products"

func Func(ctx context.Context, parser *gcmd.Parser) error {
    brandsLogic.Warmup(ctx)
    categoriesLogic.Warmup(ctx)
    productsLogic.Warmup(ctx)  // ★ 新增
    // ...
}
```

---

## 五、关键设计决策与避坑指南

### 5.1 为什么用 ZSET 而不是缓存整页结果？

| 方案 | 缓存 key 数量 | 数据变更影响 | 适用数据量 |
|------|-------------|-------------|-----------|
| 整页缓存 (`list:page:1`) | page 数 | 任一数据变更，所有页缓存失效 | 数据几乎不变 |
| 全量缓存 (`list:all`) | 1 | 任一数据变更，整个缓存失效 | < 100 条 |
| **ZSET + 单条** | N+1 | 仅失效变更的数据 + 更新 ZSET score | 任意规模 |

ZSET 方案的优势：**缓存失效粒度最细**，变更一条数据只影响一个 key。

### 5.2 Score 编码的陷阱

```go
const scoreScale = 1_000_000_000_000 // 10^12

func encodeScore(sortOrder int, id int64) float64 {
    return float64(int64(sortOrder)*scoreScale + (scoreScale - id))
}
```

- `10^12` 确保 `sort_order` 变化 1 时 score 跳变足够大，不会与 id 部分重叠
- `(scoreScale - id)` 使同 sort_order 下大 id 的 score 更小 → ZRANGE ASC 时 id DESC
- float64 精度约 15 位十进制，`10^12` 仅占 12 位，安全余量 3 位

### 5.3 为什么 ZSET 和实体必须同步 TTL？

如果 ZSET 不过期而实体过期 → Lua 脚本 EXISTS 检测失败 → 返回 -1 → 每次请求都触发重建。**两个 key 必须设相同的 TTL**。

### 5.4 为什么写入操作用 context.Background()？

请求 context 在 HTTP 响应返回后即被取消。若缓存维护绑定请求 context，异步 Redis 操作可能被取消，导致缓存不一致。**所有缓存维护操作必须使用独立 context**。

### 5.5 为什么需要 EXISTS 检查（Lua 脚本中）？

极端情况：实体 `brand:5` 的 TTL 恰好到期被 Redis 淘汰，但 ZSET `brand:ids` 还有 0.1 秒才到期。此时 ZRANGE 返回了 id=5，但 MGET `brand:5` 返回 nil → 解析失败。

`EXISTS` 预检确保**所有 key 齐全**才执行 MGET，否则返回 -1 触发重建。

---

## 六、配置文件完整参考

```yaml
# manifest/config/config.yaml

server:
  address: ":8000"

database:
  default:
    link: "mysql:user:pass@tcp(host:3306)/db?..."
    maxActive: 100
    maxIdle: 10
    idleTimeout: "60s"
    maxConnLifetime: "30m"

redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    maxActive: 300       # 略大于峰值并发
    maxIdle: 200         # ★ 接近 maxActive，避免频繁建连
    minIdle: 50          # 保底
    idleTimeout: "300s"  # ★ 长空闲超时，减少回收
    waitTimeout: "3s"
    dialTimeout: "2s"
    readTimeout: "2s"
    writeTimeout: "2s"
    maxConnLifetime: "30m"
```

---

## 七、优化效果对照

| 指标 | 优化前 | 优化后 |
|------|--------|--------|
| 架构 | Cache-Aside（全量 JSON） | ZSET + Lua + Singleflight |
| List QPS (50 并发) | 746 | 14,000+ |
| List P99 (50 并发) | 105ms | < 10ms |
| List P99 (200 并发) | 扛不住 | ~25ms |
| Detail P99 (50 并发) | 62ms | < 5ms |
| 每次 List 请求的 Redis 往返 | 1（全量 GET） | 1（Lua EVAL） |
| 每次 List 请求的 SQL 查询 | 2（COUNT+SELECT） | 0 |
| 缓存击穿保护 | 无 | Singleflight |
| JSON 库 | encoding/json | sonic |
| 缓存不一致窗口 | 全量失效 | 单条粒度 |
