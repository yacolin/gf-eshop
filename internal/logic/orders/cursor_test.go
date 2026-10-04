package orders

import (
	"testing"

	"gf-eshop/internal/errcode"
)

// 游标分页的契约：cursor = base64(id)、非法游标报错、每页条数默认/上限与 products 一致、
// 已移除的 page/page_size 必须报错而不是被静默忽略。

func TestDecodeOrderCursor(t *testing.T) {
	// 空游标代表首页
	if id, err := decodeOrderCursor(""); err != nil || id != 0 {
		t.Errorf("空游标 = (%d, %v), want (0, nil)", id, err)
	}

	// 正常往返
	c := encodeOrderCursor(123456)
	if id, err := decodeOrderCursor(c); err != nil || id != 123456 {
		t.Errorf("往返 = (%d, %v), want (123456, nil)", id, err)
	}

	// 非法游标必须报错：静默当首页会让调用方以为翻到了下一页
	//   not-base64!!  非法 base64 字符
	//   YWJj          base64("abc")，解出来不是数字
	//   LTE= / MA==   base64("-1") / base64("0")，不是正整数
	for _, bad := range []string{"not-base64!!", "YWJj", "LTE=", "MA=="} {
		if _, err := decodeOrderCursor(bad); err == nil {
			t.Errorf("非法游标 %q 应报错", bad)
		}
	}
}

func TestNormalizeListSize(t *testing.T) {
	// 默认值与 products 一致（20），上限 100
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

func TestCheckLegacyOffsetParams(t *testing.T) {
	// 没有这两个参数：放行
	if err := checkLegacyOffsetParams(func(string) bool { return false }); err != nil {
		t.Errorf("未传 page/page_size 不应报错，实际 %v", err)
	}

	// 任一出现即拒绝，且错误码是参数错误（不能让老客户端以为翻页成功）
	for _, key := range []string{"page", "page_size"} {
		err := checkLegacyOffsetParams(func(k string) bool { return k == key })
		if err == nil {
			t.Fatalf("%s 必须被拒绝", key)
		}
		if got := errcode.CodeOf(err); got != errcode.CodeInvalidParams {
			t.Errorf("%s 的错误码 = %d, want %d（%v）", key, got, errcode.CodeInvalidParams, err)
		}
	}
}
