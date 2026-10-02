package utility

import (
	"strings"
	"testing"
)

func TestGenerateNumericCode(t *testing.T) {
	for i := 0; i < 200; i++ {
		code, err := GenerateNumericCode(VerifyCodeLength)
		if err != nil {
			t.Fatalf("生成验证码失败: %v", err)
		}
		if len(code) != VerifyCodeLength {
			t.Fatalf("长度应为 %d，实际 %d（code=%q）", VerifyCodeLength, len(code), code)
		}
		for _, r := range code {
			if r < '0' || r > '9' {
				t.Fatalf("验证码应只含数字，实际 %q", code)
			}
		}
	}
}

func TestGenerateNumericCodeRandomness(t *testing.T) {
	// crypto/rand 生成，200 次采样不应只有一种结果
	seen := make(map[string]struct{})
	for i := 0; i < 200; i++ {
		code, err := GenerateNumericCode(VerifyCodeLength)
		if err != nil {
			t.Fatalf("生成验证码失败: %v", err)
		}
		seen[code] = struct{}{}
	}
	if len(seen) < 2 {
		t.Errorf("验证码缺乏随机性，200 次只产生了 %d 种结果", len(seen))
	}
}

func TestGenerateNumericCodeRejectsBadLength(t *testing.T) {
	for _, n := range []int{0, -1} {
		if _, err := GenerateNumericCode(n); err == nil {
			t.Errorf("长度 %d 应返回错误", n)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	cases := map[string]string{
		"  User@Example.COM ": "user@example.com",
		"user@example.com":    "user@example.com",
		"":                    "",
	}
	for in, want := range cases {
		if got := NormalizeEmail(in); got != want {
			t.Errorf("NormalizeEmail(%q) = %q，期望 %q", in, got, want)
		}
	}
}

func TestHashVerifyCode(t *testing.T) {
	base := HashVerifyCode("email", "login", "user@example.com", "123456")

	if base != HashVerifyCode("email", "login", "user@example.com", "123456") {
		t.Errorf("相同输入的指纹应一致")
	}
	// 指纹需绑定 channel + scene + target，保证验证码不能跨渠道/场景/收件人复用
	if base == HashVerifyCode("sms", "login", "user@example.com", "123456") {
		t.Errorf("不同渠道的指纹不应相同")
	}
	if base == HashVerifyCode("email", "register", "user@example.com", "123456") {
		t.Errorf("不同 scene 的指纹不应相同")
	}
	if base == HashVerifyCode("email", "login", "other@example.com", "123456") {
		t.Errorf("不同收件人的指纹不应相同")
	}
	if base == HashVerifyCode("email", "login", "user@example.com", "654321") {
		t.Errorf("不同验证码的指纹不应相同")
	}
	// 指纹不含明文，且大小写/空格按 NormalizeEmail 归一
	if strings.Contains(base, "123456") {
		t.Errorf("指纹中不应包含验证码明文: %s", base)
	}
	if base != HashVerifyCode("email", "login", " USER@example.com ", "123456") {
		t.Errorf("指纹应按 NormalizeEmail 归一化收件人")
	}
}

func TestConstantTimeEqual(t *testing.T) {
	if !ConstantTimeEqual("abc123", "abc123") {
		t.Errorf("相同字符串应相等")
	}
	if ConstantTimeEqual("abc123", "abc124") {
		t.Errorf("不同字符串不应相等")
	}
	if ConstantTimeEqual("abc", "abc123") {
		t.Errorf("长度不同不应相等")
	}
	if !ConstantTimeEqual("", "") {
		t.Errorf("空字符串应相等")
	}
}

func TestEncodeMailHeaderFoldsLongChineseSubject(t *testing.T) {
	subject := "【gf-eshop】邮箱验证码，请勿转发给任何人这是一段很长的中文主题用于验证折行逻辑是否正确"
	encoded := encodeMailHeader(subject, len("Subject: "))

	// 编码结果应只含 ASCII（可安全写入 SMTP 头部），且确实做了 RFC 2047 编码
	if !isASCIIString(encoded) {
		t.Fatalf("编码结果应只含 ASCII: %s", encoded)
	}
	if !strings.Contains(encoded, "=?UTF-8?") {
		t.Fatalf("中文主题未被编码: %s", encoded)
	}
	// 每个 encoded-word 不超过 75 字符
	words := strings.Fields(strings.ReplaceAll(encoded, "\r\n ", " "))
	for _, word := range words {
		if len(word) > mailEncodedWordMaxLen {
			t.Errorf("encoded-word 超长（%d 字符）: %s", len(word), word)
		}
	}
	// 除最后一个词外不应出现碎片词（行尾余量不足时应当直接折行）
	for i, word := range words {
		if i < len(words)-1 && len(word) < mailMinUsefulWordLen {
			t.Errorf("出现碎片化 encoded-word（%d 字符）: %s", len(word), word)
		}
	}
	// 折行后每行不超过行宽上限：首行含 "Subject: " 前缀，续行只含折叠空格
	for i, line := range strings.Split(encoded, "\r\n") {
		lineLen := len(line)
		if i == 0 {
			lineLen += len("Subject: ")
		}
		if lineLen > mailHeaderMaxLineLen {
			t.Errorf("折行后仍超宽（%d 字符）: %s", lineLen, line)
		}
	}
}

func TestEncodeMailHeaderKeepsASCIIReadable(t *testing.T) {
	if got := encodeMailHeader("gf-eshop verify", len("Subject: ")); got != "gf-eshop verify" {
		t.Errorf("纯 ASCII 头部应原样返回，实际 %q", got)
	}
}
