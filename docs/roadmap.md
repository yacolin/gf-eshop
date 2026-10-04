# 项目路线图 / 待办

记录 gf-eshop 当前已知的待办事项、优先级、阻塞条件与前置依赖。

> **维护约定**：完成一项就把状态改为 ✅ 并补上提交号；新增发现的问题请补上
> 「现象 / 根因 / 影响面 / 验证方式」，方便后来人直接接手。

---

## 一、待你决策（阻塞中）

### 1.1 `categories` 的 `?status=0` 同源缺陷

| 项 | 内容 |
|----|------|
| **现象** | `GET /api/v1/categories?status=0` 不带其他筛选时**不会**按「禁用」过滤，而是返回全部类目 |
| **根因** | `internal/logic/categories/categories.go:39` 的「无筛选」判定写成 `(req.Status == nil \|\| *req.Status <= 0)`。`Status` 是指针，`nil` 才代表未传该条件，`<= 0` 把「显式传 0」误判成「无筛选」 |
| **对照** | brands 已有同样缺陷，已修复（见 §三）。categories 写法比 brands 更宽（`<= 0` 而非 `== 0`），负值也会被吞掉 |
| **阻塞原因** | 需先确认**类目页的筛选下拉是否用 `status=0` 表示「全部」**。若是，直接修会把页面变成「只看禁用类目」 |
| **修法** | 与 brands 对齐，只判 `req.Status == nil` |
| **验证** | 造一条 `status=0` 的类目，确认 `?status=0` 只返回它、`?status=1` 不含它；再补 `tests/test_api.py` 回归用例 |

### 1.2 `git push`

本地领先 `origin/main`，**未推送**（无 force-push 需求）。截至订单分表 Phase 2 完成时
领先 11 个提交：`c466686`、`08b0416`，其余 9 个是订单分表相关（Phase 0/1/2、P0 修复、
测试修复、文档登记）。

> 2026-10-04 更正：本节原先写的「停在 `c450516`、领先 3 个（`fb403e0`/`7b9f12b`/`a2e630b`）」
> 已过期 —— 那三个提交早已在 `origin/main` 上。

### 1.3 订单分表方案（待决策）

| 项 | 内容 |
|----|------|
| **背景** | 订单数据预期增长，需决定分表方式与落地节奏。设计按 **10 万单/日、保留 3 年**（3 年约 2.66 亿行 / 266 GiB） |
| **方案** | [`docs/order-sharding-design.md`](order-sharding-design.md) |
| **结论** | **应用层按月分表 + `order_no` 路由**。本库单号已内嵌创建时间（实测 2000/2000 条满足），路由零映射表；且**只分表不分库**，建单四表事务保持原子 |
| **为什么不是原生分区** | 本项目订单查询全部按 `order_no`（WHERE 里没有 `created_at`）→ MySQL 分区裁剪完全失效，每次扫 36 个分区 |
| **建议节奏** | 当前仅 2000 单 / 0.4 MB，**先做 Phase 0（可切换抽象）与 Phase 1（标识生成器）** —— 两者均已完成；Phase 2~4 等触到阈值再启动 |
| **Phase 0 状态** | ✅ **已完成**（2026-10-04）。订单域 SQL 全部收敛进 `internal/logic/orders/repo_*.go`，新增路由层 `shard.go`（配置 `orderShard.mode`，默认 `single`）；`logic/payments`、`logic/dashboard` 已收回对订单表/DAO 的直接访问，改走 `service.Orders()`。**A/B 验证**：与 HEAD 逐字段对比 14 项（详情/列表筛选/翻页/错误路径/看板聚合）**零差异**；写路径另做 30+ 项 SQL 回查。路由契约测试见 `internal/logic/orders/shard_test.go` |
| **Phase 1 阻塞原因** | 需业务表态三件事：① ~~单号内嵌时间（可路由）vs 泄漏单量~~ **已定：内嵌时间**；② ~~列表 `total` / 游标分页的前端改动~~ **后端已实现（见 Phase 2 状态），剩前端适配**；③ ~~粒度按月还是按季~~ **已定：按月** |
| **Phase 1 状态** | ✅ **已完成**（`bfa5525`）。单号改为 `ORD + YYYYMMDDHHMMSS + 每秒序列(6)`（Redis `INCRBY` + 降级 + 唯一键冲突重试）；`tx_orders.id` / `tx_order_items.id` 改为显式全局唯一 ID。**实测**：100 并发建单零重号；单号内嵌时间 ↔ 主键反解时间**秒级一致 100/100**；主键全部 < 2^53-1；可控撞号自动重试成功；Redis 不可用仍可建单；读路径与 Phase 0 对比 14 项零差异 |
| **Phase 1 对设计文档的两处修正** | ① 原方案用 41 位毫秒雪花，会超过 JS 安全整数（2^53-1）导致前端**静默丢精度** → 改为 51 位「分钟(25)\|秒(6)\|序列(20)」布局；② `tx_order_items.id` 也被 `delivery_items` / `after_sales` 引用，**只换订单主键不够**，两张表都要换（原文档遗漏） |
| **Phase 2 状态** | ✅ **已完成**（`b4a3fe1` `ccd5be7`）。新增 `main shard --action=create\|migrate\|verify` 三个命令；建成 2026-08/09 共 8 张分片表，迁入 2000 单 + 5025 明细 + 2000 条老主键映射 + 32 天日汇总；**对账 8 项（行数 + 全字段 CRC32 校验和 + 金额）全部通过，重跑 migrate 新增 0 行**；日汇总与主表直查逐项一致。`tx_orders.id` / `tx_order_items.id` 的 `AUTO_INCREMENT` 已去掉 |
| **Phase 2 的 DDL 归属** | DDL 一律改在 **schema 源仓库 `std-eshop-db`**（本仓库不放 DDL，避免两个来源）：基线 `sql/tx_p0.sql`/`tx_p1.sql`/`tx_p5.sql`，存量库迁移 `sql/migrations/V001__tx_order_sharding.sql`。已过该仓库的 `sql_lint.py` 与 `schema_diff.py`，并用两个临时库验证「全新基线」与「基线 + V001」**终态 schema 逐字节一致**；`pgsql/` 侧**尚未镜像**，已登记为待补 |
| **Phase 2 的意外收获** | `decodeOrderID` 原先对任何正整数都能解出时间，**老自增主键 1..2000 会被解成 2024-01-01** 并路由到不存在的分片 → 已加 `legacyOrderIDMax` 阈值，老主键改走 `tx_order_shard_map`；`ShardOfOrderID` 是统一入口，无法定位时明确报错 |
| **Phase 3 状态** | ✅ **已完成**。双写（`orderShard.dualWrite`）+ 影子读对账（`orderShard.shadowReadPercent`），缺省全关 = 现状。**实测**：建单/状态流转/支付回写三条路径主表与分片逐列一致；影子读 `total=11 diff=0 error=0`；**反向验证**把分片改坏后差异计数 +1、接口仍返回主表值；**重放验证**只改主表复现缺口后重跑 migrate 即追平且无重复行；默认关闭下与 Phase 3 前 HEAD 对比 14 项零差异、套件 36/36 |
| **Phase 3 的实现取舍** | 双写没有逐条改 8 个写函数，而是**在事务内按订单整单镜像**（`INSERT ... SELECT ... ON DUPLICATE KEY UPDATE 全列`）——只挂 3 个编排点、与主写原子、且全列覆盖让它**自愈**（漏掉的更新下次写入自动补上）。迁移也因此获得重放能力（重跑 `--action=migrate` 即重新同步），不需要单独的 replay 命令。`tx_order_logs` 没有 `order_no` 索引，故日志按 `order_id` 定位 |
| **Phase 0 修正了文档一处错误** | GoFrame 的 `Insert` 会**无条件覆盖** `created_at`（`gdb_model_insert.go:311-321`），实测同一次建单四表相差 2~4 ms → 设计文档「三者必然同月」在月份级成立、毫秒级**不成立**，Phase 2 必须按文档 §5.3 的方案 A（`.Unscoped()`）处理跨月窗口 |
| **顺带发现（前置必修）** | `generateOrderNo()` 后 4 位是 `rand(10000)`，@10 万单/日 **日均碰撞期望 5.8 次**、峰值秒内 86% —— 撞 `uk_order_no` 会让用户下单直接失败。修法（每秒序列）与路由天然统一 |
| **改造面** | 业务代码只有 3 个文件（`logic/orders`、`logic/payments`、`logic/dashboard`）共 35 处引用，比看上去小 |

---

## 二、高优先级（无阻塞，建议尽快）

### 2.1 products 写接口的 `created_by` / `updated_by` 类型不匹配

| 项 | 内容 |
|----|------|
| **现象** | `POST /api/v1/products` 不传 `created_by` 时直接失败：`Error 1366: Incorrect integer value: '' for column 'created_by'` |
| **根因** | API 把这两个字段声明为 `string`，DB 列是 `bigint NOT NULL`，空串被原样写入 |
| **位置** | `api/products/v1/products.go` 共 4 处：`ProductsCreateReq.CreatedBy`、`ProductsCreateFullReq.CreatedBy`、`ProductsUpdateReq.UpdatedBy`、`ProductsUpdateFullReq.UpdatedBy` |
| **影响面** | 只有传数字字符串（`"created_by": "1"`）才成功；前端若省略该字段则整个写接口不可用 |
| **性质** | 接入 ES **之前就存在**（`git show HEAD:api/products/v1/products.go` 可确认），与本次改动无关 |
| **修法** | 把 4 处改为 `int64`（或 `*int64`），并同步 `do.Products` 的赋值 |
| **验证** | 不传该字段创建商品应成功；传 `"1"` 应落库为 `1` |

### 2.2 `listByIDs` 改为按入参顺序返回（相关度排序的前置）

| 项 | 内容 |
|----|------|
| **现状** | `internal/logic/products/repo_products.go` 的 `listByIDs` 固定 `ORDER BY id DESC` |
| **问题** | 它会**覆盖** ES 产出的顺序。目前 ES 也按 `id DESC` 排，所以两者一致、没问题；但一旦想用 `_score` / 销量 / 价格排序，顺序会被这行 SQL 抹掉 |
| **修法** | 查回来后按入参 `ids` 的顺序重排（建 map + 重建切片） |
| **风险** | 低。当前行为等价（ES 与 DB 都是 id 倒序），改完结果不变，但解锁了排序能力 |
| **验证** | 改后现有对账用例（16 个筛选用例 + 全量游标翻页）必须仍然全部一致 |

### 2.3 可提交的配置模板（队友拉到代码后 ES 不生效）

| 项 | 内容 |
|----|------|
| **现状** | `manifest/config/config.yaml` 被 `.gitignore` 忽略（`**/config/config.yaml`），ES 配置只存在于本地 |
| **后果** | 队友拉代码后 `elasticsearch.enabled` 取默认值 `false`，**静默回落 MySQL**（不报错，但 ES 完全没生效），排查起来很费时间 |
| **修法** | 在 `manifest/config/` 下提供一份可提交的模板（如 `config.example.yaml`），或在 README/本文档中明确列出需要新增的 `elasticsearch` 段 |
| **验证** | 按模板配置后，启动日志应出现「索引与 DB 一致，跳过重建」而非静默跳过 |

### 2.4 【P0】订单创建接口完全不可用（代码与库中库存数据不一致）✅ 已修复 `bfa5525`

| 项 | 内容 |
|----|------|
| **现象** | `POST /api/v1/orders` 一律 500：`{"code":500,"message":"sql: no rows in result set"}`。**用户根本下不了单** |
| **根因** | `internal/logic/orders/orders.go` 的 `Create` 查库存时写死 `Where(warehouse_id, 0)`，而库里 7942 行 `sp_inventories` **全部是 `warehouse_id = 1`**（0 行）→ GoFrame 的 `Scan` 在无记录时返回 `sql.ErrNoRows`，紧随其后的 `if err != nil { return err }` 把原始错误抛成 500，**作者本意的 `CodeInsufficientStock` 分支永远走不到** |
| **影响面** | 订单域写入链路全线不通。**这解释了为什么 `tx_order_logs` 是空表、2000 条订单只能直接灌库** —— 建单接口从未成功过 |
| **性质** | **接入分表之前就存在**。Phase 0 用 HEAD 建基线服务实测：同一个请求返回**逐字相同**的 500，与订单分表重构无关 |
| **修法（已实施）** | 代码侧：不再写死任何仓库 id，改为按 `sku_id` 选一个有足够可用量的仓库行（`OrderAsc(warehouse_id)`），并显式区分 `sql.ErrNoRows` 与真实 DB 错误。历史背景：库里 `sp_warehouses` 只有「默认仓库」id=1，而 `inventories` 模块里 `WarehouseId > 0` 的写法说明 **0 在本项目是「未指定/全部」的过滤约定**，不是仓库 id |
| **附带修复** | ① SKU 查询有同样的 `sql.ErrNoRows` 问题 →「SKU不存在」现在才会真的返回 `CodeSKUNotFound`；② 预占库存的 `Increment` 一直没检查 `RowsAffected`，并发下会被静默超卖 → 现在按库存不足处理 |
| **验证** | `tests/test_tx_api.py` 建单/详情/列表/状态流转全部通过；100 并发建单全部成功且库存预占正确 |

### 2.5 【P0】数据库约束与代码写入不匹配（一类缺陷，实际触发 4 处）✅ 已修复 `ea94fb2` `bfa5525`

| 项 | 内容 |
|----|------|
| **现象** | 建单过了库存那关后接着 500：`Error 3140: Invalid JSON text ... for column 'tx_order_items.sku_spec'`；创建支付、创建退款也各自失败 |
| **根因** | 表结构与代码写法不匹配，共四类：<br>① **JSON 列写纯文本/空串**：`tx_order_items.sku_spec`、`tx_payments.channel_response`、`tx_refunds.channel_response`、`tx_cart_items.sku_spec` → `Error 3140`<br>② **NOT NULL 无默认值的列整条没写**：`tx_refunds.idempotency_key` → `Error 1364`<br>③ **可空+唯一列写空串**：`tx_refunds.channel_refund_id` → 第二条退款撞 `Duplicate entry ''`（`sp_skus.barcode` 早在 `3d71668` 踩过同一个坑）<br>④ **`Scan` 无记录返回 `sql.ErrNoRows` 被当成 DB 错误**：见 §2.4 |
| **影响面** | 建单、创建支付、创建退款、加购物车**四条写入路径全部不可用**，且一直没有用例覆盖到 —— 与 §2.4 合起来就是「交易域写入链路从未跑通过」 |
| **性质** | 均为 HEAD 既有缺陷。基线服务实测返回同一批错误（`52/3140`、`1364`、`1062`） |
| **修法（已实施）** | ① 规格摘要写 `sku_spec_summary`（varchar），JSON 列留 `NULL`；`channel_response` 不写；`tx_cart_items` 没有 varchar 摘要列，暂留 `NULL` + TODO（结构化快照需从 `sp_sku_specs` 组装）<br>② `idempotency_key` 暂用退款单号占位（接口没有幂等键入参，真正幂等应由调用方传）<br>③ `channel_refund_id` 留 `NULL` 等渠道回执再填 |
| **验证** | `tests/test_tx_api.py` 28/28 通过且**连续两次运行都通过**（可重复执行）；脚本里硬编码的唯一 `transaction_id` 一并修掉 |
| **防复发** | 本次已用 `information_schema` 全量比对「可空+唯一」列与代码里的 `"col": ""` 写法；后续新增写入路径应把建支付/退款/加购物车纳入回归 |

### 2.6 前端适配订单列表的游标分页（`gf-eshop-fe` / `gf-eshop-miniprogram`）

| 项 | 内容 |
|----|------|
| **背景** | 分表后 offset 分页无法跨片归并、`COUNT(*)` 要跨片 fan-out，后端已提供 keyset 游标分页（`ccd5be7`），需要两个前端仓库跟进 |
| **新契约** | `GET /api/v1/orders?size=20[&cursor=<上一页的 next_cursor>]`；响应新增 `next_cursor`（为空表示没有更多）；**游标模式下 `total = -1` 表示未统计**，前端不要再显示「共 N 条」 |
| **兼容性** | `page` + `page_size` 行为完全不变（`total` 仍精确），所以前端可以渐进迁移，不做也可以先不动 |
| **待办** | 把订单列表从「页码 + total」改成「游标 + 上拉加载」；错误处理里补上「非法游标」这一种（后端会返回参数错误而不是空列表） |
| **验证** | 前端翻到最后一页时 `next_cursor` 为空；切换筛选条件后游标重置（不要复用旧游标） |

---

## 三、已完成（归档对照）

| 事项 | 提交 | 说明 |
|------|------|------|
| ES 接入 brands | `7b9f12b` | 走 ES 的有筛选路径 + Redis ZSET 无筛选路径保留；别名零停机重建；双写；熔断降级；启动对账自愈 |
| ES 落地 products | `fb403e0` | 名称搜索 + 价格区间走 ES；类目/品牌/状态仍走 Redis ZSET；游标分页；SKU 改价联动；价格合计对账；schema 版本机制 |
| 修复 `?status=0` 被误判为无筛选（brands） | `7b9f12b` | 修复前返回 100 条、修复后 0 条；补回归用例 4.2~4.4 |
| 修复降级兜底查询缺 `id DESC` | `7b9f12b` | 库中 23 组 `sort_order` 重复，兜底路径排序与缓存路径不一致，实测前 20 条有 6 个位置不同 |
| 修复 ngram 检索侧分析器 | `fb403e0` | `keyword` → 与索引侧同 tokenizer + `operator=and`；修好含空格的商品名（`三星e 青春版`）与含撇号的品牌名（`Arc'teryx`） |
| 建立测试基线 | `fb403e0` | products 16 用例对账、全量游标翻页、价格联动、降级、双写均已实测 |
| 订单分表 Phase 0（可切换抽象） | `d998a47` | 订单域 SQL 收敛进 `repo_*.go` + 分片路由骨架（默认 single，行为零变化）；payments/dashboard 收回直接访问；A/B 对比 14 项零差异 |
| 订单分表 Phase 1（标识生成器） | `bfa5525` | 单号 `ORD+时间+每秒序列`；订单/明细主键改 51 位全局唯一 ID（JS 安全）；Redis 降级 + 撞号重试；100 并发零重号 |
| 修复交易域 4 处写入约束冲突（P0） | `ea94fb2` `bfa5525` | 见 §2.5：JSON 列、NOT NULL、可空+唯一、`Scan` ErrNoRows 四类；建单/支付/退款/加购四条路径恢复可用 |
| 修复订单接口测试（鉴权 + 可重复执行） | `6e7c1ac` | 业务接口改用用户令牌（原先用 staff 令牌必然 401）；`transaction_id` 不再硬编码；28/28 且可重复执行 |
| 订单分表 Phase 2（建分片 + 迁移 + 对账） | `b4a3fe1` | `main shard --action=create\|migrate\|verify`；8 张分片表、2000 单 + 5025 明细迁入、32 天日汇总回填；三重对账（行数 + 全字段 CRC32 + 金额）通过且可重复执行；去掉两个 `AUTO_INCREMENT` |
| 订单列表游标分页（keyset） | `ccd5be7` | `size` + `cursor`（base64(id)，与 products 一致）、`next_cursor`、游标模式 `total=-1`；page 模式行为不变；套件 36/36，前端适配见 §2.6 |
| 订单分表 Phase 3（双写 + 影子读） | 见 §1.3 | 事务内按订单整单镜像（自愈、原子）；影子读确定性抽样 + Redis 计数；反向验证能抓到差异；migrate 变全列覆盖因而具备重放能力 |

---

## 四、中优先级（有价值，不急）

### 4.1 `skus` 业务接入 ES

- **现状**：SKU 只通过「价格区间」间接进入 products 索引，独立 SKU 检索未接入
- **场景**：按 `sku_code` / `barcode` 反查、按规格文本检索
- **做法**：完全照搬 `internal/logic/products/es.go`，在 `reindexTargets` 登记 `skus`
- **注意**：SKU 量级远大于 products（当前 7942），是全项目最适合 ES 的表

### 4.2 products 多字段检索 + 相关度排序

- **现状**：只按 `name` 匹配；`subtitle` 已建索引但刻意未参与查询
- **原因**：subtitle 是 IK 词级检索，与 `name` 的 ngram 子串语义不同，
  混进 `should` 会让结果噪声很大（实测一次错配 5 条）
- **前置**：§2.2（`listByIDs` 顺序保持），否则 `_score` 排序无效
- **附带**：可考虑叠加品牌名（`brand_id` → 品牌名）等字段

### 4.3 双写改异步（高频写入场景）

- **现状**：`writeRefresh: "wait_for"`，单次写最多多等一个 refresh_interval
- **影响**：products/skus 写入量远大于 brands 时可能拖慢写接口
- **修法**：改 `writeRefresh: ""`（最低写入延迟，代价是约 1s 内检索不同步），
  或引入异步/消息队列
- **注意**：写入与删除必须用**同一个**刷新策略，否则会出现「刚删除的记录仍能搜到」

### 4.4 ES 高可用

- **现状**：单节点 8G 堆，`_cluster/health` 为 `yellow`（副本未分配）
- **实测天花板**：单节点约 6.6k RPS（两次压测 6619.16 / 6618.93）
- **风险**：单节点故障时全部筛选查询回落 MySQL，MySQL 将独自承受全部并发
- **方向**：生产环境至少双节点 + 副本；容量按 6.6k RPS 预留余量

### 4.5 WebSocket 同账号多端互踢，前端反复断线重连

| 项 | 内容 |
|----|------|
| **现象** | 两台电脑用**同一个 admin 账号**同时连 WebSocket，两端反复断线重连、永不收敛 |
| **根因** | `internal/ws/hub.go:51-62` `handleRegister`：新连接注册时遍历该 userID 下所有旧连接，执行 `old.closed = true; close(old.Send); old.Conn.Close()`，随后把 `h.clients[userID]` 整体替换成「只含新连接」的 map —— 即强制单用户单连接 |
| **为什么是「反复」** | 踢人是**静默硬断**：`WritePump`（`internal/ws/client.go:74-79`）在 `Send` 被关闭时确实会补发 close 帧，但 payload 为空（**无 code、无 reason**），且 hub 紧接着 `Conn.Close()` 与之存在竞争，客户端有时只看到连接断开。两种情况下前端都无法区分「被顶替」与「网络抖动」，于是照常自动重连；重连又去踢对方 → 互踢死循环。**这才是"反复"的机制，不是连接不稳定** |
| **影响面** | ① 同一账号多端（两台电脑 / 浏览器 + 小程序）必然互踢；② 被踢的旧连接不经过 `handleUnregister`，其 `UpdateLastSeq`、离线事件、在线统计都不会为它更新；③ 会话状态**按用户**存（`internal/ws/session.go` 的 `ws:session:<userID>`，只有一个 `LastSeq`/`ReconnectCount`），而内存里 `Client.LastSeq`（`client.go:23`）是按连接的 —— 真要多设备，恢复位点必须改成按连接/设备维度，否则两端互相覆盖 `last_seq`、重连补发错乱 |
| **临时绕过** | 两个电脑用不同账号登录（当前可用） |
| **验证方式** | 同一账号两个客户端同时连：① 是否仍互踢；② 被顶替端是否收到明确提示且**不再自动重连**；③ `GET /api/v1/ws/stats` 的 `connections` 是否为 2；④ 多设备各自断开重连后 `last_seq` 补发是否正确 |
| **部署依赖** | 服务器上无 Go 源码，改完必须走「本地交叉编译 + 上传」链路（见 [`deploy/INSTALL.md`](../deploy/INSTALL.md)）；前端两端（`gf-eshop-fe`、`gf-eshop-miniprogram`）需同步处理 |

**方案（建议先做 B 止血，再评估 A 作为终态）**

- **A. 支持多设备并发（终态推荐）**：删掉 `handleRegister` 里 53-59 行的踢人分支，
  注册时改为 `h.clients[userID][client] = true`，而不是把 map 重置成只含新连接的。
  好消息是**数据结构本身已支持多连接** —— `clients` 就是 `map[int64]map[*Client]bool` 集合，
  `SendToUser` / `snapshotUserClients` / `GetOnlineCount` 都按集合遍历，广播侧不用改。
  **真正的工作量在会话状态**：`ws:session:<userID>` 目前只有一个 `LastSeq`，
  需改成按连接/设备存（例如 `ws:session:<userID>:<connID>`，或复用 JWT 里的 `TokenId`）。
  这一点如果漏了，会出现「A 端把 B 端的恢复位点覆盖掉」这类难查的补发错乱。
- **B. 保持单连接 + 被顶替通知（成本低、立刻止血）**：踢之前先给旧连接发一条业务消息
  （如新消息类型 `type: "kicked"`，payload 带 reason），并带自定义 code（如 4001）的 close 帧
  **优雅关闭**（先发完再关，不要 `close(Send)` 与 `Conn.Close()` 同时做）；
  前端收到该 code **不自动重连**，改为提示「账号已在其他设备登录」。
  这一步就能消掉互踢死循环，且不动会话模型。现有消息类型只有
  `stats/user/welcome/sync_required/notification/pong`（`internal/ws/message.go`），需新增一种。
- **C. 临时绕过**：不同账号登录（不入代码）。

> 决策提示：如果预期「多个运营共用同一 admin 账号同时在线」是常态，A 应升到 §二 高优先级；
> 如果只是本机双开测试图的方便，B 足够且代价最小。

---

## 五、低成本优化（顺手可做）

| 事项 | 说明 |
|------|------|
| `HasFilter()` 抽取（brands） | 把「无筛选」判定从手写枚举字段抽到一处，降低新增筛选字段时漏改的风险（`?status=0` 就是这类 bug）。**只能集中风险，不能自动防住**，但比散在业务逻辑里强 |
| 空表退化 | `getBrandPage` 把「缓存未建」与「真的没数据」都归成 `total <= 0`，品牌表为空时每请求会走 2 次 Redis + 2 次 DB。brands 不会为空，加行注释即可 |
| ngram 索引体积 | products 名 `varchar(200)`，`max_gram` 现为 32。若索引膨胀明显，可只对 `name` 开 ngram |
| 压测补充 | brands 已有 P99 < 50ms 的压测基线（见 `docs/cache-optimization-guide.md`）；products 走 ES 后尚未压测 |

---

## 六、参考

| 文档 | 内容 |
|------|------|
| [`docs/elasticsearch-search-guide.md`](elasticsearch-search-guide.md) | ES 完整方案、products 接入实录、行为变化与坑位（§10 已知问题、§11.9 行为变化必读） |
| [`docs/cache-optimization-guide.md`](cache-optimization-guide.md) | 缓存优化历程与 P99 目标 |
| [`docs/curd-workflow.md`](curd-workflow.md) | 新增 CRUD 模块流程 |
| [`CLAUDE.md`](../CLAUDE.md) | 项目开发指引（架构、约定、搜索层规则） |
