package orders

import (
	"context"
	"testing"
	"time"

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

	if mode := shardMode(ctx); mode != shardModeSingle {
		t.Fatalf("未配置 orderShard.mode 时应为 %q，实际 %q", shardModeSingle, mode)
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

	// 聚合 / 列表路径
	if sh, err = shardForScan(ctx, "test"); err != nil || !sh.isZero() {
		t.Fatalf("single 模式下跨片查询应放行，实际 shard=%q err=%v", sh.suffix, err)
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
