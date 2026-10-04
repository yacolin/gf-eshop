package orders

import (
	"encoding/base64"
	"strconv"

	"gf-eshop/internal/errcode"
)

// 本文件是订单列表的游标编解码，格式与 products 保持一致：
// cursor = base64(十进制 id)，排序固定 id DESC。
//
// 为什么要游标分页：分表后 offset 分页无法正确归并（每个分片各取 offset 再合并是错的），
// 而 COUNT(*) 需要跨 36 片 fan-out，深翻页也会越来越慢。
// 游标分页只需要「我读到哪了」，天然适配分片与 keyset。
//
// 与 products 的一个**故意差异**：非法游标返回错误而不是静默当成首页 ——
// 静默降级会让调用方以为翻到了下一页，实际又拿到第一页数据。

// encodeOrderCursor 生成游标（base64 编码的末位订单 ID）。
func encodeOrderCursor(id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

// decodeOrderCursor 解析游标；空字符串代表首页（返回 0）。
func decodeOrderCursor(cursor string) (int64, error) {
	if cursor == "" {
		return 0, nil
	}
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, errcode.Newf(errcode.CodeInvalidParams, "游标格式非法（应为 base64(id)）")
	}
	id, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil || id <= 0 {
		return 0, errcode.Newf(errcode.CodeInvalidParams, "游标内容非法: %q", string(b))
	}
	return id, nil
}
