package orders

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/errcode"
)

// 列表的时间窗口。
//
// 分表之后列表**必须**带时间范围，否则只能对全部活跃分片 fan-out：
// 实测一次翻页 = 活跃分片数次查询，按 10 万单/日、保留 3 年（36 片）算就是 36 次/页。
// 方案与取舍见 docs/order-sharding-design.md §5.5。
//
// 优先级：month（精确到月）> created_from/created_to（区间）> 默认窗口（近 N 个月）。
// 按 order_no 查询走点查、不经过这里（客服找单不能被月份默认值挡住）。

const (
	// defaultWindowMonths 未传时间参数时默认回溯的月数。
	defaultWindowMonths = 12
	// maxWindowMonths 一次查询最多覆盖的月数（防止构造一个跨十年的区间把分片全扫一遍）。
	maxWindowMonths = 240
	// monthLayout 月份参数的规范形式（也接受带 - 的 2026-10）。
	monthLayout = "200601"
	dayLayout   = "2006-01-02"
)

// orderListWindow 一次列表查询实际生效的时间窗口。
type orderListWindow struct {
	Months    []string   // 命中的月份（YYYYMM，升序）；空表示不限月份（只有单表模式会这样）
	Start     *time.Time // 闭区间起点；nil 表示不限
	End       *time.Time // 开区间终点（"到某日"会 +1 天）；nil 表示不限
	FromDay   string     // 回显用的起点日期（含当天）
	ToDay     string     // 回显用的终点日期（含当天）
	Defaulted bool       // 未传时间参数、由后端默认窗口兜底
}

// resolveListWindow 把请求参数解析成时间窗口。
//
// now 显式传入（而不是内部取 time.Now），这样默认窗口与边界行为可以单测。
func resolveListWindow(ctx context.Context, month, fromDay, toDay string, now time.Time) (orderListWindow, error) {
	month = strings.TrimSpace(month)
	fromDay = strings.TrimSpace(fromDay)
	toDay = strings.TrimSpace(toDay)

	switch {
	case month != "":
		start, err := parseMonthStart(month)
		if err != nil {
			return orderListWindow{}, err
		}
		end := start.AddDate(0, 1, 0)
		return orderListWindow{
			Months:  []string{start.Format(monthLayout)},
			Start:   &start,
			End:     &end,
			FromDay: start.Format(dayLayout),
			ToDay:   end.AddDate(0, 0, -1).Format(dayLayout),
		}, nil

	case fromDay != "" || toDay != "":
		var (
			start, end *time.Time
			fromText   string
			toText     string
		)
		if fromDay != "" {
			t, err := time.ParseInLocation(dayLayout, fromDay, time.Local)
			if err != nil {
				return orderListWindow{}, errcode.Newf(errcode.CodeInvalidParams,
					"created_from 需形如 2026-08-01，收到 %q", fromDay)
			}
			start, fromText = &t, fromDay
		}
		if toDay != "" {
			t, err := time.ParseInLocation(dayLayout, toDay, time.Local)
			if err != nil {
				return orderListWindow{}, errcode.Newf(errcode.CodeInvalidParams,
					"created_to 需形如 2026-08-31，收到 %q", toDay)
			}
			t = t.AddDate(0, 0, 1) // 闭区间转开区间
			end, toText = &t, toDay
		}
		if start != nil && end != nil && !end.After(*start) {
			return orderListWindow{}, errcode.Newf(errcode.CodeInvalidParams,
				"created_from(%s) 不能晚于 created_to(%s)", fromText, toText)
		}
		months, err := monthsBetween(start, end, now)
		if err != nil {
			return orderListWindow{}, err
		}
		return orderListWindow{
			Months: months, Start: start, End: end, FromDay: fromText, ToDay: toText,
		}, nil
	}

	// 没传时间参数：默认回溯 N 个月，并把生效区间回显给客户端
	n := g.Cfg().MustGet(ctx, "orderShard.defaultWindowMonths", defaultWindowMonths).Int()
	if n <= 0 {
		n = defaultWindowMonths
	}
	if n > maxWindowMonths {
		n = maxWindowMonths
	}
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	start := end.AddDate(0, -n, 0)
	months := make([]string, 0, n)
	for i := 0; i < n; i++ {
		months = append(months, start.AddDate(0, i, 0).Format(monthLayout))
	}
	return orderListWindow{
		Months:    months,
		Start:     &start,
		End:       &end,
		FromDay:   start.Format(dayLayout),
		ToDay:     end.AddDate(0, 0, -1).Format(dayLayout),
		Defaulted: true,
	}, nil
}

// windowShards 返回窗口命中的活跃分片（升序）。
//
//   - 窗口为空（未限月份）→ 全部分片：只有单表模式会走到这里；
//   - 窗口内的月份与活跃分片取交集 —— 查一个还没建表的月份就是「空结果」，
//     不该报错（例如查下个月、或查早已归档清理的月份）。
func windowShards(ctx context.Context, w orderListWindow) ([]shard, error) {
	shards, err := activeShards(ctx)
	if err != nil {
		return nil, err
	}
	if len(w.Months) == 0 {
		return shards, nil
	}
	inWindow := make(map[string]bool, len(w.Months))
	for _, m := range w.Months {
		inWindow[m] = true
	}
	out := make([]shard, 0, len(shards))
	for _, sh := range shards {
		if inWindow[sh.suffix] {
			out = append(out, sh)
		}
	}
	return out, nil
}

func parseMonthStart(month string) (time.Time, error) {
	normalized := strings.ReplaceAll(month, "-", "")
	t, err := time.ParseInLocation(monthLayout, normalized, time.Local)
	if err != nil {
		return time.Time{}, errcode.Newf(errcode.CodeInvalidParams,
			"month 需形如 202610（或 2026-10），收到 %q", month)
	}
	return t, nil
}

// monthsBetween 列出 [start, end) 覆盖的月份；nil 表示该侧不限。
func monthsBetween(start, end *time.Time, now time.Time) ([]string, error) {
	first := now.AddDate(0, -(maxWindowMonths - 1), 0)
	last := now.AddDate(0, maxWindowMonths, 0)
	if start != nil {
		first = *start
	}
	if end != nil && end.Before(last) {
		last = *end
	}
	if !last.After(first) {
		return nil, nil
	}
	var months []string
	for cur := time.Date(first.Year(), first.Month(), 1, 0, 0, 0, 0, time.Local); cur.Before(last); cur = cur.AddDate(0, 1, 0) {
		months = append(months, cur.Format(monthLayout))
		if len(months) > maxWindowMonths {
			return nil, errcode.Newf(errcode.CodeInvalidParams,
				"查询区间超过 %d 个月，请缩小范围", maxWindowMonths)
		}
	}
	return months, nil
}

// describe 便于日志与报错里说明本次窗口。
func (w orderListWindow) describe() string {
	if w.Defaulted {
		return fmt.Sprintf("默认窗口 %s ~ %s（近 %d 个月）", w.FromDay, w.ToDay, len(w.Months))
	}
	return fmt.Sprintf("窗口 %s ~ %s", w.FromDay, w.ToDay)
}

// listSource 列表查询的数据来源。
type listSource struct {
	Main   bool    // 读主表（单表模式，或灰度/兜底）
	Shards []shard // 读这些分片；Main=false 且为空切片 = 窗口内没有分片 ⇒ 空结果
}

// resolveListSource 决定这次列表查询读哪里。
//
// 三档优先级：
//  1. order_no（点查）：等价单分片，且**不受月份窗口限制**（客服找单不能被月份挡住）；
//  2. 分片模式：窗口命中的分片 —— 落在 1 片就能精确 offset，落在多片则 keyset 归并；
//     窗口内一个分片都没有（比如查一个还没建表的月份）⇒ 空结果，而不是回落主表
//     （回落会返回与该月份无关的行）；
//  3. 否则主表（单表模式），并保留灰度回落计数。
func resolveListSource(ctx context.Context, f orderListFilter) (listSource, error) {
	if sh, ok := listShardForPointLookup(ctx, f); ok {
		return listSource{Shards: []shard{sh}}, nil
	}
	if !readShardsForList(ctx) {
		shardFallbackCounter(ctx, "list_main")
		return listSource{Main: true}, nil
	}
	shards, err := windowShards(ctx, f.Window)
	if err != nil {
		return listSource{}, err
	}
	return listSource{Shards: shards}, nil
}
