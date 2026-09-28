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

`origin/main` 停在 `c450516`，本地领先 3 个提交，**未推送**且不需要 force-push：

```
fb403e0  feat(search): 将 ES 方案落地到 products（名称搜索 + 价格区间）
7b9f12b  feat(search): 接入 Elasticsearch，brands 筛选查询走 ES
a2e630b  refactor(products): 拆分商品逻辑为仓储/组装/规则分层   ← 你原有的未推送提交
```

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
