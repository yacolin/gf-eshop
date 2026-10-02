// Package utility 的邮件发送能力：仅依赖标准库 net/smtp，不引入第三方发信库。
//
// 设计要点：
//  1. 配置全部来自 config.yaml 的 email 段，读取时显式带默认值，
//     避免 g.Cfg().Get 在缺键时返回 nil 指针导致 panic；
//  2. 465 走隐式 TLS，587/25 交给 net/smtp 自动 STARTTLS；
//  3. 全程受 timeout 约束（连接 + 读写共用 deadline），防止半开连接挂死请求；
//  4. smtp.PlainAuth 本身拒绝在明文连接上发送凭据，因此不会静默降级泄露授权码。
package utility

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

// MailConfig 描述 SMTP 发信与验证码风控所需的完整配置（对应 config.yaml 的 email 段）。
type MailConfig struct {
	Enabled  bool          // 发信总开关
	DevMode  bool          // 开发模式：不真实发信，只把内容写入日志
	Host     string        // SMTP 服务器地址，例如 smtp.qq.com
	Port     int           // 465 隐式 TLS / 587 STARTTLS / 25 明文（不推荐）
	TLSMode  string        // TLS 模式：auto（默认，按端口判断）/ implicit / starttls
	User     string        // SMTP 登录账号
	Pass     string        // SMTP 授权码（不是邮箱登录密码）
	From     string        // 发件人地址，留空回退为 User
	FromName string        // 发件人显示名
	HeloName string        // EHLO/HELO 使用的域名，留空取 from 的域名
	Subject  string        // 验证码邮件主题
	Timeout  time.Duration // 单次发信超时
}

// ImplicitTLSPort 是使用「连接即 TLS」的 SMTP 端口，其余端口按 STARTTLS 处理。
const ImplicitTLSPort = 465

// TLS 模式取值
const (
	TLSModeAuto     = "auto"     // 按端口判断：465 隐式 TLS，其余 STARTTLS
	TLSModeImplicit = "implicit" // 强制隐式 TLS（部分服务商在非 465 端口也要求）
	TLSModeStartTLS = "starttls" // 强制 STARTTLS
)

// mailTLSConfig 构造隐式 TLS 的客户端配置。
// 声明为变量是为了在测试中注入自签 CA，生产行为与直接构造 tls.Config 完全一致。
var mailTLSConfig = func(host string) *tls.Config {
	return &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	}
}

// useImplicitTLS 判断是否走「连接即 TLS」。
func (c *MailConfig) useImplicitTLS() bool {
	switch c.TLSMode {
	case TLSModeImplicit:
		return true
	case TLSModeStartTLS:
		return false
	default:
		return c.Port == ImplicitTLSPort
	}
}

// resolveHeloName 返回 EHLO/HELO 使用的域名。
//
// net/smtp 的默认值是 "localhost"，而接收方（尤其是国内邮箱）会把
// "EHLO localhost" 视为明显的垃圾邮件特征，因此这里显式传发件域名。
func (c *MailConfig) resolveHeloName() string {
	if c.HeloName != "" {
		return c.HeloName
	}
	if at := strings.LastIndex(c.From, "@"); at >= 0 && at+1 < len(c.From) {
		return c.From[at+1:]
	}
	return "localhost"
}

const (
	defaultMailTimeout = 10 * time.Second
	defaultMailSubject = "邮箱验证码"

	// mailSendAttempts 单次发信的尝试次数（含首次）。
	mailSendAttempts = 2
	// mailRetryBackoff 重试前的等待时间。
	mailRetryBackoff = time.Second
)

// LoadMailConfig 读取 email 段配置。
//
// 所有键都带默认值：配置缺失只会得到零值/默认值，不会 panic，
// 是否「可发信」由 Ready()/Validate() 判断。
func LoadMailConfig(ctx context.Context) (*MailConfig, error) {
	cfg := &MailConfig{}
	var err error

	if cfg.Enabled, err = mailConfigBool(ctx, "email.enabled", false); err != nil {
		return nil, err
	}
	if cfg.DevMode, err = mailConfigBool(ctx, "email.dev_mode", false); err != nil {
		return nil, err
	}
	if cfg.Host, err = mailConfigText(ctx, "email.smtp"); err != nil {
		return nil, err
	}
	if cfg.Port, err = mailConfigInt(ctx, "email.port"); err != nil {
		return nil, err
	}
	if cfg.TLSMode, err = mailConfigText(ctx, "email.tls_mode"); err != nil {
		return nil, err
	}
	if cfg.TLSMode == "" {
		cfg.TLSMode = TLSModeAuto
	}
	if cfg.User, err = mailConfigText(ctx, "email.user"); err != nil {
		return nil, err
	}
	if cfg.Pass, err = mailConfigText(ctx, "email.pass"); err != nil {
		return nil, err
	}
	if cfg.From, err = mailConfigText(ctx, "email.from"); err != nil {
		return nil, err
	}
	if cfg.FromName, err = mailConfigText(ctx, "email.from_name"); err != nil {
		return nil, err
	}
	if cfg.HeloName, err = mailConfigText(ctx, "email.helo_name"); err != nil {
		return nil, err
	}
	if cfg.Subject, err = mailConfigText(ctx, "email.subject"); err != nil {
		return nil, err
	}
	if cfg.Subject == "" {
		cfg.Subject = defaultMailSubject
	}
	if cfg.From == "" {
		cfg.From = cfg.User
	}
	if cfg.Timeout, err = mailConfigDuration(ctx, "email.timeout", defaultMailTimeout); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Ready 判断当前配置能否发信。dev_mode 下不校验 SMTP 参数，因为此时不会真正连服务器。
func (c *MailConfig) Ready() error {
	if c == nil {
		return errors.New("邮件配置为空")
	}
	if c.DevMode {
		return nil
	}
	if !c.Enabled {
		return errors.New("email.enabled 为 false，邮件服务未启用")
	}
	return c.Validate()
}

// Validate 校验真实发信所需的 SMTP 参数是否齐备。
func (c *MailConfig) Validate() error {
	if c == nil {
		return errors.New("邮件配置为空")
	}
	var missing []string
	if c.Host == "" {
		missing = append(missing, "email.smtp")
	}
	if c.User == "" {
		missing = append(missing, "email.user")
	}
	if c.Pass == "" {
		missing = append(missing, "email.pass")
	}
	if c.From == "" {
		missing = append(missing, "email.from")
	}
	if len(missing) > 0 {
		return gerror.Newf("邮件配置缺失：%s（注意 email.pass 需填邮箱授权码，不是登录密码）", strings.Join(missing, ", "))
	}
	if c.Port <= 0 || c.Port > 65535 {
		return gerror.Newf("email.port 非法：%d", c.Port)
	}
	if _, err := netmail.ParseAddress(c.From); err != nil {
		return gerror.Wrapf(err, "email.from 不是合法邮箱地址：%s", c.From)
	}
	if c.Timeout <= 0 {
		return gerror.Newf("email.timeout 必须为正数，当前 %s", c.Timeout)
	}
	switch c.TLSMode {
	case "", TLSModeAuto, TLSModeImplicit, TLSModeStartTLS:
	default:
		return gerror.Newf("email.tls_mode 非法：%q（可选 auto / implicit / starttls）", c.TLSMode)
	}
	if strings.ContainsAny(c.HeloName, "\r\n") {
		return gerror.Newf("email.helo_name 含有非法换行符：%q", c.HeloName)
	}
	return nil
}

// SendMail 发送一封 HTML 邮件。
//
// 发信前会校验配置；失败会按 mailSendAttempts 重试，最终仍失败则返回错误，
// 由调用方决定是否回滚风控计数与验证码。
func SendMail(ctx context.Context, cfg *MailConfig, to []string, subject, htmlBody string) error {
	if err := cfg.Ready(); err != nil {
		return err
	}
	if cfg.DevMode {
		// 开发模式：不连接 SMTP，把内容打到日志，方便本地联调
		g.Log().Warningf(ctx, "[email.dev_mode] 邮件未真实发送 to=%v subject=%q body=%s", to, subject, htmlBody)
		return nil
	}

	raw, err := buildMailMessage(cfg, to, subject, htmlBody, time.Now())
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= mailSendAttempts; attempt++ {
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(mailRetryBackoff):
			}
		}
		if lastErr = sendMailOnce(cfg, cfg.From, to, raw); lastErr == nil {
			return nil
		}
		g.Log().Warningf(ctx, "邮件发送失败（第 %d/%d 次）：%v", attempt, mailSendAttempts, lastErr)
	}
	return lastErr
}

// sendMailOnce 完成一次完整的 SMTP 会话。
func sendMailOnce(cfg *MailConfig, from string, to []string, raw []byte) error {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	dialer := &net.Dialer{Timeout: cfg.Timeout}

	var (
		conn net.Conn
		err  error
	)
	if cfg.useImplicitTLS() {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, mailTLSConfig(cfg.Host))
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return gerror.Wrapf(err, "连接 SMTP 服务器 %s 失败", addr)
	}
	// SMTP 协议没有逐命令超时，用连接 deadline 兜底，保证整体不超时 cfg.Timeout
	_ = conn.SetDeadline(time.Now().Add(cfg.Timeout))
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return gerror.Wrap(err, "SMTP 握手失败")
	}
	defer client.Close()

	// 必须显式指定 EHLO 域名：net/smtp 默认发 "EHLO localhost"，
	// 接收方（尤其国内邮箱）会把 localhost 当作垃圾邮件特征。此调用必须早于其它命令。
	heloName := cfg.resolveHeloName()
	if err = client.Hello(heloName); err != nil {
		return gerror.Wrapf(err, "SMTP EHLO %s 失败", heloName)
	}

	// 注意：Client API 不会自动 STARTTLS（只有包级 smtp.SendMail 会），
	// 必须显式升级。否则 587 端口上 smtp.PlainAuth 会因「明文连接」拒绝发送凭据，
	// 表现为认证失败——即 587 完全发不出去。
	if !cfg.useImplicitTLS() {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(mailTLSConfig(cfg.Host)); err != nil {
				return gerror.Wrap(err, "SMTP STARTTLS 升级失败")
			}
		} else if cfg.TLSMode == TLSModeStartTLS {
			// 显式要求 STARTTLS 但服务端不支持：宁可失败，也不能明文发凭据
			return gerror.Newf("SMTP 服务器 %s 未宣告 STARTTLS，按配置拒绝在明文连接上继续", cfg.Host)
		}
	}

	// 凭据只会在 TLS 连接上发送（smtp.PlainAuth 自身保证，明文连接直接报错）
	if err = client.Auth(smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)); err != nil {
		return gerror.Wrap(err, "SMTP 认证失败（请确认使用的是邮箱授权码而非登录密码）")
	}
	if err = client.Mail(from); err != nil {
		return gerror.Wrapf(err, "SMTP MAIL FROM %s 被拒绝", from)
	}
	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			return gerror.Wrapf(err, "SMTP RCPT TO %s 被拒绝", rcpt)
		}
	}
	writer, err := client.Data()
	if err != nil {
		return gerror.Wrap(err, "SMTP DATA 失败")
	}
	if _, err = writer.Write(raw); err != nil {
		return gerror.Wrap(err, "写入邮件正文失败")
	}
	if err = writer.Close(); err != nil {
		return gerror.Wrap(err, "提交邮件失败")
	}
	return client.Quit()
}

// buildMailMessage 拼装 RFC 5322 邮件报文（纯函数，便于单测）。
//
// 主题用 RFC 2047 编码以支持中文，正文用 base64 传输以支持 UTF-8 与超长行。
func buildMailMessage(cfg *MailConfig, to []string, subject, htmlBody string, now time.Time) ([]byte, error) {
	if len(to) == 0 {
		return nil, errors.New("收件人不能为空")
	}
	from := cfg.From
	if from == "" {
		from = cfg.User
	}
	// 防止头部注入：任何进入头部的值都不允许包含 CR/LF
	if err := checkMailHeaderValue("发件人", from); err != nil {
		return nil, err
	}
	if err := checkMailHeaderValue("主题", subject); err != nil {
		return nil, err
	}
	rcpts := make([]string, 0, len(to))
	for _, addr := range to {
		if err := checkMailHeaderValue("收件人", addr); err != nil {
			return nil, err
		}
		if _, err := netmail.ParseAddress(addr); err != nil {
			return nil, gerror.Wrapf(err, "收件人地址非法：%s", addr)
		}
		rcpts = append(rcpts, addr)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", formatMailFrom(from, cfg.FromName))
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(rcpts, ", "))
	fmt.Fprintf(&buf, "Subject: %s\r\n", encodeMailHeader(subject, len("Subject: ")))
	fmt.Fprintf(&buf, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(&buf, "Message-ID: %s\r\n", newMailMessageID(from))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(wrapBase64(base64.StdEncoding.EncodeToString([]byte(htmlBody))))
	return buf.Bytes(), nil
}

func formatMailFrom(from, fromName string) string {
	if fromName == "" {
		return from
	}
	return fmt.Sprintf("%s <%s>", encodeMailHeader(fromName, len("From: ")), from)
}

// mailHeaderMaxLineLen 是 MIME 头部建议的最大行宽（RFC 5322 推荐 78）。
const mailHeaderMaxLineLen = 78

// mailEncodedWordMaxLen 是单个 encoded-word 的长度上限（RFC 2047 规定 75）。
const mailEncodedWordMaxLen = 75

// mailMinUsefulWordLen 是「值得开启一个新 encoded-word」的最小长度。
// 一个 CJK 字符 Q 编码后为 "=?UTF-8?q?" + 9 + "?=" = 21 字符；
// 行尾余量小于该值时直接折行，避免在行尾留下 1~2 字符的碎片词。
const mailMinUsefulWordLen = 21

// encodeMailHeader 按 RFC 2047 编码非 ASCII 头部；纯 ASCII 直接返回，保持可读。
//
// 中文编码后会显著变长，因此按字符切片拆成多个 encoded-word
// （每个不超过 RFC 2047 规定的 75 字符），并在接近行宽上限处折行。
// prefixLen 是头部字段名长度（如 "Subject: " 为 9），用于计算首行剩余宽度，
// 保证「字段名 + 首个 encoded-word」也不会超出行宽。
func encodeMailHeader(s string, prefixLen int) string {
	if isASCIIString(s) {
		return s
	}

	// 单个 encoded-word 的可用长度：受 RFC 2047 的 75 字符与首行剩余宽度双重约束
	maxWordLen := mailHeaderMaxLineLen - prefixLen
	if maxWordLen > mailEncodedWordMaxLen {
		maxWordLen = mailEncodedWordMaxLen
	}

	var (
		b       strings.Builder
		cur     []rune
		lineLen = prefixLen
		wrote   bool
	)
	// fold 折行：续行以空格开头（RFC 5322 的折叠空白）
	fold := func() {
		b.WriteString("\r\n ")
		lineLen = 1
	}
	// budget 返回当前行还能容纳的 encoded-word 长度
	budget := func() int {
		remaining := mailHeaderMaxLineLen - lineLen
		if wrote {
			remaining-- // 与上一个 encoded-word 之间需要一个空格
		}
		if remaining > maxWordLen {
			remaining = maxWordLen
		}
		return remaining
	}
	flush := func() {
		if len(cur) == 0 {
			return
		}
		word := mime.QEncoding.Encode("UTF-8", string(cur))
		if wrote {
			if lineLen+1+len(word) > mailHeaderMaxLineLen {
				fold()
			} else {
				b.WriteString(" ")
				lineLen++
			}
		}
		b.WriteString(word)
		lineLen += len(word)
		wrote = true
		cur = cur[:0]
	}

	for _, r := range s {
		if len(cur) > 0 && len(mime.QEncoding.Encode("UTF-8", string(cur)+string(r))) > budget() {
			flush()
		}
		if len(cur) == 0 && wrote && budget() < mailMinUsefulWordLen {
			fold()
		}
		cur = append(cur, r)
	}
	flush()
	return b.String()
}

func isASCIIString(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

func newMailMessageID(from string) string {
	domain := "localhost"
	if at := strings.LastIndex(from, "@"); at >= 0 && at+1 < len(from) {
		domain = from[at+1:]
	}
	return fmt.Sprintf("<%s@%s>", guid.S(), domain)
}

func checkMailHeaderValue(field, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return gerror.Newf("%s 含有非法换行符，可能是头部注入", field)
	}
	return nil
}

// wrapBase64 把 base64 文本按 76 字符折行，符合 MIME 行宽要求。
func wrapBase64(s string) string {
	const width = 76
	var b strings.Builder
	for len(s) > width {
		b.WriteString(s[:width])
		b.WriteString("\r\n")
		s = s[width:]
	}
	b.WriteString(s)
	b.WriteString("\r\n")
	return b.String()
}

// ---------------------------------------------------------------------------
// 配置读取辅助：全部带默认值，缺键时拿到默认值而不是 nil 指针
// ---------------------------------------------------------------------------

func mailConfigText(ctx context.Context, key string) (string, error) {
	v, err := g.Cfg().Get(ctx, key, "")
	if err != nil {
		return "", gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return strings.TrimSpace(v.String()), nil
}

func mailConfigInt(ctx context.Context, key string) (int, error) {
	return mailConfigIntDefault(ctx, key, 0)
}

func mailConfigIntDefault(ctx context.Context, key string, def int) (int, error) {
	v, err := g.Cfg().Get(ctx, key, def)
	if err != nil {
		return 0, gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return v.Int(), nil
}

func mailConfigBool(ctx context.Context, key string, def bool) (bool, error) {
	v, err := g.Cfg().Get(ctx, key, def)
	if err != nil {
		return false, gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	return v.Bool(), nil
}

// mailConfigDuration 支持 "10s" / "2m" 这类字符串写法，也兼容纯数字（按秒处理）。
func mailConfigDuration(ctx context.Context, key string, def time.Duration) (time.Duration, error) {
	v, err := g.Cfg().Get(ctx, key, def.String())
	if err != nil {
		return 0, gerror.Wrapf(err, "读取配置项 %s 失败", key)
	}
	raw := strings.TrimSpace(v.String())
	if raw == "" {
		return def, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		// 兼容 120 这种纯秒数写法
		if secs, convErr := strconv.Atoi(raw); convErr == nil {
			return time.Duration(secs) * time.Second, nil
		}
		return 0, gerror.Wrapf(err, "配置项 %s 不是合法时长（示例：10s、2m）：%s", key, raw)
	}
	return d, nil
}
