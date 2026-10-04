package orders

import (
	"context"

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
func topProductsByOrderItems(ctx context.Context, limit int) ([]orderProductRankRow, error) {
	sh, err := shardForScan(ctx, "热销商品统计")
	if err != nil {
		return nil, err
	}
	var rows []orderProductRankRow
	err = model(ctx, sh, tableOrderItems).
		Fields("product_id", "COUNT(*) AS count", "COALESCE(SUM(subtotal), 0) AS amount").
		Where("deleted_at IS NULL").
		Group(dao.OrderItems.Columns().ProductId).
		Order("count DESC, amount DESC").
		Limit(limit).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
