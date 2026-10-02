package utility

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

func testMailConfig() *MailConfig {
	return &MailConfig{
		Enabled:  true,
		Host:     "smtp.example.com",
		Port:     ImplicitTLSPort,
		User:     "noreply@example.com",
		Pass:     "auth-code",
		From:     "noreply@example.com",
		FromName: "gf-eshop",
		Subject:  "【gf-eshop】邮箱验证码",
		Timeout:  5 * time.Second,
	}
}

func TestBuildMailMessageBasic(t *testing.T) {
	cfg := testMailConfig()
	body := `<b>您的验证码为 123456</b>`
	now := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)

	raw, err := buildMailMessage(cfg, []string{"user@example.com"}, cfg.Subject, body, now)
	if err != nil {
		t.Fatalf("buildMailMessage 返回错误: %v", err)
	}
	msg := string(raw)

	for _, want := range []string{
		"To: user@example.com\r\n",
		"noreply@example.com",
		"MIME-Version: 1.0\r\n",
		`Content-Type: text/html; charset="UTF-8"`,
		"Content-Transfer-Encoding: base64\r\n",
		"Date: Wed, 01 May 2024 12:00:00 +0000\r\n",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("邮件报文缺少 %q\n实际报文:\n%s", want, msg)
		}
	}

	// 中文主题必须做 RFC 2047 编码，否则会被客户端显示为乱码
	if !strings.Contains(msg, "Subject: =?UTF-8?") {
		t.Errorf("中文主题未做 RFC 2047 编码:\n%s", msg)
	}

	// 正文必须是 base64，且解码后与原文一致
	headerEnd := strings.Index(msg, "\r\n\r\n")
	if headerEnd < 0 {
		t.Fatalf("报文缺少头部与正文的分隔空行:\n%s", msg)
	}
	encoded := strings.ReplaceAll(msg[headerEnd+4:], "\r\n", "")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("正文不是合法 base64: %v", err)
	}
	if string(decoded) != body {
		t.Errorf("正文解码结果不符\n期望: %s\n实际: %s", body, decoded)
	}

	// MIME 行宽不应超过 78 字符（base64 折行 76 + CRLF）
	for _, line := range strings.Split(msg, "\r\n") {
		if len(line) > 78 {
			t.Errorf("存在过长的 MIME 行（%d 字符）: %q", len(line), line)
		}
	}
}

func TestBuildMailMessageFromFallback(t *testing.T) {
	cfg := testMailConfig()
	cfg.From = ""     // 未配置 from 时应回退到 user
	cfg.FromName = "" // 不设显示名，便于断言头部

	raw, err := buildMailMessage(cfg, []string{"user@example.com"}, "s", "b", time.Now())
	if err != nil {
		t.Fatalf("buildMailMessage 返回错误: %v", err)
	}
	if !strings.Contains(string(raw), "From: "+cfg.User+"\r\n") {
		t.Errorf("from 未回退到 user:\n%s", raw)
	}
}

func TestBuildMailMessageFromWithDisplayName(t *testing.T) {
	cfg := testMailConfig()
	cfg.FromName = "电商平台"

	raw, err := buildMailMessage(cfg, []string{"user@example.com"}, "s", "b", time.Now())
	if err != nil {
		t.Fatalf("buildMailMessage 返回错误: %v", err)
	}
	msg := string(raw)
	// 中文显示名需做 RFC 2047 编码，地址部分保持原样
	if !strings.Contains(msg, "<"+cfg.From+">") {
		t.Errorf("发件人地址丢失:\n%s", msg)
	}
	if !strings.Contains(msg, "=?UTF-8?") {
		t.Errorf("中文发件人显示名未编码:\n%s", msg)
	}
}

func TestBuildMailMessageRejectsInjectionAndBadInput(t *testing.T) {
	cfg := testMailConfig()
	now := time.Now()

	cases := []struct {
		name    string
		to      []string
		subject string
	}{
		{"收件人含换行", []string{"user@example.com\r\nBcc: evil@example.com"}, "主题"},
		{"主题含换行", []string{"user@example.com"}, "主题\r\nBcc: evil@example.com"},
		{"收件人非法", []string{"not-an-email"}, "主题"},
		{"收件人为空", nil, "主题"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := buildMailMessage(cfg, c.to, c.subject, "body", now); err == nil {
				t.Errorf("期望返回错误，实际通过")
			}
		})
	}
}

func TestMailConfigReady(t *testing.T) {
	// dev_mode 下不校验 SMTP 参数，因为不会真正连服务器
	dev := &MailConfig{DevMode: true}
	if err := dev.Ready(); err != nil {
		t.Errorf("dev_mode 应视为可用，实际: %v", err)
	}

	// 未开启任何开关时视为未启用
	if err := (&MailConfig{}).Ready(); err == nil {
		t.Errorf("enabled=false 且 dev_mode=false 时应不可用")
	}
}

func TestMailConfigValidate(t *testing.T) {
	if err := testMailConfig().Validate(); err != nil {
		t.Fatalf("合法配置校验失败: %v", err)
	}

	cases := []struct {
		name    string
		mutate  func(*MailConfig)
		wantSub string
	}{
		{"缺少授权码", func(c *MailConfig) { c.Pass = "" }, "email.pass"},
		{"缺少发信账号", func(c *MailConfig) { c.User = "" }, "email.user"},
		{"端口非法", func(c *MailConfig) { c.Port = 0 }, "email.port"},
		{"发件人非法", func(c *MailConfig) { c.From = "not-an-email" }, "email.from"},
		{"超时非正数", func(c *MailConfig) { c.Timeout = 0 }, "email.timeout"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := testMailConfig()
			c.mutate(cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatalf("期望校验失败，实际通过")
			}
			if !strings.Contains(err.Error(), c.wantSub) {
				t.Errorf("错误信息应包含 %q，实际: %v", c.wantSub, err)
			}
		})
	}
}

func TestSendMailRejectsUnreadyConfig(t *testing.T) {
	// 配置不可用时必须在连接 SMTP 之前就失败，避免把请求挂在网络超时上
	cfg := &MailConfig{Enabled: true}
	if err := SendMail(context.Background(), cfg, []string{"user@example.com"}, "s", "b"); err == nil {
		t.Errorf("配置不完整时应返回错误")
	}
}
