package orders

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/internal/errcode"
)

// 本文件是订单号与订单主键的生成器，方案见 docs/order-sharding-design.md §5.3。
//
//	单号：ORD + YYYYMMDDHHMMSS(14) + 每秒序列(6)
//	主键：分钟位(25) | 秒位(6) | 序列位(20)
//
// 两者**共用同一个每秒序列**，所以单号内嵌的时间与主键高位解出的时间是同一时刻，
// 不会各推一份时间而在跨月瞬间分叉。
//
// 为什么主键不用标准雪花（41 位毫秒）：那会超过 2^53，前端 JS 解析 JSON number 会丢精度。
// 这里刻意压在 51 位以内，并且把「分钟」编进高位，按 id 也能反解出月份用于分片路由。

const (
	// orderNoPrefix 父订单号前缀。
	orderNoPrefix = "ORD"
	// subOrderNoPrefix 子订单号前缀。
	subOrderNoPrefix = "SUB"

	// orderNoSeqDigits 单号尾部序列位数：6 位 → 每秒最多 999999 单。
	orderNoSeqDigits = 6
	// orderNoSeqMax 序列上限。超过就报错，绝不静默截断（宁可失败也不发重号）。
	orderNoSeqMax = 999999

	// 主键位布局：25 + 6 + 20 = 51 位，最大约 2.25e15 < 2^53-1（JS 安全整数上限）。
	orderIDMinutesBits = 25
	orderIDSecondsBits = 6
	orderIDSeqBits     = 20
	orderIDSeqMask     = 1<<orderIDSeqBits - 1
	orderIDSecondMask  = 1<<orderIDSecondsBits - 1
)

// 单号生成模式（配置 orderNo.mode）。
const (
	// orderNoModeSequence Redis 每秒序列，默认。
	orderNoModeSequence = "sequence"
	// orderNoModeLegacy 回滚用：旧的「4 位随机后缀」。会有同秒重号风险，
	// 仅在需要临时关掉序列分配时使用。
	orderNoModeLegacy = "legacy"
)

const (
	// orderNoTimeGoFmt 单号内嵌时间的 Go 布局。
	orderNoTimeGoFmt = "20060102150405"
	// orderNoTimeStart/End 单号中时间段的起止下标（"ORD" 之后 14 位），与 shard.go 的解析保持一致。
	orderNoTimeStart = 3
	orderNoTimeEnd   = 17
	// seqKeyTTL 每秒序列键的存活时间，避免键无限堆积。
	seqKeyTTL = 10 * time.Second
	// createMaxAttempts 建单遇到唯一键冲突时的最大尝试次数。
	createMaxAttempts = 3
)

// orderIDEpoch 主键分钟位的起算点（UTC）。
var orderIDEpoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// ── 序列分配 ──────────────────────────────────────────────────────────────

// seqAllocator 为「某一秒」分配一段连续序列。抽成接口是为了可单测，
// 也为了 Redis 不可用时能整体换成降级实现。
type seqAllocator interface {
	// nextBlock 为 second（形如 20261004130509）分配 n 个连续序列并返回段首。
	nextBlock(ctx context.Context, second string, n int64) (int64, error)
}

// redisSeqAllocator 生产实现：Redis INCRBY，全局单调，天然跨实例唯一。
type redisSeqAllocator struct{}

func (redisSeqAllocator) nextBlock(ctx context.Context, second string, n int64) (int64, error) {
	key := "order:seq:" + second
	v, err := g.Redis().Do(ctx, "INCRBY", key, n)
	if err != nil {
		return 0, err
	}
	// 续期失败不影响正确性，忽略错误
	_, _ = g.Redis().Do(ctx, "EXPIRE", key, int(seqKeyTTL.Seconds()))
	return v.Int64() - n + 1, nil
}

// localSeqAllocator 是 Redis 不可用时的降级实现：进程内的每秒计数器。
// 单实例下唯一；多实例下**不保证**跨实例唯一，最终由 uk_order_no / 主键冲突 + 重试兜底。
type localSeqAllocator struct {
	mu     sync.Mutex
	second string
	seq    int64
}

func (l *localSeqAllocator) nextBlock(_ context.Context, second string, n int64) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.second != second {
		l.second = second
		l.seq = 0
	}
	start := l.seq + 1
	l.seq += n
	return start, nil
}

// resilientSeqAllocator 优先用主实现，失败时降级到本地计数器并告警一次。
// 订单写入不因 Redis 不可用而中断（与项目「缓存是加速器不是硬依赖」的取向一致）。
type resilientSeqAllocator struct {
	primary  seqAllocator
	fallback seqAllocator
	warnOnce sync.Once
}

func (a *resilientSeqAllocator) nextBlock(ctx context.Context, second string, n int64) (int64, error) {
	start, err := a.primary.nextBlock(ctx, second, n)
	if err == nil {
		return start, nil
	}
	a.warnOnce.Do(func() {
		g.Log().Warningf(ctx,
			"订单号序列分配失败，已降级为进程内计数器：多实例部署下不保证唯一，"+
				"将由唯一键冲突触发重试兜底。原因: %v", err)
	})
	return a.fallback.nextBlock(ctx, second, n)
}

// ── 生成器 ────────────────────────────────────────────────────────────────

// orderIdentity 一次建单所需的全套标识。
type orderIdentity struct {
	OrderNo    string
	SubOrderNo string
	OrderID    int64
	SubOrderID int64   // 子订单主键（与订单同一序列，全局唯一）
	LogID      int64   // 建单那条操作日志的主键
	ItemIDs    []int64 // 与入参明细顺序一一对应
}

// identityGen 订单标识生成器。
type identityGen struct {
	alloc seqAllocator
}

// orderIdentities 是包级单例：主实现为 Redis 序列，降级为进程内计数器。
var orderIdentities = &identityGen{
	alloc: &resilientSeqAllocator{
		primary:  redisSeqAllocator{},
		fallback: &localSeqAllocator{},
	},
}

// allocate 用同一个时间基准为一次建单分配单号、订单主键与各明细主键。
//
// 序列段一次取够（订单 1 个 + 每条明细 1 个），既省 Redis 往返，
// 也保证同一订单内的主键互不相同。legacy 模式只是把**单号字符串**换回旧格式，
// 主键仍走同一个序列，避免回滚开关把主键唯一性一起牺牲掉。
// 需要的主键个数：订单 + 子订单 + 每条明细 + 建单日志。
func (g *identityGen) allocate(ctx context.Context, mode string, now time.Time, itemCount int) (*orderIdentity, error) {
	need := int64(3 + itemCount)
	start, err := g.alloc.nextBlock(ctx, now.Format(orderNoTimeGoFmt), need)
	if err != nil {
		return nil, errcode.Newf(errcode.CodeOrderNoAllocateFailed,
			"分配订单序列失败: %v", err)
	}
	// 序列必须整段落在 [1, orderNoSeqMax]：单号补零位数与主键的序列位都依赖这个上界
	if start < 1 || start+need-1 > orderNoSeqMax {
		return nil, errcode.Newf(errcode.CodeOrderNoAllocateFailed,
			"本秒订单序列已用尽（start=%d, need=%d, max=%d）", start, need, orderNoSeqMax)
	}

	orderNo := buildOrderNo(orderNoPrefix, now, start)
	subOrderNo := buildOrderNo(subOrderNoPrefix, now, start)
	if mode == orderNoModeLegacy {
		// 回滚用：旧的 4 位随机后缀（同秒有重号风险，由唯一键 + 重试兜底）
		orderNo = fmt.Sprintf("%s%s%04d", orderNoPrefix, now.Format(orderNoTimeGoFmt), grand.Intn(10000))
		subOrderNo = fmt.Sprintf("%s%s%04d", subOrderNoPrefix, now.Format(orderNoTimeGoFmt), grand.Intn(10000))
	}

	// 段内顺序：订单 → 子订单 → 明细… → 日志。
	// 四张表的主键必须**全局唯一且互不重复**：分片之间、以及分片与主表之间都不能撞
	// （tx_sub_orders / tx_order_logs 原先各自 AUTO_INCREMENT，分片表从 1 开始计数，
	// 回灌时会按主键命中无关行 —— 见 docs/order-sharding-design.md §7 Phase 5）。
	id := &orderIdentity{
		OrderNo:    orderNo,
		SubOrderNo: subOrderNo,
		OrderID:    encodeOrderID(now, start),
		SubOrderID: encodeOrderID(now, start+1),
		LogID:      encodeOrderID(now, start+2+int64(itemCount)),
		ItemIDs:    make([]int64, itemCount),
	}
	for i := 0; i < itemCount; i++ {
		id.ItemIDs[i] = encodeOrderID(now, start+2+int64(i))
	}
	return id, nil
}

// nextLogID 单独分配一个日志主键（用于状态变更这类「只有一条日志」的写入）。
func (g *identityGen) nextLogID(ctx context.Context, mode string, now time.Time) (int64, error) {
	start, err := g.alloc.nextBlock(ctx, now.Format(orderNoTimeGoFmt), 1)
	if err != nil {
		return 0, errcode.Newf(errcode.CodeOrderNoAllocateFailed, "分配日志主键失败: %v", err)
	}
	if start < 1 || start > orderNoSeqMax {
		return 0, errcode.Newf(errcode.CodeOrderNoAllocateFailed,
			"本秒序列已用尽（start=%d, max=%d）", start, orderNoSeqMax)
	}
	return encodeOrderID(now, start), nil
}

// orderNoMode 读取 orderNo.mode，缺省 sequence；非法值按 sequence 处理并告警。
func orderNoMode(ctx context.Context) string {
	mode := g.Cfg().MustGet(ctx, "orderNo.mode", orderNoModeSequence).String()
	switch mode {
	case orderNoModeSequence:
		return orderNoModeSequence
	case orderNoModeLegacy:
		g.Log().Warningf(ctx, "orderNo.mode=legacy：单号已退回 4 位随机后缀，存在同秒重号风险")
		return orderNoModeLegacy
	default:
		g.Log().Warningf(ctx, "orderNo.mode=%q 非法（sequence/legacy），已按 sequence 处理", mode)
		return orderNoModeSequence
	}
}

// buildOrderNo 拼接业务单号。
func buildOrderNo(prefix string, now time.Time, seq int64) string {
	return fmt.Sprintf("%s%s%0*d", prefix, now.Format(orderNoTimeGoFmt), orderNoSeqDigits, seq)
}

// encodeOrderID 把「时间 + 序列」编成全局唯一主键。
func encodeOrderID(now time.Time, seq int64) int64 {
	minutes := now.Unix()/60 - orderIDEpoch.Unix()/60
	return minutes<<(orderIDSecondsBits+orderIDSeqBits) |
		int64(now.Second())<<orderIDSeqBits |
		seq
}

// decodeOrderID 反解主键中的时间，供「按 id 定位分片」使用（Phase 2/3）。
// 解出的时间按本地时区返回，与 created_at / 单号内嵌时间同一口径。
//
// 小于 legacyOrderIDMax 的 id 一律判为「不是编码主键」并返回 false：
// 迁移前的自增主键（如 1..2000）在数学上也能解出一串位，会被误判成
// 2024-01-01 的单而路由到根本不存在的分片 —— 必须靠阈值区分，
// 老数据改走 tx_order_shard_map（见 ShardOfOrderID）。
func decodeOrderID(id int64) (time.Time, bool) {
	if id < legacyOrderIDMax {
		return time.Time{}, false
	}
	seq := id & orderIDSeqMask
	second := (id >> orderIDSeqBits) & orderIDSecondMask
	minutes := id >> (orderIDSecondsBits + orderIDSeqBits)
	if seq == 0 || second > 59 || minutes < 0 {
		return time.Time{}, false
	}
	t := orderIDEpoch.Add(time.Duration(minutes)*time.Minute + time.Duration(second)*time.Second)
	return t.Local(), true
}

// isDuplicateKey 判定唯一键冲突，只在降级生成标识时才可能出现，用于触发重试。
func isDuplicateKey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
