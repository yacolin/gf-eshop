# GoFrame CRUD 开发流程

> 本文档以 `sp_brands` 表的 CRUD 为例，梳理哪些代码是**命令自动生成**的，哪些需要**手写**。

---

## 全景图

```
┌────────────────────────────────────────────────────────────┐
│                  手写代码（你的工作）                       │
│                                                            │
│  hack/config.yaml       ← 告诉 gf CLI 连哪个库、哪张表    │
│  api/brands/v1/*.go     ← 定义 API 请求/响应结构体        │
│  internal/service/*.go  ← 定义 Service 接口                │
│  internal/logic/*/*.go  ← 业务逻辑实现（CRUD 怎么写）     │
│  internal/controller/*/*.go ← 控制器（桥接 API → Service）│
│  internal/cmd/cmd.go    ← 注册路由                         │
│  main.go                ← 导入 MySQL 驱动 + logic 包      │
├────────────────────────────────────────────────────────────┤
│               命令生成（gf gen dao）                         │
│                                                            │
│  internal/dao/internal/*.go  ← DAO 底层实现（勿改）        │
│  internal/dao/*.go           ← DAO 包装（可扩展）         │
│  internal/model/entity/*.go  ← 实体结构体                  │
│  internal/model/do/*.go      ← DO 结构体（插入/更新用）   │
└────────────────────────────────────────────────────────────┘
```

---

## 一、命令生成部分

### 1. 告诉 gf CLI 连哪个数据库

**文件：`hack/config.yaml`**（只需写一次）

```yaml
gfcli:
  gen:
    dao:
      - link: "mysql:root:123456@tcp(localhost:3306)/eshop_db?..."
        tables: "sp_brands"          # 指定表名，多表用逗号分隔
        removePrefix: "sp_"           # 去掉表前缀，代码中变 brands
        descriptionTag: true          # 字段注释转 json tag
        noModelComment: true
```

> ⚠️ 注意：`gf gen dao` 的 `link` 格式是 `mysql:user:pass@tcp(...)`，**没有 `//`**。
> 运行时 `manifest/config/config.yaml` 里的格式才是 `mysql://user:pass@tcp(...)`（有 `//`）。

### 2. 运行命令

```bash
gf gen dao
```

这一条命令会生成 **4 个文件**：

| 生成文件 | 用途 | 能不能改 |
|----------|------|---------|
| `internal/dao/internal/brands.go` | DAO 底层：表名、字段名常量、Ctx()、Transaction() 等基础方法 | **❌ 不能改**（clI 重新生成会覆盖） |
| `internal/dao/brands.go` | DAO 包装：包了一个内部 DAO，给你留了扩展空间 | ✅ 可以加自定义查询方法 |
| `internal/model/entity/brands.go` | 实体结构体：一条记录的完整字段，用于查询返回 | **❌ 不能改** |
| `internal/model/do/brands.go` | DO 结构体：字段全是 `interface{}`，用于插入/更新/条件查询 | **❌ 不能改** |

这些文件和表结构一一对应。如果表改了字段，重新跑一次 `gf gen dao` 就能同步更新。

---

## 二、手写代码部分

### 1. API 接口定义（约 50 行）

```
api/
  brands/
    brands.go          ← 接口定义（IBrandsV1）
    v1/brands.go       ← Req / Res 结构体
```

**`api/brands/brands.go`** — 声明有哪些接口方法：

```go
type IBrandsV1 interface {
    List(ctx, *v1.ListReq) (*v1.ListRes, error)
    Detail(ctx, *v1.DetailReq) (*v1.DetailRes, error)
    Create(ctx, *v1.CreateReq) (*v1.CreateRes, error)
    Update(ctx, *v1.UpdateReq) (*v1.UpdateRes, error)
    Delete(ctx, *v1.DeleteReq) (*v1.DeleteRes, error)
}
```

**`api/brands/v1/brands.go`** — 每个接口的请求/响应结构体，核心是 `g.Meta` 标签定义路由：

```go
type ListReq struct {
    g.Meta `path:"/brands" tags:"Brands" method:"get" summary:"品牌列表"`
    Name string `json:"name"`
    // ...
}
```

> 关键点：`path` 决定 URL 路径，`method` 决定 HTTP 方法，
> `{id}` 表示路径参数，`v:"required"` 是校验规则。

### 2. Service 接口（约 30 行）

**`internal/service/brands.go`**

声明 Service 层的接口 + 注册/获取的模板代码：

```go
type IBrands interface {
    List(ctx, *v1.ListReq) (*v1.ListRes, error)
    // ...
}

var localBrands IBrands

func Brands() IBrands {        // 获取已注册的实现
    if localBrands == nil {
        panic("implement not found...")
    }
    return localBrands
}

func RegisterBrands(i IBrands) {  // logic 包的 init() 调用此函数
    localBrands = i
}
```

> 这个文件是标准的"服务定位器"模式，每张表一个接口。GoFrame 的 `gf gen service` 可以自动生成这个文件（后续再优化）。

### 3. 业务逻辑（最核心，约 100 行）

**`internal/logic/brands/brands.go`**

这里是 **真正的 CRUD 代码**，需要手写：

```go
func init() {
    service.RegisterBrands(&sBrands{})  // 注册到 service
}

func (s *sBrands) List(ctx, req) (res, err) {
    // 构建查询、分页、筛选
    m = dao.Brands.Ctx(ctx)
    m = m.WhereLike(...)
    total, _ = m.Count()
    m.Page(page, size).Scan(&list)
}

func (s *sBrands) Create(ctx, req) (res, err) {
    dao.Brands.Ctx(ctx).Insert(do.Brands{...})
}

func (s *sBrands) Update(ctx, req) (res, err) {
    dao.Brands.Ctx(ctx).Data(do.Brands{...}).Where(id).Update()
}

func (s *sBrands) Delete(ctx, req) (res, err) {
    dao.Brands.Ctx(ctx).Where(id).Delete()
    // 因为表有 deleted_at 字段，自动变软删除（UPDATE 而非 DELETE）
}
```

> **初学者常见问题：**
> - 查询用 `entity.Brands` 接收结果
> - 插入/更新用 `do.Brands` 传参
> - `dao.Brands.Columns().Name` 取字段名，比硬编码 `"name"` 更安全

### 4. 控制器（约 40 行）

**`internal/controller/brands/brands_new.go`**

```go
type ControllerV1 struct{}

func NewV1() brands.IBrandsV1 {
    return &ControllerV1{}
}
```

**`internal/controller/brands/brands_v1_brands.go`**

```go
func (c *ControllerV1) List(ctx, req) (res, err) {
    return service.Brands().List(ctx, req)   // 直接委托给 service
}
// ... 其他方法同理
```

> 控制器就是很薄的一层，每个方法一行，把 API 请求转给 Service 处理。

### 5. 注册路由（修改 2 行）

**`internal/cmd/cmd.go`**

```go
import "gf-eshop/internal/controller/brands"

s.Group("/api", func(group *ghttp.RouterGroup) {
    group.Middleware(ghttp.MiddlewareHandlerResponse)
    group.Bind(
        brands.NewV1(),   // ← 加这一行
    )
})
```

### 6. 注册 logic 包导入（修改 1 行）

**`main.go`**

```go
import (
    _ "github.com/gogf/gf/contrib/drivers/mysql/v2"  // ← MySQL 驱动
    _ "gf-eshop/internal/logic/brands"                 // ← 触发 logic 的 init()
)
```

> 为什么要加 `_` 导入？因为 `logic/brands` 包的 `init()` 会调用 `service.RegisterBrands()`，
> 不加这个导入，init() 不会执行，运行时就会 panic。

---

## 三、小结：下次你做一张新表的 CRUD

### 你需要手写

| 文件 | 行数 | 难度 |
|------|------|------|
| `api/xxx/v1/xxx.go` | 40-60 | ⭐ 定义 Req/Res 结构体 |
| `internal/service/xxx.go` | 25-35 | ⭐ 抄模板改方法名 |
| `internal/logic/xxx/xxx.go` | 80-120 | ⭐⭐⭐ 真正的业务逻辑 |
| `internal/controller/xxx/xxx_new.go` | 10 | ⭐ 构造函数 |
| `internal/controller/xxx/xxx_v1_xxx.go` | 30 | ⭐ 每个方法调一次 service |
| `internal/cmd/cmd.go` | +1行 | ⭐ 加一行 Bind |
| `main.go` | +1行 | ⭐ 加一行 import |

### 命令生成

| 步骤 | 文件 | 次数 |
|------|------|------|
| `hack/config.yaml` 加表名 | - | **首次**配一次数据库连接 |
| `gf gen dao` | 4 个文件（dao/internal, dao, entity, do） | **每次表结构变了**跑一次 |

### 标准操作流程

```bash
# 1. hack/config.yaml 的 tables 字段加上新表名
# 2. 生成 DAO / Model
gf gen dao

# 3. 手写 API 定义
touch api/xxx/v1/xxx.go

# 4. 手写 Service 接口
touch internal/service/xxx.go

# 5. 手写业务逻辑
touch internal/logic/xxx/xxx.go

# 6. 手写控制器
touch internal/controller/xxx/xxx_new.go
touch internal/controller/xxx/xxx_v1_xxx.go

# 7. 注册路由 + 导入 logic
vim internal/cmd/cmd.go
vim main.go

# 8. 编译验证
go build -o main .
```

---

## 四、FAQ

### Q: 改表结构了怎么办？

```bash
# 1. 执行 ALTER TABLE
# 2. 重新生成
gf gen dao
# 3. 检查 entity 和 do 是否同步了 ✅
```

### Q: dao/brands.go 可以加自定义方法吗？

可以。`dao/brands.go` 是"一次生成、后续可改"的文件，你可以加自定义查询：

```go
func (dao *brandsDao) FindActive(ctx context.Context) ([]*entity.Brands, error) {
    return dao.Ctx(ctx).Where("status", 1).All()
}
```

### Q: 如果不想用 service 注册模式？

可以跳过 `internal/service/`，让 controller 直接调用 logic：

```go
// controller 直接 new logic 结构体
func (c *ControllerV1) List(ctx, req) (res, err) {
    return (&sBrands{}).List(ctx, req)
}
```

但推荐注册模式——解耦、方便单元测试。
