# 订单分表设计（时间分片 + order_no 路由）

本文是 gf-eshop 订单域分表的**决策文档**：对比候选方案、给出推荐、拆解迁移步骤与回滚路径。
按 **10 万单/日、保留 3 年** 的量级设计。配套阅读：
[`roadmap.md`](./roadmap.md)（待办与优先级）、[`elasticsearch-search-guide.md`](./elasticsearch-search-guide.md)（订单检索索引将复用其底座）。

**状态：待决策。** 本文只做设计，未落任何代码改动。

---

## 一、决策摘要

| 项 | 结论 |
|----|------|
| 分片键 | `created_at` 的**月份** |
| 分片方式 | **应用层分表**（`tx_orders_YYYYMM` / `tx_sub_orders_YYYYMM` / `tx_order_items_YYYYMM` / `tx_order_logs_YYYYMM`），四张表同键同片 |
| 路由来源 | **`order_no` 中间 14 位时间戳**（本库单号天然内嵌时间，零映射表） |
| 分库？ | **不分库**，仅分表。同实例内建单事务仍保持原子性 |
| 落地节奏 | **分两步**：先做「可切换」抽象（配置成单表模式，行为零变化），量到了再开分片 |
| 前置必修 | ① `order_no` 生成器（当前日均碰撞期望 **5.8 次/天**）② 主键换全局唯一 ID |
| 必付代价 | 列表 `COUNT(*)`+offset 分页作废；dashboard 聚合必须改汇总表；按「业务时间+状态」的运营检索必须走 ES |

> 最关键的一条：**这个库的 `order_no` 已经内嵌了创建时间**，我拿实库验证过 2000/2000 条满足
> `order_no == 'ORD' + DATE_FORMAT(created_at,'%Y%m%d%H%i%s') + 后4位`。
> 所以「按时间分片」在这里能拿到零映射表的直接路由 —— 这是 `user_id` 哈希分片拿不到的优势，
> 也是选它的**唯一决定性理由**。

---

## 二、现状核查（先说实话）

### 2.1 数据量：当前离「必须分表」还有 3~4 个数量级

实库（`eshop_db`）实测：

| 表 | 行数（`COUNT(*)`） | 数据 | 索引 | 字节/行（估） |
|----|------|------|------|------------|
| `tx_orders` | 2,000 | 0.4 MB | 0.5 MB | 490 |
| `tx_sub_orders` | 2,000 | 0.4 MB | 0.5 MB | 487 |
| `tx_order_items` | 5,025 | 1.5 MB | 0.6 MB | 450 |
| `tx_order_logs` | **0** | — | — | — |
| `tx_payments` | 1,634 | 0.3 MB | 0.4 MB | 489 |

> 行数取自 `COUNT(*)`；体积取自 `information_schema.TABLES`（其 `TABLE_ROWS` 是估算值，
> 小表上页级固定开销会被摊到很少的行上，因此「字节/行」偏大 —— **对小表仅供参考**）。
> 第六节的容量推算改用 DDL 逻辑行宽 + 页填充估算，不直接外推这里的数字。

订单分布：2026-08（1595）、2026-09（405），覆盖 50 个用户，约 66 单/天。
比例：`sub_orders : orders = 1.0`，`items : orders = 2.51`。

**结论：分表不是当前性能问题的解药。** 本文的价值是「提前把路修好」，而不是「现在必须切」。
第十节会给出明确的**触发阈值**，到线才动手。

### 2.2 订单域不是一张表 —— 本节涉及 9 张表的联动

```
tx_orders ──1:1── tx_sub_orders ──1:N── tx_order_items
    │                                        ▲
    ├── tx_order_logs                        │
    ├── tx_payments ── tx_refunds            │
    ├── tx_deliveries ── tx_delivery_items ──┤  (delivery_items.order_item_id)
    └── tx_after_sales ──────────────────────┘  (order_id + order_item_id)
```

关键结构事实（决定方案可行性）：

- 子表**大多冗余了 `order_no`**：`tx_sub_orders.parent_order_no`、`tx_order_items.order_no`、
  `tx_order_logs.order_no`、`tx_payments.order_no`、`tx_refunds.order_no`、`tx_deliveries.order_no`
  → 都能靠单号直接路由；
- **`tx_after_sales` 只有 `order_id` / `order_item_id`，没有 `order_no`** → 它只能靠订单 ID 路由（见 §5.3）；
- `tx_deliveries` 挂在**订单**维度（`order_id` + `merchant_id`），明细在 `tx_delivery_items.order_item_id`；
- `tx_payments.order_id` / `tx_deliveries.order_id` / `tx_refunds.order_id` 引用 `tx_orders.id` → **ID 必须全局唯一**；
- 库里**没有物理外键约束**（只有逻辑引用）→ 改主键定义不会被外键挡住，这是运气好。

### 2.3 约束分片键的 5 条访问路径

| 路径 | 位置 | 对分片的要求 |
|------|------|-------------|
| 建单事务：4 张表 + 扣库存 | `internal/logic/orders/orders.go:58` | 4 张表**必须同键同片**，否则丢原子性 |
| 列表：`user_id`+`status`+`payment_status`+`order_no`，`COUNT(*)`+offset | `orders.go:244` | 无时间窗时跨片；`total` 要归并 |
| 详情 / 改状态：**按 `order_no`**（不带 user_id） | `orders.go:282`、`:341` | 必须能**直接定位单分片** |
| 支付回调：按 `order_no` 回写 orders、按 `parent_order_no` 回写 sub_orders | `logic/payments/payments.go:143`、`:153` | 同上；且不可跨片 |
| dashboard：全表 `COUNT`、`SUM(pay_amount)`、`GROUP BY status`、`GROUP BY DATE(created_at)`、items 聚合 | `logic/dashboard/dashboard.go:193`、`:224`、`:264`、`:396` | 必须改成汇总表或 fan-out，否则数字错 |

### 2.4 框架约束

- **GoFrame v2.6.1 没有分表原语**。我在 `$GOMODCACHE/github.com/gogf/gf/v2@v2.6.1/database/` 下 grep `sharding` 零命中，路由层必须自己写。
- 订单域的表/DAO 引用共 **83 处**，但其中只有 **35 处**在业务代码里（`internal/logic` + `internal/controller` + `internal/cmd`），
  其余全是生成物（`internal/dao`、`internal/model/entity`、`internal/model/do`）。
  **真正要动的只有 3 个文件**：`logic/orders/orders.go`、`logic/payments/payments.go`、`logic/dashboard/dashboard.go`。
  改造面比看上去小得多 —— 这是推荐「先做 Phase 0 抽象」的底气。
- `gf gen dao` 是按表名生成的：**绝不能为每个分片生成 DAO**，否则 3 年后是 144 个 DAO。分片表结构一致，复用 `entity.Orders` 即可，DAO 层退化为 repo 内部的表名拼接。

### 2.5 顺带发现的缺陷：单号生成器上量就会撞

```go
// internal/logic/orders/orders.go:34
func generateOrderNo() string {
    return fmt.Sprintf("ORD%s%04d", gtime.Now().Format("YmdHis"), grand.Intn(10000))
}
```

后 4 位是 `rand(10000)`，同一秒内只有 1 万个槽位。按 10 万单/日估算：

| 场景 | 同秒内至少一次碰撞的概率 |
|------|------------------------|
| 均值 1.157 单/秒（10 万/日） | 0.0009% / 秒，**日均期望 5.8 次碰撞** |
| 峰值 200 单/秒 | **86.3%** |

撞上 `uk_order_no` 的表现是**用户下单直接失败**（不是静默降级）。所以这是分表工作里
**必须先修**的前置项，而且修法与路由天然统一（见 §5.3）。

> ✅ **Phase 1 已修复**：改为 Redis 每秒序列 + 降级 + 唯一键冲突重试。
> 实测 100 并发建单零重号；可控撞号（预置 Redis 序列 + 抢先插入同一单号）下
> 自动换号重试并在同一秒内成功。

---

## 三、分表到底解决什么（把话说准）

容易犯的错是把分表当「点查变快」的手段。**点查在 1.1 亿行时并不慢**：

| 指标 | 1.1 亿行不分片 |
|------|---------------|
| PK / 唯一索引高度 | **4 层**（按 16 KB 页、扇出约 900 估） |
| `WHERE order_no = ?` | 4 次页访问，root 常驻 buffer pool → 亚毫秒级 |

真正的痛点在**运维与生命周期**，这才是分表的收益来源：

| 痛点 | 1.1 亿行 / 266 GiB 单表下 | 月分片（304 万行 / 1.13 GiB）下 |
|------|--------------------------|-------------------------------|
| 加列 / 加索引（在线 DDL） | 数小时，期间风险与主从延迟不可控 | 秒级~分钟级，可**只对当月表**执行 |
| 物理备份 / 克隆 / 恢复 | 单文件 266 GiB，恢复窗口以小时计 | 按片备份，可只备热数据 |
| 历史数据淘汰 | `DELETE` 上亿行 → 海量 undo、碎片、主从延迟 | `DROP TABLE` **秒级**，零碎片 |
| 冷热分离 | buffer pool 被冷数据挤占 | 热表小，命中率高 |
| 单实例故障爆炸半径 | 全量 | 可按片迁到不同实例（未来分库的跳板） |

> 判断标准：**如果分表不能换来「可按片运维」，那它只是把复杂度搬进了代码，不值。**
> 上面这一列才是它值的原因。

---

## 四、候选方案对比

| 方案 | 路由 | 优点 | 致命缺点 | 结论 |
|------|------|------|---------|------|
| **A. 不分表**，只做归档 + 索引优化 | 无 | 零改造、零风险 | 单表继续涨，DDL/备份/删除三座山不动 | 现状维持，不是终态 |
| **B. MySQL 原生 RANGE 分区**（`PARTITION BY RANGE (created_at)`） | 靠 `WHERE` 里的 `created_at` 做分区裁剪 | **应用零改动**（表名不变）；`DROP PARTITION` 秒级 | ① 本项目订单查询**全部按 `order_no`**，WHERE 里没有 `created_at` → **裁剪完全失效，每次查 36 个分区**；② 唯一键必须含分区列 → `uk_order_no` 变 `(order_no, created_at)`、PK 变 `(id, created_at)`；③ 仍是**一个大 `.ibd` 文件**，备份/恢复粒度不变 | ❌ 裁剪失效是决定性的 |
| **C. 应用层按月分表 + `order_no` 路由** | 解析单号 → 单分片 | 点查零代价；可按片运维；`order_no` 已内嵌时间，**零映射表**；天然支持冷热分层与未来分库 | 需自建路由层（业务侧 3 个文件 / 35 处引用）；跨片查询要 fan-out | ✅ **推荐** |
| **D. 按 `user_id` 哈希分表** | `hash(user_id) % N` | 「我的订单」单分片 | 详情/支付回调**只带 `order_no` 不带 `user_id`** → 必须维护 `order_no → user_id` 映射表，多一次写与一次查；dashboard 必须 fan-out；时间归档能力丧失 | ❌ 路由成本更高 |
| **E. 按 `order_no` 尾号取模** | 尾号 | 分布最均匀 | 与 D（映射表问题）相同，且**丢掉时间归档能力**——而订单恰恰是最典型的时间序数据 | ❌ |

**选 C 的决定性论据**：B 的分区裁剪依赖查询里带 `created_at`，而本库订单查询的入口是
`order_no`（`GET /orders/{order_no}`、`PUT /orders/{order_no}/status`、支付回调）；
C 则把「单号里的时间」直接变成路由，**查询模式一行不用改**，同时拿到可按片运维的能力。

---

## 五、推荐方案详设

### 5.1 分片键与粒度

分片键 = `created_at` 的月份，表名后缀 `YYYYMM`。

**为什么四张表必须同键同片**：建单是一个事务，写 orders + sub_orders + items + logs。
同键同片 → 四条 INSERT 落在**同一个 MySQL 实例**，事务原子性天然保持。
一旦将来跨实例分库，这个事务立刻失效（届时需要 outbox / TCC）——
**所以在设计上要明确：订单域四张表永远同库，分库时间表要单独评估。**

粒度选月还是季？两者在 10 万单/日 下都在舒适区：

| 粒度 | 分片数（3 年） | `tx_orders` 单片 | `tx_order_items` 单片 | 表总数（4 表） |
|------|--------------|-----------------|---------------------|--------------|
| **月** | 36 | 304 万行 / **1.13 GiB** | 764 万行 / 2.85 GiB | 144 |
| 季 | 12 | 912 万行 / 3.4 GiB | 2,290 万行 / 8.6 GiB | 48 |

**推荐按月**，理由：① 归档/DROP 粒度细；② DDL 窗口小，可只对热表执行；
③ 热表更小，buffer pool 命中率更高。跨片 fan-out 的代价不靠「减少分片数」来省，
而是靠 §5.6 的汇总表绕开。

> 若运维人手紧张、想要更少的表，**按季同样成立**（12 片 vs 36 片），
> 因为 10 万单/日 的量级下两种粒度都不触阈值。这是个可自由选择的取舍，不是对错。

**单片容量目标**（用于未来重新评估粒度）：

| 指标 | 目标 | 月分片 @10 万/日 | 余量 |
|------|------|-----------------|------|
| 单表行数 | < 2,000 万 | 304 万（orders） / 764 万（items） | 6.5x / 2.6x |
| 单表体积 | < 5 GiB | 1.13 GiB / 2.85 GiB | 4.4x / 1.8x |

即：**月分片可支撑到约 40 万单/日** 才触及 5 GiB 阈值。超过再考虑按天或按周。

### 5.2 路由规则

```go
// 路由层接口形状（文档示意，非最终实现）
type ShardKey struct{ Year, Month int }   // 分片标识；单表模式下恒为零值

// 三级来源，优先级从高到低
// 1) order_no：取 [3:17] 解析成 yyyymmddHHMMSS → 月份   ← 绝大多数路径走这条
// 2) created_at：写入路径直接用它
// 3) 无法解析：进入兜底（广播热分片 + 计数告警），绝不静默返回空
func ShardOfOrderNo(orderNo string) (ShardKey, bool)
func ShardOfCreatedAt(t time.Time) ShardKey

// 表名拼接：tx_orders + 202608
func Table(base string, k ShardKey) string
```

要点：

- **解析失败必须可观测**：兜底广播要打点（`order_shard_fallback_total` 计数器）。
  这个指标长期不为 0，说明有异常单号在持续消耗全片扫描，是必须清掉的债。
- 维护「**活跃分片清单**」配置（如近 13 个月），兜底只在活跃清单内广播 —— 不能真的扫 36 个片。
- 老数据兼容：现有 2000 条单号是同一格式（`ORD`+YmdHis+4 位随机），
  **路由只看前 14 位时间 → 旧格式照样能路由**，无需特殊分支。这是本方案的一个额外便宜。

### 5.3 单号与主键（前置必修）

**问题 1：单号同秒碰撞**（§2.5）。修法：把后 4 位随机数换成**每秒单调递增序列**：

```
ORD + YYYYMMDDHHMMSS(14) + 序号(6)      → 23 字符，varchar(32) 放得下
```

序列从 Redis `INCRBY order:seq:<yyyyMMddHHMMSS>` 取（TTL 10s）；
**Redis 不可用时降级为进程内计数器并告警一次**，此时多实例可能撞号，
由唯一键冲突触发**有界重试**（最多 3 次，每次重取时间基准与序列）兜底 ——
订单写入不因 Redis 故障而中断，与本项目「缓存是加速器不是硬依赖」的取向一致。

序列**一次取够**（订单 1 个 + 每条明细 1 个），既省 Redis 往返，
也让订单主键与明细主键共用同一个序列源，不会各推一份时间。
单号与主键共用同一个 `now`，因此**单号内嵌时间与主键反解时间秒级完全相同**（已实测 100/100）。

> 代价（必须知情）：单号里带时间+序号 → **泄漏单量信息**。真实电商常因此改用不可解析的
> 单号 + 映射表。这里选择「可路由优先」，接受信息暴露；若业务不接受，则要额外维护
> `order_no → 分片` 映射表，路由成本上升一档。**这是需要在决策时明确表态的一点。**
>
> 回滚开关：`orderNo.mode: legacy` 可让单号退回旧的「4 位随机后缀」格式
> （主键仍走序列，不受影响）。新旧两种格式的时间段位置相同，路由层通吃。

**问题 2：分片后 `AUTO_INCREMENT` 各片独立 → ID 冲突**，而 `tx_payments` / `tx_refunds` /
`tx_deliveries` / `tx_after_sales` 都用 `order_id` 引用，`tx_delivery_items` / `tx_after_sales`
还用 `order_item_id` 引用 `tx_order_items.id`。

修法：**给 `tx_orders.id` 与 `tx_order_items.id` 都换全局唯一主键**
（两张表都要，只换订单主键会漏掉订单项的跨域引用 —— 这是 Phase 0 文档的一处遗漏）。

> ⚠️ **不能直接用 41 位毫秒的标准雪花**：`41 位毫秒 << 22` 会得到约 9.2e18 量级，
> 远超 JS 的 `Number.MAX_SAFE_INTEGER`（2^53-1 ≈ 9.0e15）。前端 `gf-eshop-fe` /
> `gf-eshop-miniprogram` 把 `id` 当 JSON number 解析，超限会**静默丢精度**，
> 回传的 id 变成另一个数 —— 这类 bug 极难查。接口把 id 改成字符串也会破坏既有前端。
>
> 因此采用**压进 53 位**的布局，并按「分钟」编码时间（分片只需要月份，分钟精度足够）：

```
主键 51 位： 分钟位(25) | 秒位(6) | 序列位(20)
              ↑ 自 2024-01-01 UTC 起算，25 位够用约 63 年
最大约 2.25e15 < 2^53-1  →  前端解析 JSON number 不丢精度
```

**好处是三份**：ID 全局唯一解决跨域引用；高位含时间，**by `order_id` 也能反解月份**，
正好覆盖没有 `order_no` 冗余的 `tx_after_sales`；并且天然按时间递增，`ORDER BY id DESC`
仍然等价于「最新在前」。

`decodeOrderID()` 已随 Phase 1 落地并单测覆盖（含月末/年末边界）。

> ✅ **已完成（Phase 2）**：`tx_orders.id` / `tx_order_items.id` 的 `AUTO_INCREMENT` 已去掉。
> DDL 归 **schema 源仓库 `std-eshop-db`** 管理（本仓库不放 DDL，避免两个来源）：
> 基线 `sql/tx_p0.sql`、`sql/tx_p1.sql`、`sql/tx_p5.sql`；
> 存量库前向升级 `sql/migrations/V001__tx_order_sharding.sql`。
> 去掉之后，任何漏改的、不带 id 的 INSERT 会直接报 `Error 1364` —— 实测确认，
> **宁可写入失败，也不要产生一个不含时间位、无法按 id 路由的脏主键**。

**必须定一条铁律 —— 单一时间基准**：

> 建单时**只取一次 `now`**，用它同时派生 `order_no`、雪花 ID 的时间部分与**分片键**。
> 三者因此必然落在同一月份。

否则会出现 23:59:59.999 的边界订单：单号指向 8 月、ID 指向 9 月 → 按单号读和按 ID 读
命中**不同分片**，这类 bug 极难复现。

> ⚠️ **`created_at` 不在「单一时间基准」的覆盖范围内 —— 它被框架覆盖。**
> GoFrame 的 `Insert` 会**无条件**用自己取的 `gtime.Now()` 覆盖 `created_at` / `updated_at`
> （`gdb_model_insert.go:311-321`：`v[fieldNameCreate] = now`），data map 里传的值会被丢掉。
> Phase 0 实测：同一次建单的四张表，`created_at` 分别是 `.071 / .072 / .072 / .073` ——
> **月份级一致（四表 `%Y%m` 实测相同），但毫秒级不同。**
> 也就是说「四表 created_at 完全相同」这个强不变量**不成立**，只能在月份粒度上成立。

**Phase 2 必须处理这个毫秒窗口**：若分片取自代码里的 `now`，而行内 `created_at` 由框架
在几毫秒后才生成，那么**跨月瞬间（每月 1 号 00:00:00.000~005 左右）建的单会被写进上个月的
分片**，与自身 `created_at` 不符 —— 按 `created_at` 扫描的列表/看板会漏掉它。
三种处理方式，**推荐第一种**：

| 方案 | 做法 | 评价 |
|------|------|------|
| **A. 让显式值生效（推荐）** | 建单的四条 INSERT 用 `.Unscoped()`，使显式传入的 `created_at`/`updated_at` 不被覆盖；此时四表 `created_at` 与分片键同源，**可断言完全相等** | 一处改动，不变量变强且可测；`.Unscoped()` 只影响插入时的自动填充，对只做 INSERT 的模型无副作用 |
| B. 以库中值为准 | 先插主订单，回读其 `created_at`，再用该值推导分片插子表 | 多一次回读；子表自身 `created_at` 仍可能再偏几毫秒 |
| C. 接受窗口 | 记录为已知风险，靠「按单号路由」兜住点查 | 跨月那个窗口内的单在按时间扫描时会漏；不推荐 |

**老数据（id 1..2000）** 不满足雪花结构，且迁移时**不应重写 ID**（会级联破坏
payments/refunds/deliveries/after_sales 的引用）。处理方式：

- 路由时判断 `id < 阈值`（如 10^10）→ 走兜底广播或查一次性的 `order_id → shard` 小映射表；
- 老数据只占最早两个月，活跃度极低，建一张几千行的映射表成本可忽略。
  迁移脚本顺手生成即可，无需改动线上写路径。

### 5.4 事务与一致性

| 事项 | 结论 |
|------|------|
| 建单事务（4 表 + 扣库存 `sp_inventories`） | **保持原子**：分表不分库，同实例内跨表事务不受影响 |
| 订单状态更新（orders + sub_orders + logs） | 同上，单分片内 |
| 支付回调（payments + orders + sub_orders） | `tx_payments` **不分片**，orders/sub_orders 按 `order_no` 路由到同片 → 事务仍原子 |
| 跨分片事务 | **不需要**（所有子表同键同片）；若将来引入跨月写入需求，需重新评估 |
| 分库 | **不在本次范围**。分库后建单事务失效，需要 outbox/TCC —— 单独评估 |

> 这是选「分表不分库」的核心收益：**用一个不跨实例的约束，换掉整套分布式事务复杂度。**
> 代价是单实例仍是容量上限，未来必须再走一步分库。

> **Phase 0 已实测**：支付回调里回写 orders + sub_orders 的 `MarkPaidByOrderNo`
> **确实加入了调用方事务** —— 在事务中调用它后故意回滚，订单仍是 `unpaid`，
> 证明它没有另开连接或提前提交。这条必须一直成立，否则分表后跨表写入会悄悄丢掉原子性。

### 5.5 查询改造清单

| 接口 / 查询 | 改造 | 破坏性 |
|-------------|------|--------|
| `GET /orders/{order_no}` 详情 | 解析单号 → 单分片；子表同片 | 无 |
| `PUT /orders/{order_no}/status` | 同上 | 无 |
| 支付回调按 `order_no` 回写 | 同上 | 无 |
| `POST /orders` 建单 | 定单一时间基准 → 落当前月分片 | 无 |
| `GET /orders?order_no=x` | **退化为分片内条件查询**（其实等价于点查） | 无 |
| `GET /orders?page&page_size` | 保留（offset 分页，分表后**只有单表模式精确**） | 无 |
| `GET /orders?size&cursor` | ✅ **已实现**（Phase 2）：keyset 分页，`cursor=base64(末位订单ID)`、排序固定 `id DESC`、返回 `next_cursor`、`total=-1`（不统计）。分片后各片各取一页再归并即可 | 新增（与 products 的游标约定一致） |
| `GET /orders?user_id=x` 无时间窗 | 默认**强制近 N 个月**（如 12 个月），否则 fan-out 全部活跃片 | ⚠️ 行为变更 |
| dashboard 全部聚合 | 改走**日汇总表**（§5.6） | 需新增表与写入钩子 |
| 「今天待发货」这类**按业务时间+状态**的运营检索 | `created_at` 分表**无法覆盖** → **必须走 ES**（§5.7） | 新增能力 |
| 库存 `sp_inventories`、商品 `sp_products` | 不分片 | 无 |

⚠️ **API 破坏性变更需要跨仓库协调**：`OrdersListRes.Total` 与 `page` 语义变化会影响
`gf-eshop-fe`、`gf-eshop-miniprogram`，以及现有回归用例
（`tests/test_tx_api.py:152`、`:163`、`:173` 都在断言 `total`）。
建议新增 `cursor` + `size`，`total` 仅在带时间窗时返回精确值，否则返回 `-1` 表示「不精确」。

> 这印证了 CLAUDE.md 里那句话：**ES 是加速器，不是唯一真相源**。
> 分表后它同时变成「跨分片检索层」—— 详情/写走分片 MySQL，列表/运营检索走 ES。

### 5.6 汇总表（dashboard 的解）

实时 `SELECT COUNT/SUM/GROUP BY` 在 36 片上跑是不可接受的（串行估算 2~7 秒）。
用一张**日粒度汇总表**承接所有看板聚合：

```sql
CREATE TABLE tx_order_daily_stats (
  stat_date      date     NOT NULL COMMENT '统计日期（按 Asia/Shanghai）',
  order_cnt      bigint   NOT NULL DEFAULT 0 COMMENT '下单数',
  paid_cnt       bigint   NOT NULL DEFAULT 0 COMMENT '支付笔数',
  refund_cnt     bigint   NOT NULL DEFAULT 0 COMMENT '退款笔数',
  cancelled_cnt  bigint   NOT NULL DEFAULT 0 COMMENT '取消数',
  gmv            bigint   NOT NULL DEFAULT 0 COMMENT '下单金额（分）',
  paid_amount    bigint   NOT NULL DEFAULT 0 COMMENT '实收金额（分）',
  refund_amount  bigint   NOT NULL DEFAULT 0 COMMENT '退款金额（分）',
  updated_at     datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (stat_date)
) ENGINE=InnoDB COMMENT='订单日汇总（看板数据源）';
```

- **写入时机**：建单、支付成功、退款、取消时 `INSERT ... ON DUPLICATE KEY UPDATE 累加`。
  放在业务事务内会放大锁竞争，建议**异步**（复用现有通知/WS 的异步通道，或落 outbox 表）。
- **历史回填**：迁移脚本按分片逐月 `SELECT DATE(created_at), COUNT(*), SUM(...) GROUP BY 1` 灌入，
  每片一次聚合，秒级完成。
- 有了它，看板的「订单总数/营收/7 日趋势/状态分布」全部变成**单表单行/单表索引查询**，
  与分片数无关 —— 也就不需要 fan-out 归并逻辑（比 fan-out 更简单也更稳）。
- **进度**：表已建、历史已回填（Phase 2，实测 32 天，与主表直查逐项一致）；
  看板**读路径的切换**放在 Phase 4，与读路径灰度一起做。

### 5.7 订单检索索引（ES）

分表后有一类查询**结构性地**没法单分片完成：
「**按状态 + 业务时间**」—— 例如「今天要发货的订单」= `status='paid' AND paid_at >= 今天`。
8 月下的单可能 9 月才发货，它躺在 8 月分片里，按 `created_at` 分表永远定位不到。

做法：**完全照搬 `internal/logic/products/es.go` 的模式**，新增订单检索索引：

- 只存**检索投影**（`order_no`、`user_id`、`status`、`payment_status`、`created_at`、`paid_at`、
  `shipped_at`、`merchant_id`、`pay_amount`、`source`），不存收货地址等全量字段；
- 别名 + schema 版本前缀 + 零停机重建 + 双写 + 熔断降级，全部复用 `internal/search` 底座；
- 列表响应仍由 MySQL 侧按 ID 组装（与 products 一致）；
- 在 `internal/cmd/cmd.go` 的 `reindexTargets` 登记 `orders`，即可获得 `main reindex --entity=orders`。

> 注意 products 踩过的坑：**改 mapping / analyzer / 字段语义必须递增 schema 常量**，
> 否则启动对账发现不了（文档见 `elasticsearch-search-guide.md` §10）。

### 5.8 归档与生命周期

| 温度 | 范围 | 存储 | 访问方式 |
|------|------|------|---------|
| 热 | 近 3 个月 | 按月分片，主库 | 直接路由 |
| 温 | 3 个月 ~ 1 年 | 按月分片，主库（可选迁到只读实例） | 强制时间窗 |
| 冷 | 1 ~ 3 年 | 归档表 / 对象存储（Parquet/CSV） | 离线导出，不进在线链路 |
| 到期 | > 3 年 | `DROP TABLE` 秒级删除（先满足审计/对账要求） | — |

- 淘汰**必须用 `DROP TABLE`，不要 `DELETE`**：不分表时淘汰等于 `DELETE` 上亿行 → 海量 undo、
  页碎片、主从延迟同时爆；分表后是 `DROP TABLE`，**秒级且零碎片**。这正是分表最主要的收益之一。
- `tx_order_logs` 是增长最快的表（3 年 4.38 亿行 / 81.6 GiB，超过 orders）。
  建议**单独缩短保留期**（如 1 年），它是审计日志，多数场景不需要留满 3 年。

### 5.9 运维

| 事项 | 做法 |
|------|------|
| 建表 | 定时任务/启动时确保「当前月 + 下月」分片存在（幂等 `CREATE TABLE IF NOT EXISTS`） |
| DDL 一致性 | 用模板建表 + 巡检脚本比对 `SHOW CREATE TABLE`（144 张表结构漂移是主要风险） |
| 加列 | 遍历活跃分片执行，用 `gh-ost`/`pt-online-schema-change`；注意 `LIKE` 建表不会自动跟进 |
| 时区 | 分片月份推导必须与 `created_at` **同一个时区**。当前配置 `loc=Local`，建议统一固定 `Asia/Shanghai`，否则跨时区部署会出现「同一订单算出不同月份」 |
| 备份 | 可按片备份：热片每日全量，冷片一次性归档 |
| 监控 | 单片行数/体积、`order_shard_fallback_total`、跨片查询 P99、汇总表与分片对账差异 |
| 表数量 | 3 年后 144 张，MySQL 8/9 无压力；但要控制 `information_schema` 类监控查询的频率 |

---

## 六、容量规划（10 万单/日 × 3 年）

假设：`sub_orders : orders = 1.0`（当前实测）、`items : orders = 2.51`、
`logs : orders = 4`（created/paid/shipped/delivered/completed，取消退款另加）；
行宽取 orders/sub 400 B、items 400 B、logs 200 B（含索引，按 DDL 逻辑行宽 ~270 B + 页填充估）。
单位 GiB。

| 表 | 3 年行数 | 3 年体积 | 单片行数（月） | 单片体积（月） |
|----|---------|---------|--------------|--------------|
| `tx_orders` | 1.10 亿 | 40.8 | 304 万 | **1.13** |
| `tx_sub_orders` | 1.10 亿 | 40.8 | 304 万 | 1.13 |
| `tx_order_items` | 2.75 亿 | 102.5 | 764 万 | **2.85** |
| `tx_order_logs` | 4.38 亿 | 81.6 | 1,216 万 | 2.26 |
| **合计** | — | **265.7** | — | **7.4 / 月** |

对照：**不分表**时 `tx_order_items` 是一张 **1.03 亿行 / 102.5 GiB** 的单表。

敏感度：若未来引入多商家拆单使 `sub:order` 升到 1.5、`items:order` 升到 3，
月分片 `tx_order_items` 升到约 3.4 GiB —— 仍在 5 GiB 阈值内。

---

## 七、迁移步骤

设计原则：**每一步都单独可验证、可回滚；旧表在最后一步之前绝不删除。**

### Phase 0：抽象 repo 层（不改行为）

- 动作：把订单域全部 SQL 收敛进 `internal/logic/orders/repo_*`，引入路由层接口，
  配置 `order_shard.mode = "single"`（单表模式，路由恒返回空分片 → 仍读写 `tx_orders`）。
  同时收回 `logic/payments`、`logic/dashboard` 里对 `dao.Orders` / `tx_orders` 的直接引用。
- 价值：**改完行为完全一致**，但「分表」从一次重写变成改一个配置项。
- 验证：`tests/test_tx_api.py` 全绿，接口响应逐字段不变。
- 回滚：纯重构，回滚即 revert。

### Phase 1：单号与主键生成器（仍单表）✅ 已完成

- 动作：`order_no` 后 4 位随机 → 6 位每秒序列；`tx_orders.id` / `tx_order_items.id`
  改为显式指定的全局唯一 ID（51 位布局，见 §5.3）。
- 落点：`internal/logic/orders/identity.go`（生成器 + 序列分配器 + 降级）、
  `identity_test.go`（契约单测）。
- **前置必修**：Phase 1 一开工就发现建单接口在 HEAD 上根本跑不通（库存查询写死
  `warehouse_id=0`、JSON 列写入非法值、退款缺 `idempotency_key`），必须先修掉才能验证，
  详见 roadmap §2.4/§2.5。
- 验证（全部实测）：
  - 100 并发建单：单号零重复、主键零重复、同秒序列互不相同；
  - 单号内嵌时间与主键反解时间**秒级精确一致 100/100**；与 `created_at` 同月 100/100；
  - 主键全部 < 2^53-1（JS 安全）；
  - 可控撞号场景下自动重试成功（同秒换号）；
  - Redis 指向不可用端口后仍能建单（降级为进程内计数器 + 告警一次）；
  - 读路径与 Phase 0 逐字段对比 14 项依然零差异；`tests/test_tx_api.py` 28/28 通过且可重复执行。
- 回滚：`orderNo.mode=legacy` 退回旧单号格式；新老格式路由兼容，回滚不产生脏数据。

### Phase 2：建分片 + 历史迁移 ✅ 已完成

- 动作（全部落地为命令，见下方「运维命令」）：
  1. `CREATE TABLE ... LIKE` 批量建好月份分片（幂等，4 张表 × 每个月）；
  2. 按 `created_at` 分月 `INSERT ... SELECT`，**按 id 区间分批**（默认每批 5,000 行）；
  3. 登记老数据 `order_id → shard` 映射（新主键可反解，不入表）；
  4. 从**该月分片**回填 `tx_order_daily_stats`（每天必然只落在一个月分片里）。
- 与计划的差异：对账的第三项从「抽样 100 单逐字段哈希」升级为
  **全字段 CRC32 校验和**（列名从 `information_schema` 动态取，自动适配后续加列）。
  抽样只能证明「没抽到的地方没坏」，全字段校验和是 O(1) 的全量证据。
- 幂等性：写数据统一用 `ON DUPLICATE KEY UPDATE`（**不用 `INSERT IGNORE`** ——
  后者会把 `CHECK` 约束冲突降级成告警并丢行，等于静默丢数据）。
- 实测（当前库 2000 单 / 5025 明细）：
  - 建分片：8 张表（2026-08、2026-09）就绪；
  - 迁移：202608 迁 1595 单 + 4027 明细，202609 迁 405 单 + 998 明细，合计 2000 / 5025，
    与主表逐一对上；登记老主键映射 2000 行；日汇总回填 32 天；
  - 对账：8 项（2 月 × 4 表）行数 / 校验和 / 金额**全部通过**；
  - **重跑 migrate 新增 0 行**，再次对账仍全通过（可重复执行）；
  - 日汇总与主表直查一致：订单数 2000=2000、GMV 3097410823=3097410823、
    实收 2127606184=2127606184。
- 回滚：分片表、映射表、汇总表都是**新增**对象，`DROP` 即可；主表未被改动。

**运维命令**（照 `main reindex` 的写法，选项传参，不能写位置参数）：

```bash
./main shard --action=create  --from=2026-08 --to=2026-09
./main shard --action=migrate --from=2026-08 --to=2026-09 [--batch=5000]
./main shard --action=verify  --from=2026-08 --to=2026-09
```

**主键 → 分片的统一入口**：`ShardOfOrderID(ctx, id)` ——
新主键（`>= 1<<40`）直接反解月份；老主键查 `tx_order_shard_map`；
两者都不匹配则**明确报错**，绝不猜一个分片。实测三种情况都符合预期。

> ⚠️ `decodeOrderID` 对 `id < 1<<40` 一律返回失败：迁移前的自增主键在数学上也能
> 解出一串位（会得到 2024-01-01），不设阈值就会把老订单路由到不存在的分片。
>
> 📌 **Phase 4 待办**：`tx_order_shard_map` 与 `tx_order_daily_stats` 目前走原始 SQL
> （`g.DB().Exec/GetValue`）。等看板读路径切到汇总表时，把这两张表加进
> `hack/config.yaml` 的 `tables` 并跑 `gf gen dao`，换成 DAO 访问。

### Phase 3：开双写 + 影子读对账

- 动作：写入同时写主表与分片；读仍走主表。
  影子读：抽样请求同时读两边，比对结果并打点差异率。
- 竞态说明：Phase 2 迁移期间对**历史月份**的零星更新（发货、退款）不会进分片 ——
  由「迁移完成后按 `updated_at > 迁移开始时刻` 对该月重放一次（`REPLACE INTO`）」补齐。
  这个窗口很小，且被限定在历史月份。
- 验证：差异率连续 N 天为 0；双写失败只记日志不阻塞主链路（与 ES 双写同策略）。
- 回滚：关双写，主表仍是完整真相源。

### Phase 4：灰度切读

- 动作：`order_shard.enabled = true` 按流量比例（如 1% → 10% → 50% → 100%）放量。
- 验证：核心接口 P99、错误率、`order_shard_fallback_total`、跨片查询耗时。
- 回滚：**单开关回落主表**，秒级生效。

### Phase 5：停双写

- 动作：确认全量切读稳定后，停止写主表。
- 回滚：主表停写期间的新数据只在新分片 —— **必须在这之前确认无需回滚**。

### Phase 6：归档旧表

- 动作：`RENAME TABLE tx_orders TO tx_orders_legacy_202609`，观察 **30 天**后再 `DROP`。
- 验证：全链路无任何引用旧表名（`grep` + 监控空连接）。
- 回滚：RENAME 回来即可（30 天内）。

---

## 八、回滚方案总表

| 阶段 | 回滚动作 | 数据风险 | RTO |
|------|---------|---------|-----|
| Phase 0 | revert 代码 | 无 | 分钟级 |
| Phase 1 | 切回旧生成器 | 无 | 分钟级 |
| Phase 2 | 删除分片表 | 无（主表未动） | 分钟级 |
| Phase 3 | 关双写开关 | 无（主表是真相源） | 秒级 |
| Phase 4 | `order_shard.enabled=false` | 无 | **秒级** |
| Phase 5 | 重启双写 | 停写期间新单只在分片，需手工回灌 | 小时级 |
| Phase 6 | `RENAME` 回旧名 | 无（30 天窗口内） | 分钟级 |

**核心保障**：Phase 6 之前**旧表始终是完整的**，所以前五个阶段随时可以「拔开关回落」。
唯一不可逆的窗口是 Phase 5→6，因此这两步之间要留足够长的观察期。

---

## 九、风险与坑

| 风险 | 触发条件 | 缓解 | 验证方式 |
|------|---------|------|---------|
| 单号撞唯一键 | 上量后同秒并发 | Phase 1 换每秒序列 | 并发压测零重复 |
| 分片错位（月末边界） | 单号/ID/created_at 时间基准不一致 | **单一时间基准**铁律 + 一致性巡检 | 巡检脚本断言三者月份一致率 100% |
| 老数据 ID 路由失效 | `id < 阈值` | 广播兜底 + 小映射表 | `fallback_total` 指标归零 |
| 跨片查询随分片数变慢 | 36 片 fan-out | 默认强制时间窗；列表走 ES；看板走汇总表 | 跨片查询 P99 监控 |
| 时区错位 | 部署机时区与数据时区不一致 | 固定 `Asia/Shanghai` | 分片归属抽样核对 |
| 144 张表结构漂移 | 手工加列漏改分片 | 模板建表 + 结构巡检 | `SHOW CREATE TABLE` diff |
| 汇总表与分片对账不一致 | 异步累加丢事件 | 每日离线重算对账 + 差异告警 | 日终对账差异为 0 |
| 雪花时钟回拨 | NTP 校时 | 回拨检测 + 拒绝/等待 | 单元测试 |
| 单号泄漏单量 | 格式可解析 | 知情接受，或改不可解析+映射表（**待决策**） | — |
| 忽略 order_logs 增长 | 未单独设保留期 | 日志表独立保留 1 年 | 体积监控 |
| 未来分库时建单事务失效 | 跨实例 | 本设计明确不分库；分库需单独设计 outbox/TCC | — |

---

## 十、待决策 / 开放问题

1. **本次是否只做 Phase 0？** 按当前数据量（2000 单 / 0.4 MB），
   推荐**只做 Phase 0 的「可切换」抽象**：把 3 个业务文件里的 35 处引用收敛进 repo + 路由留出口子，
   等触到阈值再走 Phase 1~6。这一步不改行为、可 revert，却是后面全部工作的前提。
2. **单号是否允许泄漏单量？** 可解析 = 零映射表路由（推荐）；不可解析 = 需额外映射表。
   这一条直接决定路由成本，需要业务表态。
3. **列表 API 的 `total` 是否保留？** 需与 `gf-eshop-fe`、`gf-eshop-miniprogram` 协调，
   并同步更新 `tests/test_tx_api.py` 的断言。
4. **`tx_order_logs` 是否单独缩短保留期**（3 年 4.38 亿行，是最大的表）？
5. **粒度按月还是按季？** 两者都成立，按月运维更灵活、按季表更少。
6. **分库时间表**：一旦要跨实例，建单事务与库存扣减需要重新设计，越早评估越好。

**触发阈值**（到线才启动 Phase 1~6）：

| 信号 | 阈值 |
|------|------|
| `tx_orders` 行数 | > 5,000 万（@10 万/日 约 **1.4 年**） |
| `tx_orders` 体积 | > 20 GiB（同上，约 5,400 万行，两个信号基本同时到线） |
| 在线 DDL（加列/加索引）耗时 | > 10 分钟 |
| 归档 `DELETE` 耗时 | > 10 分钟 |
| 订单列表 P99 | > 500 ms（先排查索引与分页，再考虑分表） |

---

## 十一、参考

| 文档 / 代码 | 用途 |
|------------|------|
| [`elasticsearch-search-guide.md`](./elasticsearch-search-guide.md) | 订单检索索引将复用其别名/schema/双写/熔断模式 |
| [`roadmap.md`](./roadmap.md) | 本文对应的决策项登记位置 |
| [`cache-optimization-guide.md`](./cache-optimization-guide.md) | 缓存分层与 P99 基线 |
| `internal/logic/orders/orders.go` | 建单事务与三条查询路径 |
| `internal/logic/payments/payments.go` | 支付回调按 `order_no` 跨表回写 |
| `internal/logic/dashboard/dashboard.go` | 5 处全表聚合，汇总表的动机 |
| `internal/search/` | ES 共用底座（订单检索索引的落点） |
