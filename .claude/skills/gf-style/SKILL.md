---
name: "gf-style"
description: "GoFrame 代码规范（GF Layered Structure）。基于 gf-eshop 项目的 brands/categories 模块参考实现提炼。适用于 code review、新模块开发、规范统一。"
---

# GoFrame 代码规范（GF Layered Structure）

## 适用范围

本规范适用于本项目（gf-eshop）所有 GoFrame 模块，包括 brands、categories、category_brands 等。

## 核心原则

1. **分层而治** — Controller → Service（接口） → Logic（实现） → DAO → MySQL，每层职责清晰
2. **声明式路由** — 通过 `g.Meta` 结构体标签定义路由，不在代码中手写 URL
3. **注册模式** — Logic 包通过 `init()` 中调用 `service.RegisterX()` 注册自身
4. **字段安全** — 用 `dao.X.Columns().FieldName` 代替硬编码字符串
5. **缓存旁路** — 读走 Cache-Aside，写操作立即失效缓存

## 1. 目录结构

```
gf-eshop/
├── main.go                           # 入口：导入 driver + logic 包，运行 cmd
├── api/{module}/                     # API 接口定义
│   ├── {module}.go                   #   接口声明（IBrandsV1）
│   └── v1/{module}.go               #   Req/Res 结构体（含 g.Meta 路由标签）
├── internal/
│   ├── cmd/cmd.go                    #   HTTP 服务器 + 路由注册 + 缓存预热
│   ├── consts/consts.go              #   项目常量
│   ├── controller/{module}/          #   控制器（薄层，委托给 service）
│   │   ├── {module}_new.go          #     ControllerV1 + NewV1()
│   │   └── {module}_v1_{module}.go  #     各方法实现
│   ├── dao/                          #   数据访问层
│   │   ├── {module}.go              #     DAO 包装（可扩展自定义方法）
│   │   └── internal/{module}.go     #     DAO 内部实现（自动生成，禁止修改）
│   ├── logic/{module}/               #   业务逻辑层
│   │   ├── {module}.go              #     业务逻辑（init 注册 + CRUD 实现）
│   │   └── cache.go                  #     缓存存取 + 预热（可选）
│   ├── model/
│   │   ├── entity/{module}.go       #     实体结构体（自动生成）
│   │   └── do/{module}.go           #     DO 结构体（自动生成）
│   └── service/{module}.go          #   Service 接口 + 注册函数
├── hack/
│   └── config.yaml                  #   gf CLI 配置（数据库连接 + 表定义）
└── manifest/
    └── config/config.yaml           #   运行时配置（server、database、redis）
```

### 示例：brands 模块

```
api/brands/
├── brands.go              # IBrandsV1 接口
└── v1/brands.go           # ListReq, DetailReq, CreateReq, UpdateReq, DeleteReq

internal/
├── cmd/cmd.go              # brands.NewV1() 绑定路由
├── controller/brands/
│   ├── brands_new.go       # ControllerV1 + NewV1()
│   └── brands_v1_brands.go # List/Detail/Create/Update/Delete → service.Brands()
├── dao/
│   ├── brands.go           # Brands 全局对象（包装 internal 实现）
│   └── internal/brands.go  # 自动生成：表名、字段常量、Ctx()、Transaction()
├── logic/brands/
│   ├── brands.go           # sBrands 实现 IBrands 接口
│   └── cache.go            # 缓存存取 + Warmup 预热
├── model/
│   ├── entity/brands.go    # Brands 结构体（自动生成）
│   └── do/brands.go        # do.Brands（自动生成）
└── service/brands.go       # IBrands 接口 + RegisterBrands() + Brands()
```

## 2. 文件命名规范

| 层 | 命名规则 | 示例 |
|---|---------|------|
| API 接口声明 | `api/{module}/{module}.go` | `api/brands/brands.go` |
| API Req/Res | `api/{module}/v1/{module}.go` | `api/brands/v1/brands.go` |
| 控制器结构体 | `{module}_new.go` | `brands_new.go` |
| 控制器实现 | `{module}_v1_{module}.go` | `brands_v1_brands.go` |
| Service 接口 | `internal/service/{module}.go` | `internal/service/brands.go` |
| 业务逻辑 | `internal/logic/{module}/{module}.go` | `internal/logic/brands/brands.go` |
| 缓存模块 | `internal/logic/{module}/cache.go` | `internal/logic/brands/cache.go` |
| DAO 包装 | `internal/dao/{module}.go` | `internal/dao/brands.go` |
| DAO 内部 | `internal/dao/internal/{module}.go` | `internal/dao/internal/brands.go` |
| Entity | `internal/model/entity/{module}.go` | `internal/model/entity/brands.go` |
| DO | `internal/model/do/{module}.go` | `internal/model/do/brands.go` |

## 3. API 定义规范（Req/Res + g.Meta 路由）

### 规范

- 每个模块在 `api/{module}/` 下声明接口，在 `api/{module}/v1/` 下定义 Req/Res
- 路由通过 `g.Meta` 结构体标签声明式定义：`path`、`tags`、`method`、`summary`
- `{id}` 表示路径参数
- JSON 标签使用 snake_case：`json:"first_letter"`
- 校验规则使用 `v` 标签：`v:"required|length:1,100"`
- 路径参数用 `g.Meta` 标签定义，请求体/查询参数字段用 `json` 标签

### 示例

```go
package v1

import (
    "github.com/gogf/gf/v2/frame/g"
)

// ---------- List ----------
type ListReq struct {
    g.Meta `path:"/brands" tags:"Brands" method:"get" summary:"品牌列表"`

    Page       int    `json:"page"`        // 页码，默认1
    PageSize   int    `json:"page_size"`   // 每页条数，默认20
    Name       string `json:"name"`        // 按名称模糊搜索
    FirstLetter string `json:"first_letter"` // 按首字母筛选
    Status     int    `json:"status"`      // 按状态筛选
}
type ListRes struct {
    List  []*entity.Brands `json:"list"`
    Total int              `json:"total"`
}

// ---------- Detail ----------
type DetailReq struct {
    g.Meta `path:"/brands/{id}" tags:"Brands" method:"get" summary:"品牌详情"`
    Id     int64 `json:"id"`
}
type DetailRes struct {
    *entity.Brands
}

// ---------- Create ----------
type CreateReq struct {
    g.Meta `path:"/brands" tags:"Brands" method:"post" summary:"新增品牌"`

    Name        string `json:"name"         v:"required|length:1,100" description:"品牌名称"`
    EnglishName string `json:"english_name" description:"英文名"`
    LogoUrl     string `json:"logo_url"     description:"品牌Logo"`
    FirstLetter string `json:"first_letter" v:"length:1,1" description:"首字母"`
}
type CreateRes struct {
    Id int64 `json:"id"`
}

// ---------- Update ----------
type UpdateReq struct {
    g.Meta `path:"/brands/{id}" tags:"Brands" method:"put" summary:"更新品牌"`

    Id          int64  `json:"id"           v:"required"`
    Name        string `json:"name"         v:"length:1,100" description:"品牌名称"`
    EnglishName string `json:"english_name" description:"英文名"`
}
type UpdateRes struct{}

// ---------- Delete ----------
type DeleteReq struct {
    g.Meta `path:"/brands/{id}" tags:"Brands" method:"delete" summary:"删除品牌"`
    Id     int64 `json:"id"`
}
type DeleteRes struct{}
```

### 各操作 g.Meta 路由标签速查

| 操作 | path | method | 说明 |
|------|------|--------|------|
| 列表 | `/brands` | `get` | 支持分页、筛选参数 |
| 详情 | `/brands/{id}` | `get` | 路径参数 |
| 创建 | `/brands` | `post` | 请求体 JSON |
| 更新 | `/brands/{id}` | `put` | 路径参数 + 请求体 |
| 删除 | `/brands/{id}` | `delete` | 路径参数 |

## 4. Service 接口规范

### 规范

- 接口命名：`I{Module}`（`IBrands`、`ICategories`）
- 方法签名：`method(ctx context.Context, req *v1.Request) (*v1.Response, error)`
- 使用 Service 定位器模式（Register + 全局获取）
- 每个模块三个要素：`interface`、`localXxx` 变量、`RegisterXxx()` + `Xxx()` 函数

### 模板

```go
package service

import (
    "context"
    "gf-eshop/api/{module}/v1"
)

type I{Module} interface {
    List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
    Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
    Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
    Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
    Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
}

var local{Module} I{Module}

func {Module}() I{Module} {
    if local{Module} == nil {
        panic("implement not found for interface I{Module}, forgot register?")
    }
    return local{Module}
}

func Register{Module}(i I{Module}) {
    local{Module} = i
}
```

## 5. Logic（业务逻辑）规范

### 规范

- 结构体命名：`s{Module}`（私有，`sBrands`、`sCategories`）
- 在 `init()` 中调用 `service.RegisterXxx()` 注册
- 方法签名与 Service 接口完全一致
- 查询用 `entity.X` 接收 Scan 结果
- 插入/更新用 `do.X` 传参
- 字段引用用 `dao.X.Columns().FieldName`，不用硬编码
- 分页默认值：page=1, pageSize=20
- 空结果返回 `make([]*T, 0)` 而非 `nil`
- 记录不存在返回 `gerror.NewCode(gcode.CodeNotFound, "xxx不存在")`

### 通用写法

```go
package {module}

import (
    "context"
    "github.com/gogf/gf/v2/errors/gcode"
    "github.com/gogf/gf/v2/errors/gerror"

    "gf-eshop/api/{module}/v1"
    "gf-eshop/internal/dao"
    "gf-eshop/internal/model/do"
    "gf-eshop/internal/model/entity"
    "gf-eshop/internal/service"
)

type s{Module} struct{}

func init() {
    service.Register{Module}(&s{Module}{})
}

// List 列表查询（分页 + 筛选）
func (s *s{Module}) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
    var (
        m     = dao.{Module}.Ctx(ctx)
        list  []*entity.{Module}
        total int
        page  = req.Page
        size  = req.PageSize
    )
    if page <= 0 { page = 1 }
    if size <= 0 { size = 20 }

    // 条件筛选
    if req.Name != "" {
        m = m.WhereLike(dao.{Module}.Columns().Name, "%"+req.Name+"%")
    }

    total, err = m.Count()
    if err != nil { return nil, err }
    if total == 0 {
        return &v1.ListRes{List: make([]*entity.{Module}, 0), Total: 0}, nil
    }

    err = m.Page(page, size).
        OrderAsc(dao.{Module}.Columns().SortOrder).
        OrderDesc(dao.{Module}.Columns().Id).
        Scan(&list)
    if err != nil { return nil, err }
    return &v1.ListRes{List: list, Total: total}, nil
}

// Detail 详情（走缓存）
func (s *s{Module}) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
    // 先查缓存
    cached, err := get{Module}EntityCache(ctx, req.Id)
    if err == nil && cached != nil {
        return &v1.DetailRes{{Module}: cached}, nil
    }

    var entity *entity.{Module}
    err = dao.{Module}.Ctx(ctx).
        Where(dao.{Module}.Columns().Id, req.Id).
        Scan(&entity)
    if err != nil { return nil, err }
    if entity == nil {
        return nil, gerror.NewCode(gcode.CodeNotFound, "{Module}不存在")
    }
    // 回写缓存
    set{Module}EntityCache(ctx, entity)
    return &v1.DetailRes{{Module}: entity}, nil
}

// Create 新增
func (s *s{Module}) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
    result, err := dao.{Module}.Ctx(ctx).Insert(do.{Module}{
        Name: req.Name,
        // ...
    })
    if err != nil { return nil, err }
    id, _ := result.LastInsertId()
    del{Module}AllCache(ctx) // 缓存放失效
    return &v1.CreateRes{Id: id}, nil
}

// Update 更新
func (s *s{Module}) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
    count, err := dao.{Module}.Ctx(ctx).Where(dao.{Module}.Columns().Id, req.Id).Count()
    if err != nil { return nil, err }
    if count == 0 { return nil, gerror.NewCode(gcode.CodeNotFound, "{Module}不存在") }

    _, err = dao.{Module}.Ctx(ctx).
        Data(do.{Module}{Name: req.Name}).
        Where(dao.{Module}.Columns().Id, req.Id).
        Update()
    if err != nil { return nil, err }
    del{Module}EntityCache(ctx, req.Id)
    del{Module}AllCache(ctx)
    return &v1.UpdateRes{}, nil
}

// Delete 删除（含 deleted_at 字段自动软删除）
func (s *s{Module}) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
    _, err = dao.{Module}.Ctx(ctx).
        Where(dao.{Module}.Columns().Id, req.Id).
        Delete()
    if err != nil { return nil, err }
    del{Module}EntityCache(ctx, req.Id)
    del{Module}AllCache(ctx)
    return &v1.DeleteRes{}, nil
}
```

## 6. 缓存规范（Cache-Aside）

### 规范

- 缓存文件单独放：`internal/logic/{module}/cache.go`
- 使用 GoFrame 的 `g.Redis()` 操作，通过 `SETEX` / `GET` / `DEL`
- JSON 序列化/反序列化存取
- 缓存键格式：`{module}:<id>`、`{module}:all`
- TTL 统一 10 分钟（`10 * time.Minute`）
- 提供 Warmup 函数在启动时预加载

### 三层缓存函数

```go
// 全量列表缓存
get{Module}AllCache(ctx)          // 读取
set{Module}AllCache(ctx, list)    // 写入
del{Module}AllCache(ctx)          // 失效

// 单条实体缓存
get{Module}EntityCache(ctx, id)   // 读取
set{Module}EntityCache(ctx, b)    // 写入
del{Module}EntityCache(ctx, id)   // 失效
```

### 模板

```go
package {module}

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/gogf/gf/v2/frame/g"

    "gf-eshop/internal/dao"
    "gf-eshop/internal/model/entity"
)

const (
    {module}AllTTL    = 10 * time.Minute
    {module}EntityTTL = 10 * time.Minute
)

func cacheKey{Module}All() string       { return "{module}:all" }
func cacheKey{Module}(id int64) string  { return fmt.Sprintf("{module}:%d", id) }

// Warmup 启动预热：加载全量数据到 Redis
func Warmup(ctx context.Context) {
    var list []*entity.{Module}
    if err := dao.{Module}.Ctx(ctx).OrderAsc(dao.{Module}.Columns().SortOrder).Scan(&list); err != nil {
        g.Log().Warningf(ctx, "{module} cache warmup query failed: %v", err)
        return
    }
    if len(list) == 0 { return }
    data, _ := json.Marshal(list)
    g.Redis().Do(ctx, "SETEX", cacheKey{Module}All(), int({module}AllTTL.Seconds()), string(data))
    // 顺便预热单条缓存
    for _, item := range list {
        b, _ := json.Marshal(item)
        g.Redis().Do(ctx, "SETEX", cacheKey{Module}(item.Id), int({module}EntityTTL.Seconds()), string(b))
    }
    g.Log().Infof(ctx, "{module} cache warmed up: %d items", len(list))
}

func get{Module}AllCache(ctx context.Context) ([]*entity.{Module}, error) {
    v, err := g.Redis().Do(ctx, "GET", cacheKey{Module}All())
    if err != nil || v.IsNil() { return nil, err }
    var list []*entity.{Module}
    if err := json.Unmarshal(v.Bytes(), &list); err != nil { return nil, err }
    return list, nil
}

func set{Module}AllCache(ctx context.Context, list []*entity.{Module}) error {
    data, _ := json.Marshal(list)
    _, err := g.Redis().Do(ctx, "SETEX", cacheKey{Module}All(), int({module}AllTTL.Seconds()), string(data))
    return err
}

func del{Module}AllCache(ctx context.Context) {
    g.Redis().Do(ctx, "DEL", cacheKey{Module}All())
}

func get{Module}EntityCache(ctx context.Context, id int64) (*entity.{Module}, error) {
    v, err := g.Redis().Do(ctx, "GET", cacheKey{Module}(id))
    if err != nil || v.IsNil() { return nil, err }
    var item entity.{Module}
    if err := json.Unmarshal(v.Bytes(), &item); err != nil { return nil, err }
    return &item, nil
}

func set{Module}EntityCache(ctx context.Context, item *entity.{Module}) error {
    data, _ := json.Marshal(item)
    _, err := g.Redis().Do(ctx, "SETEX", cacheKey{Module}(item.Id), int({module}EntityTTL.Seconds()), string(data))
    return err
}

func del{Module}EntityCache(ctx context.Context, id int64) {
    g.Redis().Do(ctx, "DEL", cacheKey{Module}(id))
}
```

### 缓存失效时机

| 操作 | 失效的缓存 |
|------|-----------|
| Create | `del{Module}AllCache` |
| Update | `del{Module}EntityCache(id)` + `del{Module}AllCache` |
| Delete | `del{Module}EntityCache(id)` + `del{Module}AllCache` |

## 7. Controller 规范

### 规范

- 结构体命名：`ControllerV1`
- 每个方法一行，直接委托给 `service.Xxx()` 调用
- 不做任何业务处理，只做请求转发

### 模板

```go
// {module}_new.go
package {module}

import "gf-eshop/api/{module}"

type ControllerV1 struct{}

func NewV1() {module}.I{Module}V1 {
    return &ControllerV1{}
}
```

```go
// {module}_v1_{module}.go
package {module}

import (
    "context"
    "gf-eshop/api/{module}/v1"
    "gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
    return service.{Module}().List(ctx, req)
}
// ...Detail, Create, Update, Delete 同理
```

## 8. 路由注册规范

### 规范

- 在 `internal/cmd/cmd.go` 的 `Main.Func` 中注册
- API 路由注册在 `/api/v1` 分组下
- 所有路由使用 `ghttp.MiddlewareHandlerResponse` 中间件
- 缓存预热放在路由注册之前

### 模板

```go
s.Group("/api/v1", func(group *ghttp.RouterGroup) {
    group.Middleware(ghttp.MiddlewareHandlerResponse)
    group.Bind(
        brands.NewV1(),
        categories.NewV1(),
        categoryBrands.NewV1(),
    )
})
```

## 9. 导入规范

标准库 → 第三方库 → 内部包，三组间空行分隔：

```go
import (
    "context"
    "fmt"

    "github.com/gogf/gf/v2/errors/gcode"
    "github.com/gogf/gf/v2/errors/gerror"
    "github.com/gogf/gf/v2/frame/g"

    "gf-eshop/api/brands/v1"
    "gf-eshop/internal/dao"
    "gf-eshop/internal/model/do"
)
```

## 10. 错误处理规范

| 场景 | 写法 |
|------|------|
| 记录不存在 | `gerror.NewCode(gcode.CodeNotFound, "品牌不存在")` |
| 参数非法 | `gerror.NewCode(gcode.CodeInvalidParameter, "xxx不能为空")` |
| 数据库错误 | 直接返回 `err`，不包装 |
| 缓存操作失败 | `g.Log().Warning(ctx, "xxx failed: %v", err)` — 不阻断请求 |

## 11. 事务规范

使用 `dao.X.Transaction()` 封装多表写入：

```go
err = dao.CategoryBrands.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    // 先删旧关联
    dao.CategoryBrands.Ctx(ctx).TX(tx).Where(...).Delete()
    // 再插新关联
    tx.Model("sp_category_brands").Insert(...)
    return nil
})
```

## 12. 生成代码命令规范

```bash
# 首次或改表后：从数据库表生成 DAO/Entity/DO
# 1. 在 hack/config.yaml 的 tables 字段加上新表名
# 2. 运行：
gf gen dao

# 重新生成 Service 接口（当 API Req/Res 改变时）
gf gen service

# 重新生成 Controller（当 API 接口定义改变时）
gf gen ctrl
```

## 检查清单

### 目录结构
- [ ] `api/{module}/{module}.go` 接口声明是否存在
- [ ] `api/{module}/v1/{module}.go` Req/Res 是否定义
- [ ] `internal/service/{module}.go` Service 接口是否定义
- [ ] `internal/logic/{module}/{module}.go` 业务逻辑是否存在
- [ ] `internal/controller/{module}/` 控制器是否存在
- [ ] `internal/dao/{module}.go` 是否已通过 `gf gen dao` 生成

### API 定义
- [ ] 路由是否使用 `g.Meta` 标签声明，未在代码中手写 URL
- [ ] `path` 是否以 `/` 开头
- [ ] `tags` 是否使用首字母大写英文（`Brands`、`Categories`）
- [ ] `summary` 是否使用中文
- [ ] `method` 是否符合 RESTful（get/post/put/delete）
- [ ] 路径参数是否使用 `{id}` 格式
- [ ] 入参校验 `v` 标签是否完整（`required`、`length` 等）
- [ ] JSON 标签是否 snake_case

### Logic（业务逻辑）
- [ ] 是否在 `init()` 中调用 `service.RegisterXxx()`
- [ ] 查询是否使用 `entity.X` 作为 Scan 类型
- [ ] 插入/更新是否使用 `do.X` 传参
- [ ] 字段引用是否使用 `dao.X.Columns().FieldName`
- [ ] 是否处理了 page=0 和 pageSize=0 的默认值
- [ ] 空结果是否返回 `make([]*T, 0)` 而非 `nil`
- [ ] 记录不存在是否返回 `gerror.NewCode(gcode.CodeNotFound, ...)`
- [ ] 新增/更新/删除后是否正确失效缓存

### 缓存
- [ ] cache.go 中是否提供 Warmup 函数
- [ ] Warmup 是否在 cmd.go 的路由注册前调用
- [ ] 读取时是否先查缓存 → 未命中再查 DB → 回写缓存
- [ ] 写入操作是否失效了对应缓存（单条 + 全量列表）
- [ ] 缓存键格式是否统一（`{module}:{id}`、`{module}:all`）
- [ ] TTL 是否统一为 10 分钟

### Controller
- [ ] 是否每个方法仅一行委派给 `service.Xxx()`
- [ ] 是否不包含任何业务处理逻辑

### main.go 入口
- [ ] 是否导入了 MySQL 驱动 `_ "github.com/gogf/gf/contrib/drivers/mysql/v2"`
- [ ] 是否导入了 Redis 驱动 `_ "github.com/gogf/gf/contrib/nosql/redis/v2"`
- [ ] 是否导入了所有 logic 包（`_ "gf-eshop/internal/logic/{module}"`）
