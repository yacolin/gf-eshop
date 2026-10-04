package middleware

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// 这一组测试守两条线：
//  1. **不含时间字段的结构，输出必须与 encoding/json 逐字节一致**（否则就是改了全局行为）；
//  2. 含时间字段时，只把时间换成带毫秒的字符串，其它字段一个字节都不许动。

type inner struct {
	Name string `json:"name"`
	TS   int64  `json:"ts"`
}

type noTime struct {
	ID       int64             `json:"id"`
	Title    string            `json:"title"`
	Empty    string            `json:"empty,omitempty"`
	Zero     int               `json:"zero,omitempty"`
	NilPtr   *string           `json:"nil_ptr"`
	OmitNil  *int              `json:"omit_nil,omitempty"`
	List     []inner           `json:"list"`
	Map      map[string]string `json:"map"`
	Blob     []byte            `json:"blob"`
	Skip     string            `json:"-"`
	Nested   *inner            `json:"nested"`
	FlagTrue bool              `json:"flag_true"`
	FlagOmit bool              `json:"flag_omit,omitempty"`
}

type embeddedNoTime struct {
	noTime
	Extra string `json:"extra"`
}

func ptr[T any](v T) *T { return &v }

func mkNoTime() noTime {
	return noTime{
		ID:       9007199254740993, // 超过 2^53-1，验证整数不被浮点化
		Title:    `标题 "x" <b>`,     // 验证转义不被改动
		NilPtr:   nil,
		List:     []inner{{Name: "a", TS: 1}},
		Map:      map[string]string{"k": "v"},
		Blob:     []byte{1, 2, 3},
		Skip:     "不该出现",
		FlagTrue: true,
	}
}

// TestNoTimeFieldsByteIdentical 不含时间字段时，毫秒化必须零影响。
func TestNoTimeFieldsByteIdentical(t *testing.T) {
	seven := 7
	cases := map[string]interface{}{
		"普通结构":     mkNoTime(),
		"指针结构":     ptr(mkNoTime()),
		"嵌入结构":     embeddedNoTime{noTime: mkNoTime(), Extra: "e"},
		"切片":       []noTime{mkNoTime(), mkNoTime()},
		"字符串键 map": map[string]inner{"a": {Name: "n", TS: 2}},
		"嵌套 map":   map[string][]noTime{"x": {mkNoTime()}},
		"指针切片":     []*noTime{ptr(mkNoTime())},
		"带值的可选字段": struct {
			P *int `json:"p,omitempty"`
		}{P: &seven},
		"标量":   123,
		"字符串":  "2026-10-04 14:57:50",
		"nil":  nil,
		"空切片":  []int{},
		"接口结构": interface{}(mkNoTime()),
	}
	for name, v := range cases {
		want, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("%s: 标准序列化失败: %v", name, err)
		}
		got, err := marshalResponse(v)
		if err != nil {
			t.Fatalf("%s: marshalResponse 失败: %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s: 输出被改动了\n  want %s\n  got  %s", name, want, got)
		}
	}
}

// ── 含时间字段的情形 ──────────────────────────────────────────────────────

type orderLike struct {
	ID        int64       `json:"id"`
	OrderNo   string      `json:"order_no"`
	CreatedAt *gtime.Time `json:"created_at"`
	PaidAt    *gtime.Time `json:"paid_at"`
	ClosedAt  *gtime.Time `json:"closed_at,omitempty"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

type detailLike struct {
	Order *orderLike   `json:"order"`
	Items []*orderLike `json:"items"`
	Meta  interface{}  `json:"meta"`
}

func TestMillisPrecision(t *testing.T) {
	created := gtime.NewFromStr("2026-10-04 14:57:50.560")
	paid := gtime.NewFromStr("2026-10-04 14:57:51.001")
	updated := gtime.New(time.Date(2026, 10, 4, 14, 57, 52, 0, time.Local)) // 毫秒为 0 → 补 .000

	got, err := marshalResponse(orderLike{
		ID:        1,
		OrderNo:   "ORD20261004145750000001",
		CreatedAt: created,
		PaidAt:    paid,
		ClosedAt:  nil, // omitempty：应被省略
		UpdatedAt: updated,
	})
	if err != nil {
		t.Fatalf("marshalResponse 失败: %v", err)
	}
	want := `{"id":1,"order_no":"ORD20261004145750000001",` +
		`"created_at":"2026-10-04 14:57:50.560",` +
		`"paid_at":"2026-10-04 14:57:51.001",` +
		`"updated_at":"2026-10-04 14:57:52.000"}`
	if string(got) != want {
		t.Errorf("毫秒输出不符\n  want %s\n  got  %s", want, got)
	}
}

func TestMillisEmbeddedAndNested(t *testing.T) {
	created := gtime.NewFromStr("2026-08-27 06:38:15.123")

	// 匿名嵌入指针（响应结构里大量使用：type XxxRes struct { *entity.Orders }）
	embedded := struct {
		*orderLike
		Total int `json:"total"`
	}{orderLike: &orderLike{ID: 2, CreatedAt: created}, Total: 1}

	got, err := marshalResponse(embedded)
	if err != nil {
		t.Fatalf("marshalResponse 失败: %v", err)
	}
	want := `{"id":2,"order_no":"","created_at":"2026-08-27 06:38:15.123",` +
		`"paid_at":null,"updated_at":null,"total":1}`
	if string(got) != want {
		t.Errorf("嵌入结构输出不符\n  want %s\n  got  %s", want, got)
	}

	// 深层嵌套：切片 + 接口字段
	d, err := marshalResponse(detailLike{
		Order: &orderLike{ID: 3, CreatedAt: created},
		Items: []*orderLike{{ID: 4, CreatedAt: created}},
		Meta:  map[string]interface{}{"at": created, "n": 1},
	})
	if err != nil {
		t.Fatalf("marshalResponse 失败: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(d, &parsed); err != nil {
		t.Fatalf("输出不是合法 JSON: %v", err)
	}
	order := parsed["order"].(map[string]interface{})
	if order["created_at"] != "2026-08-27 06:38:15.123" {
		t.Errorf("嵌套 order.created_at = %v", order["created_at"])
	}
	meta := parsed["meta"].(map[string]interface{})
	if meta["at"] != "2026-08-27 06:38:15.123" {
		t.Errorf("接口字段里的时间未被处理: %v", meta["at"])
	}
	item := parsed["items"].([]interface{})[0].(map[string]interface{})
	if item["created_at"] != "2026-08-27 06:38:15.123" {
		t.Error("切片元素里的时间未被处理")
	}
}

// TestZeroTimeKeepsOldBehavior 零值沿用 GoFrame 的既有行为（空字符串），不额外变化。
func TestZeroTimeKeepsOldBehavior(t *testing.T) {
	got, err := marshalResponse(orderLike{CreatedAt: &gtime.Time{}})
	if err != nil {
		t.Fatalf("marshalResponse 失败: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatalf("输出不是合法 JSON: %v", err)
	}
	if parsed["created_at"] != "" {
		t.Errorf("零值时间应输出空字符串（与改造前一致），实际 %v", parsed["created_at"])
	}
}

// TestFallbackKeepsStdlibOutput 遇到没把握的形状（,string 等）必须整棵回退、零改动。
func TestFallbackKeepsStdlibOutput(t *testing.T) {
	type tricky struct {
		ID int64       `json:"id,string"` // 不支持的 tag 选项
		At *gtime.Time `json:"at"`
	}
	v := tricky{ID: 5, At: gtime.NewFromStr("2026-10-04 14:57:50.560")}
	want, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("标准序列化失败: %v", err)
	}
	got, err := marshalResponse(v)
	if err != nil {
		t.Fatalf("marshalResponse 失败: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("回退分支应保持标准输出\n  want %s\n  got  %s", want, got)
	}
}
