package search

// 本文件封装可跨业务复用的 analysis 配置与字段映射。
//
// 为什么需要 ngram 子字段：MySQL 侧原本用 LIKE '%x%' 做子串匹配，而 IK 分词
// 是「按词」切分的 —— 搜「苹」不会命中「苹果」。因此每个需要子串能力的
// text 字段都额外挂一个 ngram 子字段：
//
//	索引侧 ngram_index  : 把 "苹果" 切成 苹 / 苹果 / 果，把 "Apple" 切成 a/ap/app/...
//	检索侧 ngram_search : 用 keyword tokenizer，把用户输入当成一个整体 token
//
// 于是搜索「苹」「app」「苹果」都能命中，语义回到 LIKE '%x%'，
// 同时保留 IK 的词级相关度（name 主字段）。
const (
	// DefaultNGramMax ngram 上限。品牌名 varchar(100)，取 32 足以覆盖
	// 常见品牌名，同时避免索引过度膨胀。
	DefaultNGramMax = 32
)

// IndexSettings 返回一套通用索引 settings：IK 分词 + ngram 子字段分析器。
// maxGram 控制子串匹配的最大长度（输入超过该长度将无法命中）。
func IndexSettings(maxGram int) map[string]any {
	if maxGram <= 0 {
		maxGram = DefaultNGramMax
	}
	return map[string]any{
		// ngram tokenizer 默认要求 max_gram-min_gram <= 1，必须显式放开
		"max_ngram_diff": maxGram,
		"analysis": map[string]any{
			"tokenizer": map[string]any{
				"ngram_tok": map[string]any{
					"type":        "ngram",
					"min_gram":    1,
					"max_gram":    maxGram,
					"token_chars": []string{"letter", "digit"},
				},
			},
			"analyzer": map[string]any{
				// 索引侧：产出 1..maxGram 的全部 n-gram
				"ngram_index": map[string]any{
					"type":      "custom",
					"tokenizer": "ngram_tok",
					"filter":    []string{"lowercase"},
				},
				// 检索侧：整体作为一个 token，避免把「苹果」拆成 苹/果 造成过度命中
				"ngram_search": map[string]any{
					"type":      "custom",
					"tokenizer": "keyword",
					"filter":    []string{"lowercase"},
				},
			},
		},
	}
}

// TextWithNGram 返回支持「中文分词 + 子串/前缀」的 text 字段映射。
//
// 三路子字段用法：
//
//	name          → match 查询，IK 词级相关度（可 _score 排序）
//	name.ngram    → match 查询，等价 LIKE '%x%'，子串/前缀/英文部分词
//	name.kw       → term 查询，精确值
func TextWithNGram() map[string]any {
	return map[string]any{
		"type":            "text",
		"analyzer":        "ik_max_word",
		"search_analyzer": "ik_smart",
		"fields": map[string]any{
			"ngram": map[string]any{
				"type":            "text",
				"analyzer":        "ngram_index",
				"search_analyzer": "ngram_search",
			},
			"kw": map[string]any{
				"type":         "keyword",
				"ignore_above": 256,
			},
		},
	}
}

// TextWithNGramNotIndexed 返回只存储、不参与检索的 text 字段映射，
// 用于需要原样回显但不需要搜索的长文本（如商品详情、品牌故事）。
func TextWithNGramNotIndexed() map[string]any {
	return map[string]any{"type": "text", "index": false}
}

// StoredOnly 返回只存储不解析的字段映射。
// 用于 created_at/deleted_at 这类时间字段：写进 _source 供接口原样回显，
// 但完全不参与索引与解析，因此不受 JSON 时间格式差异影响。
func StoredOnly() map[string]any {
	return map[string]any{"type": "object", "enabled": false}
}
