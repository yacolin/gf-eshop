// Package search 提供 Elasticsearch 的共用基础设施：客户端、熔断降级、
// 索引与别名管理、批量重建。
//
// 各业务模块（brands / products / skus）只需提供自己的 mapping 与文档转换，
// 即可复用同一套检索与重建底座：
//
//	search.Alias(ctx, "brands")               // eshop_brands（检索走别名）
//	search.EnsureIndex(ctx, idx, settings, mapping)
//	search.BulkIndex(ctx, idx, docs)
//	search.SwapAlias(ctx, alias, idx)          // 原子切换，无空窗
//
// 设计要点：ES 是「加速器」而非「唯一真相源」。任何一步失败都会被熔断并
// 记录日志，由调用方回落到 MySQL，因此 ES 挂掉不会影响接口可用性。
package search

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	defaultTimeout         = 3 * time.Second
	defaultPrefix          = "eshop_"
	defaultMaxResultWindow = 10000
	// breakerCooldown 熔断冷却时间：ES 故障后在这段时间内直接走 DB 降级，
	// 避免每个请求都白等一次超时。
	breakerCooldown = 5 * time.Second
)

var (
	// ErrDisabled 表示配置里关闭了 ES。
	ErrDisabled = errors.New("elasticsearch disabled")
	// ErrUnavailable 表示 ES 处于熔断冷却期。
	ErrUnavailable = errors.New("elasticsearch unavailable")
	// ErrDeepPaging 表示 from+size 超出 max_result_window，需要回落 DB。
	ErrDeepPaging = errors.New("elasticsearch deep paging out of range")
	// ErrIndexNotFound 表示别名/索引尚未建立（首次重建之前）。
	ErrIndexNotFound = errors.New("elasticsearch index not found")
)

// Config 是 manifest/config/config.yaml 中 elasticsearch 段的映射。
type Config struct {
	Enabled         bool
	Addresses       []string
	Username        string
	Password        string
	IndexPrefix     string
	Timeout         time.Duration
	MaxResultWindow int
	// WriteRefresh 单条双写后的刷新策略：
	// "wait_for"（默认）= 写完立即可检索，代价是最多等一个 refresh_interval；
	// "" = 立即返回，文档在下次 refresh 后才可见（写入延迟最低）。
	WriteRefresh string
}

var (
	cfgOnce sync.Once
	cfgVal  Config
)

// Cfg 返回（首次调用时加载并缓存）ES 配置。
func Cfg(ctx context.Context) Config {
	cfgOnce.Do(func() {
		c := Config{
			Enabled:         g.Cfg().MustGet(ctx, "elasticsearch.enabled", false).Bool(),
			Addresses:       g.Cfg().MustGet(ctx, "elasticsearch.addresses").Strings(),
			Username:        g.Cfg().MustGet(ctx, "elasticsearch.username", "").String(),
			Password:        g.Cfg().MustGet(ctx, "elasticsearch.password", "").String(),
			IndexPrefix:     g.Cfg().MustGet(ctx, "elasticsearch.indexPrefix", defaultPrefix).String(),
			Timeout:         g.Cfg().MustGet(ctx, "elasticsearch.timeout", defaultTimeout).Duration(),
			MaxResultWindow: g.Cfg().MustGet(ctx, "elasticsearch.maxResultWindow", defaultMaxResultWindow).Int(),
			WriteRefresh:    g.Cfg().MustGet(ctx, "elasticsearch.writeRefresh", "wait_for").String(),
		}
		if len(c.Addresses) == 0 {
			c.Addresses = []string{"http://127.0.0.1:9200"}
		}
		if c.Timeout <= 0 {
			c.Timeout = defaultTimeout
		}
		if c.IndexPrefix == "" {
			c.IndexPrefix = defaultPrefix
		}
		if c.MaxResultWindow <= 0 {
			c.MaxResultWindow = defaultMaxResultWindow
		}
		cfgVal = c
	})
	return cfgVal
}

var (
	clientOnce sync.Once
	clientVal  *elasticsearch.Client
	clientErr  error
)

// Client 返回惰性初始化的 ES 客户端；未启用时返回 ErrDisabled。
func Client(ctx context.Context) (*elasticsearch.Client, error) {
	cfg := Cfg(ctx)
	if !cfg.Enabled {
		return nil, ErrDisabled
	}
	clientOnce.Do(func() {
		clientVal, clientErr = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
			Transport: &http.Transport{
				MaxIdleConns:          64,
				MaxIdleConnsPerHost:   32,
				IdleConnTimeout:       90 * time.Second,
				ResponseHeaderTimeout: cfg.Timeout,
				DialContext:           (&net.Dialer{Timeout: cfg.Timeout}).DialContext,
			},
		})
	})
	return clientVal, clientErr
}

// --- 熔断 ---

var breakerDownUntil atomic.Int64

// Available 判断当前是否值得尝试访问 ES。处于冷却期时返回 false，
// 调用方应立即回落 DB 而不是白等一次超时。
func Available(ctx context.Context) bool {
	if !Cfg(ctx).Enabled {
		return false
	}
	return time.Now().UnixNano() >= breakerDownUntil.Load()
}

func markDown(ctx context.Context, err error) {
	prev := breakerDownUntil.Swap(time.Now().Add(breakerCooldown).UnixNano())
	if prev == 0 {
		g.Log().Warningf(ctx, "elasticsearch 故障，%s 内降级到 MySQL: %v", breakerCooldown, err)
	}
}

func markUp() { breakerDownUntil.Store(0) }

// --- 索引命名 ---

// Alias 返回业务实体对外检索用的别名，如 eshop_brands。
// 检索与单条双写都走别名，全量重建时只切换别名指向，检索侧无空窗。
func Alias(ctx context.Context, entity string) string {
	return Cfg(ctx).IndexPrefix + entity
}

// VersionedIndex 返回物理索引名，如 eshop_brands_20240928180000。
func VersionedIndex(ctx context.Context, entity, suffix string) string {
	return fmt.Sprintf("%s%s_%s", Cfg(ctx).IndexPrefix, entity, suffix)
}
