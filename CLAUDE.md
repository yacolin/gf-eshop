# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供本仓库代码工作的指导。

## 项目概述

基于 GoFrame (v2.6.1) 的电商后端 API，使用 **MySQL + Redis + Elasticsearch**。
单仓库 GoFrame 项目，`internal/logic/` 下现有 **43 个业务模块**（品牌、类目、商品、SKU、
库存、订单、支付、购物车、用户、商家、营销、评价、权限等）。

查询侧是**两层架构**，不要把两者混为一谈：

| 层 | 承担 | 说明 |
|----|------|------|
| **Redis** | 热点 / 简单查询 | 单条实体缓存、列表 ZSET、products 的 L1+Bloom+L2 |
| **Elasticsearch** | 组合筛选 / 文本检索 | 多字段子串检索、价格区间、排序、翻页 |

ES 是**加速器，不是唯一真相源**：任何一步失败都会熔断并自动回落 MySQL 原有查询，
`elasticsearch.enabled: false` 可一键关停。完整方案见
[`docs/elasticsearch-search-guide.md`](docs/elasticsearch-search-guide.md)。

## 开发命令

```bash
# 开发服务器，支持热加载（需要 gf CLI）
make run

# 生产构建
make build          # 通过 gf build 构建二进制文件
go build -o main .  # 或使用纯 go build

# 构建后运行服务器
./main

# 停止正在运行的服务器
make stop           # 杀死 :8000 端口上的进程

# 重建 Elasticsearch 索引
./main reindex                     # 重建全部已接入实体（brands + products）
./main reindex --entity=brands     # 只重建 brands
./main reindex --entity=products   # 只重建 products
# 注意：GoFrame gcmd 会把多余的位置参数当成「多级命令名」，
#       实体名只能用 --entity= 传，`main reindex brands` 会报 command not found

# 代码生成（需要 gf CLI）
gf gen dao          # 根据数据库表重新生成 DAO/DO/Entity
gf gen service      # 重新生成服务接口文件
gf gen ctrl         # 根据 API 定义重新生成控制器

# GoFrame CLI 检查/安装
make cli            # 下载并安装 gf CLI
```

## 架构概览

### 层级结构

```
请求 → ghttp.Server → Controller（薄层）→ Service 接口 → Logic → DAO → MySQL
                                                          ↕        ↕
                                                    Redis（缓存）  Elasticsearch（检索）
                                                                    ↑
                                        internal/search（ES 共用底座）
```

- **`api/<模块>/`** — API 接口定义：`IBrandsV1` 接口 + `v1/` 包下的 Req/Res 结构体。通过 `g.Meta` 结构体标签声明式定义路由（path、method、summary）。
- **`internal/controller/`** — 薄控制器层，每个方法直接委托给对应的 service 方法，不做业务处理。
- **`internal/service/`** — 服务接口 + 注册模式（`RegisterX()` 注册，`X()` 获取）。
- **`internal/logic/`** — 业务逻辑实现。每个模块在 `init()` 中通过 `service.RegisterX()` 注册自身实现。
- **`internal/dao/`** — 数据访问对象。`dao/internal/` 为自动生成（禁止手动修改）。`dao/` 包装层可扩展自定义查询方法。
- **`internal/model/entity/`** — 实体结构体（`gf gen dao` 自动生成），用于查询结果扫描。
- **`internal/model/do/`** — DO 结构体（自动生成），用于 Insert/Update 传参。
- **`internal/search/`** — **ES 共用底座，与业务无关**：客户端与熔断、索引/别名管理、
  分批 bulk、索引统计、schema 版本前缀、ngram 检索构件。业务模块只需提供自己的
  mapping 与文档转换（参考 `internal/logic/brands/es.go`、`internal/logic/products/es.go`）。
- **`internal/verifycode/`** — **验证码共用底座，与业务无关**：生成、`sha256` 指纹存储、
  重发冷却、单收件人/单 IP 日限额、失败次数上限、防邮箱枚举的「静默跳过」策略，
  以及 `Sender` 渠道接口。邮件与短信只是渠道实现（`sender_email.go` / `sender_sms.go`）；
  **接新渠道只需实现 `Sender`**，风控与校验完全复用。业务侧只提供场景合法性、
  收件人↔账号映射与会话签发（参考 `internal/logic/user_auth/verification_code.go`）。
  配置分三段：`email.*`（渠道参数）、`verifycode.*`（渠道无关风控）、`sms.*`。
  详见 [`docs/email-verify-code-guide.md`](docs/email-verify-code-guide.md)。
- **`main.go`** — 入口：导入 MySQL + Redis 驱动，导入 logic 包（触发 `init()`），运行 `cmd.Main.Run()`。
- **`internal/cmd/cmd.go`** — 路由设置、中间件、**缓存与 ES 索引预热**、`reindex` 子命令。

### 路由注册

```
/        → Hello（示例接口）
/api/v1/ → 无鉴权：brands、categories、category_brands、products、skus、
                     attributes、attribute_values、inventories、merchants…
         → authMiddleware：staff、departments、permissions、notification、
                            ws、user_admin、user_levels、user_points
         → authMiddleware + RequireAdmin：roles、points_rules、level_rules
         → UserAuthMiddleware：user、user_auth、address、carts、orders、
                               payments、marketing
         → /api/v1/ws（WebSocket 升级，token 走查询参数，不走上面的中间件）
```

统一响应与错误处理用 **`middleware.ErrorHandler`**（`internal/middleware/response.go`，
它替代了 `ghttp.MiddlewareHandlerResponse`）。

### 添加新 CRUD 模块

完整流程参考 `docs/curd-workflow.md`。

1. 在 `hack/config.yaml` 的 `tables` 字段添加新表名
2. 运行 `gf gen dao`（生成 DAO、entity、DO）
3. 创建 `api/<模块>/<模块>.go` 接口定义
4. 创建 `api/<模块>/v1/<模块>.go` Req/Res 结构体（使用 `g.Meta` 标签）
5. 创建 `internal/service/<模块>.go` 接口 + Register/X 函数
6. 创建 `internal/logic/<模块>/<模块>.go`，包含 `init()` 注册 + 业务逻辑
7. 创建 `internal/controller/<模块>/`（new.go + v1 实现文件）
8. 在 `internal/cmd/cmd.go` 中注册路由
9. 在 `main.go` 中导入 logic 包

### 缓存模式

**brands / categories（Cache-Aside + 列表 ZSET）**
- **缓存键**：`brand:<id>`、`category:<id>`；列表用 ZSET（`brand:ids`），分数编码排序权重
- **TTL**：10 分钟
- **列表读取**：一次 Lua 往返完成 `ZCARD + ZRANGE + EXISTS + MGET`；
  缓存不完整时经 singleflight 单次重建
- **写入/删除**：使单条缓存与列表 ZSET 均失效
- **预热**：`cmd.go` 启动时调用 `Warmup()`

**products（多级缓存）**
- **L1 本地缓存**（60s，带抖动）→ **布隆过滤器** → **L2 Redis**（`product:<id>`）→ DB 兜底
- **列表 ZSET 按筛选组合分 key**：`product:list:ids:cat=<x>:brand=<y>:status=<z>`
  —— 因此类目/品牌/状态的组合查询**已经被 Redis 覆盖**，不需要走 ES

### 搜索层（Elasticsearch）

已接入 **brands** 与 **products**。统一模式：

- **索引 + 别名**：检索与双写都走别名 `eshop_<entity>`；物理索引名为
  `<别名>_<schema>_<时间戳>`，全量重建时「新建 → 分批 bulk → refresh → 原子切别名 → 删旧索引」，
  检索侧无空窗
- **双写**：Create/Update/Delete 同步写 ES，失败只记日志（DB 才是唯一真相源）
- **启动对账自愈**：仅当「结构版本 / 文档数」不一致才重建；
  products 额外比对**价格合计**（SKU 改价不改变文档数，只比计数发现不了漂移）
- **降级**：ES 未启用、熔断冷却中、索引缺失、深分页 → 自动回落 MySQL 原有查询
- **ngram 子字段**：用 `search.NGramMatch(field, query)` 构造查询，
  保证「查询串作为连续子串出现」，语义等价 `LIKE '%x%'`。
  检索侧分析器必须与索引侧使用**同一个 tokenizer** 并加 `operator: and`
- **索引只存检索投影**（参与筛选/排序的字段），不存全量实体；列表响应统一由 MySQL 侧组装

> ⚠️ **修改 mapping / analyzer / 字段语义时，必须递增该模块的 schema 常量**
> （`brandIndexSchema` / `productIndexSchema`）。这类变更**不会改变文档数**，
> 不递增的话启动对账发现不了，旧索引会被一直沿用，改动看起来「没生效」。

### 数据模块

| 模块 | 数据库表 | 接口 | ES |
|------|----------|------|-----|
| brands | `sp_brands` | 列表（分页、筛选）、详情、新增、更新、删除 | ✅ |
| products | `sp_products` | 列表（**游标分页**、名称/价格/类目/品牌/状态筛选）、详情、创建、全量创建/更新、批量生成 SKU、删除 | ✅ |
| skus | `sp_skus` | 列表、详情、按编码查询、新增、更新、删除 | 间接（价格进 products 索引） |
| categories | `sp_categories` | 列表、全部、根类目、子类目、层级、树形、详情、增删改 | — |
| category_brands | `sp_category_brands` | 列表（带品牌 left join）、更新（事务内先删后插） | — |

其余模块（inventories、orders、payments、carts、user、merchants、marketing、reviews、
roles/permissions 等共 40+）见 `internal/logic/`，目前未接入 ES。

### 关键约定

- **查询与写入**：Scan 用 `entity.X`，Insert/Update 用 `do.X`
- **字段安全**：使用 `dao.X.Columns().FieldName` 代替硬编码字符串
- **软删除**：含 `deleted_at` 字段的表自动实现软删除（DELETE 转为 UPDATE）
- **错误处理**：记录不存在用 `gerror.NewCode(gcode.CodeNotFound, ...)`，参数非法用 `gcode.CodeInvalidParameter`
- **空结果**：返回 `make([]*T, 0)` 而非 `nil`，避免 JSON 序列化为 `null`
- **JSON 字段**：全部使用 snake_case（在 `hack/config.yaml` 中配置）
- **分页有两套，不要混用**：
  - brands / categories：`page` + `page_size`（offset，默认 page=1、pageSize=20）
  - products：`cursor` + `size`（keyset；`cursor` = base64(id)，排序固定 `id DESC`）
- **SKU 价格变更必须同步商品索引**：索引里的 `price_min`/`price_max` 来自 SKU 聚合，
  从 `logic/skus` 改价后要调用 `service.Products().SyncSearchDoc(ctx, productId)`，
  否则价格区间筛选会用到过期数据

## 相关文档

| 文档 | 内容 |
|------|------|
| [`docs/curd-workflow.md`](docs/curd-workflow.md) | 新增 CRUD 模块的完整流程 |
| [`docs/cache-optimization-guide.md`](docs/cache-optimization-guide.md) | 缓存优化历程（ZSET+Lua+singleflight、P99 目标） |
| [`docs/elasticsearch-search-guide.md`](docs/elasticsearch-search-guide.md) | ES 接入方案、products 接入实录、行为变化与坑位 |
| [`docs/roadmap.md`](docs/roadmap.md) | 待办事项与优先级 |
| [`docs/perf-workflow.md`](docs/perf-workflow.md) | 压测流程 |
| [`docs/email-verify-code-guide.md`](docs/email-verify-code-guide.md) | 验证码底座：渠道分层、配置、接口、风控与安全设计、到达率 |
