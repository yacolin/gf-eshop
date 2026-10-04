package orders

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

// 本文件是 tx_sub_orders 子订单表的数据访问层。
// 子订单与父订单一同分片，分片由调用方传入（见 shard.go 与 repo_orders.go 的说明）。

// listSubOrdersByParentOrderID 查询父订单下的全部子订单。
// 改造前后都按 parent_order_id 查询，SQL 逐字一致；分片由父订单链路传入。
func listSubOrdersByParentOrderID(ctx context.Context, sh shard, parentOrderID int64) ([]*entity.SubOrders, error) {
	var list []*entity.SubOrders
	err := model(ctx, sh, tableSubOrders).
		Where(dao.SubOrders.Columns().ParentOrderId, parentOrderID).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// insertSubOrder 写入子订单并返回自增主键。
func insertSubOrder(ctx context.Context, sh shard, data g.Map) (int64, error) {
	return model(ctx, sh, tableSubOrders).InsertAndGetId(data)
}

// updateSubOrdersByParentOrderID 更新父订单下的全部子订单。
func updateSubOrdersByParentOrderID(ctx context.Context, sh shard, parentOrderID int64, data g.Map) error {
	_, err := model(ctx, sh, tableSubOrders).
		Where(dao.SubOrders.Columns().ParentOrderId, parentOrderID).
		Data(data).
		Update()
	return err
}

// updateSubOrdersByParentOrderNo 按父订单号更新全部子订单（支付回调路径使用）。
func updateSubOrdersByParentOrderNo(ctx context.Context, sh shard, parentOrderNo string, data g.Map) error {
	_, err := model(ctx, sh, tableSubOrders).
		Where(dao.SubOrders.Columns().ParentOrderNo, parentOrderNo).
		Data(data).
		Update()
	return err
}
