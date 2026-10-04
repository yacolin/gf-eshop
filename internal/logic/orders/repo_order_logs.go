package orders

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// 本文件是 tx_order_logs 订单操作日志表的数据访问层。
// 日志与父订单一同分片（表内冗余了 order_no）；目前只有写入，没有读取路径。

// insertOrderLog 写入一条订单操作日志。
func insertOrderLog(ctx context.Context, sh shard, data g.Map) error {
	_, err := model(ctx, sh, tableOrderLogs).Insert(data)
	return err
}
