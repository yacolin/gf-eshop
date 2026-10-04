package orders

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gtime"
)

// 本文件锁定分片路由的契约（见 docs/order-sharding-design.md §5.2）。
// 这些规则一旦被改动，读写路径就会落到不同月份的表上，所以必须有测试兜住。

func TestParseOrderNoTime(t *testing.T) {
	cases := []struct {
		name    string
		orderNo string
		want    time.Time
		wantOK  bool
	}{
		{
			name:    "历史格式（后 4 位随机数）",
			orderNo: "ORD202608202105157688",
			want:    time.Date(2026, 8, 20, 21, 5, 15, 0, time.Local),
			wantOK:  true,
		},
		{
			name:    "Phase 1 新格式（后 6 位序列）",
			orderNo: "ORD20260820210515000042",
			want:    time.Date(2026, 8, 20, 21, 5, 15, 0, time.Local),
			wantOK:  true,
		},
		{
			name:    "跨年边界",
			orderNo: "ORD202512312359599999",
			want:    time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local),
			wantOK:  true,
		},
		{
			name:    "子订单号同样可解析（前缀不作校验）",
			orderNo: "SUB202608202105157688",
			want:    time.Date(2026, 8, 20, 21, 5, 15, 0, time.Local),
			wantOK:  true,
		},
		{name: "非订单号", orderNo: "ORD_NOT_EXIST", wantOK: false},
		{name: "空串", orderNo: "", wantOK: false},
		{name: "时间段不是数字", orderNo: "ORD2026AB202105157688", wantOK: false},
		{name: "月份非法", orderNo: "ORD202613202105157688", wantOK: false},
		{name: "长度不足 17", orderNo: "ORD202608202105", wantOK: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := parseOrderNoTime(c.orderNo)
			if ok != c.wantOK {
				t.Fatalf("parseOrderNoTime(%q) ok=%v, want %v", c.orderNo, ok, c.wantOK)
			}
			if ok && !got.Equal(c.want) {
				t.Fatalf("parseOrderNoTime(%q) = %v, want %v", c.orderNo, got, c.want)
			}
		})
	}
}

func TestShardTable(t *testing.T) {
	cases := []struct {
		name  string
		shard shard
		base  string
		want  string
	}{
		{name: "不分表（零值分片）", shard: shard{}, base: tableOrders, want: "tx_orders"},
		{name: "按月分片", shard: shard{suffix: "202608"}, base: tableOrders, want: "tx_orders_202608"},
		{name: "子表同样加后缀", shard: shard{suffix: "202612"}, base: tableOrderItems, want: "tx_order_items_202612"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.shard.table(c.base); got != c.want {
				t.Fatalf("shard%v.table(%q) = %q, want %q", c.shard, c.base, got, c.want)
			}
		})
	}
}

// TestSingleModeRoutesToBaseTable 锁定 Phase 0 的核心不变量：
// 默认配置下路由恒为空分片，物理表名与分表前一致 —— 也就是「行为零变化」。
func TestSingleModeRoutesToBaseTable(t *testing.T) {
	ctx := context.Background()
	setShardConfig(t, shardModeSingle, false, 0, 0)

	if mode := shardMode(ctx); mode != shardModeSingle {
		t.Fatalf("显式配置 single 后应读回 %q，实际 %q", shardModeSingle, mode)
	}

	// 写入路径
	if sh := shardFromCreatedAt(ctx, gtime.Now()); !sh.isZero() {
		t.Fatalf("single 模式下写入分片应为零值，实际 %q", sh.suffix)
	}

	// 读取路径：合法单号
	sh, err := shardFromOrderNo(ctx, "ORD202608202105157688")
	if err != nil || !sh.isZero() {
		t.Fatalf("single 模式下按单号路由应为零值且无错误，实际 shard=%q err=%v", sh.suffix, err)
	}

	// 读取路径：无法解析的单号在 single 模式下也应放行（不改变既有行为）
	if sh, err = shardFromOrderNo(ctx, "ORD_NOT_EXIST"); err != nil || !sh.isZero() {
		t.Fatalf("single 模式下非法单号不应报错，实际 shard=%q err=%v", sh.suffix, err)
	}

	// offset 分页（page/page_size）：single 模式下回落主表且不报错
	if sh, err = offsetListShard(ctx, orderListFilter{}); err != nil || !sh.isZero() {
		t.Fatalf("single 模式下 offset 分页应回落主表，实际 shard=%q err=%v", sh.suffix, err)
	}

	// 未配置灰度时，点查与列表都不走分片
	if readFromShardByKey(ctx, "ORD202608202105157688") {
		t.Error("未配置 readShardsPercent 时点查不应走分片")
	}
	if readShardsForList(ctx) {
		t.Error("未配置 readShardsPercent 时列表不应走分片")
	}
	if !allowShardFallback(ctx) {
		t.Error("single 模式应允许分片读回落主表")
	}
}

// TestShardSuffixFormat 锁定分片后缀格式：必须是 YYYYMM（gtime 的 Ym 与 Go 的 200601 必须一致）。
func TestShardSuffixFormat(t *testing.T) {
	at := gtime.New(time.Date(2026, 8, 20, 21, 5, 15, 0, time.Local))
	if got := at.Format("Ym"); got != "202608" {
		t.Fatalf(`gtime Format("Ym") = %q, want "202608"`, got)
	}
	parsed, ok := parseOrderNoTime("ORD202608202105157688")
	if !ok {
		t.Fatal("解析示例单号失败")
	}
	if got := parsed.Format(shardSuffixGoFmt); got != "202608" {
		t.Fatalf("time.Format(%q) = %q, want 202608", shardSuffixGoFmt, got)
	}
}

// setShardConfig 在内存里覆盖分表配置。
//
// 必要性：这些用例读的是 g.Cfg()，而测试进程会加载开发机本地的
// manifest/config/config.yaml —— 本地一旦切到 monthly（Phase 5 dogfood），
// 不显式设定的用例就会跟着变，测试结果取决于「谁的机器」。所以每个模式相关的
// 用例都先自己声明要测的模式。
func setShardConfig(t *testing.T, mode string, dualWrite bool, readPercent, shadowPercent int) {
	t.Helper()
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		t.Fatalf("配置适配器不是 *gcfg.AdapterFile，无法在内存里改写配置")
	}
	for k, v := range map[string]interface{}{
		"orderShard.mode":              mode,
		"orderShard.dualWrite":         dualWrite,
		"orderShard.readShardsPercent": readPercent,
		"orderShard.shadowReadPercent": shadowPercent,
	} {
		if err := adapter.Set(k, v); err != nil {
			t.Fatalf("设置配置 %s 失败: %v", k, err)
		}
	}
	shardWarnOnce = sync.Once{} // 一次性告警要让每个用例都能重新触发
}

// TestMonthlyModeRouting 锁定 monthly（Phase 5 终态）的路由矩阵。
func TestMonthlyModeRouting(t *testing.T) {
	ctx := context.Background()
	setShardConfig(t, shardModeMonthly, false, 100, 100)

	const orderNo = "ORD202608202105157688"

	// 写：按创建时间落到当月分片
	if sh := shardFromCreatedAt(ctx, gtime.New(time.Date(2026, 8, 20, 21, 5, 15, 0, time.Local))); sh.suffix != "202608" {
		t.Errorf("monthly 下写入分片 = %q, want 202608", sh.suffix)
	}
	// 读：按单号内嵌时间落到同一个月
	sh, err := shardFromOrderNo(ctx, orderNo)
	if err != nil || sh.suffix != "202608" {
		t.Errorf("monthly 下按单号路由 = (%q, %v), want (202608, nil)", sh.suffix, err)
	}
	// 无法解析的单号在 monthly 下必须报错（不能静默读主表）
	if _, err = shardFromOrderNo(ctx, "ORD_NOT_EXIST"); err == nil {
		t.Error("monthly 下非法单号应报错")
	}
	// 终态不再镜像、不允许回落
	if !mirrorShardOfOrderNo(ctx, orderNo).isZero() {
		t.Error("monthly 下不应有双写镜像目标")
	}
	if allowShardFallback(ctx) {
		t.Error("monthly 下不允许回落主表（主表已停写）")
	}
	if !readFromShardByKey(ctx, orderNo) {
		t.Error("monthly 下点查必须走分片")
	}
	// offset 分页：普通条件无法跨片 → 报错；带 order_no 等价点查 → 放行
	if _, err = offsetListShard(ctx, orderListFilter{}); err == nil {
		t.Error("monthly 下无 order_no 的 offset 分页应报错")
	}
	if sh, err = offsetListShard(ctx, orderListFilter{OrderNo: orderNo}); err != nil || sh.suffix != "202608" {
		t.Errorf("monthly 下带 order_no 的 offset 分页应路由到 202608，实际 (%q, %v)", sh.suffix, err)
	}
}

// TestDualWriteAndGraySwitches 锁定双写与读灰度的开关语义（single 模式）。
func TestDualWriteAndGraySwitches(t *testing.T) {
	ctx := context.Background()
	const orderNo = "ORD20261004130509000001"
	at := gtime.New(time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local))

	// 双写开：主写仍在主表（零值），但要额外镜像到当月分片
	setShardConfig(t, shardModeSingle, true, 0, 0)
	if sh := shardFromCreatedAt(ctx, at); !sh.isZero() {
		t.Errorf("single 模式主写应在主表，实际分片 %q", sh.suffix)
	}
	if sh := mirrorShardOfCreatedAt(ctx, at); sh.suffix != "202610" {
		t.Errorf("双写开启时镜像分片 = %q, want 202610", sh.suffix)
	}
	if !readShardsForList(ctx) == (readShardsPercent(ctx) > 0) {
		t.Error("列表读灰度的开关语义不一致")
	}

	// 双写关：没有镜像目标
	setShardConfig(t, shardModeSingle, false, 0, 0)
	if sh := mirrorShardOfOrderNo(ctx, orderNo); !sh.isZero() {
		t.Errorf("双写关闭时不应有镜像分片，实际 %q", sh.suffix)
	}

	// 读灰度 100%：命中分片；0%：不命中
	setShardConfig(t, shardModeSingle, false, 100, 0)
	if !readFromShardByKey(ctx, orderNo) || !readShardsForList(ctx) {
		t.Error("读灰度 100% 时点查与列表都应走分片")
	}
	setShardConfig(t, shardModeSingle, false, 0, 0)
	if readFromShardByKey(ctx, orderNo) || readShardsForList(ctx) {
		t.Error("读灰度 0% 时点查与列表都不应走分片")
	}
}
