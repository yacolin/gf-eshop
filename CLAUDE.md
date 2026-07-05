# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供本仓库代码工作的指导。

## 项目概述

基于 GoFrame (v2.6.1) 的电商后端 API，使用 MySQL + Redis。单仓库 GoFrame 模板项目，支持品牌、类目、类目品牌关联管理，采用 Cache-Aside 缓存模式。

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
                                                              ↕
                                                            Redis（缓存旁路）
```

- **`api/<模块>/`** — API 接口定义：`IBrandsV1` 接口 + `v1/` 包下的 Req/Res 结构体。通过 `g.Meta` 结构体标签声明式定义路由（path、method、summary）。
- **`internal/controller/`** — 薄控制器层，每个方法直接委托给对应的 service 方法，不做业务处理。
- **`internal/service/`** — 服务接口 + 注册模式（`RegisterX()` 注册，`X()` 获取）。
- **`internal/logic/`** — 业务逻辑实现。每个模块在 `init()` 中通过 `service.RegisterX()` 注册自身实现。
- **`internal/dao/`** — 数据访问对象。`dao/internal/` 为自动生成（禁止手动修改）。`dao/` 包装层可扩展自定义查询方法。
- **`internal/model/entity/`** — 实体结构体（`gf gen dao` 自动生成），用于查询结果扫描。
- **`internal/model/do/`** — DO 结构体（自动生成），用于 Insert/Update 传参。
- **`main.go`** — 入口：导入 MySQL + Redis 驱动，导入 logic 包（触发 `init()`），运行 `cmd.Main.Run()`。
- **`internal/cmd/cmd.go`** — 路由设置：路由分组、中间件、缓存预热。

### 路由注册

```
/              → Hello（示例接口）
/api/v1/       → Brands、Categories、CategoryBrands（业务接口）
```

所有 API 路由使用 `ghttp.MiddlewareHandlerResponse` 中间件（统一结构化 JSON 响应）。

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

品牌和类目均采用相同的 Cache-Aside 模式：
- **缓存键**：`brand:<id>`、`category:<id>`、`brand:all`、`category:all`
- **TTL**：10 分钟
- **读取**：先查缓存 → 未命中 → 查数据库 → 回写缓存
- **写入/删除**：使单条缓存和全量列表缓存均失效
- **预热**：`cmd.go` 中启动时调用 `Warmup()`，预加载全部记录及单条缓存

### 数据模块

| 模块 | 数据库表 | 接口 |
|------|----------|------|
| brands | `sp_brands` | 列表（分页、筛选）、详情、新增、更新、删除 |
| categories | `sp_categories` | 列表、全部、根类目、子类目、层级、树形、详情、新增、更新、删除 |
| category_brands | `sp_category_brands` | 列表（带品牌 left join）、更新（事务内先删后插） |

### 关键约定

- **查询与写入**：Scan 用 `entity.X`，Insert/Update 用 `do.X`
- **字段安全**：使用 `dao.X.Columns().FieldName` 代替硬编码字符串
- **软删除**：含 `deleted_at` 字段的表自动实现软删除（DELETE 转为 UPDATE）
- **错误处理**：记录不存在用 `gerror.NewCode(gcode.CodeNotFound, ...)`，参数非法用 `gcode.CodeInvalidParameter`
- **空结果**：返回 `make([]*T, 0)` 而非 `nil`，避免 JSON 序列化为 `null`
- **分页默认值**：page=1，pageSize=20
- **JSON 字段**：全部使用 snake_case（在 `hack/config.yaml` 中配置）
