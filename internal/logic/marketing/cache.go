package marketing

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/bytedance/sonic"
	"github.com/hashicorp/golang-lru/v2"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	promotionEntityTTL   = 10 * time.Minute
	promotionDetailTTL   = 5 * time.Minute
	promotionLocalSize   = 1024
	promotionLocalTTL    = 60 * time.Second
	promotionLocalJitter = 0.2
	promotionBloomN      = 50000
	promotionBloomP      = 0.01
)

type cacheEntry[T any] struct {
	val       T
	expiresAt time.Time
}

func (e *cacheEntry[T]) expired() bool {
	return time.Now().After(e.expiresAt)
}

type multiLevelCache struct {
	local      *lru.Cache[int64, *cacheEntry[*entity.Promotions]]
	bloom      *bloom.BloomFilter
	bloomMu    sync.RWMutex
	bloomReady bool
}

var campaignCache *multiLevelCache

func initCache() *multiLevelCache {
	local, _ := lru.New[int64, *cacheEntry[*entity.Promotions]](promotionLocalSize)
	return &multiLevelCache{
		local: local,
		bloom: bloom.NewWithEstimates(promotionBloomN, promotionBloomP),
	}
}

func ttlWithJitter(base time.Duration, jitter float64) time.Duration {
	delta := time.Duration(float64(base) * jitter * (2*rand.Float64() - 1))
	return base + delta
}

func (mc *multiLevelCache) bloomAdd(id int64) {
	mc.bloomMu.Lock()
	mc.bloom.AddString(fmt.Sprintf("%d", id))
	mc.bloomReady = true
	mc.bloomMu.Unlock()
}

func (mc *multiLevelCache) bloomMayExist(id int64) bool {
	mc.bloomMu.RLock()
	defer mc.bloomMu.RUnlock()
	if !mc.bloomReady {
		return true
	}
	return mc.bloom.TestString(fmt.Sprintf("%d", id))
}

func (mc *multiLevelCache) getLocal(id int64) (*entity.Promotions, bool) {
	e, ok := mc.local.Get(id)
	if !ok || e.expired() {
		return nil, false
	}
	return e.val, true
}

func (mc *multiLevelCache) setLocal(id int64, p *entity.Promotions) {
	mc.local.Add(id, &cacheEntry[*entity.Promotions]{
		val:       p,
		expiresAt: time.Now().Add(ttlWithJitter(promotionLocalTTL, promotionLocalJitter)),
	})
}

func (mc *multiLevelCache) delLocal(id int64) {
	mc.local.Remove(id)
}

type promotionDetail struct {
	Promotion *entity.Promotions          `json:"promotion"`
	Rule      *entity.PromotionRules      `json:"rule"`
	Products  []*entity.PromotionProducts `json:"products"`
}

func cacheKeyPromotion(id int64) string {
	return fmt.Sprintf("promotion:%d", id)
}

func cacheKeyPromotionDetail(id int64) string {
	return fmt.Sprintf("promotion:detail:%d", id)
}

func getRedisEntity(ctx context.Context, id int64) (*entity.Promotions, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyPromotion(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var p entity.Promotions
	if err := sonic.Unmarshal(v.Bytes(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func setRedisEntity(ctx context.Context, p *entity.Promotions) error {
	data, err := sonic.Marshal(p)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyPromotion(p.Id), int(promotionEntityTTL.Seconds()), string(data))
	return err
}

func delRedisEntity(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyPromotion(id))
}

func getRedisDetail(ctx context.Context, id int64) (*promotionDetail, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyPromotionDetail(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var d promotionDetail
	if err := sonic.Unmarshal(v.Bytes(), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func setRedisDetail(ctx context.Context, id int64, d *promotionDetail) error {
	data, err := sonic.Marshal(d)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyPromotionDetail(id), int(promotionDetailTTL.Seconds()), string(data))
	return err
}

func delRedisDetail(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyPromotionDetail(id))
}

func mustMarshal(v any) string {
	data, err := sonic.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Warmup(ctx context.Context) {
	var list []*entity.Promotions
	if err := dao.Promotions.Ctx(ctx).Scan(&list); err != nil {
		g.Log().Warningf(ctx, "promotion warmup failed: scan: %v", err)
		return
	}

	for _, p := range list {
		g.Redis().Do(ctx, "SETEX", cacheKeyPromotion(p.Id), int(promotionEntityTTL.Seconds()), mustMarshal(p))
		campaignCache.bloomAdd(p.Id)
		campaignCache.setLocal(p.Id, p)
	}

	g.Log().Infof(ctx, "promotion warmup complete: %d items", len(list))
}
