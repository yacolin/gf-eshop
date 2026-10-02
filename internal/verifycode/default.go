package verifycode

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	defaultMu  sync.Mutex
	defaultSvc *Service
)

// Default 返回全局验证码服务，按需装配各渠道发送器。
//
// 只缓存成功结果：配置读取失败时不缓存错误，下次调用会重试，
// 避免一次瞬时抖动让整个进程后续一直不可用。
func Default(ctx context.Context) (*Service, error) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultSvc != nil {
		return defaultSvc, nil
	}

	cfg, err := LoadConfig(ctx)
	if err != nil {
		return nil, err
	}
	emailSender, err := NewEmailSender(ctx)
	if err != nil {
		return nil, err
	}
	smsSender, err := NewSMSSender(ctx)
	if err != nil {
		return nil, err
	}

	defaultSvc = New(cfg, emailSender, smsSender)
	return defaultSvc, nil
}

// WarnNotReady 在启动时逐渠道自检配置，只告警不阻断启动。
//
// 验证码属于可选能力：某个渠道没配好，只会让使用该渠道的接口返回
// CodeVerifyChannelNotReady，不影响其它功能。
func WarnNotReady(ctx context.Context) {
	svc, err := Default(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "读取 verifycode 配置失败，验证码功能不可用：%v", err)
		return
	}

	for _, sender := range svc.senders {
		// 日志模式：功能可用但不会真正外发，必须显著告警，避免生产环境误判
		switch s := sender.(type) {
		case *EmailSender:
			if s.cfg.DevMode {
				g.Log().Warning(ctx, "email.dev_mode 已开启：验证码只会写入日志、不会真实发信，生产环境请关闭")
				continue
			}
		case *SMSSender:
			if s.cfg.DevMode {
				g.Log().Warning(ctx, "sms.dev_mode 已开启：验证码只会写入日志、不会真实外发，生产环境请关闭")
				continue
			}
		}
		if err = sender.Ready(); err != nil {
			g.Log().Warningf(ctx, "验证码渠道「%s」未就绪，使用该渠道的接口将不可用：%v", sender.Name(), err)
		}
	}
}
