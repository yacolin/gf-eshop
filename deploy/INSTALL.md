# 部署说明（linux/amd64）

面向 `gf-eshop`（GoFrame 电商后端）。流程参考同目录体系的 `campus_express`：
**本地交叉编译 → 打包 → scp → 服务器 systemd 托管**。

线上不做现场编译：小内存机器上 `go build` 的链接阶段内存峰值 1G+，容易被 OOM killer 干掉；
本地编译成静态二进制（`CGO_ENABLED=0`，不依赖目标机 glibc 版本）再上传。

---

## 0. 服务器需要什么

| 依赖 | 必需 | 说明 |
|------|------|------|
| MySQL（**云 RDS**，不在本机） | ✅ | 业务库 `eshop_db`。用 RDS 内网地址连接，本机不装 MySQL（建库与白名单见 §4） |
| Redis | ✅ | 缓存、列表 ZSET、验证码风控、refresh token 白名单 |
| Elasticsearch 7.17.x | ⭕ 可选 | 组合筛选/文本检索的加速器。**必须装同版本 IK 分词插件**（见下）。未安装或内存紧张时设 `elasticsearch.enabled: false`，商品/品牌的筛选查询会自动回落 MySQL（这是设计好的降级，不是故障） |
| 开放端口 | — | 默认监听 `:8000`。直接对外就放行 8000；前面挂 Nginx 则只对内网开放 |

机器规格参考：MySQL 在 RDS 上，本机只跑应用 + Redis，**2G 内存相当宽裕**（实测应用 36MB）。
要再加 Elasticsearch 建议升到 4G（见 §5.1 的说明）。

### ES 安装要点（踩过坑，务必看）

项目索引用的分词器是 `ik_max_word`（索引侧）/ `ik_smart`（检索侧），来自
[medcl/elasticsearch-analysis-ik](https://github.com/medcl/elasticsearch-analysis-ik)。
**不装这个插件，索引创建就会失败**，启动日志里只有 `warmup stage "brands_es" failed` 之类的 WARN，
业务会自动回落 MySQL —— 看起来"能跑"，但 ES 其实完全没生效。插件版本必须与 ES 版本**完全一致**：

```bash
# 以 7.17.4 为例（本机开发环境就是这个版本 + analysis-ik 7.17.4）
sudo /usr/share/elasticsearch/bin/elasticsearch-plugin install \
  https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v7.17.4/elasticsearch-analysis-ik-7.17.4.zip
sudo systemctl restart elasticsearch

# 验证插件
curl -s localhost:9200/_cat/plugins
# 验证分词器可用（返回 tokens 即为正常）
curl -s -X POST "localhost:9200/_analyze" -H 'Content-Type: application/json' \
  -d '{"analyzer":"ik_max_word","text":"苹果手机"}'
```

另外两点：

- **单节点集群显示 yellow 是正常的**：项目没有显式设置副本数，ES 默认 1 副本，
  单节点上副本分片无处分配，于是索引是 `yellow`（主分片正常、搜索可用）。
  想变绿可以在服务器上把副本设为 0：
  `curl -X PUT "localhost:9200/eshop_*/_settings" -H 'Content-Type: application/json' -d '{"index":{"number_of_replicas":0}}'`
- **ES 版本别乱升**：Go 客户端是 `go-elasticsearch/v7`，请保持 7.17.x；升到 8.x 需要同时升级客户端并改造索引 API。


---

## 1. 本地打包

```bash
cd gf-eshop

# 首次部署：导出本地库结构（只导结构不导数据），随包一起带去服务器
make db-schema

# 交叉编译 + 打包（默认 linux/amd64）
make release
# ARM 服务器（如部分云主机/树莓派）：
make release GOARCH=arm64
```

产物 `bin/gf-eshop-linux-amd64.tar.gz`，内含：

```
gf-eshop                        静态二进制
manifest/config/config.yaml     配置（权限已置 600，内含 DB 密码与邮箱授权码）
db/schema.sql                   建库脚本（执行过 make db-schema 才有）
start.sh                        启动脚本（切工作目录 + 配置体检）
gf-eshop.service                systemd 单元
INSTALL.md                      本文件
```

`make release` 会打印 `file` 结果（确认是 ELF x86-64）、包大小、sha256 和 git revision。

> ⚠️ **包里带着你的密钥**（MySQL 密码、邮箱 SMTP 授权码）。`bin/` 已在 `.gitignore` 里，
> 不会进 git；但别把 tarball 传到公开位置，也别随手发给别人。

## 2. 上传

```bash
make upload DEPLOY_HOST=root@<你的服务器IP>
# 一步到位：make deploy DEPLOY_HOST=root@<你的服务器IP>
# 需要换落地目录：make upload DEPLOY_HOST=... DEPLOY_DIR=/data
```

## 3. 服务器上解包

```bash
sudo tar -xzf /tmp/gf-eshop-linux-amd64.tar.gz -C /opt
sudo ls -l /opt/gf-eshop
```

首次部署建账号（二选一，Ubuntu/Debian 惯例）：

```bash
sudo useradd -r -s /sbin/nologin gfeshop
sudo chown -R gfeshop:gfeshop /opt/gf-eshop
```

## 4. 首次建库（MySQL 在 RDS 上）

本项目的 MySQL 跑在**云 RDS**，不在应用服务器上。所以在 RDS 上建好库，再让应用连过去。

```bash
# 把 <RDS内网地址> / <账号> 换成你 RDS 控制台上的「内网地址」与高权限账号
RDS=<RDS内网地址>

mysql -h $RDS -P 3306 -u <高权限账号> -p \
  -e "CREATE DATABASE IF NOT EXISTS eshop_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# 导入结构（schema.sql 由 make db-schema 生成，随包带过来；只含表结构不含数据）
mysql -h $RDS -P 3306 -u <高权限账号> -p eshop_db < /opt/gf-eshop/db/schema.sql
```

> ⚠️ **`schema.sql` 含 74 条 `DROP TABLE IF EXISTS`**：它是「重建结构」语义，
> 在**已有数据**的库上执行会先删表再重建 = 清空数据，只用于首次建库。
> 好消息是它已刻意避开云 RDS 导入最常见的三个拦路虎 —— `DEFINER`、`SQL_LOG_BIN`、`GTID_PURGED`
> （这三个都需要 SUPER 权限，RDS 不给）。已实测本包的 schema.sql 中这三者均为 0 条，
> 74 张表全部 `ENGINE=InnoDB` + `utf8mb4`，可直接导入。

**RDS 侧的三件事（顺序别错，否则应用连不上）**：

1. **白名单 / 安全组**：把应用服务器的**内网 IP** 加进 RDS 白名单（或与应用同 VPC + 同一安全组）。
   这一步漏了，应用日志只会看到 `connection refused` / `i/o timeout`，很容易误判成代码问题。
2. **用内网地址，不要用公网地址**：内网延迟低、免流量费，也不暴露在公网。应用与 RDS 在同一地域时
   `link` 里就写内网 endpoint。
3. **建一个专用账号给应用**（不要用 RDS 主账号）：给 `eshop_db.*` 的 SELECT/INSERT/UPDATE/DELETE 权限即可，
   DDL 只在你导入 schema 时用高权限账号。

```sql
-- 在 RDS 上执行
CREATE USER 'eshop'@'%' IDENTIFIED BY '<强密码>';
GRANT SELECT, INSERT, UPDATE, DELETE ON eshop_db.* TO 'eshop'@'%';
FLUSH PRIVILEGES;
```

> 想连本地开发数据一起带过去：`mysqldump -h127.0.0.1 -uroot -p123456 eshop_db | mysql -h $RDS -u<高权限账号> -p eshop_db`
>
> **时区**：`link` 里用了 `parseTime=True&loc=Local`，RDS 的 `time_zone` 必须与应用服务器一致
> （建议都把服务器 `timedatectl set-timezone Asia/Shanghai`，RDS 参数组里设 `time_zone = +08:00`）。
> 不一致会让 `created_at` 之类由 MySQL 生成的时间整体偏移 8 小时。

## 5. 改配置（必做）

### 5.1 只有 2G 内存（且 MySQL 在 RDS）？先看这份清单

因为 **MySQL 不在本机**，2G 机器其实相当宽裕 —— 本机只需要跑应用和 Redis：

| 组件 | 预计占用 | 说明 |
|------|----------|------|
| gf-eshop 应用 | **实测 36MB**（峰值同为 36MB） | 已含全部缓存预热；单元里给 512M 上限纯粹是安全网 |
| Redis | ~100MB | 数据集很小（缓存 + 验证码 + refresh 白名单） |
| 系统 / sshd / journald | ~250MB | |
| MySQL | 0（在 RDS） | 内存由 RDS 规格承担，本机不用管 |
| **合计** | **< 400MB** | 2G 有大量余量 |

必做的两件事：

```yaml
# 1) /opt/gf-eshop/manifest/config/config.yaml
elasticsearch:
  enabled: false        # 见下方说明；2G 上装 ES 会把它自己吃到 swap
database:
  default:
    # 换成 RDS 内网地址；账号用上一步建的专用账号
    link: "mysql:eshop:<强密码>@tcp(<RDS内网地址>:3306)/eshop_db?charset=utf8mb4&parseTime=True&loc=Local"
    maxActive: 50       # 建议 ≤ RDS 规格的 max_connections 留出余量（默认 100 偏大）
    maxIdle: 10
    idleTimeout: "60s"
```

```bash
# 2) Redis 加个内存上限（数据集小，256M 足够；allkeys-lru 避免写满后写失败）
redis-cli CONFIG SET maxmemory 256mb
redis-cli CONFIG SET maxmemory-policy allkeys-lru
# 写进 /etc/redis/redis.conf 才能重启后保留

# 建议再挂 2G swap 兜底（有 swap 总比被 OOM killer 杀掉强）
sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile && sudo mkswap /swapfile && sudo swapon /swapfile
```

**关于 ES：2G 上仍然不建议装。** 因为 MySQL 挪走了，机器确实空出 1G 左右，理论上可以给 ES 配
512M 堆跑起来；但 ES 除堆之外还有 Lucene 段缓存与 mmap、段合并时的瞬时开销，2G 上基本没有
抗突发余量。你现在(101 品牌 / 2025 商品)走 MySQL 完全够用，等数据量或搜索需求上来**直接升到 4G 再加 ES**更稳。

ES 关掉后，启动日志里会有这几行 **WARN，属预期降级，不用管**：

```
warmup stage "brands_es" failed: elasticsearch unavailable
warmup stage "products_es" failed: elasticsearch unavailable
```

> 以后想加 ES：装好 ES 7.17.x + IK 插件，把 `enabled` 改回 `true` 重启即可 ——
> 索引会由启动对账自愈自动重建（100 个品牌 / 2025 个商品，几秒就完），无需手动 reindex。

### 5.2 配置逐项

```bash
sudo vi /opt/gf-eshop/manifest/config/config.yaml
sudo chmod 600 /opt/gf-eshop/manifest/config/config.yaml   # 含密钥，别留 644
```

要点：

| 配置段 | 要改什么 |
|--------|----------|
| `database.default.link` | **RDS 内网地址**、专用账号、密码、库名；`maxActive` 建议 ≤ RDS 的 `max_connections`（默认 100 偏大） |
| `redis.default.address` | Redis 地址（默认 `127.0.0.1:6379`）；有密码要加 `pass` |
| `elasticsearch` | 装了 ES 才保持 `enabled: true`；**没装或内存紧张就改 `false`** |
| `jwt.secret` | **务必换掉**，别沿用仓库里的默认值 |
| `email.*` | 邮箱验证码用。不使用时 `enabled: false` 即可（接口返回 1021，不影响其它功能） |
| `server.address` | 默认 `:8000`。前面挂 Nginx 可改成 `127.0.0.1:8000` 只监听本机 |

> 配置是**按工作目录**找的（`manifest/config/config.yaml`），所以必须用 `start.sh` 启动，
> 或保证 systemd 的 `WorkingDirectory=/opt/gf-eshop`。

### 5.3 推荐：把线上配置放 `/etc/gf-eshop.env`（而不是改 yaml）

> **为什么**：发布包里的 `manifest/config/config.yaml` 是**打包机上的本地文件**
> （含开发库密码、邮箱授权码），而且**每次解包都会覆盖服务器上那一份**。
> 所以线上的真实值应该只放在 `/etc/gf-eshop.env`（600），
> 由 `start.sh` 与 systemd 注入进程环境，**运行期覆盖** yaml 的同名配置 —— 解包不再影响它。

```bash
sudo cp /opt/gf-eshop/env.example /etc/gf-eshop.env
sudo chmod 600 /etc/gf-eshop.env
sudo vi /etc/gf-eshop.env          # 填 DB 密码、JWT 密钥等
sudo systemctl restart gf-eshop
```

命名规则：**配置键 → 大写蛇形，`.` 与驼峰边界都变成 `_`**，大小写不敏感：

| 配置键（yaml） | 环境变量 |
|---|---|
| `server.address` | `SERVER_ADDRESS` |
| `database.default.link` | `DATABASE_DEFAULT_LINK`（整条 DSN，别省 `parseTime=True`） |
| `jwt.secret` | `JWT_SECRET` |
| `elasticsearch.enabled` | `ELASTICSEARCH_ENABLED` |
| `elasticsearch.addresses` | `ELASTICSEARCH_ADDRESSES`（逗号分隔） |
| `orderShard.mode` | `ORDER_SHARD_MODE` |
| `orderShard.defaultWindowMonths` | `ORDER_SHARD_DEFAULT_WINDOW_MONTHS` |

语义与坑：

- **注释掉 = 用 yaml 里的默认值**（不是「关」）；设置成**空串也算设置**（可用来显式清空）；
- 可覆盖项是**显式清单**（`internal/cmd/envconfig.go` 的 `EnvOverridableKeys`）——
  不是任意环境变量都能改配置；清单与 `env.example` 的一致性有单测兜着；
- 语法：`KEY=value`，**不要写 `export`**，**注释必须单独占一行**
  （systemd 不认行内 `#`，会把注释当成值的一部分）；
- 启动日志会打印「环境变量覆盖了 N 项配置：…」（只列键名，**不打值**，避免泄密）；
- 没配 env 文件也能跑：全用 yaml 的值。

## 6. 装 systemd 并启动

```bash
sudo cp /opt/gf-eshop/gf-eshop.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now gf-eshop
systemctl status gf-eshop --no-pager
journalctl -u gf-eshop -n 50 --no-pager     # 启动自检都在这里
```

启动日志里应能看到：路由表、`brand cache warmed up`、各 `Warmup stage` 结果，
以及验证码渠道的就绪情况（例如 `验证码渠道「短信」未就绪…sms.enabled 为 false`，属正常）。
**没有 ES 时会看到 `warmup stage "brands_es" failed: elasticsearch unavailable` 之类的 WARN —— 这是预期降级，不是错误。**

## 7. 验证

```bash
curl -s localhost:8000/hello
curl -s "localhost:8000/api/v1/brands?page=1&page_size=5" | head -c 300
curl -s localhost:8000/api.json | head -c 200     # OpenAPI
# 浏览器：http://<IP>:8000/swagger
```

## 8. 更新（改完代码后）

```bash
# 本地
make release && make upload DEPLOY_HOST=root@<IP>
# 服务器
sudo tar -xzf /tmp/gf-eshop-linux-amd64.tar.gz -C /opt      # 覆盖二进制与配置模板
sudo chmod 600 /opt/gf-eshop/manifest/config/config.yaml    # 解包会重置权限
sudo systemctl restart gf-eshop
```

> ⚠️ 解包会**覆盖** `manifest/config/config.yaml`。因此**不要把线上配置改在 yaml 里** ——
> 按 §5.3 放到 `/etc/gf-eshop.env`，解包就不影响它（env 在运行期覆盖 yaml）。
> 如果你确实改过 yaml，解包前先备份：`sudo cp .../config.yaml /root/gf-eshop-config.bak`。
> 新版本可能新增可覆盖项，用 `env.example` 对照补齐（缺的项回落到 yaml 默认值）。

ES 索引无需手动重建：启动时的对账自愈会按「结构版本 / 文档数」自动重建
（改过 mapping 会因 schema 版本变化触发，见 `internal/logic/brands/es.go`）。
确实需要手动全量重建时：`sudo -u gfeshop /opt/gf-eshop/gf-eshop reindex --entity=brands`。

## 9. 注意事项

- **内存**：`gf-eshop.service` 里默认 `MemoryMax=1G / MemoryHigh=768M`。机器更大可以调高；
  1~2G 的机器建议关掉 ES，并把 `MemoryMax` 调到 512M，同时给 MySQL 留出内存。
- **时区**：容器/服务器若是 UTC，`created_at`（`datetime(3)`，由 MySQL 生成）与 Go 侧
  `time.Now()` 会不一致。建议 `timedatectl set-timezone Asia/Shanghai`，并保持
  `database.link` 里的 `loc=Local`。
- **防火墙/安全组**：只用 8000 直连就放行 8000；更推荐 Nginx 反代 + HTTPS，
  并把 `server.address` 改成 `127.0.0.1:8000`。
- **WebSocket**：`/api/v1/ws` 走 Nginx 反代时需要 `Upgrade`/`Connection` 头透传。
- **验证码邮件**：`email.dev_mode: true` 时验证码只写日志不外发，**生产必须 false**；
  用个人邮箱发信有日限额（新浪/QQ 都有），量大要换带 SPF/DKIM 的事务邮件服务。
- **备份**：至少定期 `mysqldump eshop_db`；Redis 里只有缓存与短周期凭证（refresh token 白名单、
  验证码），丢失只影响登录态，不需要备份。

---

## 附 A：在 2G 机器上试装 ES（实验性）

> 这是**实验**，不是推荐配置。因为 MySQL 在 RDS 上，本机只剩应用（实测 36MB）+ Redis（~100MB）
> 可用内存约 1.6G，ES 配 512M 堆在**当前数据量**（101 品牌 / 2025 商品）下是有可能撑住的。
> 但「能启动」≠「能用」：内存压力大时 JVM 堆被换到 swap 会让 GC 停顿到秒级，
> 应用侧 3s 超时一到就**静默降级回 MySQL**（日志里只有一行 WARN），表面看服务还活着。
> 所以必须按下面的判据观测，而不是看 ES 有没有起来。

### A.1 内核参数（ES 必需的，漏了起不来）

```bash
# ES 要求 vm.max_map_count >= 262144，否则启动直接失败：
#   max virtual memory areas vm.max_map_count [65530] is too low
sudo sysctl -w vm.max_map_count=262144
# 内存不足时优先丢 page cache，而不是把 ES 堆换出去（GC 停顿才是致命的）
sudo sysctl -w vm.swappiness=1
# Redis 做 BGSAVE 要 fork，建议开；对 ES 也无害
sudo sysctl -w vm.overcommit_memory=1

# 持久化
printf 'vm.max_map_count=262144\nvm.swappiness=1\nvm.overcommit_memory=1\n' | sudo tee /etc/sysctl.d/99-gf-eshop.conf
```

### A.2 安装 ES 7.17.x + IK（版本必须一致）

```bash
# Debian/Ubuntu：加 Elastic 官方 7.x 源后安装
curl -fsSL https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elastic.gpg
echo "deb [signed-by=/usr/share/keyrings/elastic.gpg] https://artifacts.elastic.co/packages/7.x/apt stable main" \
  | sudo tee /etc/apt/sources.list.d/elastic-7.x.list
sudo apt update && sudo apt install -y elasticsearch=7.17.4
# 若源里已经没有 7.17.4，就装 7.17.x 里的最新补丁 —— 但 IK 插件版本必须与 ES 版本**完全一致**

# IK 分词插件（不装则索引创建失败，业务静默回落 MySQL —— 见 §0「ES 安装要点」）
sudo /usr/share/elasticsearch/bin/elasticsearch-plugin install \
  https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v7.17.4/elasticsearch-analysis-ik-7.17.4.zip
```

### A.3 小内存配置

```bash
# 堆固定 512M：2G 机器上不要超过 512M（JVM 堆 + Lucene 段缓存 + off-heap 合计约 700~900M）
sudo tee /etc/elasticsearch/jvm.options.d/heap.options >/dev/null <<'JVM'
-Xms512m
-Xmx512m
JVM

# 只监听本机、单节点、关安全模块（应用连的是 http://127.0.0.1:9200，无账号密码）
sudo tee -a /etc/elasticsearch/elasticsearch.yml >/dev/null <<'YML'

# === gf-eshop 小内存部署追加 ===
node.name: node-1
network.host: 127.0.0.1
http.port: 9200
discovery.type: single-node
xpack.security.enabled: false
# 2G 机器不要锁内存（锁了反而更容易 OOM/启动失败）
bootstrap.memory_lock: false
YML

sudo systemctl enable --now elasticsearch
sleep 20 && curl -s localhost:9200 | head -12
```

> 索引建出来之后（见 A.5）会显示 `yellow`：单节点上副本分片无处分配，属正常，不影响检索。
> 想变绿可以设副本为 0（可选，纯美化），命令见 A.5。

### A.4 观测：到底撑住了没有

把 `deploy/es-watch.sh` 传到服务器（release 包里已含），跑 10 分钟：

```bash
chmod +x es-watch.sh
./es-watch.sh 120 5 /tmp/es-watch.log      # 每 5s 采一次，共 120 次
```

它每轮打印：`内存 used/free/swap_used`、`ES 堆 used/max 与占比`、`si/so`（swap 换入换出）、ES 健康状态，
并自动检查内核 OOM kill 与应用日志里的「ES 检索失败，降级 DB 查询」。最后给结论：

| 结论 | 含义 | 处置 |
|------|------|------|
| ✅ 通过 | 全程 si/so=0、无 OOM、无降级，堆峰值有富余 | 可以留着用，但别再往上加负载 |
| ⚠️ 临界 | 偶发 swap 换出，或堆峰值 >80% | 能用别硬撑；数据量涨了就加内存 |
| ❌ 不通过 | 持续 swap 换出 / 出现 OOM kill / 应用已在降级 | 把 ES 关掉（见 A.5），老实走 MySQL |

> 观测脚本是 Linux 专用（依赖 `free`/`vmstat`/`journalctl`），只在 macOS 上做过语法检查、
> **未经真实 Linux 机器验证**；它只读不写配置，跑起来发现不对 Ctrl-C 即可。

### A.5 打开应用侧并验收 / 失败的退路

```bash
# 打开 ES：手工编辑，只改 elasticsearch 段那一处（配置里有多个 enabled，别用 sed 批量替换）
sudo vi /opt/gf-eshop/manifest/config/config.yaml     # elasticsearch.enabled: false -> true
sudo systemctl restart gf-eshop

# 验收 1：索引被自动建出来（100 品牌 / 2025 商品，几秒）
curl -s "localhost:9200/_cat/indices/eshop_*?h=index,health,docs.count"

# 可选：让索引变绿（单节点副本无处分配才 yellow，设 0 即绿；不影响功能）
curl -X PUT "localhost:9200/eshop_*/_settings" -H 'Content-Type: application/json' \
  -d '{"index":{"number_of_replicas":0}}'

# 验收 2：走 ES 的筛选查询正常返回（对比是否与 DB 顺序一致）
curl -s "localhost:8000/api/v1/brands?status=1&page=1&page_size=5" | head -c 200

# 验收 3：日志里没有降级告警（有就说明 ES 不可用）
journalctl -u gf-eshop -n 200 --no-pager | grep -c "ES 检索失败，降级 DB 查询"
```

**退路（一分钟回到稳定态）**：把 `elasticsearch.enabled` 改回 `false` 并
`sudo systemctl stop elasticsearch && sudo systemctl disable elasticsearch` 重启应用即可 ——
筛选查询自动回落 MySQL，功能不缺，且此时机器会立刻松开约 700~900M 内存。

> 压测建议：真要定论，按 [`docs/perf-workflow.md`](../docs/perf-workflow.md) 跑一轮并发检索，
> 观察「并发检索 + 段合并」同时发生时是否触发降级。稳态空闲不换出，不代表高峰不换出。
