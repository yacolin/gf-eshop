package middleware

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/os/gtime"
)

// 让响应里的时间字段统一输出「带毫秒」的字符串。
//
// 为什么要在响应出口做：GoFrame 的 gtime.Time.MarshalJSON() 走
// gtime_time_wrapper.go 的 wrapper.String()，那里**写死**了
// t.Format("2006-01-02 15:04:05")（秒级）。第三方方法改不了、gtime 没有全局格式开关，
// 而 gf gen dao 只支持自定义 dao 模板、不支持 entity 模板 —— 逐字段改要动
// 200+ 个 gtime 字段 / 72 个实体文件，不现实。而我们的响应出口
// （ErrorHandler 里的 r.GetHandlerResponse()）拿到的是**原始 Go 值**，能做一次
// 「类型驱动」的重写：只把 gtime 值换成毫秒字符串，其余照旧交给 encoding/json。
//
// 三条硬性设计约束（都有测试兜住）：
//  1. **没有改写就一定回退**：任何子树只要没发生时间改写，就整段交回 encoding/json，
//     因此「不含时间字段的响应」输出与改造前**逐字节一致**；
//  2. **保持字段顺序**：重写时按声明顺序拼装（含匿名嵌入展开），不做 map 中转，
//     否则键序会变，全接口 A/B 会全是噪声；
//  3. **没把握就回退**：字段名冲突、,string 等不支持的 tag、非字符串 key 的 map ——
//     一律整棵子树交回标准 encoder，宁可那里保持秒级，也不改变结构。

// millisLayout 是 **GoFrame 的格式串**（Y-m-d H:i:s），不是 Go 的 layout：
// gtime.Format 收前者，写成 "2006-01-02..." 会原样输出字面量。
const millisLayout = "Y-m-d H:i:s.u"

var (
	gtimeType    = reflect.TypeOf(gtime.Time{})
	gtimePtrType = reflect.TypeOf((*gtime.Time)(nil))
	marshalerTyp = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

	typeHasGTimeCache sync.Map // reflect.Type -> bool
)

// marshalWithMillis 序列化响应体；changed=false 表示没有任何时间改写，
// 此时返回的字节与 json.Marshal(v) 完全相同。
func marshalWithMillis(v interface{}) (data []byte, changed bool, err error) {
	if v == nil {
		b, err := json.Marshal(nil)
		return b, false, err
	}
	rv := reflect.ValueOf(v)
	if !typeHasGTime(rv.Type()) {
		b, err := json.Marshal(v)
		return b, false, err
	}
	b, ch, err := encodeValue(rv)
	if err != nil {
		return nil, false, err
	}
	if !ch {
		// 类型上可能有、实际没有：仍保证逐字节一致
		b, err = json.Marshal(v)
		return b, false, err
	}
	return b, true, nil
}

// marshalResponse 供测试使用：等价于 json.Marshal，但时间带毫秒。
func marshalResponse(v interface{}) ([]byte, error) {
	b, _, err := marshalWithMillis(v)
	return b, err
}

// isEmbeddedStruct 判断匿名字段是否（指向）结构体。
//
// 这类字段即使**类型未导出**，encoding/json 也会展开它的导出字段
// （例如 type outer struct{ inner } 里的 inner），所以不能按「未导出」跳过。
func isEmbeddedStruct(f reflect.StructField) bool {
	if !f.Anonymous {
		return false
	}
	ft := f.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	return ft.Kind() == reflect.Struct
}

// typeHasGTime 判断类型（含嵌套）是否可能出现 gtime 值；结果按类型缓存。
func typeHasGTime(t reflect.Type) bool {
	if v, ok := typeHasGTimeCache.Load(t); ok {
		return v.(bool)
	}
	res := computeTypeHasGTime(t)
	typeHasGTimeCache.Store(t, res)
	return res
}

func computeTypeHasGTime(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface:
		return true // 里面可能装任何东西：先认为有，真正处理时看具体值
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		return typeHasGTime(t.Elem())
	case reflect.Struct:
		if t == gtimeType {
			return true
		}
		// 自带 MarshalJSON 的类型自己决定输出，我们不进去（保持原行为）
		if t.Implements(marshalerTyp) || reflect.PointerTo(t).Implements(marshalerTyp) {
			return false
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			// 未导出字段会被 encoding/json 忽略；但匿名嵌入的结构体例外（它的导出字段会被展开）
			if f.PkgPath != "" && !isEmbeddedStruct(f) {
				continue
			}
			if name, opts := parseJSONTag(f.Tag.Get("json")); name == "-" || opts.unsupported {
				continue
			}
			if typeHasGTime(f.Type) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// encodeValue 把一个值编码成 JSON 片段，并报告是否发生了时间改写。
func encodeValue(rv reflect.Value) ([]byte, bool, error) {
	if !rv.IsValid() {
		return []byte("null"), false, nil
	}
	t := rv.Type()

	switch t {
	case gtimePtrType:
		if rv.IsNil() {
			return []byte("null"), false, nil
		}
		b, err := json.Marshal(formatGTime(rv.Interface().(*gtime.Time)))
		return b, true, err
	case gtimeType:
		gt := rv.Interface().(gtime.Time)
		b, err := json.Marshal(formatGTime(&gt))
		return b, true, err
	}

	// 不含时间的类型：整体交给标准 encoder（逐字节一致）
	if !typeHasGTime(t) {
		b, err := json.Marshal(rv.Interface())
		return b, false, err
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return []byte("null"), false, nil
		}
		return encodeValue(rv.Elem())
	case reflect.Interface:
		if rv.IsNil() {
			return []byte("null"), false, nil
		}
		return encodeValue(rv.Elem())
	case reflect.Slice, reflect.Array:
		return encodeSlice(rv)
	case reflect.Map:
		return encodeMap(rv)
	case reflect.Struct:
		return encodeStruct(rv)
	default:
		b, err := json.Marshal(rv.Interface())
		return b, false, err
	}
}

func encodeSlice(rv reflect.Value) ([]byte, bool, error) {
	var (
		buf     bytes.Buffer
		changed bool
	)
	buf.WriteByte('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		b, ch, err := encodeValue(rv.Index(i))
		if err != nil {
			return nil, false, err
		}
		buf.Write(b)
		changed = changed || ch
	}
	buf.WriteByte(']')
	if !changed {
		return mustMarshal(rv.Interface())
	}
	return buf.Bytes(), true, nil
}

func encodeMap(rv reflect.Value) ([]byte, bool, error) {
	if rv.Type().Key().Kind() != reflect.String {
		return mustMarshal(rv.Interface()) // 非字符串 key：交回标准 encoder
	}
	// encoding/json 对 map 按 key 排序输出，这里保持一致
	keys := rv.MapKeys()
	strs := make([]string, 0, len(keys))
	for _, k := range keys {
		strs = append(strs, k.String())
	}
	sort.Strings(strs)

	var (
		buf     bytes.Buffer
		changed bool
	)
	buf.WriteByte('{')
	for i, k := range strs {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, false, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		b, ch, err := encodeValue(rv.MapIndex(reflect.ValueOf(k)))
		if err != nil {
			return nil, false, err
		}
		buf.Write(b)
		changed = changed || ch
	}
	buf.WriteByte('}')
	if !changed {
		return mustMarshal(rv.Interface())
	}
	return buf.Bytes(), true, nil
}

// encodeStruct 按**声明顺序**拼装对象（含匿名嵌入展开），与 encoding/json 的字段顺序一致。
func encodeStruct(rv reflect.Value) ([]byte, bool, error) {
	fields, ok := visibleFields(rv.Type())
	if !ok {
		return mustMarshal(rv.Interface()) // 有歧义：整棵回退
	}
	var (
		buf     bytes.Buffer
		changed bool
		first   = true
	)
	buf.WriteByte('{')
	for _, f := range fields {
		fv, ok := fieldByIndexSafe(rv, f.index)
		if !ok {
			return mustMarshal(rv.Interface())
		}
		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		if !first {
			buf.WriteByte(',')
		}
		first = false
		kb, err := json.Marshal(f.name)
		if err != nil {
			return nil, false, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		b, ch, err := encodeValue(fv)
		if err != nil {
			return nil, false, err
		}
		buf.Write(b)
		changed = changed || ch
	}
	buf.WriteByte('}')
	if !changed {
		return mustMarshal(rv.Interface())
	}
	return buf.Bytes(), true, nil
}

// fieldByIndexSafe 等价于 reflect.Value.FieldByIndex，但遇到 nil 的匿名指针时返回 ok=false
// （标准库此时会 panic，我们选择回退）。
func fieldByIndexSafe(rv reflect.Value, index []int) (reflect.Value, bool) {
	for i, x := range index {
		if i > 0 && rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return reflect.Value{}, false
			}
			rv = rv.Elem()
		}
		rv = rv.Field(x)
	}
	return rv, true
}

func mustMarshal(v interface{}) ([]byte, bool, error) {
	b, err := json.Marshal(v)
	return b, false, err
}

// ── 字段可见性（与 encoding/json 对齐的常用子集）──────────────────────────

type visibleField struct {
	index     []int
	name      string
	omitEmpty bool
}

// visibleFields 收集结构体的可见字段，展开匿名嵌入。
//
// 只支持本项目实际用到的形状：单层匿名嵌入 + 无同名字段冲突。
// 一旦出现歧义就返回 ok=false，调用方整棵回退 —— 宁可保持秒级，也不猜。
func visibleFields(t reflect.Type) ([]visibleField, bool) {
	var (
		out  []visibleField
		seen = map[string]bool{}
	)
	var walk func(rt reflect.Type, prefix []int) bool
	walk = func(rt reflect.Type, prefix []int) bool {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.PkgPath != "" && !isEmbeddedStruct(f) {
				continue
			}
			name, opts := parseJSONTag(f.Tag.Get("json"))
			if name == "-" {
				continue
			}
			if opts.unsupported {
				return false
			}
			index := append(append([]int{}, prefix...), i)

			if f.Anonymous {
				// 带显式 json 名的匿名字段：交给标准 encoder
				if name != "" {
					return false
				}
				ft := f.Type
				if ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				if ft.Kind() != reflect.Struct {
					return false
				}
				if !walk(ft, index) {
					return false
				}
				continue
			}
			if name == "" {
				name = f.Name
			}
			if seen[name] { // 冲突：不猜
				return false
			}
			seen[name] = true
			out = append(out, visibleField{index: index, name: name, omitEmpty: opts.omitEmpty})
		}
		return true
	}
	if !walk(t, nil) {
		return nil, false
	}
	return out, true
}

type jsonTagOpts struct {
	omitEmpty   bool
	unsupported bool // 例如 ,string：不支持就整棵回退
}

func parseJSONTag(tag string) (name string, opts jsonTagOpts) {
	if tag == "" {
		return "", opts
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, p := range parts[1:] {
		switch p {
		case "omitempty":
			opts.omitEmpty = true
		case "":
		default:
			opts.unsupported = true
		}
	}
	return name, opts
}

// isEmptyValue 与 encoding/json 的 omitempty 判定一致。
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// formatGTime 把 gtime 值格式化成带毫秒的字符串。
// 零值沿用 GoFrame 的既有行为（空字符串），不额外制造变化。
func formatGTime(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(millisLayout)
}
