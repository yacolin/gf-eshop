package brands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	brandAllTTL     = 10 * time.Minute
	brandEntityTTL  = 10 * time.Minute
)

func cacheKeyBrandAll() string          { return "brand:all" }
func cacheKeyBrand(id int64) string     { return fmt.Sprintf("brand:%d", id) }

// Warmup 启动预热：加载全量品牌到 Redis
func Warmup(ctx context.Context) {
	var list []*entity.Brands
	if err := dao.Brands.Ctx(ctx).OrderAsc(dao.Brands.Columns().SortOrder).Scan(&list); err != nil {
		g.Log().Warningf(ctx, "brand cache warmup query failed: %v", err)
		return
	}
	if len(list) == 0 {
		return
	}
	data, err := json.Marshal(list)
	if err != nil {
		g.Log().Warningf(ctx, "brand cache warmup marshal failed: %v", err)
		return
	}
	if _, err := g.Redis().Do(ctx, "SETEX", cacheKeyBrandAll(), int(brandAllTTL.Seconds()), string(data)); err != nil {
		g.Log().Warningf(ctx, "brand cache warmup set failed: %v", err)
		return
	}
	// 顺便预热单条缓存
	for _, b := range list {
		item, _ := json.Marshal(b)
		g.Redis().Do(ctx, "SETEX", cacheKeyBrand(b.Id), int(brandEntityTTL.Seconds()), string(item))
	}
	g.Log().Infof(ctx, "brand cache warmed up: %d items", len(list))
}

// getBrandAllCache 获取全量品牌列表缓存
func getBrandAllCache(ctx context.Context) ([]*entity.Brands, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyBrandAll())
	if err != nil || v.IsNil() {
		return nil, err
	}
	var list []*entity.Brands
	if err := json.Unmarshal(v.Bytes(), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func setBrandAllCache(ctx context.Context, list []*entity.Brands) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyBrandAll(), int(brandAllTTL.Seconds()), string(data))
	return err
}

func delBrandAllCache(ctx context.Context) {
	g.Redis().Do(ctx, "DEL", cacheKeyBrandAll())
}

// getBrandEntityCache 获取单个品牌缓存
func getBrandEntityCache(ctx context.Context, id int64) (*entity.Brands, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyBrand(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var b entity.Brands
	if err := json.Unmarshal(v.Bytes(), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func setBrandEntityCache(ctx context.Context, b *entity.Brands) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyBrand(b.Id), int(brandEntityTTL.Seconds()), string(data))
	return err
}

func delBrandEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyBrand(id))
}
