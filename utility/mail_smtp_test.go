package utility

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"math/big"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// fakeSMTP 是一个最小可用的 SMTP 服务端，用于在无外网环境下验证客户端的真实协议行为。
//
// 它记录收到的 DATA 报文体、AUTH 命令行与 EHLO 域名，测试据此断言发信流程与报文内容。
// startTLS 为 true 时会宣告 STARTTLS 并在收到命令后升级为 TLS（模拟 587）。
type fakeSMTP struct {
	listener net.Listener
	messages chan string
	authCmds chan string
	heloArgs chan string
	tlsUsed  chan bool
	startTLS bool
	cert     tls.Certificate
	upgraded chan bool
}

// newFakeSMTP 启动假 SMTP 服务；useTLS 为 true 时监听隐式 TLS（模拟 465）。
func newFakeSMTP(t *testing.T, useTLS bool, cert tls.Certificate) *fakeSMTP {
	t.Helper()
	return newFakeSMTPWithOptions(t, useTLS, cert, false)
}

// newFakeSMTPWithOptions 额外支持「明文监听 + 宣告 STARTTLS」的 587 形态。
func newFakeSMTPWithOptions(t *testing.T, useTLS bool, cert tls.Certificate, startTLS bool) *fakeSMTP {
	t.Helper()

	var (
		ln  net.Listener
		err error
	)
	if useTLS {
		ln, err = tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	} else {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		t.Fatalf("启动假 SMTP 监听失败: %v", err)
	}

	f := &fakeSMTP{
		listener: ln,
		messages: make(chan string, 1),
		authCmds: make(chan string, 1),
		heloArgs: make(chan string, 4),
		tlsUsed:  make(chan bool, 1),
		startTLS: startTLS,
		cert:     cert,
		upgraded: make(chan bool, 1),
	}
	go func() {
		// 循环 Accept：SendMail 失败会重试，第二次连接必须也有人处理
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go f.serve(conn)
		}
	}()
	t.Cleanup(func() { _ = ln.Close() })
	return f
}

func (f *fakeSMTP) port() int {
	_, portStr, _ := net.SplitHostPort(f.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func (f *fakeSMTP) serve(conn net.Conn) {
	defer conn.Close()

	_, isTLS := conn.(*tls.Conn)
	if isTLS {
		// 显式完成握手，确认客户端确实走了隐式 TLS
		if err := conn.(*tls.Conn).Handshake(); err != nil {
			return
		}
	}
	select {
	case f.tlsUsed <- isTLS:
	default:
	}

	tlsActive := isTLS
	reader := bufio.NewReader(conn)
	write := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
	write("220 fake ESMTP ready")

	var (
		body   strings.Builder
		inData bool
	)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")

		if inData {
			if line == "." {
				inData = false
				write("250 OK queued")
				select {
				case f.messages <- body.String():
				default:
				}
				continue
			}
			body.WriteString(line)
			body.WriteString("\r\n")
			continue
		}

		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "EHLO"):
			select {
			case f.heloArgs <- line:
			default:
			}
			// 多行响应；未升级且支持 STARTTLS 时宣告之，并声明支持 PLAIN 认证
			resp := "250-fake\r\n250-AUTH PLAIN LOGIN\r\n250 OK\r\n"
			if f.startTLS && !tlsActive {
				resp = "250-fake\r\n250-STARTTLS\r\n250-AUTH PLAIN LOGIN\r\n250 OK\r\n"
			}
			_, _ = conn.Write([]byte(resp))
		case strings.HasPrefix(upper, "HELO"):
			select {
			case f.heloArgs <- line:
			default:
			}
			write("250 fake")
		case upper == "STARTTLS":
			if !f.startTLS {
				write("502 STARTTLS not supported")
				continue
			}
			write("220 Ready to start TLS")
			tlsConn := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{f.cert}})
			if err := tlsConn.Handshake(); err != nil {
				return
			}
			conn = tlsConn
			reader = bufio.NewReader(conn)
			tlsActive = true
			select {
			case f.upgraded <- true:
			default:
			}
		case strings.HasPrefix(upper, "AUTH"):
			select {
			case f.authCmds <- line:
			default:
			}
			write("235 2.7.0 Authentication successful")
		case strings.HasPrefix(upper, "MAIL FROM"), strings.HasPrefix(upper, "RCPT TO"), upper == "RSET":
			write("250 OK")
		case upper == "DATA":
			write("354 End data with <CR><LF>.<CR><LF>")
			inData = true
		case upper == "QUIT":
			write("221 Bye")
			return
		default:
			write("250 OK")
		}
	}
}

// selfSignedCert 生成一张绑定 127.0.0.1 的自签证书，供隐式 TLS 测试使用。
func selfSignedCert(t *testing.T) tls.Certificate {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("签发自签证书失败: %v", err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// trustCert 把自签证书注入隐式 TLS 的根证书池（仅测试期间生效）。
func trustCert(t *testing.T, cert tls.Certificate) {
	t.Helper()

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("解析自签证书失败: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(leaf)

	original := mailTLSConfig
	mailTLSConfig = func(host string) *tls.Config {
		return &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
			RootCAs:    pool,
		}
	}
	t.Cleanup(func() { mailTLSConfig = original })
}

// decodeMailBody 取出报文中的正文并做 base64 解码。
func decodeMailBody(t *testing.T, msg string) string {
	t.Helper()

	idx := strings.Index(msg, "\r\n\r\n")
	if idx < 0 {
		t.Fatalf("报文缺少头部与正文分隔空行:\n%s", msg)
	}
	encoded := strings.ReplaceAll(msg[idx+4:], "\r\n", "")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("正文不是合法 base64: %v", err)
	}
	return string(decoded)
}

func TestSendMailOverPlainSMTP(t *testing.T) {
	server := newFakeSMTP(t, false, tls.Certificate{})

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port()
	cfg.TLSMode = TLSModeAuto // 服务端不宣告 STARTTLS → 明文连接，由 PlainAuth 决定是否发凭据
	cfg.Subject = "【gf-eshop】邮箱验证码"

	body := `<b>验证码 246810</b>`
	if err := SendMail(context.Background(), cfg, []string{"to@example.com"}, cfg.Subject, body); err != nil {
		t.Fatalf("SendMail 失败: %v", err)
	}

	select {
	case msg := <-server.messages:
		if !strings.Contains(msg, "Subject: =?UTF-8?") {
			t.Errorf("主题未编码:\n%s", msg)
		}
		if got := decodeMailBody(t, msg); got != body {
			t.Errorf("正文不符\n期望: %s\n实际: %s", body, got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("假 SMTP 未收到报文体")
	}

	// 127.0.0.1 属于 PlainAuth 允许的明文例外，应完成一次 PLAIN 认证
	select {
	case authLine := <-server.authCmds:
		if !strings.HasPrefix(strings.ToUpper(authLine), "AUTH PLAIN ") {
			t.Errorf("认证命令应为 AUTH PLAIN，实际: %s", authLine)
		}
		if strings.Contains(authLine, cfg.Pass) {
			t.Errorf("授权码不应以明文形式出现在命令行: %s", authLine)
		}
	case <-time.After(time.Second):
		t.Error("假 SMTP 未收到 AUTH 命令")
	}

	if used := <-server.tlsUsed; used {
		t.Error("该用例不应使用 TLS")
	}
}

// TestSendMailUsesRealHeloName 断言不会发出 "EHLO localhost"。
// net/smtp 的默认 EHLO 域名是 localhost，国内邮箱普遍把它当垃圾邮件特征，
// 因此必须显式传参（见 MailConfig.resolveHeloName）。
func TestSendMailUsesRealHeloName(t *testing.T) {
	server := newFakeSMTP(t, false, tls.Certificate{})

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port()
	cfg.TLSMode = TLSModeAuto
	cfg.From = "noreply@qq.com"
	cfg.User = "noreply@qq.com"
	cfg.HeloName = "" // 留空 → 取 from 的域名

	if err := SendMail(context.Background(), cfg, []string{"to@sina.com"}, cfg.Subject, "<b>1</b>"); err != nil {
		t.Fatalf("SendMail 失败: %v", err)
	}

	select {
	case line := <-server.heloArgs:
		if !strings.EqualFold(line, "EHLO qq.com") {
			t.Errorf("EHLO 应为发件域名 qq.com，实际: %q", line)
		}
		if strings.Contains(strings.ToLower(line), "localhost") {
			t.Errorf("不应发送 EHLO localhost，实际: %q", line)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("假 SMTP 未收到 EHLO")
	}
}

func TestSendMailHonoursExplicitHeloName(t *testing.T) {
	server := newFakeSMTP(t, false, tls.Certificate{})

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port()
	cfg.TLSMode = TLSModeAuto
	cfg.HeloName = "mail.eshop.dev"

	if err := SendMail(context.Background(), cfg, []string{"to@sina.com"}, cfg.Subject, "<b>1</b>"); err != nil {
		t.Fatalf("SendMail 失败: %v", err)
	}

	select {
	case line := <-server.heloArgs:
		if !strings.EqualFold(line, "EHLO mail.eshop.dev") {
			t.Errorf("EHLO 应为显式配置值，实际: %q", line)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("假 SMTP 未收到 EHLO")
	}
}

func TestResolveHeloName(t *testing.T) {
	cases := []struct {
		from     string
		heloName string
		want     string
	}{
		{"noreply@qq.com", "", "qq.com"},
		{"noreply@qq.com", "mail.eshop.dev", "mail.eshop.dev"},
		{"", "", "localhost"},
		{"broken", "", "localhost"},
	}
	for _, c := range cases {
		cfg := &MailConfig{From: c.from, HeloName: c.heloName}
		if got := cfg.resolveHeloName(); got != c.want {
			t.Errorf("resolveHeloName(from=%q, helo=%q) = %q，期望 %q", c.from, c.heloName, got, c.want)
		}
	}
}

// TestSendMailUpgradesStartTLS 覆盖 587 场景：服务端宣告 STARTTLS 时必须显式升级。
// 早期实现依赖「smtp.NewClient 自动 STARTTLS」的错误假设，会导致 587 上认证被拒。
func TestSendMailUpgradesStartTLS(t *testing.T) {
	cert := selfSignedCert(t)
	trustCert(t, cert)
	server := newFakeSMTPWithOptions(t, false, cert, true)

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port() // 非 465 → 走 STARTTLS 分支
	cfg.TLSMode = TLSModeAuto
	cfg.From = "noreply@qq.com"
	cfg.User = "noreply@qq.com"

	body := "<b>验证码 888888</b>"
	if err := SendMail(context.Background(), cfg, []string{"to@sina.com"}, cfg.Subject, body); err != nil {
		t.Fatalf("STARTTLS 发信失败: %v", err)
	}

	select {
	case <-server.upgraded:
	case <-time.After(2 * time.Second):
		t.Fatal("客户端未执行 STARTTLS 升级")
	}

	// 升级后仍需完成一次真正的 TLS 会话（DATA 内容可读）
	select {
	case msg := <-server.messages:
		if got := decodeMailBody(t, msg); got != body {
			t.Errorf("正文不符\n期望: %s\n实际: %s", body, got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("假 SMTP 未收到报文体")
	}

	// 至少发生过两次 EHLO：明文一次 + STARTTLS 升级后一次（RFC 3207 要求）
	if len(server.heloArgs) < 2 {
		t.Errorf("STARTTLS 后应重新 EHLO，实际只收到 %d 次", len(server.heloArgs))
	}
}

// TestSendMailRefusesPlaintextAuthWhenStartTLSForced 保证服务端不支持 STARTTLS 时
// 宁可失败，也绝不把授权码以明文发出去。
func TestSendMailRefusesPlaintextAuthWhenStartTLSForced(t *testing.T) {
	server := newFakeSMTP(t, false, tls.Certificate{}) // 不宣告 STARTTLS

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port()
	cfg.TLSMode = TLSModeStartTLS
	cfg.From = "noreply@qq.com"
	cfg.User = "noreply@qq.com"

	err := SendMail(context.Background(), cfg, []string{"to@sina.com"}, cfg.Subject, "<b>1</b>")
	if err == nil {
		t.Fatal("服务端不支持 STARTTLS 时应发送失败")
	}
	if !strings.Contains(err.Error(), "STARTTLS") {
		t.Errorf("错误信息应指明 STARTTLS 问题，实际: %v", err)
	}
	select {
	case line := <-server.authCmds:
		t.Fatalf("不应在明文连接上发送 AUTH，实际发出: %s", line)
	case <-time.After(300 * time.Millisecond):
		// 预期：没有发出任何认证命令
	}
}

func TestMailConfigRejectsHeloInjection(t *testing.T) {
	cfg := testMailConfig()
	cfg.HeloName = "qq.com\r\nMAIL FROM:<evil@x.com>"
	if err := cfg.Validate(); err == nil {
		t.Errorf("helo_name 含换行应校验失败")
	} else if !strings.Contains(err.Error(), "email.helo_name") {
		t.Errorf("错误信息应包含 email.helo_name，实际: %v", err)
	}
}

func TestSendMailOverImplicitTLS(t *testing.T) {
	cert := selfSignedCert(t)
	trustCert(t, cert)
	server := newFakeSMTP(t, true, cert)

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = server.port() // 随机端口：必须靠 tls_mode=implicit 强制隐式 TLS
	cfg.TLSMode = TLSModeImplicit

	body := "<b>验证码 135790</b>"
	if err := SendMail(context.Background(), cfg, []string{"to@example.com"}, cfg.Subject, body); err != nil {
		t.Fatalf("隐式 TLS 发信失败: %v", err)
	}

	select {
	case msg := <-server.messages:
		if got := decodeMailBody(t, msg); got != body {
			t.Errorf("正文不符\n期望: %s\n实际: %s", body, got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("假 SMTP 未收到报文体")
	}

	if used := <-server.tlsUsed; !used {
		t.Error("客户端未使用隐式 TLS")
	}
}

func TestUseImplicitTLSDecision(t *testing.T) {
	cases := []struct {
		mode string
		port int
		want bool
	}{
		{TLSModeAuto, ImplicitTLSPort, true},
		{TLSModeAuto, 587, false},
		{TLSModeAuto, 25, false},
		{TLSModeImplicit, 587, true},  // 部分服务商在非 465 端口也要求隐式 TLS
		{TLSModeStartTLS, 465, false}, // 显式覆盖端口推断
		{"", ImplicitTLSPort, true},   // 未配置按 auto 处理
	}
	for _, c := range cases {
		cfg := &MailConfig{TLSMode: c.mode, Port: c.port}
		if got := cfg.useImplicitTLS(); got != c.want {
			t.Errorf("useImplicitTLS(mode=%q, port=%d) = %v，期望 %v", c.mode, c.port, got, c.want)
		}
	}
}

func TestMailConfigRejectsBadTLSMode(t *testing.T) {
	cfg := testMailConfig()
	cfg.TLSMode = "ssl"
	if err := cfg.Validate(); err == nil {
		t.Errorf("非法 tls_mode 应校验失败")
	} else if !strings.Contains(err.Error(), "email.tls_mode") {
		t.Errorf("错误信息应包含 email.tls_mode，实际: %v", err)
	}
}

func TestSendMailRetriesThenSucceeds(t *testing.T) {
	// 第一次连接失败（端口未监听）不被吃掉：SendMail 会重试，
	// 这里通过「第一次拒绝、第二次接受」的服务端验证重试确实发生。
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听失败: %v", err)
	}
	defer ln.Close()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)

	attempts := make(chan int, 2)
	go func() {
		for i := 1; i <= 2; i++ {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			attempts <- i
			// 第一次直接断开，迫使客户端重试；第二次走正常会话
			if i == 1 {
				_ = conn.Close()
				continue
			}
			f := &fakeSMTP{
				messages: make(chan string, 1),
				authCmds: make(chan string, 1),
				heloArgs: make(chan string, 2),
				tlsUsed:  make(chan bool, 1),
			}
			f.serve(conn)
			return
		}
	}()

	cfg := testMailConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = port
	cfg.TLSMode = TLSModeAuto
	cfg.Timeout = 3 * time.Second

	if err := SendMail(context.Background(), cfg, []string{"to@example.com"}, cfg.Subject, "<b>1</b>"); err != nil {
		t.Fatalf("重试后应发送成功，实际: %v", err)
	}

	first := <-attempts
	second := <-attempts
	if first != 1 || second != 2 {
		t.Errorf("期望两次连接尝试，实际顺序: %d, %d", first, second)
	}
}
