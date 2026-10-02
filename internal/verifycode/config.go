package verifycode

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Config 是渠道无关的验证码风控配置（对应 config.yaml 的 verifycode 段）。
//
// 放在这里而不是各渠道配置里：额度、冷却、失败上限属于「验证码」这件事本身，
// 不应该随着渠道增加而各配一份。
type Config struct {
	CodeExpire   time.Duration // 验证码有效期
	ResendAfter  time.Duration // 同一收件人的重发冷却
	DailyLimit   int           // 单收件人每日发送上限
	IPDailyLimit int           // 单 IP 每日发送上限
	MaxAttempts  int           // 单个验证码最多校验次数
}

const (
	defaultCodeExpire  = 2 * time.Minute
	defaultResendAfter = 60 * time.Second
	defaultDailyLimit  = 10
	defaultIPLimit     = 30
	defaultMaxAttempt  = 5
)

// LoadConfig 读取 verifycode 段配置。缺键时使用内置默认值，不会 panic。
func LoadConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	var err error

	if cfg.CodeExpire, err = configDuration(ctx, "verifycode.code_expire", defaultCodeExpire); err != nil {
		return nil, err
	}
	if cfg.ResendAfter, err = configDuration(ctx, "verifycode.resend_interval", defaultResendAfter); err != nil {
		return nil, err
	}
	if cfg.DailyLimit, err = configInt(ctx, "verifycode.daily_limit", defaultDailyLimit); err != nil {
		return nil, err
	}
	if cfg.IPDailyLimit, err = configInt(ctx, "verifycode.ip_daily_limit", defaultIPLimit); err != nil {
		return nil, err
	}
	if cfg.MaxAttempts, err = configInt(ctx, "verifycode.max_attempts", defaultMaxAttempt); err != nil {
		return nil, err
	}

	if cfg.CodeExpire <= 0 {
		return nil, gerror.Newf("verifycode.code_expire 必须为正数，当前 %s", cfg.CodeExpire)
	}
	if cfg.ResendAfter < 0 {
		return nil, gerror.Newf("verifycode.resend_interval 不能为负数，当前 %s", cfg.ResendAfter)
	}
	if cfg.MaxAttempts <= 0 {
		return nil, gerror.Newf("verifycode.max_attempts 必须为正数，当前 %d", cfg.MaxAttempts)
	}
	return cfg, nil
}

func configText(ctx context.Context, key string) (string, error) {
	v, err := g.Cfg().Get(ctx, key, "")
	if err != nil {
		return "", gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return v.String(), nil
}

func configBool(ctx context.Context, key string, def bool) (bool, error) {
	v, err := g.Cfg().Get(ctx, key, def)
	if err != nil {
		return false, gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return v.Bool(), nil
}

func configInt(ctx context.Context, key string, def int) (int, error) {
	v, err := g.Cfg().Get(ctx, key, def)
	if err != nil {
		return 0, gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return v.Int(), nil
}

// configDuration 支持 "2m" / "60s" 写法，也兼容纯数字（按秒处理）。
func configDuration(ctx context.Context, key string, def time.Duration) (time.Duration, error) {
	raw, err := configText(ctx, key)
	if err != nil {
		return 0, err
	}
	if raw == "" {
		return def, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		if secs, convErr := time.ParseDuration(raw + "s"); convErr == nil {
			return secs, nil
		}
		return 0, gerror.Wrapf(err, "配置项 %s 不是合法时长（示例：2m、60s）：%s", key, raw)
	}
	return d, nil
}
