package verifycode

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"gf-eshop/internal/errcode"
)

// fakeSender 记录投递内容，用于验证核心与渠道的解耦。
type fakeSender struct {
	ch       Channel
	name     string
	ready    error
	sent     []string
	lastCode string
	sendErr  error
}

func (f *fakeSender) Channel() Channel { return f.ch }
func (f *fakeSender) Name() string     { return f.name }
func (f *fakeSender) Ready() error     { return f.ready }
func (f *fakeSender) Send(_ context.Context, to, code string, _ time.Duration) error {
	f.sent = append(f.sent, to)
	f.lastCode = code
	return f.sendErr
}

func testConfig() *Config {
	return &Config{
		CodeExpire:   2 * time.Minute,
		ResendAfter:  60 * time.Second,
		DailyLimit:   10,
		IPDailyLimit: 30,
		MaxAttempts:  5,
	}
}

// TestSendRejectsNotReadyChannel：渠道配置不全时必须直接失败，且错误文案带上渠道名。
func TestSendRejectsNotReadyChannel(t *testing.T) {
	sender := &fakeSender{ch: ChannelSMS, name: "短信", ready: errors.New("sms.enabled 为 false")}
	svc := New(testConfig(), sender)

	_, err := svc.Send(context.Background(), ChannelSMS, "login", "13800000000", DeliverAlways)
	if err == nil {
		t.Fatal("渠道不可用时 Send 应返回错误")
	}
	if errcode.CodeOf(err) != errcode.CodeVerifyChannelNotReady {
		t.Errorf("错误码应为 %d，实际 %d", errcode.CodeVerifyChannelNotReady, errcode.CodeOf(err))
	}
	// 关键：文案必须带渠道名，便于前端直接展示
	if !strings.Contains(err.Error(), "短信服务未启用或配置不完整") {
		t.Errorf("错误文案应带渠道名，实际：%v", err)
	}
	if len(sender.sent) != 0 {
		t.Errorf("渠道不可用时不应尝试投递，实际投递了 %d 次", len(sender.sent))
	}
}

// TestSendUnknownChannel：未注册的渠道要给出明确错误而不是空指针。
func TestSendUnknownChannel(t *testing.T) {
	svc := New(testConfig())
	if _, err := svc.Send(context.Background(), Channel("wechat"), "login", "x", DeliverAlways); err == nil {
		t.Fatal("未注册渠道应返回错误")
	}
}

// TestConsumeRejectsEmptyInput：空验证码/空收件人要在访问 Redis 之前就返回，避免无谓往返。
func TestConsumeRejectsEmptyInput(t *testing.T) {
	svc := New(testConfig(), &fakeSender{ch: ChannelEmail, name: "邮件"})

	for _, c := range []struct{ target, code string }{
		{"", "123456"},
		{"user@example.com", ""},
	} {
		err := svc.Consume(context.Background(), ChannelEmail, "login", c.target, c.code)
		if errcode.CodeOf(err) != errcode.CodeVerifyCodeInvalid {
			t.Errorf("target=%q code=%q 应返回 %d，实际 %v",
				c.target, c.code, errcode.CodeVerifyCodeInvalid, err)
		}
	}
}

// TestSenderLookup：按渠道取发送器，未注册渠道返回 nil。
func TestSenderLookup(t *testing.T) {
	email := &fakeSender{ch: ChannelEmail, name: "邮件"}
	sms := &fakeSender{ch: ChannelSMS, name: "短信"}
	svc := New(testConfig(), email, sms)

	if svc.Sender(ChannelEmail) != email {
		t.Error("邮件渠道发送器未正确注册")
	}
	if svc.Sender(ChannelSMS) != sms {
		t.Error("短信渠道发送器未正确注册")
	}
	if svc.Sender(Channel("wechat")) != nil {
		t.Error("未注册渠道应返回 nil")
	}
}

// TestKeyNamespaceIsolation：键空间必须带渠道维度，否则同一收件人在不同渠道上会互相覆盖。
func TestKeyNamespaceIsolation(t *testing.T) {
	codeEmail := codeKey(ChannelEmail, "login", "a@b.com")
	codeSMS := codeKey(ChannelSMS, "login", "a@b.com")
	if codeEmail == codeSMS {
		t.Error("不同渠道的验证码 key 不应相同")
	}
	if codeKey(ChannelEmail, "login", "a@b.com") == codeKey(ChannelEmail, "reset", "a@b.com") {
		t.Error("不同场景的验证码 key 不应相同")
	}
	if cooldownKey(ChannelEmail, "login", "a@b.com") == codeKey(ChannelEmail, "login", "a@b.com") {
		t.Error("冷却 key 不应与验证码 key 相同")
	}
	if dailyKey(ChannelEmail, "a@b.com") == ipDailyKey(ChannelEmail, "a@b.com") {
		t.Error("收件人维度与 IP 维度的计数 key 不应相同")
	}
	for _, k := range []string{
		codeKey(ChannelEmail, "login", "a@b.com"),
		cooldownKey(ChannelEmail, "login", "a@b.com"),
		dailyKey(ChannelEmail, "a@b.com"),
		ipDailyKey(ChannelEmail, "1.2.3.4"),
	} {
		if !strings.HasPrefix(k, "verify:code:") {
			t.Errorf("键应统一以 verify:code: 开头，实际 %s", k)
		}
	}
}

// TestNormalizeTarget：邮箱归一化（域名大小写不敏感），其它渠道只去空格。
func TestNormalizeTarget(t *testing.T) {
	if got := normalizeTarget(ChannelEmail, "  User@Example.COM "); got != "user@example.com" {
		t.Errorf("邮箱应归一化为小写，实际 %q", got)
	}
	if got := normalizeTarget(ChannelSMS, " 13800000000 "); got != "13800000000" {
		t.Errorf("手机号应仅去空格，实际 %q", got)
	}
}

func TestFormatDuration(t *testing.T) {
	cases := map[time.Duration]string{
		45 * time.Second: "45 秒",
		time.Minute:      "1 分钟",
		2 * time.Minute:  "2 分钟",
		90 * time.Second: "1 分钟",
	}
	for d, want := range cases {
		if got := FormatDuration(d); got != want {
			t.Errorf("FormatDuration(%s) = %q，期望 %q", d, got, want)
		}
	}
}

// TestBuildEmailBody：正文必须包含验证码与有效期，且只含数字验证码（无注入面）。
func TestBuildEmailBody(t *testing.T) {
	body := buildEmailBody("123456", 2*time.Minute)
	if !strings.Contains(body, "123456") {
		t.Error("邮件正文缺少验证码")
	}
	if !strings.Contains(body, "2 分钟") {
		t.Error("邮件正文缺少有效期")
	}
	if !strings.Contains(body, "text/html") && !strings.Contains(body, "<div") {
		t.Error("邮件正文应为 HTML")
	}
}

// TestSMSSenderDevMode：短信渠道在 dev_mode 下视为就绪，且不外发。
func TestSMSSenderDevMode(t *testing.T) {
	s := &SMSSender{cfg: &SMSConfig{DevMode: true}, provider: logSMSProvider{}}

	if err := s.Ready(); err != nil {
		t.Fatalf("dev_mode 下短信渠道应就绪，实际：%v", err)
	}
	if err := s.Send(context.Background(), "13800000000", "123456", time.Minute); err != nil {
		t.Fatalf("dev_mode 下 Send 不应报错，实际：%v", err)
	}
	if s.Channel() != ChannelSMS || s.Name() != "短信" {
		t.Error("短信渠道的 Channel/Name 不正确")
	}
}

// TestSMSSenderNotReady：开启了渠道但没接真实 provider / 没报备签名模板时，
// 必须明确判定为不可用，绝不能假装发送成功。
func TestSMSSenderNotReady(t *testing.T) {
	// 未开启
	if err := (&SMSSender{cfg: &SMSConfig{}, provider: logSMSProvider{}}).Ready(); err == nil {
		t.Error("sms.enabled=false 时应不可用")
	}
	// 开启了但仍是日志 provider
	err := (&SMSSender{cfg: &SMSConfig{Enabled: true}, provider: logSMSProvider{}}).Ready()
	if err == nil || !strings.Contains(err.Error(), "未接入真实短信通道") {
		t.Errorf("已启用但未接真实通道时应给出明确提示，实际：%v", err)
	}
	// 接了真实 provider 但缺签名/模板
	real := &fakeSMSProvider{name: "aliyun"}
	err = (&SMSSender{cfg: &SMSConfig{Enabled: true}, provider: real}).Ready()
	if err == nil || !strings.Contains(err.Error(), "sign_name") {
		t.Errorf("缺签名/模板时应不可用，实际：%v", err)
	}
}

// TestSMSSenderWithProvider：注入真实 provider 后应走 provider 投递（这就是接短信的全部增量）。
func TestSMSSenderWithProvider(t *testing.T) {
	real := &fakeSMSProvider{name: "aliyun"}
	s := &SMSSender{
		cfg:      &SMSConfig{Enabled: true, SignName: "gf-eshop", TemplateID: "123456"},
		provider: real,
	}
	if err := s.Ready(); err != nil {
		t.Fatalf("配置齐备时应就绪，实际：%v", err)
	}
	if err := s.Send(context.Background(), "13800000000", "654321", time.Minute); err != nil {
		t.Fatalf("Send 失败：%v", err)
	}
	if real.calls != 1 || real.lastCode != "654321" || real.lastPhone != "13800000000" {
		t.Errorf("provider 未收到正确调用：%+v", real)
	}
}

type fakeSMSProvider struct {
	name      string
	calls     int
	lastPhone string
	lastCode  string
}

func (f *fakeSMSProvider) Name() string { return f.name }

func (f *fakeSMSProvider) Send(_ context.Context, _ *SMSConfig, phone, code string, _ time.Duration) error {
	f.calls++
	f.lastPhone = phone
	f.lastCode = code
	return nil
}
