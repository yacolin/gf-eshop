package products

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	productEntityTTL = 10 * time.Minute
)

func cacheKeyProduct(id int64) string { return fmt.Sprintf("product:%d", id) }

// --- 单条缓存 ---

func getProductEntityCache(ctx context.Context, id int64) (*entity.Products, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyProduct(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var p entity.Products
	if err := sonic.Unmarshal(v.Bytes(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func setProductEntityCache(ctx context.Context, p *entity.Products) error {
	data, err := sonic.Marshal(p)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyProduct(p.Id), int(productEntityTTL.Seconds()), string(data))
	return err
}

func delProductEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyProduct(id))
}

// --- 缓存重建（启动预热） ---

var rebuildOnce sync.Once

func Warmup(ctx context.Context) {
	rebuildOnce.Do(func() {
		var list []*entity.Products
		if err := dao.Products.Ctx(ctx).OrderDesc(dao.Products.Columns().Id).Scan(&list); err != nil {
			g.Log().Warningf(ctx, "product cache warmup failed: %v", err)
			return
		}
		for _, p := range list {
			data, _ := sonic.Marshal(p)
			g.Redis().Do(ctx, "SETEX", cacheKeyProduct(p.Id), int(productEntityTTL.Seconds()), string(data))
		}
		g.Log().Infof(ctx, "product cache warmed up: %d items", len(list))
	})
}
