package products

import (
	"encoding/base64"
	"strconv"
	"testing"

	"gf-eshop/internal/model/entity"
)

// 列表游标分页的契约。字段名与行为都必须与 orders 列表一致：
// 响应 next_cursor + has_more + total(-1)、请求 cursor + size、
// 非法游标报错、每页条数默认 20 / 上限 100。

func TestDecodeCursor(t *testing.T) {
	if id, err := decodeCursor(""); err != nil || id != 0 {
		t.Errorf("空游标 = (%d, %v), want (0, nil)", id, err)
	}

	c := encodeCursor(998877)
	if id, err := decodeCursor(c); err != nil || id != 998877 {
		t.Errorf("往返 = (%d, %v), want (998877, nil)", id, err)
	}

	// 与 orders 一致：非法游标报错，不静默降级成首页
	//   not-base64!!  非法 base64 字符
	//   YWJj          base64("abc")，解出来不是数字
	//   LTE= / MA==   base64("-1") / base64("0")，不是正整数
	for _, bad := range []string{"not-base64!!", "YWJj", "LTE=", "MA=="} {
		if _, err := decodeCursor(bad); err == nil {
			t.Errorf("非法游标 %q 应报错", bad)
		}
	}

	// 裸 id（未 base64）不能通过：base64("7") 才是合法游标
	if _, err := decodeCursor(strconv.FormatInt(7, 10)); err == nil {
		t.Error("未 base64 的 id 应报错")
	}
	if id, err := decodeCursor(base64.StdEncoding.EncodeToString([]byte("7"))); err != nil || id != 7 {
		t.Errorf("base64(7) = (%d, %v), want (7, nil)", id, err)
	}
}

func TestNormalizeListSize(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{0, 20}, {-1, 20}, {1, 1}, {20, 20}, {100, 100}, {101, 100}, {1000, 100},
	}
	for _, c := range cases {
		if got := normalizeListSize(c.in); got != c.want {
			t.Errorf("normalizeListSize(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestBuildListResponseCursorAndHasMore(t *testing.T) {
	// 有下一页：游标指向末位条目
	products := []*entity.Products{{Id: 30}, {Id: 20}, {Id: 10}}
	res := buildListResponse(products, nil, true)
	if !res.HasMore {
		t.Error("has_more 应为 true")
	}
	if want := encodeCursor(10); res.NextCursor != want {
		t.Errorf("next_cursor = %q, want %q", res.NextCursor, want)
	}

	// 没有下一页：next_cursor 必须为空 ——「next_cursor 为空 ⟺ 没有更多」
	last := buildListResponse(products, nil, false)
	if last.NextCursor != "" {
		t.Errorf("末页 next_cursor = %q, want 空（末页不能给一个翻不到东西的游标）", last.NextCursor)
	}
	if last.HasMore {
		t.Error("末页 has_more 应为 false")
	}

	// 无数据：空切片（不是 nil）、无游标
	empty := emptyListResponse()
	if empty.List == nil || len(empty.List) != 0 {
		t.Error("空结果应返回空切片而非 nil")
	}
	if empty.NextCursor != "" || empty.HasMore {
		t.Error("空结果不应有 next_cursor / has_more")
	}
}
