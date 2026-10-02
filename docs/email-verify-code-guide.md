# 验证码（邮箱 / 短信）接入说明

面向 `user_auth` 模块的验证码能力：**用 GoFrame 管配置与风控，用标准库 `net/smtp` 发信**，
不引入第三方发信库（`gomail` 等），因此没有额外依赖、也没有供应链风险。

## 0. 分层：渠道无关的核心 + 可插拔渠道

```
internal/verifycode/            渠道无关的验证码底座（与 internal/search 同类定位）
  verifycode.go                 Send / Consume、指纹存储、冷却、额度、防爆破
  config.go                     verifycode.* 风控配置
  sender_email.go               邮件渠道（正文模板；SMTP 细节在 utility/mail.go）
  sender_sms.go                 短信渠道（日志模式 + SMSProvider 扩展点）
  default.go                    全局装配与启动自检
internal/logic/user_auth/       业务胶水：场景合法性、收件人↔账号映射、会话签发
  verification_code.go          发码通用逻辑（按 channel 分派）+ 邮箱验证码登录
  password_reset.go             忘记密码重置
```

**接一个新渠道（例如短信）只需要实现 `verifycode.Sender`**，额度、冷却、防爆破、防枚举、
指纹存储全部复用 —— 不需要再抄一遍风控。`verifycode_test.go` 里用假 sender 验证了这一点。

## 1. 配置

`manifest/config/config.yaml` 被 `.gitignore` 排除，需要**手动**补上这几段（完整示例）：

```yaml
# 邮件渠道
email:
  enabled: false          # 邮件渠道总开关
  dev_mode: false         # 仅本地开发：不真实发信，把验证码写进日志（生产必须 false）
  smtp: "smtp.qq.com"     # SMTP 服务器地址
  port: 465               # 465=隐式 TLS（推荐） 587=STARTTLS 25=明文（不推荐）
  tls_mode: "auto"        # auto=按端口判断 implicit=强制隐式 TLS starttls=强制 STARTTLS
  user: ""                # SMTP 登录账号
  pass: ""                # 邮箱授权码（不是登录密码！）
  from: ""                # 发件人地址；必须与 user 同一个邮箱，否则 SPF 不过会被拒收
  from_name: "gf-eshop"   # 发件人显示名
  helo_name: ""           # EHLO 域名，留空取 from 的域名
  subject: "【gf-eshop】邮箱验证码"
  timeout: "10s"          # 单次发信超时（连接 + 读写共用）

# 验证码风控（渠道无关：email / sms 共用同一套额度、冷却与失败上限）
verifycode:
  code_expire: "2m"       # 验证码有效期
  resend_interval: "60s"  # 同一收件人的重发冷却
  daily_limit: 10         # 单收件人每日发送上限
  ip_daily_limit: 30      # 单 IP 每日发送上限
  max_attempts: 5         # 单个验证码最多校验次数

# 短信渠道（资质到位前保持 enabled=false）
sms:
  enabled: false          # 短信渠道总开关
  dev_mode: false         # 仅本地开发：不真实外发，把验证码写进日志
  sign_name: ""           # 短信签名（报备通过后获得）
  template_id: ""         # 短信模板 ID（报备通过后获得）
  timeout: "5s"
```

要点：

- **风控参数在 `verifycode.*`**（不再是 `email.*`）：它们属于「验证码」本身，
  不应该随渠道增加而各配一份。
- **`pass` 填授权码**：QQ 邮箱 / 163 邮箱都需在网页端开启 SMTP 服务并获取授权码，
  直接用登录密码会认证失败（错误信息里也会提示这一点）。
- 所有键都有代码内置默认值，缺失不会 panic。
- `dev_mode: true`（邮件/短信各自独立）时不连外部服务，验证码会打一行结构化日志：

  ```
  [email.dev_mode] 验证码已生成（未真实发信）to=you@x.com code=123456 expire=2m0s
  ```

  取码：`grep -oE 'code=[0-9]{6}' 日志文件 | tail -1 | cut -d= -f2`
  （不要从邮件正文 HTML 里抠码：日志是多行写入的，并发读取时可能拿到半行）。
  启动日志会给出醒目告警。生产环境必须关闭。
- 服务启动时**逐渠道自检**，未就绪只在日志告警、**不阻断启动**；此时该渠道接口返回 `1021`。
  例如短信未接入时会打印：`验证码渠道「短信」未就绪，使用该渠道的接口将不可用：sms.enabled 为 false...`

### 接短信要做什么

1. 实现 `verifycode.SMSProvider`（一个 HTTPS 调用云厂商 API 即可，官方 SDK 只需几行）；
2. 在 `NewSMSSender` 里注入它（当前默认注入的是只写日志的 `logSMSProvider`）；
3. 在 `verifycode.LoadSMSConfig` 或 provider 内部读取自己的 AK/SK 配置；
4. 业务侧按需加一个 `phone/code` 接口，`verification_code.go` 里按 `ChannelSMS` 分派即可
   （`findUserByTarget` 已预留 `users.phone` 分支）。

**注意**：`sms.enabled: true` 但仍用日志 provider 时，`Ready()` 会判定为不可用并报错，
不会假装发送成功 —— 避免线上配置写错却"看起来正常"。

## 2. 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/auth/email/code` | 发送验证码（三种 scene 共用） |
| POST | `/api/v1/user/auth/email/login` | 邮箱验证码登录 |
| POST | `/api/v1/user/auth/password/reset` | 邮箱验证码重置密码（忘记密码） |

三个接口都在 `internal/middleware/user_auth.go` 的 `publicPaths` 白名单里，无需登录态。

**发送验证码**

```jsonc
// 请求
{ "email": "user@example.com", "scene": "login" }   // scene 可省略，默认 login
// 响应
{ "code": 0, "data": { "success": true, "expire_in": 120, "resend_after": 60 } }
```

`scene` 是验证码用途，三种**互相隔离**（注册码不能用于登录，登录码不能用于重置）：

- `login` —— 邮箱验证码登录
- `register` —— 注册前验证邮箱归属
- `reset` —— 忘记密码重置

**邮箱验证码登录**

```jsonc
{ "email": "user@example.com", "code": "123456" }
```

**邮箱验证码重置密码**

```jsonc
// 请求
{ "email": "user@example.com", "code": "123456", "new_password": "至少8位" }
// 响应
{ "code": 0, "data": { "success": true, "revoked_sessions": 2 } }
```

重置成功后会**吊销该用户全部 refresh token**（所有设备需重新登录），`revoked_sessions` 是吊销数量。

> ⚠️ `access_token` 是无状态 JWT，无法回收，最长会在 `jwt.access_expire`（默认 30 分钟）后自然失效。
> 需要「改密后立刻失效」就必须引入 access token 黑名单或 token 版本号。

**注册时验证邮箱**（可选字段，兼容旧行为）

```jsonc
{
  "username": "colin", "password": "***",
  "email": "user@example.com",
  "email_code": "123456"        // 传了才校验，通过后 email_verified = 1
}
```

## 3. 错误码

| code | 含义 |
|------|------|
| 1017 | 验证码错误 / 已过期 / **邮箱未注册**（故意合并，见下） |
| 1018 | 发送过于频繁（消息里带剩余秒数） |
| 1019 | 验证码错误次数过多，请重新获取 |
| 1020 | 该邮箱今日发送次数已达上限 |
| 1021 | `<渠道>服务未启用或配置不完整`（文案按渠道拼装，如「短信服务未启用…」）|
| 1022 | `验证码<渠道>发送失败`（已回滚冷却与额度，可立即重试）|
| 1023 | 该邮箱已被绑定（注册场景） |
| 51 | 新密码长度不是 8-64 位（框架校验） |

**为什么「邮箱未注册」不单独给一个错误码**：登录与重置都采用「先验证码、后查账号」的顺序，
未注册邮箱不可能持有验证码，因此它和「验证码错误」返回的是**完全相同的 1017**。
若给未注册单独一个码，攻击者就能用错误码差异批量探测某邮箱是否在本站注册过 —— 这正是
最初实现踩过的坑。只有验证码校验**通过之后**（意味着请求方已能收信）才会返回账号状态类错误。

## 4. 安全设计

| 风险 | 处理方式 |
|------|----------|
| 缓存泄露导致验证码被冒用 | Redis 只存 `sha256(scene + email + code)` 指纹，不存明文 |
| 验证码被重放 | 校验通过立即 `DEL`，一次性使用 |
| 暴力猜解 | 失败次数累加，达到 `max_attempts` 直接销毁；比较用 `subtle.ConstantTimeCompare` |
| 被当成垃圾邮件发射器 | 单邮箱日限额 + 单 IP 日限额 + 重发冷却（`SET NX` 原子抢占，防并发重复发） |
| 邮箱枚举 | 未注册邮箱同样返回成功但**不发信**；且额度/冷却/验证码写入照常执行，两条路径无观测差异；验证类接口统一返回 1017 |
| 未注册邮箱绕过限流探测 | 风控对「发信」与「静默跳过」一视同仁，静默路径同样计数、同样受冷却约束 |
| 连点「重发」耗光当天额度 | 冷却抢占放在额度扣减**之前**，被冷却拒绝的请求不消耗额度；额度已满时释放已抢到的冷却 |
| 跨场景/跨渠道复用验证码 | 指纹与 Redis key 都绑定 `channel` 与 `scene`：`verify:code:<channel>:<scene>:<收件人>` |
| 改密后旧会话继续可用 | 重置成功即按 `SCAN` 游标删除该用户全部 `user:refresh:<uid>:*`（不用 `KEYS`，避免阻塞 Redis） |
| 头部注入 | 所有进入邮件头的值拒绝 CR/LF；收件人经 `net/mail.ParseAddress` 校验 |
| 请求被 SMTP 拖死 | 连接 + 读写共用 `timeout` deadline；失败重试 1 次（间隔 1s） |
| 明文连接泄露授权码 | 依赖 `smtp.PlainAuth` 自带保护：非 TLS 且非 localhost 时拒绝发送凭据 |

发信失败时会回滚：删除验证码 key、释放重发冷却、归还每日额度，用户可立即重试。

## 5. 验证码生成

用 `crypto/rand`（`utility.GenerateNumericCode`），**不是** `math/rand`：
验证码属于身份凭证，`math/rand` 的序列可被预测，且 Go 1.20+ 已废弃全局 `rand.Seed`。

## 6. 到达率（能不能真的进收件箱）

代码层面能保证的是「握手、TLS、协议、报文格式正确」，**能否进收件箱由发件域的信誉与认证决定**。

### 6.1 铁律：`from` 必须与 `user` 是同一个邮箱

国内免费邮箱只授权自己的域发信，实测域名 SPF：

| 域 | SPF | 含义 |
|----|-----|------|
| `qq.com` | `v=spf1 include:spf.mail.qq.com ~all` | 只有 QQ 的服务器能用 `@qq.com` 发信 |
| `sina.com` / `sina.cn` | `v=spf1 include:spf.sinamail.sina.com.cn -all` | `-all` 是**硬失败**：其它服务器冒充 `@sina.com` 直接判定不合格 |

所以「用 QQ 邮箱的授权码认证，却把 `from` 写成自己的 `@sina.com` 地址」会被新浪判为 SPF 失败
（`sina.com` 的 DMARC 是 `p=none`，不会立刻退信，但大概率进垃圾箱，或被以「发件人已被拒绝」退信）。
**结论：`from` 就填你拿来认证的那个邮箱地址。**

### 6.2 已验证的协议事实

对着真实服务器探测过（仅握手，不发信）：

| 服务器 | 端口 | 结果 |
|--------|------|------|
| `smtp.qq.com` | 465 | 隐式 TLS，宣告 `AUTH LOGIN PLAIN XOAUTH XOAUTH2`、`SIZE 70MB` |
| `smtp.qq.com` | 587 | 宣告 STARTTLS，升级后同理支持 PLAIN |
| `smtp.sina.com` | 465 / 587 | 隐式 TLS / STARTTLS，宣告 `AUTH LOGIN PLAIN` |
| `smtp.163.com` | 465 | 隐式 TLS，宣告 `AUTH LOGIN PLAIN XOAUTH2` |

两个由此发现并修掉的坑（都有回归测试）：

1. **`net/smtp` 的 Client 不会自动 STARTTLS**（只有包级 `smtp.SendMail` 会），必须显式调用
   `client.StartTLS(...)`。否则 587 上 `PlainAuth` 会以 `unencrypted connection` 拒绝发送凭据，
   即 **587 端口完全发不出去**。实测升级前拒绝、升级后允许。
2. **`net/smtp` 默认发 `EHLO localhost`**，国内邮箱普遍把它当垃圾邮件特征，
   现已显式传发件域名（`HELO_NAME`/`email.helo_name`）。

### 6.3 新浪邮箱侧的准备

- `@sina.com` → `smtp.sina.com:465`；`@sina.cn` → `smtp.sina.cn:465`（VIP 域另有 `smtp.vip.sina.*`）；
  官方说明见[新浪帮助：如何在 iPhone/iPad 上设置新浪免费邮箱](http://help.sina.com.cn/comquestiondetail/view/798/)。
- 必须先在网页版「设置区 → 客户端 POP/IMAP/SMTP」开启服务，并生成 16 位**授权码**填到 `email.pass`；
  官方说明见[如何获取授权码登录第三方邮箱](http://help.sina.com.cn/comquestiondetail/view/1566/)。
- 用个人邮箱（QQ/163/新浪）当发信通道，**每日发信量很小且容易触发风控**，
  只适合开发/演示；正式环境请换成带 SPF + DKIM 的域名，或事务邮件服务（阿里云邮件推送、腾讯云 SES、SendGrid 等）。

### 6.4 开发阶段建议

- 只想跑通流程：把 `dev_mode` 打开，验证码直接进日志，不用碰真实邮箱。
- 想验证真实到达率：用**收件方就是发件方**的方式自测（`user`/`from` 都填你的新浪邮箱，
  `smtp.sina.com:465`），避免 SPF 不一致带来的干扰；记得去垃圾箱看一眼。

## 7. 键空间与迁移说明

重构后 Redis 键统一为：

```
verify:code:<channel>:<scene>:<收件人>            验证码指纹（TTL=有效期，校验通过即删）
verify:code:cooldown:<channel>:<scene>:<收件人>   重发冷却
verify:code:daily:<channel>:<收件人>              收件人维度日计数（24h）
verify:code:daily:ip:<channel>:<ip>               IP 维度日计数（24h）
```

重构前是 `email:code:*`。旧键不迁移（最长 24h 自然过期）；副作用只有一个：
**切换当天各收件人/本机的日计数会从 0 重新开始**，等于当天多了一点额度，无安全影响。

## 8. 顺带修复的既有缺陷

`usr_users.email` / `usr_users.phone` 有唯一索引且允许为 NULL，但 `Register` 原先插入的是
**空串**，于是「第一个不带手机号的用户注册成功后，后续所有同类注册都会撞 `uk_phone` 报 SQL 错误」。
现改为未填写时传 `nil`（GoFrame 的 DO 结构体会自动 `OmitNilData`，落 `DEFAULT NULL`），
与字段注释「NULL 表示未绑定」的语义一致。
