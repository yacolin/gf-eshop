package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Doc 一条待写入的文档。
type Doc struct {
	ID     string
	Source any
}

// SearchResult 是 _search 的精简结果：命中总数 + 原始 _source 列表。
type SearchResult struct {
	Total int64
	Hits  []json.RawMessage
}

// --- 内部工具 ---

// finish 统一收尾 esapi 响应：读取并关闭 body、维护熔断状态。
// 404 视为「索引/别名尚未建立」的正常状态，不触发熔断。
func finish(ctx context.Context, res *esapi.Response, err error) ([]byte, error) {
	if err != nil {
		markDown(ctx, err)
		return nil, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		markDown(ctx, readErr)
		return nil, readErr
	}
	if res.IsError() {
		detail := strings.TrimSpace(string(body))
		if res.StatusCode == http.StatusNotFound {
			return body, fmt.Errorf("%w: %s", ErrIndexNotFound, detail)
		}
		e := fmt.Errorf("elasticsearch %s: %s", res.Status(), detail)
		markDown(ctx, e)
		return body, e
	}
	markUp()
	return body, nil
}

// --- 索引生命周期 ---

// IndexExists 判断索引或别名是否存在。
func IndexExists(ctx context.Context, name string) (bool, error) {
	es, err := Client(ctx)
	if err != nil {
		return false, err
	}
	res, err := es.Indices.Exists([]string{name}, es.Indices.Exists.WithContext(ctx))
	if err != nil {
		markDown(ctx, err)
		return false, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		markUp()
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		e := fmt.Errorf("elasticsearch indices.exists %s", res.Status())
		markDown(ctx, e)
		return false, e
	}
}

// EnsureIndex 幂等地创建物理索引（已存在则跳过）。
func EnsureIndex(ctx context.Context, index string, settings, mappings map[string]any) error {
	exists, err := IndexExists(ctx, index)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	es, err := Client(ctx)
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]any{"settings": settings, "mappings": mappings})
	if err != nil {
		return err
	}
	res, err := es.Indices.Create(
		index,
		es.Indices.Create.WithContext(ctx),
		es.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	_, err = finish(ctx, res, err)
	if err != nil {
		return fmt.Errorf("create index %s: %w", index, err)
	}
	return nil
}

// DeleteIndex 删除物理索引；不存在时静默返回。
func DeleteIndex(ctx context.Context, index string) error {
	es, err := Client(ctx)
	if err != nil {
		return err
	}
	res, err := es.Indices.Delete(
		[]string{index},
		es.Indices.Delete.WithContext(ctx),
		es.Indices.Delete.WithIgnoreUnavailable(true),
	)
	_, err = finish(ctx, res, err)
	return err
}

// Refresh 强制刷新索引，让刚写入的文档立即可被检索。
func Refresh(ctx context.Context, index string) error {
	es, err := Client(ctx)
	if err != nil {
		return err
	}
	res, err := es.Indices.Refresh(es.Indices.Refresh.WithContext(ctx), es.Indices.Refresh.WithIndex(index))
	_, err = finish(ctx, res, err)
	return err
}

// Count 返回索引/别名下的文档数。
func Count(ctx context.Context, index string) (int64, error) {
	es, err := Client(ctx)
	if err != nil {
		return 0, err
	}
	res, err := es.Count(es.Count.WithContext(ctx), es.Count.WithIndex(index))
	body, err := finish(ctx, res, err)
	if err != nil {
		return 0, err
	}
	var parsed struct {
		Count int64 `json:"count"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, err
	}
	return parsed.Count, nil
}

// --- 别名 ---

// AliasedIndices 返回别名当前指向的物理索引列表。
func AliasedIndices(ctx context.Context, alias string) ([]string, error) {
	es, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	res, err := es.Indices.GetAlias(
		es.Indices.GetAlias.WithContext(ctx),
		es.Indices.GetAlias.WithName(alias),
	)
	body, err := finish(ctx, res, err)
	if err != nil {
		// 别名尚未建立属于「首次重建」的正常情况，视为当前没有已存在索引
		if errors.Is(err, ErrIndexNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(parsed))
	for name := range parsed {
		out = append(out, name)
	}
	return out, nil
}

// SwapAlias 把别名原子地指向 newIndex，并返回被摘除的旧索引。
// 先 add 后 remove，保证任意时刻别名都至少指向一个索引，检索侧不出现空窗。
func SwapAlias(ctx context.Context, alias, newIndex string) (removed []string, err error) {
	es, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	current, err := AliasedIndices(ctx, alias)
	if err != nil {
		return nil, err
	}
	actions := make([]map[string]any, 0, len(current)+1)
	actions = append(actions, map[string]any{"add": map[string]any{"index": newIndex, "alias": alias}})
	for _, idx := range current {
		if idx == newIndex {
			continue
		}
		actions = append(actions, map[string]any{"remove": map[string]any{"index": idx, "alias": alias}})
		removed = append(removed, idx)
	}
	body, err := json.Marshal(map[string]any{"actions": actions})
	if err != nil {
		return nil, err
	}
	res, err := es.Indices.UpdateAliases(
		bytes.NewReader(body),
		es.Indices.UpdateAliases.WithContext(ctx),
	)
	_, err = finish(ctx, res, err)
	if err != nil {
		return nil, fmt.Errorf("swap alias %s -> %s: %w", alias, newIndex, err)
	}
	return removed, nil
}

// --- 写入 ---

// BulkIndex 批量写入文档，返回成功条数。
func BulkIndex(ctx context.Context, index string, docs []Doc) (int, error) {
	if len(docs) == 0 {
		return 0, nil
	}
	es, err := Client(ctx)
	if err != nil {
		return 0, err
	}
	var buf bytes.Buffer
	for _, d := range docs {
		meta, err := json.Marshal(map[string]any{"index": map[string]any{"_id": d.ID}})
		if err != nil {
			return 0, err
		}
		src, err := json.Marshal(d.Source)
		if err != nil {
			return 0, err
		}
		buf.Write(meta)
		buf.WriteByte('\n')
		buf.Write(src)
		buf.WriteByte('\n')
	}
	res, err := es.Bulk(
		bytes.NewReader(buf.Bytes()),
		es.Bulk.WithContext(ctx),
		es.Bulk.WithIndex(index),
	)
	body, err := finish(ctx, res, err)
	if err != nil {
		return 0, err
	}
	var parsed struct {
		Errors bool `json:"errors"`
		Items  []map[string]struct {
			Status int `json:"status"`
			Error  any `json:"error"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, err
	}
	ok := 0
	var firstErr string
	for _, item := range parsed.Items {
		for _, r := range item {
			if r.Status >= 200 && r.Status < 300 {
				ok++
			} else if firstErr == "" {
				firstErr = fmt.Sprintf("status=%d error=%v", r.Status, r.Error)
			}
		}
	}
	if parsed.Errors {
		return ok, fmt.Errorf("bulk index %s: %d/%d failed: %s", index, len(docs)-ok, len(docs), firstErr)
	}
	return ok, nil
}

// IndexDoc 写入单条文档（双写路径）。刷新策略取 cfg.WriteRefresh：
// "wait_for" 保证写完立即可检索（代价是最多等一个 refresh_interval）。
func IndexDoc(ctx context.Context, index, id string, source any) error {
	es, err := Client(ctx)
	if err != nil {
		return err
	}
	body, err := json.Marshal(source)
	if err != nil {
		return err
	}
	opts := []func(*esapi.IndexRequest){
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(id),
	}
	if refresh := Cfg(ctx).WriteRefresh; refresh != "" {
		opts = append(opts, es.Index.WithRefresh(refresh))
	}
	res, err := es.Index(index, bytes.NewReader(body), opts...)
	_, err = finish(ctx, res, err)
	return err
}

// DeleteDoc 按 id 删除文档（不存在时静默成功）。
// 与 IndexDoc 使用同一刷新策略，保证「删完立刻查不到」，
// 否则刚删除的品牌会在下次 refresh 前继续出现在检索结果里。
func DeleteDoc(ctx context.Context, index, id string) error {
	es, err := Client(ctx)
	if err != nil {
		return err
	}
	opts := []func(*esapi.DeleteRequest){es.Delete.WithContext(ctx)}
	if refresh := Cfg(ctx).WriteRefresh; refresh != "" {
		opts = append(opts, es.Delete.WithRefresh(refresh))
	}
	res, err := es.Delete(index, id, opts...)
	if err != nil {
		markDown(ctx, err)
		return err
	}
	defer res.Body.Close()
	// 文档本就不存在（404）对删除语义来说就是成功
	if res.StatusCode == http.StatusNotFound {
		markUp()
		return nil
	}
	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		e := fmt.Errorf("elasticsearch delete %s/%s %s: %s", index, id, res.Status(), strings.TrimSpace(string(b)))
		markDown(ctx, e)
		return e
	}
	markUp()
	return nil
}

// Search 执行 _search 查询，返回命中总数与原始 _source。
func Search(ctx context.Context, index string, body map[string]any) (*SearchResult, error) {
	es, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(index),
		es.Search.WithBody(bytes.NewReader(raw)),
	)
	respBody, err := finish(ctx, res, err)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	out := &SearchResult{
		Total: parsed.Hits.Total.Value,
		Hits:  make([]json.RawMessage, 0, len(parsed.Hits.Hits)),
	}
	for _, h := range parsed.Hits.Hits {
		out.Hits = append(out.Hits, h.Source)
	}
	return out, nil
}
