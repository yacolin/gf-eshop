package merchants

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	merchantEntityTTL = 10 * time.Minute
)

func cacheKeyMerchant(id int64) string { return fmt.Sprintf("merchant:%d", id) }

func getMerchantEntityCache(ctx context.Context, id int64) (*entity.Merchants, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyMerchant(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var m entity.Merchants
	if err := sonic.Unmarshal(v.Bytes(), &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func setMerchantEntityCache(ctx context.Context, m *entity.Merchants) error {
	data, err := sonic.Marshal(m)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyMerchant(m.Id), int(merchantEntityTTL.Seconds()), string(data))
	return err
}

func delMerchantEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyMerchant(id))
}

func rebuildMerchantCache(ctx context.Context) error {
	var list []*entity.Merchants
	if err := dao.Merchants.Ctx(ctx).OrderDesc(dao.Merchants.Columns().Id).Scan(&list); err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	for _, m := range list {
		data, _ := sonic.Marshal(m)
		g.Redis().Do(ctx, "SETEX", cacheKeyMerchant(m.Id), int(merchantEntityTTL.Seconds()), string(data))
	}
	return nil
}

func Warmup(ctx context.Context) {
	if err := rebuildMerchantCache(ctx); err != nil {
		g.Log().Warningf(ctx, "merchant cache warmup failed: %v", err)
		return
	}
	g.Log().Infof(ctx, "merchant cache warmed up")
}
