package orders

import (
	"context"
	"testing"
	"time"
)

// 列表时间窗口的契约：优先级（month > 区间 > 默认窗口）、闭开区间转换、回显字段。
// 这些规则决定了「一次列表要查几个分片」，所以必须有测试兜住。

func TestResolveListWindowByMonth(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 15, 0, 0, 0, time.Local)

	w, err := resolveListWindow(ctx, "202608", "", "", now)
	if err != nil {
		t.Fatalf("按月份解析失败: %v", err)
	}
	if len(w.Months) != 1 || w.Months[0] != "202608" {
		t.Errorf("命中月份 = %v, want [202608]", w.Months)
	}
	if w.FromDay != "2026-08-01" || w.ToDay != "2026-08-31" {
		t.Errorf("回显区间 = %s ~ %s, want 2026-08-01 ~ 2026-08-31", w.FromDay, w.ToDay)
	}
	if w.Defaulted {
		t.Error("显式传月份不应标记为默认窗口")
	}
	// 开区间终点应是下月 1 日（含当月最后一毫秒）
	if w.End == nil || !w.End.Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.Local)) {
		t.Errorf("开区间终点 = %v, want 2026-09-01", w.End)
	}

	// 带 - 的写法也接受（前端 DatePicker 的默认输出）
	if w2, err := resolveListWindow(ctx, "2026-08", "", "", now); err != nil || w2.Months[0] != "202608" {
		t.Errorf("month=2026-08 应被接受，实际 %v err=%v", w2.Months, err)
	}
	if _, err := resolveListWindow(ctx, "20268", "", "", now); err == nil {
		t.Error("非法 month 应报错")
	}
}

func TestResolveListWindowByRange(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 15, 0, 0, 0, time.Local)

	w, err := resolveListWindow(ctx, "", "2026-08-15", "2026-10-03", now)
	if err != nil {
		t.Fatalf("按区间解析失败: %v", err)
	}
	want := []string{"202608", "202609", "202610"}
	if len(w.Months) != len(want) {
		t.Fatalf("命中月份 = %v, want %v", w.Months, want)
	}
	for i := range want {
		if w.Months[i] != want[i] {
			t.Fatalf("命中月份 = %v, want %v", w.Months, want)
		}
	}
	// created_to 是**闭区间**：内部终点要 +1 天，否则当天 00:00 之后的下单会漏
	if w.End == nil || !w.End.Equal(time.Date(2026, 10, 4, 0, 0, 0, 0, time.Local)) {
		t.Errorf("开区间终点 = %v, want 2026-10-04", w.End)
	}
	if w.ToDay != "2026-10-03" {
		t.Errorf("回显终点 = %q, want 2026-10-03（回显用用户给的闭区间值）", w.ToDay)
	}

	if _, err := resolveListWindow(ctx, "", "2026-10-03", "2026-08-15", now); err == nil {
		t.Error("起点晚于终点应报错")
	}
	if _, err := resolveListWindow(ctx, "", "2026/08/15", "", now); err == nil {
		t.Error("非法日期格式应报错")
	}
}

func TestResolveListWindowDefault(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 4, 15, 0, 0, 0, time.Local)

	w, err := resolveListWindow(ctx, "", "", "", now)
	if err != nil {
		t.Fatalf("默认窗口解析失败: %v", err)
	}
	if !w.Defaulted {
		t.Error("未传时间参数必须标记 Defaulted（前端要据此提示「默认只看近 N 个月」）")
	}
	if len(w.Months) != defaultWindowMonths {
		t.Errorf("默认窗口月数 = %d, want %d", len(w.Months), defaultWindowMonths)
	}
	// 默认窗口必须以「当前月」结尾（含当月，否则今天的单会被自己过滤掉）
	if last := w.Months[len(w.Months)-1]; last != "202610" {
		t.Errorf("默认窗口最后一个月 = %s, want 202610（必须含当月）", last)
	}
	if w.FromDay != "2025-11-01" {
		t.Errorf("默认窗口起点 = %s, want 2025-11-01（近 12 个月）", w.FromDay)
	}
}
