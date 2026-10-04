package orders

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// 本文件锁定 Phase 1 的标识生成契约（见 docs/order-sharding-design.md §5.3）：
// 单号可路由、主键全局唯一且 JS 安全、序列用尽要报错而不是发重号、Redis 故障要能降级。

// fakeAllocator 返回可控的序列段，用于精确断言生成结果。
type fakeAllocator struct {
	start int64
	err   error
	calls int
}

func (f *fakeAllocator) nextBlock(_ context.Context, _ string, _ int64) (int64, error) {
	f.calls++
	if f.err != nil {
		return 0, f.err
	}
	return f.start, nil
}

func TestAllocateOrderNoAndID(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)

	gen := &identityGen{alloc: &fakeAllocator{start: 12345}}
	idn, err := gen.allocate(ctx, orderNoModeSequence, now, 2)
	if err != nil {
		t.Fatalf("allocate 失败: %v", err)
	}

	// 单号格式：ORD + YYYYMMDDHHMMSS + 6 位序列
	if want := "ORD20261004130509012345"; idn.OrderNo != want {
		t.Errorf("OrderNo = %q, want %q", idn.OrderNo, want)
	}
	if want := "SUB20261004130509012345"; idn.SubOrderNo != want {
		t.Errorf("SubOrderNo = %q, want %q", idn.SubOrderNo, want)
	}
	// 主键低 20 位就是分配到的序列，锁住位布局：
	// 段内顺序是 订单 → 子订单 → 明细… → 日志（四张表的主键必须互不重复）
	if got := idn.OrderID & orderIDSeqMask; got != 12345 {
		t.Errorf("订单主键低 20 位 = %d, want 12345（位布局被改动？）", got)
	}
	if got := idn.SubOrderID & orderIDSeqMask; got != 12346 {
		t.Errorf("子订单主键低 20 位 = %d, want 12346", got)
	}
	if got := idn.ItemIDs[0] & orderIDSeqMask; got != 12347 {
		t.Errorf("第 1 条明细主键低 20 位 = %d, want 12347", got)
	}
	if got := idn.ItemIDs[1] & orderIDSeqMask; got != 12348 {
		t.Errorf("第 2 条明细主键低 20 位 = %d, want 12348", got)
	}
	if got := idn.LogID & orderIDSeqMask; got != 12349 {
		t.Errorf("日志主键低 20 位 = %d, want 12349", got)
	}
	// 四个主键必须两两不同，且都能按 id 反解出同一个秒
	ids := []int64{idn.OrderID, idn.SubOrderID, idn.ItemIDs[0], idn.ItemIDs[1], idn.LogID}
	idSeen := map[int64]bool{}
	for _, id := range ids {
		if idSeen[id] {
			t.Errorf("主键重复: %d", id)
		}
		idSeen[id] = true
		if id < legacyOrderIDMax {
			t.Errorf("主键 %d 小于 legacyOrderIDMax，会被判为老自增主键", id)
		}
		got, ok := decodeOrderID(id)
		if !ok || !got.Equal(now) {
			t.Errorf("主键 %d 反解时间 = %v (ok=%v), want %v", id, got, ok, now)
		}
	}

	// 单号必须能被路由层解析，且解析出的时间与主键解出的时间同月
	parsed, ok := parseOrderNoTime(idn.OrderNo)
	if !ok {
		t.Fatalf("生成的单号无法被 parseOrderNoTime 解析: %q", idn.OrderNo)
	}
	if parsed.Format(shardSuffixGoFmt) != now.Format(shardSuffixGoFmt) {
		t.Errorf("单号内嵌时间落到了 %s，want %s", parsed.Format(shardSuffixGoFmt), now.Format(shardSuffixGoFmt))
	}
	decoded, ok := decodeOrderID(idn.OrderID)
	if !ok {
		t.Fatalf("主键无法反解: %d", idn.OrderID)
	}
	if decoded.Format(shardSuffixGoFmt) != parsed.Format(shardSuffixGoFmt) {
		t.Errorf("主键月份 %s != 单号月份 %s", decoded.Format(shardSuffixGoFmt), parsed.Format(shardSuffixGoFmt))
	}

	// 主键必须 JS 安全（前端按 JSON number 解析不丢精度）
	const jsMaxSafeInteger = int64(1)<<53 - 1
	if idn.OrderID > jsMaxSafeInteger {
		t.Errorf("主键 %d 超过 JS 安全整数上限 %d", idn.OrderID, jsMaxSafeInteger)
	}

	// 订单主键与各明细主键互不相同
	seen := map[int64]bool{idn.OrderID: true}
	for i, itemID := range idn.ItemIDs {
		if itemID <= 0 {
			t.Errorf("第 %d 条明细主键非正: %d", i, itemID)
		}
		if seen[itemID] {
			t.Errorf("第 %d 条明细主键与已有主键重复: %d", i, itemID)
		}
		seen[itemID] = true
	}
}

func TestAllocateSequenceExhausted(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)

	// 段首 999999 + 需要 2 个 → 越过 orderNoSeqMax，必须报错而不是补零截断
	gen := &identityGen{alloc: &fakeAllocator{start: orderNoSeqMax}}
	if _, err := gen.allocate(ctx, orderNoModeSequence, now, 1); err == nil {
		t.Fatal("序列越界时应返回错误")
	}

	gen = &identityGen{alloc: &fakeAllocator{start: 0}}
	if _, err := gen.allocate(ctx, orderNoModeSequence, now, 0); err == nil {
		t.Fatal("序列为 0 时应返回错误")
	}
}

func TestAllocateLegacyStillHasUniqueIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)

	// 用真实的每秒计数器，验证回滚模式**仍然**向序列分配器要号
	gen := &identityGen{alloc: &localSeqAllocator{}}
	a, err := gen.allocate(ctx, orderNoModeLegacy, now, 1)
	if err != nil {
		t.Fatalf("legacy allocate 失败: %v", err)
	}
	// 回滚模式：单号退回「ORD + 时间 + 4 位随机」= 21 字符
	if len(a.OrderNo) != 21 {
		t.Errorf("legacy 单号长度 = %d (%q), want 21", len(a.OrderNo), a.OrderNo)
	}
	if _, ok := parseOrderNoTime(a.OrderNo); !ok {
		t.Errorf("legacy 单号仍应可被解析路由: %q", a.OrderNo)
	}
	b, err := gen.allocate(ctx, orderNoModeLegacy, now, 1)
	if err != nil {
		t.Fatalf("legacy allocate 失败: %v", err)
	}
	if a.OrderID == b.OrderID {
		t.Error("legacy 模式两次分配的主键相同（说明主键也退回了随机，不应如此）")
	}
	// 主键低 20 位应来自序列。每次建单取 3+itemCount=4 个号（订单/子订单/明细/日志），
	// 所以两次分配的订单主键序列位分别是 1 和 5。
	if got := a.OrderID & orderIDSeqMask; got != 1 {
		t.Errorf("首次主键序列位 = %d, want 1", got)
	}
	if got := b.OrderID & orderIDSeqMask; got != 5 {
		t.Errorf("第二次主键序列位 = %d, want 5", got)
	}
}

func TestLocalSeqAllocatorResetsPerSecond(t *testing.T) {
	ctx := context.Background()
	l := &localSeqAllocator{}

	if got, _ := l.nextBlock(ctx, "20261004130509", 3); got != 1 {
		t.Errorf("同秒首个段首 = %d, want 1", got)
	}
	if got, _ := l.nextBlock(ctx, "20261004130509", 3); got != 4 {
		t.Errorf("同秒第二个段首 = %d, want 4", got)
	}
	if got, _ := l.nextBlock(ctx, "20261004130510", 3); got != 1 {
		t.Errorf("跨秒后应重置为 1, got %d", got)
	}
}

func TestResilientAllocatorFallsBack(t *testing.T) {
	ctx := context.Background()
	fallback := &localSeqAllocator{}
	a := &resilientSeqAllocator{
		primary:  &fakeAllocator{err: errors.New("redis down")},
		fallback: fallback,
	}
	got, err := a.nextBlock(ctx, "20261004130509", 1)
	if err != nil {
		t.Fatalf("主实现失败时应降级而不是报错: %v", err)
	}
	if got != 1 {
		t.Errorf("降级后的段首 = %d, want 1", got)
	}
}

// TestAllocateConcurrentUnique 是最关键的一条：同一秒并发分配时单号与主键都不得重复。
func TestAllocateConcurrentUnique(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)

	const goroutines, perGoroutine, itemCount = 50, 20, 3
	gen := &identityGen{
		alloc: &resilientSeqAllocator{
			primary:  &fakeAllocator{err: errors.New("redis down")}, // 强制走降级分配器
			fallback: &localSeqAllocator{},
		},
	}

	var (
		mu       sync.Mutex
		orderNos = make(map[string]bool, goroutines*perGoroutine)
		ids      = make(map[int64]bool, goroutines*perGoroutine*(itemCount+1))
		failed   error
		wg       sync.WaitGroup
	)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				idn, err := gen.allocate(ctx, orderNoModeSequence, now, itemCount)
				if err != nil {
					mu.Lock()
					failed = err
					mu.Unlock()
					return
				}
				mu.Lock()
				if orderNos[idn.OrderNo] {
					failed = errors.New("单号重复: " + idn.OrderNo)
				}
				orderNos[idn.OrderNo] = true
				if ids[idn.OrderID] {
					failed = errors.New("订单主键重复")
				}
				ids[idn.OrderID] = true
				for _, itemID := range idn.ItemIDs {
					if ids[itemID] {
						failed = errors.New("明细主键重复")
					}
					ids[itemID] = true
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if failed != nil {
		t.Fatal(failed)
	}
	if want := goroutines * perGoroutine; len(orderNos) != want {
		t.Errorf("唯一单号数 = %d, want %d", len(orderNos), want)
	}
	if want := goroutines * perGoroutine * (itemCount + 1); len(ids) != want {
		t.Errorf("唯一主键数 = %d, want %d", len(ids), want)
	}
}

func TestDecodeOrderIDRoundTrip(t *testing.T) {
	// 覆盖月末/年末边界：解码月份必须与编码时间一致
	times := []time.Time{
		time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local),
		time.Date(2026, 1, 31, 23, 59, 59, 0, time.Local),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local),
		time.Date(2026, 12, 31, 23, 59, 59, 0, time.Local),
		time.Date(2027, 1, 1, 0, 0, 1, 0, time.Local),
	}
	for _, at := range times {
		for _, seq := range []int64{1, 500, orderNoSeqMax} {
			id := encodeOrderID(at, seq)
			got, ok := decodeOrderID(id)
			if !ok {
				t.Fatalf("decodeOrderID(%d) 失败 (at=%v seq=%d)", id, at, seq)
			}
			if got.Format("20060102150405") != at.Format("20060102150405") {
				t.Errorf("at=%v seq=%d 解码得到 %v", at, seq, got)
			}
			if id > int64(1)<<53-1 {
				t.Errorf("at=%v seq=%d 编码出 %d，超过 JS 安全整数上限", at, seq, id)
			}
		}
	}
	if _, ok := decodeOrderID(0); ok {
		t.Error("0 不是合法主键，应返回 false")
	}
	if _, ok := decodeOrderID(-1); ok {
		t.Error("负数不是合法主键，应返回 false")
	}
	// 迁移前的自增主键（1..2000）不能被当成编码主键，
	// 否则会被解成 2024-01-01 并路由到不存在的分片
	for _, legacy := range []int64{1, 2000, legacyOrderIDMax - 1} {
		if decoded, ok := decodeOrderID(legacy); ok {
			t.Errorf("老自增主键 %d 不应被反解成功（得到 %v）", legacy, decoded)
		}
	}
}

// TestNextLogID 锁定「状态变更时单独分配日志主键」的行为：
// 与批量分配同源同布局，且与订单主键同秒可反解。
func TestNextLogID(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)
	// 用真实的每秒计数器：要验证的核心是「两次分配会各自推进序列」
	gen := &identityGen{alloc: &localSeqAllocator{}}

	first, err := gen.nextLogID(ctx, orderNoModeSequence, now)
	if err != nil {
		t.Fatalf("nextLogID 失败: %v", err)
	}
	second, err := gen.nextLogID(ctx, orderNoModeSequence, now)
	if err != nil {
		t.Fatalf("nextLogID 失败: %v", err)
	}
	if first == second {
		t.Error("两次分配的日志主键相同")
	}
	if first&orderIDSeqMask != 1 || second&orderIDSeqMask != 2 {
		t.Errorf("日志主键序列位 = %d/%d, want 1/2", first&orderIDSeqMask, second&orderIDSeqMask)
	}
	if got, ok := decodeOrderID(first); !ok || !got.Equal(now) {
		t.Errorf("日志主键反解时间 = %v (ok=%v), want %v", got, ok, now)
	}
}
