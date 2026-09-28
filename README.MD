# gf-eshop

基于 [GoFrame v2.6.1](https://goframe.org) 的电商后端 API。单仓库结构，
`internal/logic/` 下现有 **43 个业务模块**，涵盖商品、SKU、库存、订单、支付、购物车、
用户、商家、营销、评价、权限等。

## 技术栈

| 组件 | 版本 | 用途 |
|------|------|------|
| Go | 1.25 | 语言运行时 |
| GoFrame | v2.6.1 | Web 框架 / ORM / 配置 / 日志 |
| MySQL | 8.x | 主存储（库名 `eshop_db`，74 张表） |
| Redis | 7.x | 缓存（单条实体、列表 ZSET、多级缓存 L1/L2） |
| Elasticsearch | 7.17.4 | 检索（组合筛选、文本子串、价格区间） |
| analysis-ik | 7.17.4 | ES 中文分词插件（**必须安装**） |

## 架构

### 分层

```
请求 → ghttp.Server → Controller（薄层）→ Service 接口 → Logic → DAO → MySQL
                                                          ↕        ↕
                                                    Redis（缓存）  Elasticsearch（检索）
                                                                    ↑
                                                        internal/search（ES 底座）
```

- `api/<模块>/` — 接口定义，用 `g.Meta` 标签声明式定义路由
- `internal/controller/` — 薄控制器，只做转发
- `internal/service/` — 服务接口 + 注册模式（`RegisterX()` / `X()`）
- `internal/logic/` — 业务逻辑，各模块在 `init()` 中注册自身
- `internal/dao/` — 数据访问（`dao/internal/` 自动生成，禁止手改）
- `internal/search/` — ES 共用底座（客户端、熔断、索引/别名、分批 bulk）

### 查询侧是两层，别混淆

| 层 | 承担 | 覆盖范围 |
|----|------|----------|
| **Redis** | 热点 / 简单查询 | 单条实体缓存、列表 ZSET、products 的 L1+Bloom+L2 |
| **Elasticsearch** | 组合筛选 / 文本检索 | 多字段子串检索、价格区间、排序、翻页 |

**ES 是加速器，不是唯一真相源**：任何一步失败都会熔断并自动回落 MySQL 原有查询，
配置 `elasticsearch.enabled: false` 可一键关停。详见
[docs/elasticsearch-search-guide.md](docs/elasticsearch-search-guide.md)。

## 快速开始

### 1. 启动依赖服务

```bash
# MySQL（库名 eshop_db）
brew services start mysql

# Redis
brew services start redis

# Elasticsearch 7.17.4
brew services start elasticsearch-full
```

### 2. 安装 IK 分词插件（ES 必须）

IK 已从 `medcl/elasticsearch-analysis-ik` 迁移到 `infinilabs/analysis-ik`，
GitHub Releases 不再提供插件包，走官方 CDN：

```bash
curl -sSL -o /tmp/ik.zip https://get.infini.cloud/elasticsearch/analysis-ik/7.17.4

# 用 ES 自带插件 CLI 安装（会校验版本，别手工解压）
/opt/homebrew/opt/elasticsearch-full/bin/elasticsearch-plugin install --batch "file:///tmp/ik.zip"
brew services restart elasticsearch-full

# 验证
curl -s "http://localhost:9200/_nodes/plugins?pretty" | grep -i analysis-ik
curl -s -X POST "http://localhost:9200/_analyze" -H 'Content-Type: application/json' \
  -d '{"analyzer":"ik_max_word","text":"苹果 Apple 华为"}'
# => ["苹果","apple","华为"]
```

### 3. 创建配置文件

> ⚠️ `manifest/config/config.yaml` 被 `.gitignore` 忽略（`**/config/config.yaml`），
> **新克隆的仓库里没有这个文件**，服务无法启动，需要手动创建。

```yaml
server:
  address: ":8000"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"
  logStdout: true
  accessLogEnabled: false

database:
  default:
    link: "mysql:root:123456@tcp(localhost:3306)/eshop_db?charset=utf8mb4&parseTime=True&loc=Local"
    maxActive: 100
    maxIdle: 10
    idleTimeout: "60s"
    maxConnLifetime: "30m"

redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    maxActive: 300
    maxIdle: 200
    minIdle: 50
    idleTimeout: "300s"
    waitTimeout: "3s"
    dialTimeout: "2s"
    readTimeout: "2s"
    writeTimeout: "2s"
    maxConnLifetime: "30m"

jwt:
  secret: "eshop-jwt-secret-key-2024"
  access_expire: 1800
  refresh_expire: 604800

logger:
  level: "all"
  stdout: true

elasticsearch:
  enabled: true                      # 总开关：false 时全部查询回落 MySQL
  addresses:
    - "http://127.0.0.1:9200"
  username: ""
  password: ""
  indexPrefix: "eshop_"
  timeout: "3s"
  maxResultWindow: 10000             # from+size 超过此值直接回落 DB
  writeRefresh: "wait_for"           # 双写刷新策略
  bulkChunk: 500                     # 重建时单次 _bulk 的文档数上限
```

> 上面的数据库口令是本地开发默认值（与 `hack/config.yaml` 一致），部署时请替换。

### 4. 构建运行

```bash
make run            # 开发模式，热加载（需要 gf CLI，首次会自动安装）
# 或
make run.prod       # go build -o main . && ./main
```

启动时会并行预热缓存与 ES 索引，日志中出现
`商品 ES 索引与 DB 一致（N 条…），跳过重建` 即表示 ES 侧正常。

访问：

- Swagger UI：<http://127.0.0.1:8000/swagger>
- OpenAPI JSON：<http://127.0.0.1:8000/api.json>

## 开发命令

```bash
# 运行
make run                     # 热加载开发服务器（gf run main.go）
make run.prod                # 构建并运行
make build                   # gf build -ew
make stop                    # 杀掉 :8000 上的进程
make cli                     # 下载安装 GoFrame CLI

# 代码生成（均依赖 gf CLI，会自动安装）
make dao                     # 根据数据库表生成 DAO / DO / Entity
make service                 # 生成服务接口
make ctrl                    # 根据 api 定义生成控制器
make enums                   # 生成枚举

# 重建 Elasticsearch 索引
./main reindex                     # 重建全部已接入实体（brands + products）
./main reindex --entity=brands     # 只重建 brands
./main reindex --entity=products   # 只重建 products
# 注意：GoFrame gcmd 会把多余的位置参数当成「多级命令名」，
#       实体名只能用 --entity= 传，`main reindex brands` 会报 command not found
```

## 目录结构

```
gf-eshop/
├── api/                      # 接口定义（<模块>/<模块>.go + v1/Req-Res）
├── internal/
│   ├── cmd/cmd.go            # 路由、中间件、预热、reindex 子命令
│   ├── controller/           # 薄控制器
│   ├── service/              # 服务接口 + 注册
│   ├── logic/                # 业务逻辑（43 个模块）
│   ├── dao/                  # 数据访问（internal/ 自动生成）
│   ├── model/entity|do/      # 实体与 DO（自动生成）
│   ├── middleware/           # ErrorHandler、鉴权等
│   └── search/               # ES 共用底座
├── docs/                     # 方案与流程文档
├── hack/                     # gf CLI 配置与 Makefile 片段
├── manifest/config/          # 运行时配置（被 gitignore，需自行创建）
└── tests/                    # API 自动化测试脚本
```

## 数据表

按前缀分组，共 74 张：

| 前缀 | 数量 | 说明 |
|------|------|------|
| `sp_` | 15 | 商品域（品牌、类目、商品、SKU、属性、库存…） |
| `tx_` | 15 | 交易域（订单、支付、退款、购物车、售后、物流…） |
| `mch_` | 11 | 商家域（商家、结算、提现、资质…） |
| `sys_` | 9 | 系统域（员工、角色、权限、部门、操作日志…） |
| `usr_` | 8 | 用户域（用户、地址、等级、积分…） |
| `mkt_` | 6 | 营销域（促销、规则、库存、用量日志…） |
| `rev_` | 6 | 评价域（评价、媒体、回复、审核、统计…） |
| `base_` | 3 | 基础（通知、模板、已读） |

## 文档

| 文档 | 内容 |
|------|------|
| [CLAUDE.md](CLAUDE.md) | 项目开发指引：架构、约定、搜索层规则 |
| [docs/elasticsearch-search-guide.md](docs/elasticsearch-search-guide.md) | ES 接入方案、products 接入实录、行为变化与坑位 |
| [docs/cache-optimization-guide.md](docs/cache-optimization-guide.md) | 缓存优化历程（ZSET+Lua+singleflight、P99 目标） |
| [docs/curd-workflow.md](docs/curd-workflow.md) | 新增 CRUD 模块的完整流程 |
| [docs/roadmap.md](docs/roadmap.md) | 待办事项、优先级与已知问题 |
| [docs/perf-workflow.md](docs/perf-workflow.md) | 压测流程 |

## 已知问题

见 [docs/roadmap.md](docs/roadmap.md)。其中需要留意的两条：

- **ES 多词查询语义与 MySQL 不同**：ES 按词 AND 匹配（`三星 智能` 返回 4 条），
  而 MySQL `LIKE '%三星 智能%'` 要求字面子串（返回 0 条）。ES 降级期间会退回字面匹配。
- **改 ES mapping / analyzer 时必须递增该模块的 schema 常量**
  （`brandIndexSchema` / `productIndexSchema`），否则结构变更不会生效。
