package verifycode

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/utility"
)

// EmailSender 是邮件渠道：SMTP 协议细节全部封装在 utility/mail.go，
// 本文件只负责「把验证码渲染成邮件正文」这一件事。
type EmailSender struct {
	cfg *utility.MailConfig
}

// NewEmailSender 读取 email 段配置构造邮件渠道。
func NewEmailSender(ctx context.Context) (*EmailSender, error) {
	cfg, err := utility.LoadMailConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &EmailSender{cfg: cfg}, nil
}

func (s *EmailSender) Channel() Channel { return ChannelEmail }

func (s *EmailSender) Name() string { return "邮件" }

// Ready 判断邮件渠道是否可用。dev_mode 下不连 SMTP，因此视为可用。
func (s *EmailSender) Ready() error {
	if s == nil || s.cfg == nil {
		return errors.New("邮件配置为空")
	}
	return s.cfg.Ready()
}

// MailConfig 暴露底层邮件配置，便于启动自检时打印具体缺哪些项。
func (s *EmailSender) MailConfig() *utility.MailConfig { return s.cfg }

func (s *EmailSender) Send(ctx context.Context, to, code string, expire time.Duration) error {
	if s.cfg.DevMode {
		// 结构化一行，便于 `grep -oE 'code=[0-9]{6}'` 取码；
		// 邮件正文本身也会被 utility.SendMail 在 dev_mode 下打进日志，供检查模板。
		g.Log().Warningf(ctx, "[email.dev_mode] 验证码已生成（未真实发信）to=%s code=%s expire=%s",
			to, code, expire)
	}
	return utility.SendMail(ctx, s.cfg, []string{to}, s.cfg.Subject, buildEmailBody(code, expire))
}

// buildEmailBody 渲染验证码邮件正文（HTML）。
func buildEmailBody(code string, expire time.Duration) string {
	// code 由 crypto/rand 生成且只含数字，直接内插不存在注入风险
	return `<div style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;">
  <h2 style="margin:0 0 16px;font-size:18px;color:#111;">邮箱验证码</h2>
  <p style="margin:0 0 12px;color:#444;font-size:14px;">您正在进行身份验证，验证码为：</p>
  <p style="margin:0 0 12px;font-size:30px;font-weight:700;letter-spacing:6px;color:#1677ff;">` + code + `</p>
  <p style="margin:0 0 12px;color:#444;font-size:14px;">验证码 ` + FormatDuration(expire) + ` 内有效，请勿转发给他人。</p>
  <p style="margin:0;color:#999;font-size:12px;">若非本人操作，请忽略本邮件。</p>
</div>`
}
