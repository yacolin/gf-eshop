# Elasticsearch 搜索接入指南（brands Pilot + products 落地）

本文记录 gf-eshop 接入 Elasticsearch 的完整过程：**brands** 作为 Pilot 跑通模式，
**products** 作为正式落地兑现价值；`skus` 与相关度排序留作后续。
配套阅读：[缓存性能优化路径](./cache-optimization-guide.md)。

---

## 一、要解决的问题

改造前的品牌列表接口是「两条腿走路」：

| 查询场景 | 走哪条路 | 体验 |
|----------|----------|------|
| 无筛选（仅翻页） | Redis ZSET + Lua 一次往返 | 快，已验证 |
| **有筛选**（name / first_letter / status） | **直查 MySQL** | `LIKE '%x%'` 无法走索引，全靠 DB 扛 |

多级缓存 + 布隆过滤器只优化了「按 ID 取单条」和「无筛选分页」，**一旦带上筛选条件就全部落回 DB**，
这正是本次要补的能力。

## 二、适用判断（先说实话）

| 维度 | 现状 | 结论 |
|------|------|------|
| `sp_brands` 数据量 | **100 行** | 远低于 ES 适用门槛 |
| `sp_products` | 2025 行 | 同样偏小 |
| `sp_skus` | 7942 行 | 同样偏小 |
| 索引情况 | `first_letter` / `status` 已有单列索引 | 当前筛选在 MySQL 上其实很快 |

**所以：brands 接入 ES 的性能收益在当前数据量下接近于零**，它的价值是
**把「索引建模 / 零停机重建 / 双写 / 熔断降级」这套模式跑通并固化成可复用底座**，
等 products / skus 数据量涨到几十万级、筛选维度变复杂（规格、价格区间、类目、库存）时直接复用。

> 判断标准：数据量大（几十万+）、查询模式复杂、且是核心业务 —— 品牌字典表一般用不上。

### 2.1 为什么保留 Redis 层：成本差与爆炸半径

一个自然的疑问：既然 ES 也能做无筛选查询，那条 Redis ZSET + Lua 路径是不是多余了？
实测结论是**不多余**。保留两条路径的依据不是「路径不同」，而是**成本差一个数量级 +
降级时爆炸半径不同**。如果只是路径不同而没有成本差，就该合并成一条、少维护一套。

压测环境：100 条品牌、单节点 ES（8G 堆）、`ab -k`。两条路径返回结果**完全一致**
（`total=100`、ids 顺序相同），所以差距纯粹是代价。为消除顺序偏差，交替顺序各跑两轮：

| 无筛选查询 | RPS | P50 | P95 | P99 |
|---|---|---|---|---|
| **Redis Lua**（当前路径） | 11001 / 11757 | 7–8ms | 13–15ms | 24–51ms |
| **ES**（对照：无筛选也走 ES） | 6619 / 6619 | 8–9ms | 47–55ms | 110–132ms |

ES 两轮 6619.16 / 6618.93 几乎完全相同，是**撞到吞吐天花板**的典型特征。并发 200 时差距进一步放大：

| 并发 200 | RPS | P50 | P95 | P99 |
|---|---|---|---|---|
| Redis Lua | 11898 | 13ms | 24ms | **73ms** |
| ES | 4838 | 18ms | 157ms | **461ms** |

原因是结构性的，不是调优能解决的：Redis Lua 是**一次内存操作**
（`EVAL` 内完成 ZCARD+ZRANGE+EXISTS+MGET），而 ES 每次都要 HTTP 往返 + DSL 解析 +
`QUERY_THEN_FETCH` 两阶段 + 分片聚合 + JVM。

**这组数据直接对着 [缓存优化路径](./cache-optimization-guide.md) 的目标**：那里定的是
P99 < 50ms，而 ES 单独扛无筛选查询时是 110–132ms（并发 200 时 461ms）会击穿目标，
Redis 路径保住了它。

**比延迟更要命的是降级时的爆炸半径**：

| 架构 | ES 挂掉时会发生什么 |
|------|---------------------|
| Redis + ES（当前） | 只有筛选查询回落 MySQL，最热的无筛选流量仍由 Redis 吸收 |
| 全走 ES | **所有**列表流量一起砸向 MySQL，MySQL 独自承受全部并发 |

后者正是 200 并发缓存优化要消灭的故障模式。ES 在这里是单节点 8G 堆，
把最热查询压上去等于让一个 JVM 成为关键路径。

### 2.2 分工判断的两个轴

**切勿把「C 端 / 管理端」当成分工依据**，两条轴各管一件事：

| | 无筛选 | 有筛选 |
|---|---|---|
| C 端首页品牌墙 / 品牌导航 | **Redis** | — |
| C 端筛选 / 搜索品牌 | — | **ES** |
| 管理端列表 | Redis | ES |

- **轴 1「查询形态」决定谁来回这个请求**：无筛选 → Redis，有筛选 → ES。
  C 端与管理端都会发出这两种形态的请求，所以分工与「端」无关。
- **轴 2「流量画像」决定值不值得为 Redis 这层付维护成本**：
  C 端高频 → 值得保留；纯管理端低频（几十 QPS）→ 不值得，
  可统一走 ES，省掉 Lua 脚本 + ZSET 维护 + singleflight 重建约 80 行。

> 即使砍掉 Lua 路径，**Redis 本身也去不掉** —— `brand:{id}` 实体缓存、categories、
> products 的 L1/L2/Bloom、marketing、dashboard 都在用它。砍掉的只是这一条路径。

### 2.3 对 products / skus 的含义

现有这套「ZSET 存全部 ID + MGET 全部实体」的写法，**Redis 内存是 O(N)**：
100 个品牌无所谓，但 10 万 SKU 的全量 JSON 缓存会直接撑不住。所以不要照抄：

| 模块 | 建议 |
|------|------|
| brands（100 条字典表、C 端高频） | 保留 Redis 路径，成本极低 |
| products（2025 条，组合筛选多） | **已落地**：类目/品牌/状态仍走 ZSET，名称 + 价格区间走 ES。索引只存检索投影，不存全量实体（见 §11） |
| skus（几十万+） | **不要照抄全量缓存**；ES 扛组合条件，Redis 只做热点查询缓存 + 单条缓存 |

ES 单节点实测天花板约 6.6k RPS。C 端筛选流量涨上去后，Redis 层可以顺势扩展成
**热点筛选组合的查询缓存**（按查询签名缓存 ES 结果集），而不只是无筛选时的兜底 ——
这个两层结构本身是可生长的，不是 brands 专用。

---

## 三、环境准备

### 3.1 ES 版本与客户端

| 组件 | 版本 | 说明 |
|------|------|------|
| Elasticsearch | 7.17.4 | 本机 `elasticsearch-full`（Homebrew），`ES_HOME=/opt/homebrew/Cellar/elasticsearch-full/7.17.4/libexec` |
| Go 客户端 | `github.com/elastic/go-elasticsearch/v7 v7.17.10` | 必须用 **v7** 分支匹配 ES 7.x |
| 分词插件 | `analysis-ik 7.17.4` | 中文分词 |

### 3.2 安装 IK 分词插件

IK 已从 `medcl/elasticsearch-analysis-ik` 迁移到 `infinilabs/analysis-ik`，
**GitHub Releases 不再提供插件包**，且部分网络环境下 `github.com` 不可达。可用官方 CDN：

```bash
# get.infini.cloud 会按版本号返回对应插件包（替换成自己的 ES 版本）
curl -sSL -o /tmp/ik.zip https://get.infini.cloud/elasticsearch/analysis-ik/7.17.4

# 用 ES 自带插件 CLI 安装（不要手工解压，CLI 会校验版本与权限）
/opt/homebrew/opt/elasticsearch-full/bin/elasticsearch-plugin install --batch "file:///tmp/ik.zip"

# 重启 ES（brands Pilot 用的是 brew services）
brew services restart elasticsearch-full

# 验证插件已加载
curl -s "http://localhost:9200/_nodes/plugins?pretty" | grep -i analysis-ik
curl -s -X POST "http://localhost:9200/_analyze" -H 'Content-Type: application/json' \
  -d '{"analyzer":"ik_max_word","text":"苹果 Apple 华为"}'
# => ["苹果","apple","华为"]
```

**注意**：插件包内的 `plugin-descriptor.properties` 必须声明
`elasticsearch.version=7.17.4`，版本不匹配 ES 会拒绝启动。

### 3.3 配置项

`manifest/config/config.yaml`（该文件被 `.gitignore` 忽略，属本地配置）：

```yaml
elasticsearch:
  enabled: true                      # 总开关：false 时全部查询回落 MySQL
  addresses:
    - "http://127.0.0.1:9200"
  username: ""
  password: ""
  indexPrefix: "eshop_"              # 索引/别名前缀
  timeout: "3s"                      # 单次请求超时
  maxResultWindow: 10000             # from+size 超过此值直接回落 DB
  writeRefresh: "wait_for"           # 双写刷新策略，见 §6.3
  bulkChunk: 500                     # 全量重建时单次 _bulk 的文档数上限（几十万级必须分批）
```

---

## 四、索引建模

### 4.1 为什么单靠 IK 不够 —— 必须叠加 ngram

MySQL 侧原语义是 `name LIKE '%x%'`（**任意子串**），而 IK 是**按词**切分：

```
ik_max_word 切 "苹果"  => ["苹果"]        # 搜 "苹" => ["苹"]，匹配不上！
```

如果只用 IK，`?name=苹` 会从「有结果」变成「无结果」，属于**功能回退**。
因此每个需要子串能力的 text 字段都挂一个 ngram 子字段：

| 子字段 | 索引侧分析器 | 检索侧分析器 | 作用 |
|--------|--------------|--------------|------|
| `name` | `ik_max_word` | `ik_smart` | 中文词级相关度（可 `_score` 排序） |
| `name.ngram` | `ngram_index`（1..32 全 n-gram） | `ngram_search`（keyword 整体成 token） | **等价 `LIKE '%x%'`**，兼顾英文前缀/部分词 |
| `name.kw` | — | — | `term` 精确匹配 |

实测切分效果：

```
索引侧 ngram_index "苹果"      => ['苹', '苹果', '果']
索引侧 ngram_index "Apple"     => ['a','ap','app','appl','apple','p','pp','ppl']
检索侧 ngram_search "苹"       => ['苹']      → 命中「苹果」
检索侧 ngram_search "app"      => ['app']     → 命中「Apple」
```

**检索侧用 keyword tokenizer 是关键**：若用 standard，`"苹果"` 会被拆成 `苹 OR 果`，
命中范围过宽；用 keyword 则整体成一个 token，精确对应索引里的 n-gram，语义干净。

> `max_ngram_diff` 默认只允许 1，`min_gram=1 / max_gram=32` 必须在 settings 里显式放开，否则建索引直接报错。

### 4.2 文档结构：`_source` 存完整实体

ES 里存的 `_source` 就是**完整的 `entity.Brands` JSON**，因此 ES 返回的记录与 DB 查询**字段完全一致**
（含 `created_at` / `updated_at` / `deleted_at`），接口层无需做任何裁剪或补齐。

```go
"dynamic": false,                      // 未登记字段只进 _source，不动态建索引，避免 mapping 漂移
"properties": {
    "id":           {"type":"long"},
    "name":         search.TextWithNGram(),    // ik + ngram + kw 三路子字段
    "english_name": search.TextWithNGram(),
    "first_letter": {"type":"keyword"},
    "sort_order":   {"type":"integer"},        // 仅回显/备用，列表排序不再使用
    "status":       {"type":"integer"},
    "created_at":   search.StoredOnly(),       // {"type":"object","enabled":false}
    ...
}
```

**时间字段为什么用 `enabled:false`**：GoFrame `gtime` 序列化成 `"2026-09-06 22:40:11"`，
不是 ES 默认接受的 ISO8601。若按 `date` 映射，动态映射可能解析失败并拒收文档。
用 `object + enabled:false` 表示「存进 `_source` 原样回显，但完全不索引不解析」，
彻底规避格式问题——这些字段本来也不参与检索。

### 4.3 查询构造

```go
filters: first_letter / status  → term（精确筛选，不参与算分）
must:    name                   → bool.should[ name.ngram, english_name.ngram ]
                                  minimum_should_match: 1
sort:    id asc                   // 与 DB 分支 ORDER BY id ASC 完全一致
```

---

## 五、零停机重建（索引 + 别名）

**不直接往检索中的索引写全量数据**，而是：

```
1. 新建 eshop_brands_<时间戳>        ← 全新物理索引
2. 全量 bulk 写入
3. refresh
4. 原子切换别名：add 新索引 → remove 旧索引   ← 先 add 后 remove，任意时刻别名都有效
5. 删除旧索引
```

检索侧**始终通过别名 `eshop_brands`**，切换期间不出现空窗。

| 名称 | 示例 | 用途 |
|------|------|------|
| 别名 | `eshop_brands` | 检索 + 单条双写都走它 |
| 物理索引 | `eshop_brands_20260928182648` | 实际存储，重建时整体替换 |

单条双写也写**别名**（别名指向单一索引时可直接写入），因此重建后新索引自动接管写入。

---

## 六、数据一致性与可用性

### 6.1 双写

`Create` / `Update` / `Delete` 在完成 DB 写入与 Redis 失效后，同步写一次 ES：

```go
addBrandToIndex(...)        // 既有：Redis ZSET
syncBrandDoc(...)           // 新增：ES 双写（Update 会重新取库，保证 ES 与 DB 完全一致）
```

**DB 是唯一真相源**，双写失败只记 WARN 日志，不影响主流程返回。

### 6.2 启动自愈

启动预热管线新增 `brands_es` 阶段，逻辑是「**对账后才重建**」：

```
ES 文档数 == DB 行数  →  跳过重建（日志：与 DB 一致，跳过重建）
索引/别名不存在        →  首次全量重建
文档数与 DB 不一致     →  全量重建
读 ES 失败            →  全量重建（并进入熔断）
```

这是双写丢数据、ES 被清空、mapping 升级的统一修复入口。实测启动日志：

```
品牌 ES 索引与 DB 一致（100 条），跳过重建
warmup stage "brands_es" done: 0 items
```

### 6.3 熔断 + DB 降级

```
请求 → search.Available()?
        ├─ false（未启用 / 冷却期内） → 直接 listBrandsFromDB()
        └─ true  → searchBrandsES()
                     ├─ 成功 → 返回
                     └─ 失败 → markDown()（熔断 5s）+ listBrandsFromDB()
```

**ES 不是唯一真相源，因此不是硬依赖**：`enabled:false`、ES 进程挂掉、网络不通、索引没建，
任一情况下接口都会自动回落 DB，返回结果依然正确。

实测（ES 指向不存在的 `127.0.0.1:9399`）：

| 请求 | HTTP | 耗时 | 结果 |
|------|------|------|------|
| `?first_letter=A` | 200 | 8.6ms | total=5，与 ES 路径一致 |
| `?name=苹果` | 200 | 4.5ms | total=1 |
| `?status=1&first_letter=S` | 200 | 3.1ms | total=8 |
| 后续请求（熔断生效） | 200 | 2~3ms | 不再等 ES 超时 |

日志确认了三层降级：

```
[WARN] elasticsearch 故障，5s 内降级到 MySQL: connection refused        ← 熔断
[WARN] 读取品牌 ES 索引状态失败，执行全量重建: connection refused      ← 预热
[WARN] 品牌 ES 检索失败，降级 DB 查询: connection refused              ← 单次请求
```

### 6.4 深分页保护

`from + size > max_result_window(10000)` 时返回 `ErrDeepPaging` 并回落 DB，
避免 ES 深分页报错。实测 `?first_letter=A&page_size=20000` 正常返回 5 条（走 DB）。

---

## 七、代码结构

```
internal/search/                    ← 可复用底座（与业务无关）
├── client.go      客户端、配置、熔断、别名/索引命名、schema 版本前缀
├── analysis.go    IndexSettings / TextWithNGram / NGramMatch / StoredOnly 等构件
└── index.go       EnsureIndex / BulkIndex(分批) / Search / CollectStats /
                   SwapAlias / IndexDoc / DeleteDoc ...

internal/logic/brands/es.go         ← brands Pilot
├── brandMapping()      索引结构（_source 存完整实体，直接用于响应）
├── ReindexBrands()     全量重建（新索引 → bulk → 切别名 → 删旧索引）
├── WarmupES()          启动自愈对账（结构版本 + 文档数）
├── syncBrandDoc()      双写
└── searchBrandsES()    检索

internal/logic/products/es.go       ← products 正式落地
├── productMapping()      索引结构（只存检索投影）
├── buildProductDocs()    文档构建（复用 listProductStats 取价格区间）
├── ReindexProducts()     分批读库 + 分批 bulk 重建
├── WarmupES()            对账：结构版本 + 文档数 + 价格合计
├── SyncSearchDoc()       供 skus 模块回调的导出方法
├── syncProductDoc()      双写
└── searchProductIDsES()  检索（只产出有序 ID，响应仍由 MySQL 组装）
```

`internal/logic/brands/brands.go` 的改动点：

```go
// List：有筛选时优先 ES，失败降级
if search.Available(ctx) {
    if esRes, esErr := searchBrandsES(ctx, req, page, size); esErr == nil {
        return esRes, nil
    }
    g.Log().Warningf(ctx, "品牌 ES 检索失败，降级 DB 查询: %v", esErr)
}
return listBrandsFromDB(ctx, req, page, size)   // 原筛选逻辑原样保留
```

**无筛选路径没有被替换**，继续走既有的 Redis ZSET + Lua 一次往返（为什么保留见 §2.1）。
唯一的后续改动是它的「无筛选」判定条件，即 §10.2 的 `?status=0` 修复。

`internal/logic/products/products.go` 的改动点：

```go
// List：类目/品牌/状态组合仍走 ZSET；名称/价格筛选优先 ES，失败降级
if req.Name == "" && req.PriceMin == 0 && req.PriceMax == 0 {
    if result, err := s.listFromZSET(ctx, req, cursorId, size); err == nil {
        return result, nil
    }
}
if search.Available(ctx) {
    if ids, hasMore, err := searchProductIDsES(ctx, req, cursorId, size); err == nil {
        return buildListFromIDs(ctx, ids, hasMore)   // 响应组装与 DB 路径共用
    }
}
return s.listFromDB(ctx, req, cursorId, size)
```

---

## 八、运维命令

```bash
./main reindex                     # 重建全部已接入实体（brands + products）
./main reindex --entity=brands     # 只重建 brands
./main reindex --entity=products   # 只重建 products
```

> GoFrame `gcmd` 会把多余的位置参数当成**多级命令名**，
> 所以实体名只能用 `--entity=` 选项传，`main reindex brands` 会报 command not found。

不重启也能触发重建的另一种方式：直接删掉 ES 索引/别名，重启时对账逻辑会自动重建。

```bash
curl -X DELETE "http://localhost:9200/eshop_brands"
curl -X DELETE "http://localhost:9200/eshop_products"
```

查看索引状态（注意物理索引名带结构版本，如 `eshop_products_v2_2026...`）：

```bash
curl -s "http://localhost:9200/_cat/indices?h=index,docs.count" | grep -E "brand|product"
curl -s "http://localhost:9200/_cat/aliases?h=alias,index"
```

---

## 九、验证结论

### 9.1 与 MySQL 基线逐条对账（ES 路径 vs 直查 DB）

| 用例 | ES total | DB total | ES ids | 结论 |
|------|----------|----------|--------|------|
| 首字母 A | 5 | 5 | `[50,36,1,15,11]` | 一致 |
| 首字母 S | 8 | 8 | 一致 | 一致 |
| 状态 1 | 100 | 100 | 一致 | 一致 |
| 首字母 A + 状态 1 | 5 | 5 | 一致 | 一致 |
| 中文整词 苹果 | 1 | 1 | `[1]` | 一致 |
| 中文子串 苹 | 1 | 1 | `[1]` | 一致 |
| 中文子串 小米 / 华为 | 1 | 1 | 一致 | 一致 |
| 翻页 A 第 1/2 页 size3 | 5 | 5 | 一致 | 一致 |

**10/10 用例的 ids 顺序与 total 全部一致**（当时排序为 `sort_order asc, id desc`；
后续已统一改为 `id asc`，见 §11.1）。

### 9.2 多字段检索（本次新增能力）

| 查询 | ES total | DB total | 说明 |
|------|----------|----------|------|
| `name=Apple` | 1 | 0 | ES 命中 `english_name` |
| `name=app` | 1 | 0 | 英文部分词（ngram） |
| `name=Huawei` | 1 | 0 | 中英混合命中 |
| `name=Samsun` | 1 | 0 | 前缀不完整也能命中 |

### 9.3 写入路径

| 操作 | 验证点 | 结果 |
|------|--------|------|
| Create | ES 文档立即可见、可检索 | 通过 |
| Update | ES 文档被替换（旧名查不到、新名查得到） | 通过 |
| Update 失败 | DB 唯一键冲突时 ES 保留旧文档（与 DB 一致） | 通过 |
| Delete | ES 文档立即移除，查不到 | 通过 |

---

## 十、行为变化与已知问题（务必知悉）

### 10.1 行为变化：检索范围变宽

`?name=` 现在**同时匹配 `name` 与 `english_name`**，而改造前只匹配 `name`。
这是「多字段检索」的目标，但属于**兼容性变化**：例如 `?name=Apple` 以前返回 0 条，现在返回「苹果」。

如需严格保持旧语义，把 `searchBrandsES` 里 `english_name.ngram` 那个 should 子句删掉即可。

### 10.2 `?status=0` 筛选失效（已修复）

改造前的无筛选判定是 `req.Status == nil || *req.Status == 0`，而 `Status` 是指针，
**nil 才代表未传该条件**，于是 `?status=0` 被判成「无筛选」，走进缓存分支返回全部品牌：

```go
// 修复前：?status=0 => total=100（库里 100 条全为 status=1）
if req.Name == "" && req.FirstLetter == "" && (req.Status == nil || *req.Status == 0) {

// 修复后：只判 nil，?status=0 正确进入筛选分支
if req.Name == "" && req.FirstLetter == "" && req.Status == nil {
```

这是 ES 接入**之前就存在**的缺陷，已单独提一个修复（未混在 ES 接入提交里），并补了
`tests/test_api.py` 回归用例 4.2~4.4。实测：修复前 `?status=0` 返回 100 条，
修复后返回 0 条；造一条 `status=0` 记录后返回 1 条且只命中该条，`?status=1` 不包含它，
DB 降级路径结果一致。

> `categories` 存在同一处缺陷（`categories.go:39` 写作 `*req.Status <= 0`，比 brands 更宽），
> **尚未修复**，待确认类目页筛选下拉是否用 `status=0` 表示「全部」后再动。

### 10.3 双写延迟

`writeRefresh: wait_for` 会让 Create/Update/Delete **最多多等一个 refresh_interval（默认 1s）**，
换来「写完立即可检索」。若不能接受写延迟，改成 `writeRefresh: ""`（默认 async refresh），
代价是新增/改名后约 1s 内检索结果可能不同步。

> 这里删和改必须用**同一个**刷新策略。最初只在写入用了 `wait_for`、删除用了默认异步，
> 结果「刚删除的品牌」会在约 1s 内继续出现在检索结果里——已修复。

### 10.4 深分页上限

`from + size` 不得超过 10000（ES `max_result_window`）。超出会回落 DB。
products/skus 铺开时若需要深翻页，应改用 `search_after`。

### 10.5 ES 是加速器，不是真相源

任何时候都可以 `enabled: false` 一键关停，业务不受影响。ES 数据丢失也可通过
`./main reindex` 或启动自愈从 MySQL 完整恢复。

### 10.6 顺带发现：products 写接口的 `created_by` / `updated_by` 类型不匹配（未修）

接入 products 时发现，`POST /api/v1/products` 在不传 `created_by` 时会直接失败：

```
Error 1366 (HY000): Incorrect integer value: '' for column 'created_by' at row 1
```

根因是 API 把这两个字段声明成了 `string`，而 DB 列是 `bigint NOT NULL`：

| 位置 | 类型 |
|------|------|
| `api/products/v1/products.go`（`ProductsCreateReq.CreatedBy`、`ProductsUpdateReq.UpdatedBy` 等 4 处） | `string` |
| `sp_products.created_by` / `updated_by` | `bigint NOT NULL` |

空串被原样写进 bigint 列即报错；只有传数字字符串（如 `"created_by": "1"`）才成功。
这是接入 ES **之前就存在**的缺陷（`git show HEAD` 可确认写法相同），
与本次改动无关，已单独记录，未混在 ES 改动里修。

---

## 十一、products 接入实录（已完成）

brands 是 Pilot，products 才是这套方案真正兑现价值的地方：**它原本唯一落 DB 的部分（名称搜索 + 价格区间）现在由 ES 承接**，
而类目/品牌/状态因为 ZSET 是按筛选组合分 key 的，早已被 Redis 覆盖。

### 11.1 与 brands 的关键差异

| 维度 | brands（100 行） | products（2025 行） |
|------|------------------|---------------------|
| 分页方式 | offset（`page` / `page_size`） | **游标 keyset**（请求 `cursor` = base64(id) + `size`，响应 `next_cursor` + `has_more`，不返回 `total`，与 orders 列表同一份契约） |
| 排序 | **`id ASC`**（与 DB、缓存 ZSET 三路统一） | **`id DESC`**（游标分页固定倒序） |
| 缓存分工 | 无筛选走 Redis ZSET | **类目/品牌/状态组合**走 ZSET（每个组合一个 key） |
| 真正落 DB 的 | 全部筛选 | **只有 `name` 搜索 + 价格区间** |
| 响应字段 | ES 直接返回实体 | 需 SKU 聚合出 `price_min`/`price_max`/`total_stock` |
| 同步面 | 仅商品写入 | 商品写入 **+ SKU 价格变更** |

因为两条路径都是「先产出有序 ID，再统一 hydrate」（`listByIDs` + `listProductStats` + `buildListResponse`），
ES 只需接管「选出哪些 ID、什么顺序」，响应组装完全复用既有逻辑。

### 11.2 索引只存「检索投影」，不存全量实体

products 字段多且含 `images` / `seo_*` 等大字段，而列表响应统一由 MySQL 侧组装，
因此 ES 文档只保留参与筛选/排序的字段（`id`、`name`、`subtitle`、类目/品牌/状态、
销量/评分、`price_min`/`price_max`）。

这样索引体积与「可检索字段数」成正比，而不是与实体宽度成正比 —— 这是 products 与 brands 的重要区别
（brands 的 `_source` 存的是完整实体，因为它直接用于响应）。

### 11.3 游标分页 = `range(id < cursor)`，不需要 `search_after`

DB 侧是 `WHERE id < cursorId ORDER BY id DESC`。因为**排序键本身就是 id**，
ES 里用 `range: {id: {lt: cursorId}}` + `sort: [{id: desc}]` 即可，
无需 `search_after` 的排序值数组。取 `size + 1` 条判断 `has_more`，与 DB 路径一致。

### 11.4 价格区间：反范式 SKU 聚合，语义精确对齐

DB 的价格筛选是：

```sql
id IN (SELECT product_id FROM sp_skus WHERE price >= PriceMin)
AND id IN (SELECT product_id FROM sp_skus WHERE price <= PriceMax)
```

即「存在一个 SKU ≥ min」且「存在一个 SKU ≤ max」。把 SKU 聚合出的区间反范式进文档后，
这两条等价于：

```
price_max >= PriceMin  AND  price_min <= PriceMax
```

**语义精确一致**（都是「区间有交集」），不是近似。两个字段复用 `listProductStats` 计算，
保证与列表响应里的 `price_min`/`price_max` 同源。

### 11.5 价格漂移：为什么只比文档数不够

`price_min`/`price_max` 来自 `sp_skus` 聚合，而 **SKU 改价不会改变商品文档数**。
只做计数对账的话，改价后的价格区间筛选会一直用过期数据且无人察觉。

因此 products 的启动对账额外比对**价格合计指纹**：

```sql
SELECT COUNT(*), SUM(mn), SUM(mx) FROM (
  SELECT product_id, MIN(price) mn, MAX(price) mx
  FROM sp_skus WHERE deleted_at IS NULL GROUP BY product_id
) t
```

与 ES 侧 `sum(price_min)` / `sum(price_max)` 对比，任一不符即全量重建。
实测日志：`商品 ES 索引与 DB 一致（2025 条，价格合计 159728226/210931209），跳过重建`。

SKU 价格可以从**两个模块**修改，都必须触发重新同步：

| 来源 | 处理方式 |
|------|----------|
| `internal/logic/products/repo_skus.go`（商品编辑器的批量 SKU） | 由商品侧 `CreateFull` / `UpdateFull` / `BatchCreateSKUs` 在事务提交后统一同步 |
| `internal/logic/skus`（独立 SKU 管理） | 通过 `service.Products().SyncSearchDoc(ctx, productId)` 回调 |

`logic/skus` 的 Update/Delete 会先取出 `product_id`（Delete 必须在软删除**之前**取），
再回调同步；不这样做的话，从 SKU 管理页改价不会反映到价格筛选上。

### 11.6 索引结构版本（schema version）

**analyzer / mapping 变更不会改变文档数**，因此计数对账与价格合计都发现不了 ——
旧索引会被一直沿用，改动看起来「没生效」。这类问题排查成本很高。

解决办法：把结构版本编进**物理索引名**。

```
别名:      eshop_products
物理索引:  eshop_products_v2_20260928200053
                     ^^^ 结构版本
启动对账先比对「别名指向的索引前缀是否为 eshop_products_v2_」，
不符即全量重建。
```

> **改 mapping / analyzer / 字段语义时，必须递增对应模块的 schema 常量**
> （`brandIndexSchema` / `productIndexSchema`），否则改动不会生效。

实测有效：把 ngram 检索侧分析器从 `keyword` 改为「与索引侧同一 tokenizer + `operator=and`」后，
重启自动把 `eshop_products_2026...` 换成了 `eshop_products_v2_2026...`，无需人工干预。

### 11.7 一个必须修正的坑：ngram 检索侧分析器不能是 keyword

最初的实现里，`ngram` 子字段的**检索侧**分析器用了 `keyword`（整串一个 token）。
这只对「不含空格/标点的单段查询」有效：

| 查询 | keyword 检索侧（错误） | 同一 tokenizer + and（正确） |
|------|------------------------|------------------------------|
| `手机` | 命中 ✓ | 命中 ✓ |
| `三星e 青春版` | **搜不到**，且会返回一堆只命中副标题无关词的商品 | 精确命中 1 条 |
| `Arc'teryx`（brands） | **搜不到**（撇号把索引侧切成了 arc + teryx） | 命中 ✓ |

原因是索引侧的 ngram tokenizer 会在空格/标点处切断，整串 token 在索引里根本不存在。
改为**检索侧与索引侧使用同一个 tokenizer，并在查询里加 `operator: and`**，
即可还原「查询串作为连续子串出现」的语义。

这个修复同时修好了 brands 的潜在缺陷（含撇号的品牌名）。

### 11.8 验证结论

| 验证项 | 结果 |
|--------|------|
| 16 个筛选/搜索用例 vs MySQL 基线 | ids 完全一致（含多词、价格区间、类目+品牌+状态组合） |
| 游标翻页走完全量 2025 条 | 41 页、无重复、无遗漏，与 DB 全量 id 序列一致 |
| SKU 改价联动 | 改价 → ES 文档价格区间同步 → 价格筛选命中集合翻转，还原后复原 |
| ES 不可用降级 | 结果与 ES 路径一致（6~17ms），三层降级日志齐全 |
| 双写 Create/Update/Delete | ES 文档即时可见/更新/移除 |
| 启动对账 | 计数 + 价格合计一致时正确跳过重建 |

### 11.9 行为变化（products，务必知悉）

1. **多词查询语义变化（预期且更优）**：ES 按词 AND 匹配，不再要求字面子串。
   例如 `三星 智能`：DB `LIKE '%三星 智能%'` 返回 **0** 条，
   ES 返回 **4** 条（三星智能手机 / 三星智能彩光灯 / 三星智能灯泡 …）。
   用户输入空格分词时，ES 的行为才符合直觉。
2. **只按 `name` 匹配，刻意不含 `subtitle`**：subtitle 是 IK 词级检索，
   与 name 的 ngram 子串语义不同，混进 `should` 会让「按商品名搜」返回
   一堆只命中副标题无关词的商品（实测一次错配 5 条）。多字段检索应配合
   按 `_score` 排序，属后续工作。
3. **ES 降级期间**，多词查询会退化为 MySQL 字面 `LIKE`（返回更少，甚至 0 条）。
   这是降级模式的预期差异。

### 11.10 未来方向

- **`skus` 业务接入**：目前 SKU 只通过商品间接进 ES（价格区间），独立 SKU 检索（按 sku_code/条码）尚未接入。
- **相关度/销量/价格排序**：需要先把 `listByIDs` 改成**按入参顺序**返回。
  它现在固定 `ORDER BY id DESC`，会覆盖 ES 算出的相关性顺序 —— 不先改这里，`_score` 排序无法生效。
- **高频写入改异步**：products/skus 写入量远大于 brands，同步双写（`wait_for`）可能需要换成异步或消息队列。

---

## 十二、一句话总结

**brands**：无筛选走 Redis（实测保留，不是 ES 的过渡态），有筛选走 ES（多字段/子串检索），
ES 挂了自动回落 MySQL，写入双写 + 启动自愈对账，重建靠别名切换零停机。

**products**：类目/品牌/状态仍走 Redis ZSET（按筛选组合分 key），
`name` 搜索 + 价格区间交给 ES —— 这正是它原本唯一落 DB 的部分。
索引只存检索投影，游标分页用 `range(id < cursor)`，价格区间由 SKU 聚合反范式且语义精确对齐；
对账除文档数外还比对价格合计以发现改价漂移；结构版本编进物理索引名，mapping/analyzer 变更可自动重建。

分工记两条轴：**查询形态决定走谁，流量画像决定要不要留 Redis 层**（§2.2）。
改 mapping/analyzer 时**必须递增 schema 常量**，否则索引结构变更不会生效（§11.6）。
