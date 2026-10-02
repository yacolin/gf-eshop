package utility

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// VerifyCodeLength 是邮箱验证码的默认位数。
const VerifyCodeLength = 6

// verifyCodeHashPrefix 作为哈希前缀，避免与其他用途的摘要混淆（域分隔）。
const verifyCodeHashPrefix = "gf-eshop:email-code:v1"

// GenerateNumericCode 生成指定长度的纯数字验证码。
//
// 使用 crypto/rand 而非 math/rand：验证码是身份凭证，
// math/rand 的序列在观察到足够样本后可被预测，且 Go 1.20+ 已废弃全局 Seed。
func GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("验证码长度必须为正数，当前为 %d", length)
	}
	buf := make([]byte, length)
	ten := big.NewInt(10)
	for i := range buf {
		n, err := rand.Int(rand.Reader, ten)
		if err != nil {
			return "", fmt.Errorf("生成验证码失败: %w", err)
		}
		buf[i] = byte('0' + n.Int64())
	}
	return string(buf), nil
}

// NormalizeEmail 统一邮箱写法：去空格 + 转小写。
//
// 邮箱域名大小写不敏感，本项目 MySQL 默认排序规则同样大小写不敏感，
// 统一归一化可以让「验证码 key」与「用户查询」保持一致。
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// HashVerifyCode 计算验证码的存储指纹。
//
// Redis 中只落指纹不落明文，即使缓存内容被读到也无法直接冒用。
// 指纹绑定 channel + scene + target（邮箱/手机号），保证同一串验证码
// 不能跨渠道、跨场景、跨收件人复用。
func HashVerifyCode(channel, scene, target, code string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		verifyCodeHashPrefix, channel, scene, NormalizeEmail(target), code,
	}, ":")))
	return hex.EncodeToString(sum[:])
}

// ConstantTimeEqual 以常量时间比较两个字符串，避免通过响应耗时逐位猜解验证码。
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
