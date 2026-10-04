package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// 环境变量覆盖配置：/etc/gf-eshop.env → manifest/config/config.yaml。
//
// 为什么需要它：发布包里的 manifest/config/config.yaml 是**打包机上的本地文件**
//（`**/config/config.yaml` 在 .gitignore 里），里面有开发库密码、邮箱授权码；
//而每次解包都会覆盖服务器上的那一份（deploy/INSTALL.md 原先就警告过这点）。
//于是服务器的真实配置只能放 /etc/gf-eshop.env（chmod 600），
//由 deploy/start.sh 与 systemd 的 EnvironmentFile 注入进程环境，再由这里覆盖到配置上。
//
// 命名约定（与隔壁 campus_express 一致）：配置键 → 大写蛇形，
// "." 与驼峰边界都变成 "_"：
//
//	orderShard.defaultWindowMonths → ORDER_SHARD_DEFAULT_WINDOW_MONTHS
//	database.default.link          → DATABASE_DEFAULT_LINK
//	elasticsearch.addresses        → ELASTICSEARCH_ADDRESSES   （逗号分隔，与 .Strings() 一致）
//
// 语义：
//   - **未设置 = 用 yaml 的值**（不会把配置清空）；
//   - 设置成空串也算「设置了」（可用来显式清空某一项）；
//   - 只覆盖 EnvOverridableKeys 里列出的键 —— 显式清单便于审阅，
//     避免「某个环境变量意外改了线上配置」；
//   - 匹配**大小写不敏感**（systemd/shell 里写成 Server_Address 也认）。
//
// 注意：GoFrame 的 gcfg 没有 viper 那种 AutomaticEnv，所以这里是显式实现。

// EnvOverridableKeys 允许被环境变量覆盖的配置键。
//
// 新增可覆盖项 = 在这里加一行 + 在 deploy/env.example 里加一条。
var EnvOverridableKeys = []string{
	// 接入层
	"server.address",
	"server.accessLogEnabled",
	"server.logStdout",

	// 依赖
	"database.default.link",
	"redis.default.address",
	"redis.default.pass",
	"redis.default.db",

	// 安全
	"jwt.secret",
	"jwt.access_expire",
	"jwt.refresh_expire",

	// 日志
	"logger.level",
	"logger.stdout",

	// 检索
	"elasticsearch.enabled",
	"elasticsearch.addresses",
	"elasticsearch.username",
	"elasticsearch.password",

	// 邮件 / 短信（渠道开关与凭据按环境不同）
	"email.enabled",
	"email.dev_mode",
	"email.smtp",
	"email.port",
	"email.user",
	"email.pass",
	"email.from",
	"email.from_name",
	"sms.enabled",
	"sms.dev_mode",
	"sms.sign_name",
	"sms.template_id",

	// 订单分表（灰度/切读开关，按环境不同）
	"orderShard.mode",
	"orderShard.dualWrite",
	"orderShard.shadowReadPercent",
	"orderShard.readShardsPercent",
	"orderShard.defaultWindowMonths",
}

// ConfigKeyToEnvName 把配置键转成环境变量名（大写蛇形，驼峰边界也断开）。
//
//	orderShard.defaultWindowMonths → ORDER_SHARD_DEFAULT_WINDOW_MONTHS
//	database.default.link          → DATABASE_DEFAULT_LINK
func ConfigKeyToEnvName(key string) string {
	var b strings.Builder
	for i, r := range key {
		switch {
		case r == '.':
			b.WriteByte('_')
		case r >= 'A' && r <= 'Z':
			// 驼峰边界：前一个字符是小写/数字时插下划线（本项目的键都是小驼峰，够用）
			if i > 0 {
				prev := rune(key[i-1])
				if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
					b.WriteByte('_')
				}
			}
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return strings.ToUpper(b.String())
}

// ApplyEnvOverrides 把环境变量里的覆盖值写进配置，返回被覆盖的配置键（供启动日志）。
//
// getenv 为 nil 时读取真实进程环境（大小写不敏感）；测试可注入假环境。
func ApplyEnvOverrides(ctx context.Context, getenv func(string) (string, bool)) []string {
	if getenv == nil {
		getenv = caseInsensitiveEnv()
	}
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		g.Log().Warningf(ctx, "环境变量覆盖跳过：配置适配器不是文件适配器（%T）", g.Cfg().GetAdapter())
		return nil
	}

	applied := make([]string, 0, len(EnvOverridableKeys))
	for _, key := range EnvOverridableKeys {
		name := ConfigKeyToEnvName(key)
		value, ok := getenv(name)
		if !ok {
			continue // 未设置：保留 yaml 的值
		}
		if err := adapter.Set(key, value); err != nil {
			g.Log().Warningf(ctx, "环境变量 %s 覆盖配置 %s 失败（已忽略）: %v", name, key, err)
			continue
		}
		applied = append(applied, key)
	}
	return applied
}

// caseInsensitiveEnv 返回一个大小写不敏感的环境变量查询函数。
func caseInsensitiveEnv() func(string) (string, bool) {
	folded := make(map[string]string, 64)
	for _, kv := range os.Environ() {
		if i := strings.IndexByte(kv, '='); i > 0 {
			folded[strings.ToUpper(kv[:i])] = kv[i+1:]
		}
	}
	return func(name string) (string, bool) {
		v, ok := folded[strings.ToUpper(name)]
		return v, ok
	}
}
