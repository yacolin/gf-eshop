package orders

import (
	"context"
	"strings"
	"testing"
	"time"
)

// 归档相关的纯逻辑与守卫：命名、后缀解析、30 天窗口、以及「非终态不许归档」。

func TestLegacyTableNaming(t *testing.T) {
	if got := legacyTableName(tableOrders, "202610"); got != "tx_orders_legacy_202610" {
		t.Errorf("归档表名 = %q", got)
	}
	suffix, ok := parseLegacySuffix("tx_sub_orders_legacy_202610")
	if !ok || suffix != "202610" {
		t.Errorf("后缀解析 = (%q, %v), want (202610, true)", suffix, ok)
	}
	// 分片表名不能被误判成归档表
	if _, ok := parseLegacySuffix("tx_orders_202610"); ok {
		t.Error("分片表名不应被当成归档表")
	}
	if _, ok := parseLegacySuffix("tx_orders_legacy_xxxxxx"); ok {
		t.Error("非法后缀应判为不是归档表")
	}
}

func TestLegacyPurgeWindow(t *testing.T) {
	// 规则：归档后缀所在月再过完一整月才能清理（比 30 天更保守，只看表名即可判断）
	at, err := legacyPurgeAllowedAt("202610")
	if err != nil {
		t.Fatalf("计算窗口失败: %v", err)
	}
	if want := time.Date(2026, 12, 1, 0, 0, 0, 0, time.Local); !at.Equal(want) {
		t.Errorf("legacy_202610 最早可清理时间 = %v, want %v", at, want)
	}
	if _, err := legacyPurgeAllowedAt("2026-10"); err == nil {
		t.Error("非法后缀应报错")
	}
}

func TestArchiveRefusedOutsideMonthly(t *testing.T) {
	ctx := context.Background()
	// 非终态（主表还在被写）下归档会让写入直接失败 —— 必须在访问数据库前就拒绝
	setShardConfig(t, shardModeSingle, true, 0, 0)
	if _, err := ArchiveBaseTables(ctx, "202610"); err == nil {
		t.Fatal("single 模式下归档应被拒绝")
	} else if !strings.Contains(err.Error(), "monthly") {
		t.Errorf("拒绝原因应说明只在 monthly 下可用，实际: %v", err)
	}
}
