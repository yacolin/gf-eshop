package categories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

const (
	categoryAllTTL    = 10 * time.Minute
	categoryEntityTTL = 10 * time.Minute
)

func cacheKeyCategoryAll() string          { return "category:all" }
func cacheKeyCategory(id int64) string     { return fmt.Sprintf("category:%d", id) }

// getCategoryAllCache 获取全量类目列表缓存
func getCategoryAllCache(ctx context.Context) ([]*entity.Categories, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyCategoryAll())
	if err != nil || v.IsNil() {
		return nil, err
	}
	var list []*entity.Categories
	if err := json.Unmarshal(v.Bytes(), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func setCategoryAllCache(ctx context.Context, list []*entity.Categories) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyCategoryAll(), int(categoryAllTTL.Seconds()), data)
	return err
}

func delCategoryAllCache(ctx context.Context) {
	g.Redis().Do(ctx, "DEL", cacheKeyCategoryAll())
}

// getCategoryEntityCache 获取单个类目缓存
func getCategoryEntityCache(ctx context.Context, id int64) (*entity.Categories, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyCategory(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var c entity.Categories
	if err := json.Unmarshal(v.Bytes(), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func setCategoryEntityCache(ctx context.Context, c *entity.Categories) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyCategory(c.Id), int(categoryEntityTTL.Seconds()), data)
	return err
}

func delCategoryEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyCategory(id))
}
