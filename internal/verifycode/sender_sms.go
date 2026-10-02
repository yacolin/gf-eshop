package verifycode

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// SMSConfig 是短信渠道配置（对应 config.yaml 的 sms 段）。
//
// 国内短信需要企业资质 + 签名报备 + 模板审核（详见 docs/email-verify-code-guide.md），
// 因此这里只保留渠道无关的开关与签名/模板标识；各服务商的 AK/SK 等参数
// 由具体 SMSProvider 实现自行读取，避免把某一家服务商的字段固化进核心包。
type SMSConfig struct {
	Enabled    bool          // 短信渠道总开关
	DevMode    bool          // 开发模式：不真正外发，把验证码写进日志
	SignName   string        // 短信签名（报备通过后获得）
	TemplateID string        // 短信模板 ID（报备通过后获得）
	Timeout    time.Duration // 调用服务商接口的超时
}

const (
	defaultSMSTimeout = 5 * time.Second
	// smsProviderLogName 是内置日志 provider 的名字：
	// 若渠道已启用却仍在用日志 provider，说明真实通道没接上，Ready 应判定为不可用。
	smsProviderLogName = "log"
)

// SMSProvider 是真实短信通道的接入点。
//
// 接入云厂商（阿里云/腾讯云/火山等）时，只需实现该接口，并在 NewSMSSender 里
// 替换默认的日志 provider —— 验证码的生成、指纹存储、冷却、限额、校验全部复用，
// 不需要改动本包以外的任何代码。
type SMSProvider interface {
	// Name 返回服务商名称，用于就绪判断与日志
	Name() string
	// Send 调用服务商接口下发一条验证码短信
	Send(ctx context.Context, cfg *SMSConfig, phone, code string, expire time.Duration) error
}

// SMSSender 是短信渠道。
type SMSSender struct {
	cfg      *SMSConfig
	provider SMSProvider
}

// NewSMSSender 读取 sms 段配置构造短信渠道。
//
// 当前默认使用日志 provider：在拿到企业资质、完成签名与模板报备之前，
// 可以先用它把「手机号验证码」的完整流程（含风控与校验）跑通。
func NewSMSSender(ctx context.Context) (*SMSSender, error) {
	cfg, err := LoadSMSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &SMSSender{cfg: cfg, provider: logSMSProvider{}}, nil
}

func (s *SMSSender) Channel() Channel { return ChannelSMS }

func (s *SMSSender) Name() string { return "短信" }

// Ready 判断短信渠道是否可用。
func (s *SMSSender) Ready() error {
	if s == nil || s.cfg == nil {
		return errors.New("短信配置为空")
	}
	if s.cfg.DevMode {
		return nil // 日志模式可用，只是不外发
	}
	if !s.cfg.Enabled {
		return errors.New("sms.enabled 为 false，短信服务未启用")
	}
	if s.provider == nil || s.provider.Name() == smsProviderLogName {
		return errors.New("sms.enabled 为 true 但未接入真实短信通道：" +
			"请实现 verifycode.SMSProvider 并在 NewSMSSender 中注入")
	}
	if strings.TrimSpace(s.cfg.SignName) == "" || strings.TrimSpace(s.cfg.TemplateID) == "" {
		return errors.New("sms.sign_name / sms.template_id 未配置（需先完成签名与模板报备）")
	}
	return nil
}

func (s *SMSSender) Send(ctx context.Context, phone, code string, expire time.Duration) error {
	if err := s.Ready(); err != nil {
		return err
	}
	if s.cfg.DevMode {
		g.Log().Warningf(ctx, "[sms.dev_mode] 验证码已生成（未真实外发）to=%s code=%s expire=%s",
			phone, code, expire)
		return nil
	}
	return s.provider.Send(ctx, s.cfg, phone, code, expire)
}

// SMSConfig 暴露底层配置，便于启动自检与测试。
func (s *SMSSender) SMSConfig() *SMSConfig { return s.cfg }

// LoadSMSConfig 读取 sms 段配置，缺键时使用默认值。
func LoadSMSConfig(ctx context.Context) (*SMSConfig, error) {
	cfg := &SMSConfig{}
	var err error

	if cfg.Enabled, err = configBool(ctx, "sms.enabled", false); err != nil {
		return nil, err
	}
	if cfg.DevMode, err = configBool(ctx, "sms.dev_mode", false); err != nil {
		return nil, err
	}
	if cfg.SignName, err = configText(ctx, "sms.sign_name"); err != nil {
		return nil, err
	}
	if cfg.TemplateID, err = configText(ctx, "sms.template_id"); err != nil {
		return nil, err
	}
	if cfg.Timeout, err = configDuration(ctx, "sms.timeout", defaultSMSTimeout); err != nil {
		return nil, err
	}
	cfg.SignName = strings.TrimSpace(cfg.SignName)
	cfg.TemplateID = strings.TrimSpace(cfg.TemplateID)
	return cfg, nil
}

// logSMSProvider 是默认 provider：只记日志，不调用任何外部接口。
//
// 它存在的意义是让「手机号验证码」这条链路可以先行开发与联调，
// 而不是在拿到资质前留下一个假装成功的空实现。
type logSMSProvider struct{}

func (logSMSProvider) Name() string { return smsProviderLogName }

func (logSMSProvider) Send(ctx context.Context, cfg *SMSConfig, phone, code string, expire time.Duration) error {
	g.Log().Warningf(ctx, "[sms.provider=%s] 短信未真实外发 sign=%s template=%s to=%s code=%s expire=%s（写入日志）",
		smsProviderLogName, cfg.SignName, cfg.TemplateID, phone, code, FormatDuration(expire))
	return nil
}

// SetProvider 注入真实的短信服务商实现（供接入时调用，也便于单测替换）。
func (s *SMSSender) SetProvider(p SMSProvider) {
	if p == nil {
		return
	}
	s.provider = p
}
