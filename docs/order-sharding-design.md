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

> ⚠️ **选 C 的代价（当时这一栏写漏了）**：原生分区（B）的**生命周期天生是自动的** ——
> 一个定时任务 `ADD/DROP PARTITION` 就够；而应用层分表把
> 「建表 / 迁移 / 归档 / 清理 / 对账 / 回滚」全部变成**应用自己负责的运维动作**。
> 这些动作绝大多数是「每月一次」的批量作业，因此必须**例行自动化、破坏性动作留审批** ——
> 现状与补齐清单见 §5.9「生命周期运维的自动化」。

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
> ✅ **Phase 5 修订**：`tx_sub_orders.id` / `tx_order_logs.id` **也必须**去掉 —— 分片运维工具
> （`migrate` / `restore`）是按主键做全列 upsert 的复制，要求主键在**主表与每个分片之间都不重复**；
> 而分片表由 `CREATE TABLE ... LIKE` 建出、各自从 1 计数，实测 202610 分片的子订单拿到 id `1..6`
> 与主表 8 月种子数据撞主键。四张表现在统一由应用生成（std-eshop-db 迁移 `V002`）。
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
| `GET /orders?order_no=x` | ✅ **已实现**（Phase 5）：等价于点查，直接路由到单号内嵌时间对应的分片，单分片内 keyset 精确；两种模式行为一致 | 无 |
| `GET /orders?page&page_size` | ❌ **已移除**（原为「single 回落主表 / monthly 报 `7004`」）：整条 offset 路径删掉了，传这两个参数统一报 1002。理由见下方 2026-10-05 说明 | ⚠️ 破坏性（老客户端必改） |
| `GET /orders?size&cursor` | ✅ **已实现**（Phase 2 定义、Phase 4 跨片）：keyset 分页，`cursor=base64(末位订单ID)`、排序固定 `id DESC`、返回 `next_cursor` + `has_more`（**不返回 `total`**）。分片模式下各活跃片各取 `size+1` 条 → 按 id 倒序归并 → 取前 `size` 条（keyset 的正确性来自「全局前 N 条必在某片的局部前 N 条里」；多取的一条判 `has_more`） | 新增（与 products 的游标契约**逐字段一致**：请求 `cursor`+`size`，响应 `list`+`next_cursor`+`has_more`） |
| `GET /orders?user_id=x` 无时间窗 | ✅ **已实现**：`month=YYYYMM`（路由单分片）/ `created_from`+`created_to`（区间）/ 都不传则**默认近 12 个月**（`orderShard.defaultWindowMonths`）。响应**回显**生效区间与命中月份（`applied_from`/`applied_to`/`applied_months`/`window_defaulted`） | ⚠️ 行为变更（不带时间参数时只看近 12 个月） |
| dashboard 全部聚合 | ✅ **已实现**（Phase 4）：5 项聚合**跨片 fan-out 归并**（订单数/营收求和；趋势按日、状态按状态、热销按商品累加）。**没有**改走日汇总表，原因见 §5.6 | 需跨片归并 |
| 「今天待发货」这类**按业务时间+状态**的运营检索 | `created_at` 分表**无法覆盖** → **必须走 ES**（§5.7） | 新增能力 |
| 库存 `sp_inventories`、商品 `sp_products` | 不分片 | 无 |

⚠️ **API 破坏性变更需要跨仓库协调**：`OrdersListRes.Total` 与 `page` 语义变化会影响
`gf-eshop-fe`、`gf-eshop-miniprogram`，以及现有回归用例
（`tests/test_tx_api.py:152`、`:163`、`:173` 都在断言 `total`）。
建议新增 `cursor` + `size`，`total` 仅在带时间窗时返回精确值，否则返回 `-1` 表示「不精确」。

> 2026-10-05 已按上面的建议落地，并把 products / orders 的游标字段统一成同一份：
> 请求 `cursor` + `size`（默认 20、上限 100），响应 `list` + `next_cursor` + `has_more`。
> 对**已有前端**的影响：products 列表响应里的 `cursor` 字段改名为 `next_cursor`（旧字段不再返回），
> orders 新增 `has_more`（`next_cursor` 语义不变，前端可继续只判它）。
>
> **同日进一步收紧**（见 [`roadmap.md`](roadmap.md) §3 归档「游标契约收紧」）：
> ① 响应**不再返回 `total`** —— 游标分页不做 COUNT，`-1` 这类占位值会被误读成「总共就这么多单」；
> ② orders 的 `page`/`page_size` **彻底移除**（连同 `pageOrders` / `countOrdersByFilter` /
> `offsetListSource`，以及 `7004` 那条分支），传这两个参数报 1002 —— 静默忽略会让 `?page=3`
> 每次都返回第一页，调用方却以为翻页成功。此前的「硬前置」随之消解：列表只剩游标分页，
> `monthly` 下不需要再靠报错兜底。

> 这印证了 CLAUDE.md 里那句话：**ES 是加速器，不是唯一真相源**。
> 分表后它同时变成「跨分片检索层」—— 详情/写走分片 MySQL，列表/运营检索走 ES。

**为什么列表必须有时间范围（以及 month 顺带带来的好处）**

分片按 `created_at` 分月，时间范围就是分片选择器；不带范围就只能对**全部活跃分片** fan-out。
实测（用 general log 数真实的分片查询次数，当时 3 个活跃分片）：

| 请求 | 分片查询次数 |
|------|--------------|
| `?month=202608&size=20` | **1**（只查 202608） |
| `?created_from=2026-08-01&created_to=2026-09-30&size=5` | **2**（8/9 两月各一次） |
| `?size=5`（默认窗口近 12 个月） | **3**（3 个分片都落在窗口内） |

按目标规模（10 万单/日 × 保留 3 年 = 36 片）算，就是 **1 次 vs 36 次/页**。

三条由此确定的语义：

1. **优先级**：`month` > `created_from/to` > 默认窗口。`month` 最省，直接路由单分片。
2. **`order_no` 点查不套时间窗**：客服拿单号找单不能被「默认只看近 12 个月」挡住
   （单号里含时间，路由本来就能定位到单分片）。这条是实测中发现的 bug ——
   最初只在**路由**上绕开了窗口、**过滤条件**仍在生效，结果查不到非当月的单。
3. ~~**带 `month` 后 offset 分页恢复精确**：单分片内的 `page/page_size` + `total` 是准的，
   所以没改游标的老页面只要带上月份就能继续用（不带则跨片，返回 `7004`）。~~
   **已作废**（2026-10-05）：`page`/`page_size` 与整条 offset 路径已删除，带不带 `month`
   都走游标分页 —— `month` 的价值只剩「把翻页收窄到单个分片」（上表：1 次 vs 36 次查询）。
   「窗口内没有分片不能回落主表」这条依然成立 —— 那会返回与该月份无关的行，所以空结果就是空结果。

> 口径提醒：月份按**下单时间**（`created_at`，即分片键）。若要按「支付/发货时间」筛，
> 那不是分片键、无法路由，只能全片 fan-out 或走 ES（§5.7）。

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
- ~~有了它，看板的「订单总数/营收/7 日趋势/状态分布」全部变成单表查询~~
  —— **这条结论在 Phase 4 落地时被推翻，实际没这么做**，原因是核对后发现两个硬伤：
  1. **覆盖不全**：表里只有 `order_cnt` / `cancelled_cnt` / `gmv` / `paid_amount` / `refund_amount`，
     而看板要的是 5 项 —— 订单总数与营收能对上，但**趋势的金额口径**（`SUM(pay_amount)` 全部订单）、
     **状态分布**（6 种状态）、**热销商品**（按商品维度）在这张表里都没有对应列；
  2. **会读到过期数据**：表只有 Phase 2 的一次性回填，**事件驱动的累加还没实现**
     （写入时机那一节写的是设计，不是现状）。今天新建的单不在表里，
     拿它当看板数据源会少数据。
  所以 Phase 4 改为**跨片 fan-out 归并**（实测与主表直查逐项一致，且分片数不大时成本很低）。
- **这张表的定位调整**：作为**可选优化**保留 —— 要用它必须先把「写入时机」那一节的事件驱动
  累加做出来（建单/支付/退款/取消四个钩子），并给它补上趋势金额与状态维度（或接受降级口径）。
  分片数涨到几十片、fan-out 成本变高时再启用。另注：热销商品的 fan-out 必须
  **每片返回全部分组**再归并，不能每片各取 TOP N（商品按时间分片，局部 TOP N 只能得到近似结果）。
- **进度**：表已建、历史已回填（Phase 2，实测 32 天，与主表直查逐项一致）；
  读路径在 Phase 4 **没有**切到它（见上），当前由 fan-out 承担。

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
| 建表 | 目前是**写路径自动兜底当月**（`EnsureShardTables`，幂等 + 进程内缓存 + 模板回退）；**尚未**做「启动时/定时预建次月」——见下方「生命周期运维的自动化」 |
| DDL 一致性 | 用模板建表 + 巡检脚本比对 `SHOW CREATE TABLE`（144 张表结构漂移是主要风险） |
| 加列 | 遍历活跃分片执行，用 `gh-ost`/`pt-online-schema-change`；注意 `LIKE` 建表不会自动跟进 |
| 时区 | 分片月份推导必须与 `created_at` **同一个时区**。当前配置 `loc=Local`，建议统一固定 `Asia/Shanghai`，否则跨时区部署会出现「同一订单算出不同月份」 |
| 备份 | 可按片备份：热片每日全量，冷片一次性归档 |
| 监控 | 单片行数/体积、`order_shard_fallback_total`、跨片查询 P99、汇总表与分片对账差异 |
| 表数量 | 3 年后 144 张，MySQL 8/9 无压力；但要控制 `information_schema` 类监控查询的频率 |

#### 生命周期运维的自动化（现状 vs 生产应有）

选应用层分表（§四 C）意味着**生命周期由应用自己负责**：原生分区只需要
`ADD/DROP PARTITION`，而这里要自己管建表、迁移、归档、清理、对账、回滚。
当前实现只自动化了其中一件（写路径兜底当月分片），其余靠人工执行 CLI ——
生产环境应按「**例行自动化、破坏性审批**」补齐：

| 动作 | 现状 | 生产应有 | 理由 |
|------|------|---------|------|
| 建当月分片 | ✅ **写路径自动兜底**（幂等 + 缓存；主表缺失时按「最新分片 → 归档表」回退模板） | 再加**定时预建次月** | 避免月初第一单踩 DDL；结构漂移能被巡检提前发现 |
| 历史迁移 `migrate` | 人工 CLI（幂等、全列覆盖） | CI/定时任务触发，**对账通过前不放量** | 可重复的批量作业，人工触发可接受，但要有进度与差异指标 |
| 日汇总累加 | ❌ **只有一次性回填，没有累加** | **每晚聚合前一日** | 比事件钩子简单且能自愈；dashboard 不再只依赖回填 |
| 归档 `archive` | 人工 CLI（无损性预检 + 打印回滚语句） | **提案 + 审批**：任务算好保留期、打印待执行的 RENAME，人/CI 批准后执行 | RENAME 会改变真相源，需要人工确认 |
| 清理 `purge` | 人工 CLI + 观察窗口 + `--force` | **保留「窗口 + 审批」**，不要全自动 | 不可逆；自动化 bug 的代价太大 |
| 对账 `verify` | 人工 CLI | **定时跑 + 差异告警** | 这是让「人工运维」变安全的关键一环 |
| 回滚 `restore` | 人工 CLI | **保持人工** | 这是人的决策，不是例行任务 |
| 备份门禁 | ❌ 目前没有 | `archive` / `purge` 之前断言「存在近期备份」 | 生产必需 |

**建议的定时任务形态**（一个子命令 + cron / systemd timer），只做四件**不破坏数据**的事：

1. 预建次月分片；
2. 对账上月，有差异就告警；
3. 聚合前一日的 `tx_order_daily_stats`（把缺失的累加补上）；
4. 打印本月归档提案（列出该归档哪些分片、以及那条 RENAME，**不执行**）。

`archive` / `purge` / `restore` 继续由人工或 CI 审批触发 —— 既符合生产实践，
也避免「自动化 bug 顺手把数据删了」。

> 📌 相关已知缺口：日汇总**累加**未实现（§5.6 已说明该表当前只有回填），
> 因此看板目前走的是跨片 fan-out；上面第 3 项落地后，看板才有条件改走汇总表。

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

### Phase 3：开双写 + 影子读对账 ✅ 已完成

- 动作：写入同时写主表与分片；读仍走主表。
  影子读：抽样请求同时读两边，比对结果并打点差异率。
- 回滚：关双写，主表仍是完整真相源。

**实现取舍（与最初设想的差异，都有理由）**

1. **双写不逐条改写入语句，而是「按订单整单镜像」**。
   最初的想法是「每条 insert/update 都写两遍」，但那要改 8 个写函数
   （4 张表 × insert/update），任何一条漏改都会让分片永久缺数，且以后每加一个写接口
   都要记得改。现在改成在**同一个事务内**按订单把 4 张表的当前行从主表复制到分片：

   ```sql
   INSERT INTO <分片> SELECT * FROM <主表> WHERE <订单键> = ?
   ON DUPLICATE KEY UPDATE <全列覆盖>
   ```

   好处：① 只挂在 3 个编排点（建单 / 改状态 / 支付回写），不会漏；
   ② 同事务 ⇒ 与主写**原子**（`gdb` 的 `Exec` 会取 ctx 里的事务，见 `Core.DoExec`）；
   ③ 全列覆盖 ⇒ 幂等且**自愈**：上一次漏掉的更新会在下一次写入时被带上。

2. **取键按「有索引」选，而不是一律用 `order_no`**：`tx_order_logs` **没有 `order_no` 索引**
   （只有 `idx_order_id`），所以日志表按 `order_id` 定位，避免大表全扫。
   因此 `MarkPaidByOrderNo` 增加了 `orderID` 入参。

3. **迁移缺口用「重跑 migrate」重放，不需要单独的 replay 命令**：
   把迁移的 `ON DUPLICATE KEY UPDATE` 从空更新（`id = VALUES(id)`）改成**全列覆盖**后，
   重跑 `--action=migrate` 就不只是「不重复」，而是把主表当前内容重新同步到分片 ——
   正好覆盖设计里说的「迁移期间落在历史月份的零星更新」。实测：直接改主表不写分片
   （复现缺口）→ 重跑 migrate → 分片追平，且不产生重复行。

**开关（都在 `orderShard` 段，缺省即关闭 = 现状）**

| 配置 | 缺省 | 作用 |
|------|------|------|
| `orderShard.mode` | `single` | `single` 读写主表；`monthly` 写分片、读分片（Phase 4/5） |
| `orderShard.dualWrite` | `false` | 写主表的同时把同一订单镜像进分片 |
| `orderShard.shadowReadPercent` | `0` | 影子读抽样比例（0~100），抽样比对主表与分片 |
| `orderShard.readShardsPercent` | `0` | 读侧灰度比例（0~100），见 Phase 4 |
| `orderShard.defaultWindowMonths` | `12` | 列表未传时间参数时默认回溯的月数（决定 fan-out 多少个分片） |

影子读抽样是**确定性**的（`crc32(order_no) % 100`），同一条订单每次都会命中同一样本，
便于复现；且它只记日志与计数，绝不改变接口返回值。

**可观测性**（Redis 计数，直接看差异率）：

```bash
redis-cli MGET order:shard:shadow:total order:shard:shadow:diff order:shard:shadow:error
redis-cli MGET order:shard:mirror:mirror_ok order:shard:mirror:mirror_failed
redis-cli KEYS 'order:shard:fallback:*'   # 读灰度回落主表的次数（应很小且可解释）
```

**实测**（`dualWrite=true` + `shadowReadPercent=100`，当月分片 202610）

- 建单：主表与分片的订单/子订单/明细/日志四张表都落到 202610，逐列一致；
- 状态流转 paid→shipped→delivered→completed：每一步主表与分片的订单与子订单状态都一致，
  日志条数一致（1 + 4）；
- 支付回写（跨模块路径）：回调后主表与分片的 `status/payment_status/paid_at` 一致；
- 影子读：计数 `total=11`、`diff=0`、`error=0`；双写 `mirror_ok≥2`、`mirror_failed=0`；
- **反向验证**：把分片行改坏（`status`/`pay_amount`）后读同一订单 → 差异计数 +1、
  日志打出两侧指纹、**接口仍返回主表的正确值**（证明影子读不会影响结果，且这套对账不是摆设）；
- **重放验证**：被改坏的分片行重跑 migrate 后追平；只改主表复现缺口后重跑 migrate 同样追平；
- **默认关闭回归**：与 Phase 3 之前的 HEAD 逐字段对比 14 项零差异；`tests/test_tx_api.py` 36/36。

> Phase 4 才会把读路径切到分片（`mode=monthly`），届时列表/看板聚合需要跨片 fan-out 或
> 日汇总表，目前 `shardForScan` 在 monthly 下仍是明确报错而不是静默读主表。

### Phase 4：灰度切读 ✅ 已完成

**开关与语义**

| 配置 | 作用 |
|------|------|
| `orderShard.readShardsPercent` | 读侧灰度比例（0~100）。点查按 `crc32(order_no)%100` **确定性**抽样，列表/看板按请求抽样 |
| `orderShard.mode` | `single`：主表仍是真相源 ⇒ 分片未命中/出错**回落主表**并计数；`monthly`：恒定读分片、**不回落**（主表已停写，回落只会读到过期数据） |

- 切读阶段对账方向**反过来**：主读走分片时再读一次主表比对（`compareShardAgainstMain`，
  只在 `single` + `shadowReadPercent>0` 时做），差异率与回落次数都进 Redis 计数。
- **列表**：游标分页跨片 fan-out —— 各活跃片各取一页，按 id 倒序归并取前 `size` 条；
  keyset 的正确性来自「全局前 N 条必然出现在某片的局部前 N 条里」，所以每片取 `size + 1` 条就够
  （多出的那条只用于判 `has_more`）。
  ~~offset 分页无法跨片正确归并：`single` 下回落主表（老客户端不受影响），`monthly` 下明确报 `7004`。~~
  **已作废**（2026-10-05）：offset 分页整条路径已移除，列表只有游标分页。
- **看板**：5 项聚合全部跨片 fan-out（详见 §5.6 的修正说明）。
- **自动建分片**：写入前用 `CREATE TABLE IF NOT EXISTS ... LIKE` 幂等建表（进程内缓存，
  每片只做一次），于是**双写会自动吸入当月分片**，不会再因为「当月分片还没建」而静默失败。
  ⚠️ DDL 必须放在**事务外**（MySQL 的 DDL 隐式提交，放进写事务会毁掉原子性），
  因此建分片发生在写事务之前 —— 建单、改状态、支付回调三个入口各一次。

**一处必须提防的坑：读源 ≠ 写目标**

灰度期间「父单从哪读」和「这条单写去哪」是两件事：`single` 下落回主表的订单，
读源是主表但写目标仍是主表（+ 镜像）；`monthly` 下两者都是分片。
改造前 `UpdateStatus` 直接拿读路径返回的分片去写 —— 灰度命中分片时会**绕过主表写入**。
现在读路径返回的分片只用于「同源读子表」，写目标由 `shardFromOrderNo`（由 `mode` 决定）单独算。

**实测**

| 项 | 结果 |
|----|------|
| 点查走分片 | 详情与主表/分片一致；反向对账 `total≥1`、`diff=0` |
| 分片缺行回落 | 只在主表存在的订单仍能点到；`fallback:read_order_miss` +1 |
| 未镜像的行 | 列表/看板**看不见它**（分片 2001 / 主表 2002）—— 这正是切读前必须保证分片完整的原因 |
| 列表跨片 | 翻完全部 2001 条与主表 id 集合完全一致、严格 id 倒序、无重复 |
| ~~offset 分页~~ | ~~`single` 下回落主表且 total 正确；`monthly` 下返回 `7004`~~ → 该能力已于 2026-10-05 整体移除 |
| 看板 5 项聚合 | 订单总数/营收/状态分布/趋势/热销 TOP1 与主表直查逐项一致 |
| 反向验证 | 改坏分片行 → 切读确实读到坏值，但反向对账立刻 `diff+1` 并打出两侧指纹（灰度的风险面与安全网同时被证实） |
| `monthly` 终态 | 建单/改状态/支付回写只落分片、主表行数不变；老数据（202608 迁移分片）仍可读 |
| 自动建分片 | 删掉 202610 分片后开双写建单：4 张表自动建出（日志仅 1 条），`mirror_ok=2` / `mirror_failed=0` |
| 默认关闭回归 | 与 Phase 4 前 HEAD 逐字段对比 14 项零差异；`tests/test_tx_api.py` 36/36；游标回归 0 失败 |

- 回滚：`readShardsPercent` 置 0（或 `mode` 改回 `single`）即秒级回落主表。

### Phase 5：停双写 ✅ 已完成

- 动作：`mode: monthly` + `dualWrite: false`，停止写主表。
- **终态语义**（与灰度期的区别）：
  - 写：只落分片，不再镜像；
  - 读：只走分片，**不回落主表** —— 主表已停写，回落只会读到过期数据；
  - ~~无 `order_no` 的 offset 分页返回 `7004`~~ → 2026-10-05 起 offset 分页已移除，传 `page`/`page_size` 统一报 1002；**带 `order_no` 的列表等价于点查，仍精确**（见下）。
- ⚠️ **前置已解除**：前端订单页已改用游标分页（`gf-eshop-fe` 提交 `6398ed6`）。
  另一个前端仓库 miniprogram **目前没有订单页面**（只有 index/product），无需改动 ——
  这修正了 §2.6 原先「两个仓库都要跟进」的说法。

**加固：终态下的运维动作**

- `shard --action=migrate` **禁止执行**（返回 `7006`）：迁移方向是「主表 → 分片」，
  主表停写后再跑就会拿冻结的旧快照把分片里的新数据覆盖掉 —— 不可逆的数据损坏。
  要回灌必须先 `mode` 改回 `single`，并明确那是回滚动作。
- `shard --action=verify` 仍放行（只读），但会明确告警：切换之后发生的状态变更只落分片，
  所以「主表 vs 分片」出现差异是**预期**的，不代表损坏 —— 此时它的价值是
  「一眼看出哪些行在切换后被改过」。

**带 `order_no` 的列表查询要单独处理**

`GET /orders?order_no=x`（不带分页参数）在 `monthly` 下原本会落到 offset 分支被拒，
但它其实等价于点查：现在直接路由到单号内嵌时间对应的那一个月，
单分片内 offset 与 keyset 都精确，两种模式行为一致。

**测试与本地配置解耦（顺带修掉的隐患）**

单测原先隐式依赖开发机本地的 `manifest/config/config.yaml`：本地一切到 `monthly`，
`TestSingleModeRoutesToBaseTable` 就失败 —— 测试结果取决于「谁的机器」。
现在模式相关的用例都用 `gcfg.AdapterFile.Set` 在**内存里**声明自己要测的模式，
并补了 monthly 路由矩阵、双写/灰度开关的用例；接口套件也改为**游标优先**
（当时新增 3.7 覆盖「single 回落主表 / monthly 报 `7004`」—— 该用例已于 2026-10-05
随 offset 分页移除，改为断言 `page`/`page_size` 被拒绝），因此两种模式下都能跑。

**实测**

| 项 | 结果 |
|----|------|
| 接口套件 | `single` 与 `monthly` 两种模式**都是 37/37** |
| `monthly` 写入 | 建单/改状态/支付回写只落分片、主表行数不变（Phase 4 已验证） |
| `migrate` 守卫 | `monthly` 下执行 → 拒绝并返回 `7006`，退出码 1 |
| `verify` 守卫 | `monthly` 下执行 → 告警说明差异属预期，然后正常出对账结果 |
| `?order_no=x` | 两种模式下都能精确查到（monthly 下走单号对应的单分片） |
| 单测 | 新增模式矩阵用例；显式设定模式，不再受本地配置影响 |

- 回滚：主表停写期间的新数据只在新分片 —— 用下面的 `restore` 回灌，主表在停写后仍是完整的历史快照。

**回滚工具：`shard --action=restore`（分片 → 主表）**

```bash
# ① 仍在 monthly 下回灌（读路径不受影响，主表逐步追平）
./main shard --action=restore --from=2026-08 --to=2026-10
# ② 确认主表与分片一致
./main shard --action=verify --from=2026-08 --to=2026-10
# ③ 再把 orderShard.mode 改回 single
```

顺序不能反：**先切 single 会让停写期间产生的新单在主表里查不到**（single 下点查不回落到分片）。
因此 `restore` 只在 `monthly` 下允许、`migrate` 只在非 `monthly` 下允许 —— 两个方向互斥，各有守卫。

**⚠️ 回灌暴露出的主键重复问题（Phase 5 的遗留阻塞）**

`restore` 与 `migrate` 都是「按主键做全列 upsert」的复制工具，因此它们隐含了一个
**比文档原先写的更强的前提**：主键不只要在「被引用的表」里全局唯一，
而是要在**主表与每个分片之间都不重复**。否则 upsert 会按主键命中并**覆盖无关的行**。

实测踩到了：`tx_sub_orders.id` 是 `AUTO_INCREMENT`，而分片表由 `CREATE TABLE ... LIKE` 建出、
**各自从 1 开始计数** —— 202610 分片的第一批子订单拿到 id `1..6`，
正好与主表的 8 月种子数据撞主键；一次 `restore` 就把那 6 行 8 月数据覆盖成了 10 月测试单的内容。
（`tx_order_logs.id` 同理。`tx_orders` / `tx_order_items` 在 Phase 2 已改为应用生成的全局唯一 ID，没有这个问题。）

- **已加的护栏**：两个复制工具在复制前做**主键冲突预检**（按业务键判断「同一 id 是否同一行」：
  订单看 `order_no`、子订单看 `sub_order_no`、日志看 `order_id + created_at`），
  命中即拒绝并报出冲突行数，不再静默覆盖。
- **根因已修**：这两张表的主键也改成应用生成的全局唯一 ID（与订单/明细同一序列、同一布局；
  段内顺序 订单 → 子订单 → 明细… → 日志），DDL 见 std-eshop-db 的基线 `sql/tx_p0.sql` /
  `sql/tx_p1.sql` 与前向迁移 `sql/migrations/V002__sub_orders_logs_global_id.sql`。
  ⚠️ `V002` 不改已存在的分片表，**旧分片要另跑一次同样的 ALTER**（迁移文件头部给了生成语句）。
- **修完后的实测**（202610，monthly）：建单后四张表主键分别为 `…129/130/131/132`（应用生成、
  互不相同、与主表无冲突）→ `restore` 预检放行、回灌 5 行 → `verify` 4 项**校验和全对**
  → 改回 `single` 后这张「只曾在分片里存在」的订单能从主表读到、列表（当时仍是 offset 分页）恢复可用。

### Phase 6：归档旧表 ✅ 已完成

```bash
# 归档（无损性预检通过才动手；一条 RENAME 改四张表，并打印回滚语句）
./main shard --action=archive --suffix=202610
# 观察窗口内回滚（把打印出来的那条 RENAME 贴回去即可）
# 窗口满后清理：不加 --force 只预览 DROP 语句
./main shard --action=purge --suffix=202610 [--force]
```

**为什么是改名而不是删**：30 天窗口内一条 `RENAME` 就能回滚；同时旧名字留着会让
人和工具误以为它还是真相源，所以必须改名成 `tx_*_legacy_<YYYYMM>`。

**归档前的无损性预检**（这是唯一真正重要的检查）：
主表里的每一行都必须已经存在于某个分片里 —— 否则改名后这些行就查不到了
（应用在 monthly 下只读分片）。逐表用 `LEFT JOIN ... IS NULL` 统计，比「行数相等」严格。
另外还要求 `mode=monthly`（主表还在被写的话，改名会让写入直接失败）。

**清理窗口**：设计文档要求「观察 30 天」，但 `RENAME` 的时刻无法从 schema 反推
（RENAME 不改变 `CREATE_TIME`），所以规则定为**归档后缀所在月再过完一整月**：
`legacy_202610` 最早 2026-12-01 可清理 —— 比 30 天更保守，且只看表名就能判断。

**归档之后的守卫**（避免出现「表不存在」这类底层报错）：
`migrate` / `restore` / `verify` 都会明确回一句「已随 Phase 6 归档，不再适用：… 如需执行请先 RENAME 回来」；
服务启动时也会自检：`mode=monthly` 下提示「符合 Phase 6 的预期状态」，
而 `mode=single` 下主表缺失会**告警**（这种组合下订单接口会持续失败）。

**⚠️ 演练中发现并修掉的缺口**：建分片表是 `CREATE TABLE <分片> LIKE <模板>`，
原先模板固定是主表 —— 主表一归档，**月初滚动要建新月份分片时直接报 Error 1146**
（`tx_orders` 不存在），monthly 模式跑不下去。现在模板按「主表 → 最新分片表 → 最新归档表」
回退，并排除「目标表自己」（否则会 `LIKE tx_orders_202610` 去建 `tx_orders_202610`，
报 Not unique table/alias）；已存在的分片直接跳过 DDL。

**实测**（202610 归档）

| 项 | 结果 |
|----|------|
| 无损性预检 | 通过（主表 2000/2000/5025/0 行都在分片里）后一条 RENAME 改四张表 |
| 归档后应用 | `tests/test_tx_api.py` **37/37**；建单/详情/游标列表全部正常；归档表行数保持不变（写入没有落回旧表） |
| 归档后建新月份分片 | `--action=create --from=2026-11` 正常，新分片结构与既有分片**逐字一致** |
| 运维工具守卫 | `restore`/`verify` 明确报「已随 Phase 6 归档」；`migrate` 先被终态守卫拦下 |
| 清理窗口 | `purge` 预览出 4 条 DROP 语句并提示窗口最早 2026-12-01；`--force` 提前执行被拒绝 |
| **回滚** | 执行打印出的 RENAME → 主表 2000/2000/5025 回来、归档表残留 0 → 切回 `single` 后详情/列表/建单全部正常 |
| 启动自检 | 归档后 monthly 下提示「符合 Phase 6 的预期状态」 |

- 回滚：RENAME 回来即可（观察窗口内）；窗口过后表已删，只能从头重建（此时分片已是唯一真相源）。

---

## 八、回滚方案总表

| 阶段 | 回滚动作 | 数据风险 | RTO |
|------|---------|---------|-----|
| Phase 0 | revert 代码 | 无 | 分钟级 |
| Phase 1 | 切回旧生成器 | 无 | 分钟级 |
| Phase 2 | 删除分片表 | 无（主表未动） | 分钟级 |
| Phase 3 | 关双写开关 | 无（主表是真相源） | 秒级 |
| Phase 4 | `orderShard.readShardsPercent` 置 0（或 `mode` 改回 `single`） | 无 | **秒级** |
| Phase 5 | `shard --action=restore` 回灌 + 改回 `single` | 停写期间新单只在分片，回灌后主表追平（**前置：主键不重复**，见 §7 Phase 5） | 分钟级 |
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
