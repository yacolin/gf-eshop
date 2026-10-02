// Package verifycode 提供与渠道无关的验证码基础设施：生成、指纹存储、
// 频次风控与校验。邮件 / 短信只是「投递渠道」，通过 Sender 接口接入。
//
// 为什么要独立成包（与 internal/search 同类定位）：
//   - user_auth 只是第一个使用者，后续 staff 登录、商家入驻、换绑手机号都会用到；
//   - 渠道无关的部分（验证码生命周期、防爆破、防枚举、限流）逻辑复杂且容易出现
//     安全缺陷，集中一处才能只写一次、只审一次；
//   - 接新渠道的增量应该只有「实现 Sender」，而不是再抄一遍风控。
package verifycode

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"gf-eshop/internal/errcode"
	"gf-eshop/utility"
)

// Channel 是验证码的投递渠道。
type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
)

// Sender 是渠道发送器：接入新渠道只需实现该接口并注册到 Service。
type Sender interface {
	// Channel 返回该发送器所属渠道
	Channel() Channel
	// Name 返回面向用户的中文渠道名（用于错误文案，如「邮件」「短信」）
	Name() string
	// Ready 判断该渠道当前是否可用；不可用时 Send 会返回 NotReady 类错误
	Ready() error
	// Send 投递验证码。开发模式（不真正外发）由实现方自行处理。
	Send(ctx context.Context, to, code string, expire time.Duration) error
}

// DeliverPolicy 决定「这条验证码是否真的需要投递」。
//
// 业务侧用它注入存在性检查：例如邮箱未注册时返回 deliver=false。
// 关键点：deliver=false 时服务**仍会照常执行**额度扣减、重发冷却、写入验证码，
// 因此「不投递」与「已投递」两条路径的响应和副作用完全一致 —— 这是防邮箱枚举的前提。
type DeliverPolicy func(ctx context.Context) (deliver bool, err error)

// DeliverAlways 是默认策略：总是投递。
func DeliverAlways(context.Context) (bool, error) { return true, nil }

// SendResult 是发送成功的回执。
type SendResult struct {
	ExpireIn    int // 验证码有效期（秒）
	ResendAfter int // 多少秒后可重新发送
}

// Service 是渠道无关的验证码服务。
type Service struct {
	cfg     *Config
	senders map[Channel]Sender
}

// New 组装验证码服务。
func New(cfg *Config, senders ...Sender) *Service {
	m := make(map[Channel]Sender, len(senders))
	for _, s := range senders {
		m[s.Channel()] = s
	}
	return &Service{cfg: cfg, senders: m}
}

// Sender 返回指定渠道的发送器（可能为 nil）。
func (s *Service) Sender(ch Channel) Sender {
	if s == nil {
		return nil
	}
	return s.senders[ch]
}

// Config 返回风控配置。
func (s *Service) Config() *Config { return s.cfg }

// ---------------------------------------------------------------------------
// 发送
// ---------------------------------------------------------------------------

// Send 生成并（按策略）投递一条验证码。
//
// 执行顺序（顺序本身有安全含义，勿随意调整）：
//  1. 渠道就绪检查 —— 配置不全时直接失败，不消耗任何额度；
//  2. 业务策略 —— 存在性检查，可返回业务错误（如邮箱已绑定）；此时同样不消耗额度；
//  3. 抢占重发冷却 —— 被冷却拒绝的请求不消耗每日额度；
//  4. 扣每日额度（邮箱/手机号维度 + IP 维度），额度不足则释放冷却；
//  5. 写入验证码指纹；
//  6. 按策略决定是否真正投递；投递失败则整体回滚（删码 + 释放冷却 + 归还额度）。
func (s *Service) Send(ctx context.Context, ch Channel, scene, target string, policy DeliverPolicy) (*SendResult, error) {
	if s == nil || s.cfg == nil {
		return nil, errcode.ErrVerifyChannelNotReady
	}
	sender := s.senders[ch]
	if sender == nil {
		return nil, errcode.Newf(errcode.CodeVerifyChannelNotReady, "不支持的验证码渠道：%s", ch)
	}
	if err := sender.Ready(); err != nil {
		g.Log().Warningf(ctx, "验证码渠道 %s 不可用：%v", ch, err)
		return nil, errcode.Newf(errcode.CodeVerifyChannelNotReady,
			"%s服务未启用或配置不完整，请联系管理员", sender.Name())
	}

	target = normalizeTarget(ch, target)
	if target == "" || strings.TrimSpace(scene) == "" {
		return nil, errcode.ErrInvalidParams
	}
	scene = strings.TrimSpace(scene)

	deliver := true
	if policy != nil {
		var err error
		if deliver, err = policy(ctx); err != nil {
			return nil, err
		}
	}

	// 抢占重发冷却（SET NX 原子操作，避免并发重复发送）；被拒绝的请求不消耗额度
	cdKey := cooldownKey(ch, scene, target)
	claimed, err := g.Redis().Do(ctx, "SET", cdKey, "1", "NX", "EX", int(s.cfg.ResendAfter.Seconds()))
	if err != nil {
		return nil, err
	}
	// 注意：SET NX 未抢到时 GoFrame 返回的是空值而不是 nil 指针，因此判断 "OK"，不能用 IsNil()
	if claimed == nil || claimed.String() != "OK" {
		return nil, errcode.Newf(
			errcode.CodeVerifyCodeTooFrequent,
			"验证码发送过于频繁，请 %d 秒后再试",
			cooldownRemain(ctx, cdKey, s.cfg.ResendAfter),
		)
	}
	releaseCooldown := func() {
		if _, delErr := g.Redis().Do(ctx, "DEL", cdKey); delErr != nil {
			g.Log().Warningf(ctx, "释放验证码重发冷却失败：%v", delErr)
		}
	}

	ip := clientIp(ctx)
	if err = s.consumeQuota(ctx, ch, target, ip); err != nil {
		releaseCooldown()
		return nil, err
	}

	code, err := utility.GenerateNumericCode(utility.VerifyCodeLength)
	if err != nil {
		s.rollback(ctx, ch, scene, target, ip)
		return nil, err
	}
	if err = s.storeCode(ctx, ch, scene, target, code); err != nil {
		s.rollback(ctx, ch, scene, target, ip)
		return nil, err
	}

	result := &SendResult{
		ExpireIn:    int(s.cfg.CodeExpire.Seconds()),
		ResendAfter: int(s.cfg.ResendAfter.Seconds()),
	}

	// 策略判定不投递：副作用已经全部落地，与「已投递」路径不可区分
	if !deliver {
		g.Log().Infof(ctx, "验证码按策略跳过投递（channel=%s scene=%s）", ch, scene)
		return result, nil
	}

	if err = sender.Send(ctx, target, code, s.cfg.CodeExpire); err != nil {
		g.Log().Errorf(ctx, "验证码投递失败（channel=%s scene=%s）：%v", ch, scene, err)
		s.rollback(ctx, ch, scene, target, ip)
		return nil, errcode.Newf(errcode.CodeVerifySendFailed,
			"验证码%s发送失败，请稍后重试", sender.Name())
	}

	g.Log().Infof(ctx, "验证码已发送（channel=%s scene=%s 有效期=%s）", ch, scene, s.cfg.CodeExpire)
	return result, nil
}

// consumeQuota 扣减目标维度与 IP 维度的每日额度，任一超限即返回错误。
func (s *Service) consumeQuota(ctx context.Context, ch Channel, target, ip string) error {
	if ip != "" && s.cfg.IPDailyLimit > 0 {
		key := ipDailyKey(ch, ip)
		n, err := incrCounter(ctx, key)
		if err != nil {
			return err
		}
		if n > int64(s.cfg.IPDailyLimit) {
			decCounter(ctx, key)
			return errcode.ErrVerifyCodeDailyLimit
		}
	}
	if s.cfg.DailyLimit > 0 {
		key := dailyKey(ch, target)
		n, err := incrCounter(ctx, key)
		if err != nil {
			return err
		}
		if n > int64(s.cfg.DailyLimit) {
			decCounter(ctx, key)
			return errcode.ErrVerifyCodeDailyLimit
		}
	}
	return nil
}

// rollback 投递失败后的回滚：删除验证码、释放重发冷却、归还每日额度。
func (s *Service) rollback(ctx context.Context, ch Channel, scene, target, ip string) {
	if _, err := g.Redis().Do(ctx, "DEL", codeKey(ch, scene, target)); err != nil {
		g.Log().Warningf(ctx, "回滚验证码失败：%v", err)
	}
	if _, err := g.Redis().Do(ctx, "DEL", cooldownKey(ch, scene, target)); err != nil {
		g.Log().Warningf(ctx, "释放验证码重发冷却失败：%v", err)
	}
	decCounter(ctx, dailyKey(ch, target))
	if ip != "" {
		decCounter(ctx, ipDailyKey(ch, ip))
	}
}

// ---------------------------------------------------------------------------
// 校验
// ---------------------------------------------------------------------------

// Consume 校验并消费验证码：通过即删除（一次性使用），失败累加次数，
// 超过上限直接销毁，防暴力枚举。
func (s *Service) Consume(ctx context.Context, ch Channel, scene, target, code string) error {
	if s == nil || s.cfg == nil {
		return errcode.ErrVerifyCodeInvalid
	}
	target = normalizeTarget(ch, target)
	if target == "" || code == "" {
		return errcode.ErrVerifyCodeInvalid
	}
	scene = strings.TrimSpace(scene)

	key := codeKey(ch, scene, target)
	v, err := g.Redis().Do(ctx, "HGETALL", key)
	if err != nil {
		return err
	}
	if v == nil {
		return errcode.ErrVerifyCodeInvalid
	}
	fields := v.Map()
	if len(fields) == 0 {
		return errcode.ErrVerifyCodeInvalid
	}

	attempts := gconv.Int(fields["attempts"])
	if attempts >= s.cfg.MaxAttempts {
		_, _ = g.Redis().Do(ctx, "DEL", key)
		return errcode.ErrVerifyCodeAttemptsExceed
	}

	expected := utility.HashVerifyCode(string(ch), scene, target, code)
	if !utility.ConstantTimeEqual(gconv.String(fields["code"]), expected) {
		n, incrErr := g.Redis().Do(ctx, "HINCRBY", key, "attempts", 1)
		if incrErr == nil && n != nil && n.Int() >= s.cfg.MaxAttempts {
			_, _ = g.Redis().Do(ctx, "DEL", key)
			return errcode.ErrVerifyCodeAttemptsExceed
		}
		return errcode.ErrVerifyCodeInvalid
	}

	if _, err = g.Redis().Do(ctx, "DEL", key); err != nil {
		g.Log().Warningf(ctx, "删除已使用的验证码失败 key=%s: %v", key, err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// 内部工具
// ---------------------------------------------------------------------------

const dailyCounterTTL = 24 * time.Hour

// 键空间统一带渠道维度：同一收件人在不同渠道上的验证码、冷却、额度互不影响。
func codeKey(ch Channel, scene, target string) string {
	return fmt.Sprintf("verify:code:%s:%s:%s", ch, scene, target)
}

func cooldownKey(ch Channel, scene, target string) string {
	return fmt.Sprintf("verify:code:cooldown:%s:%s:%s", ch, scene, target)
}

func dailyKey(ch Channel, target string) string {
	return fmt.Sprintf("verify:code:daily:%s:%s", ch, target)
}

func ipDailyKey(ch Channel, ip string) string {
	return fmt.Sprintf("verify:code:daily:ip:%s:%s", ch, ip)
}

// normalizeTarget 归一化收件人：邮箱域大小写不敏感，统一小写便于查询与比对。
func normalizeTarget(ch Channel, target string) string {
	switch ch {
	case ChannelEmail:
		return utility.NormalizeEmail(target)
	default:
		return strings.TrimSpace(target)
	}
}

// storeCode 写入验证码指纹与失败次数，并设置有效期。
func (s *Service) storeCode(ctx context.Context, ch Channel, scene, target, code string) error {
	key := codeKey(ch, scene, target)
	if _, err := g.Redis().Do(ctx, "HSET", key,
		"code", utility.HashVerifyCode(string(ch), scene, target, code),
		"attempts", 0,
	); err != nil {
		return err
	}
	_, err := g.Redis().Do(ctx, "EXPIRE", key, int(s.cfg.CodeExpire.Seconds()))
	return err
}

// incrCounter 自增日计数器，首次自增时设置 24 小时过期。
func incrCounter(ctx context.Context, key string) (int64, error) {
	v, err := g.Redis().Do(ctx, "INCR", key)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, gerror.Newf("Redis INCR 返回空值：%s", key)
	}
	n := v.Int64()
	if n == 1 {
		if _, err = g.Redis().Do(ctx, "EXPIRE", key, int(dailyCounterTTL.Seconds())); err != nil {
			g.Log().Warningf(ctx, "设置验证码日计数过期时间失败 key=%s: %v", key, err)
		}
	}
	return n, nil
}

// decCounter 归还一次日计数，失败只记日志（计数偏差不影响主流程）。
func decCounter(ctx context.Context, key string) {
	if _, err := g.Redis().Do(ctx, "DECR", key); err != nil {
		g.Log().Warningf(ctx, "归还验证码日计数失败 key=%s: %v", key, err)
	}
}

// cooldownRemain 返回冷却剩余秒数，异常时回退为配置值。
func cooldownRemain(ctx context.Context, key string, fallback time.Duration) int {
	v, err := g.Redis().Do(ctx, "TTL", key)
	if err != nil || v == nil {
		return int(fallback.Seconds())
	}
	if ttl := v.Int(); ttl > 0 {
		return ttl
	}
	return int(fallback.Seconds())
}

func clientIp(ctx context.Context) string {
	if r := g.RequestFromCtx(ctx); r != nil {
		return r.GetClientIp()
	}
	return ""
}

// FormatDuration 把有效期格式化成面向用户的中文描述（邮件/短信文案共用）。
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d 秒", int(d.Seconds()))
	}
	return fmt.Sprintf("%d 分钟", int(d.Minutes()))
}
