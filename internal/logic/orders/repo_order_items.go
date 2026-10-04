package orders

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

// 本文件是 tx_order_items 订单明细表的数据访问层。
// 明细与父订单一同分片，分片由调用方传入（见 shard.go 与 repo_orders.go 的说明）。

// listItemsByOrderID 查询订单下的全部明细。
func listItemsByOrderID(ctx context.Context, sh shard, orderID int64) ([]*entity.OrderItems, error) {
	var list []*entity.OrderItems
	err := model(ctx, sh, tableOrderItems).
		Where(dao.OrderItems.Columns().OrderId, orderID).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// insertOrderItem 写入一条订单明细。
func insertOrderItem(ctx context.Context, sh shard, data g.Map) error {
	_, err := model(ctx, sh, tableOrderItems).Insert(data)
	return err
}

// orderProductRankRow 按商品聚合的原始行。
type orderProductRankRow struct {
	ProductId int64 `orm:"product_id"`
	Count     int64 `orm:"count"`
	Amount    int64 `orm:"amount"`
}

// topProductsByOrderItems 按订单明细聚合商品销量与金额（TOP N）。
//
// 跨片时**每片返回全部分组**再归并，而不是每片各取 TOP N：
// 商品是按时间分片的（同一个商品分散在多个月），局部 TOP N 归并只能得到近似结果
// —— 一个销量分散在各月的商品可能进不了任何单片的 TOP N。全量分组归并才是精确的，
// 代价是这个查询最值得用汇总表/ES 优化（见 docs/order-sharding-design.md §5.6）。
func topProductsByOrderItems(ctx context.Context, limit int) ([]orderProductRankRow, error) {
	merged, err := aggregate(ctx, "top_products",
		func(ctx context.Context, sh shard) (map[int64]orderProductRankRow, error) {
			var rows []orderProductRankRow
			err := model(ctx, sh, tableOrderItems).
				Fields("product_id", "COUNT(*) AS count", "COALESCE(SUM(subtotal), 0) AS amount").
				Where("deleted_at IS NULL").
				Group(dao.OrderItems.Columns().ProductId).
				Scan(&rows)
			if err != nil {
				return nil, err
			}
			byProduct := make(map[int64]orderProductRankRow, len(rows))
			for _, r := range rows {
				byProduct[r.ProductId] = r
			}
			return byProduct, nil
		},
		func(acc *map[int64]orderProductRankRow, part map[int64]orderProductRankRow) {
			if *acc == nil {
				*acc = make(map[int64]orderProductRankRow, len(part))
			}
			for id, r := range part {
				cur := (*acc)[id]
				cur.ProductId = id
				cur.Count += r.Count
				cur.Amount += r.Amount
				(*acc)[id] = cur
			}
		})
	if err != nil {
		return nil, err
	}
	out := make([]orderProductRankRow, 0, len(merged))
	for _, r := range merged {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Amount > out[j].Amount
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
