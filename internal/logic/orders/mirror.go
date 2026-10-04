package orders

import (
	"context"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

// 本文件是订单分表 Phase 3 的能力：**双写镜像**与**影子读对账**。
// 方案见 docs/order-sharding-design.md §7。
//
// 双写为什么不做在每一条写入语句里：那要改 8 个写函数（4 张表 × insert/update），
// 任何一条漏改都会让分片永久缺数，而且每加一个写接口都要记得改。
// 这里改为**在同一个事务内**按订单把 4 张表的当前行从主表复制到分片：
//
//	INSERT INTO <分片> SELECT * FROM <主表> WHERE <订单键> = ?
//	ON DUPLICATE KEY UPDATE <全列覆盖>
//
// 好处：
//   - 只挂在 3 个编排点（建单 / 改状态 / 支付回写），新增写接口只要走同一入口就不会漏；
//   - 同事务 ⇒ 与主写原子（gdb 的 Exec 会取 ctx 里的事务，见 Core.DoExec）；
//   - 全列覆盖 ⇒ 幂等且自愈：上一次漏掉的更新会在下一次写入时被带上。

// mirrorTables 双写覆盖的四张表及其定位一张订单的键。
//
// 选键依据「该键上有索引」，避免在大表上全扫：
//   - tx_orders       → uk_order_no
//   - tx_sub_orders   → idx_parent_order_no
//   - tx_order_items  → idx_order_no
//   - tx_order_logs   → **没有 order_no 索引**（只有 idx_order_id），故按 order_id 定位
var mirrorTables = []struct {
	base string
	key  string
}{
	{tableOrders, "order_no"},
	{tableSubOrders, "parent_order_no"},
	{tableOrderItems, "order_no"},
	{tableOrderLogs, "order_id"},
}

// syncOrderToShard 把某订单在四张表里的当前行从主表镜像进分片。
//
// orderID 只用于 tx_order_logs（该表按 order_id 定位），传 0 时跳过日志表。
// 必须在写事务内调用才能与主写原子。返回错误由调用方决定如何处理——
// 双写失败只记日志，主表始终是唯一真相源（与 ES 双写同策略）。
func syncOrderToShard(ctx context.Context, orderNo string, orderID int64, sh shard) error {
	if sh.isZero() {
		return nil
	}
	for _, t := range mirrorTables {
		var keyVal interface{}
		if t.key == "order_id" {
			if orderID <= 0 {
				continue // 日志表按 id 定位，没有 id 就跳过（不影响订单主数据）
			}
			keyVal = orderID
		} else {
			keyVal = orderNo
		}

		upsert, err := upsertAllColumns(ctx, t.base)
		if err != nil {
			return err
		}
		sqlStr := "INSERT INTO `" + sh.table(t.base) + "` SELECT * FROM `" + t.base + "` " +
			"WHERE `" + t.key + "` = ?" + upsert
		if _, err = g.DB().Exec(ctx, sqlStr, keyVal); err != nil {
			return fmt.Errorf("镜像 %s → %s 失败: %w", t.base, sh.table(t.base), err)
		}
	}
	return nil
}

// mirrorOrderBestEffort 双写镜像：失败只记日志与计数，绝不阻断主链路。
func mirrorOrderBestEffort(ctx context.Context, orderNo string, orderID int64, sh shard) {
	if sh.isZero() {
		return
	}
	if err := syncOrderToShard(ctx, orderNo, orderID, sh); err != nil {
		shardMirrorCounter(ctx, "mirror_failed")
		g.Log().Warningf(ctx,
			"订单双写分片失败（不影响主表，主表是唯一真相源）: order_no=%s shard=%s err=%v",
			orderNo, sh.suffix, err)
		return
	}
	shardMirrorCounter(ctx, "mirror_ok")
}

// ── 影子读对账 ────────────────────────────────────────────────────────────

// orderFingerprint 把订单压成一行可读字符串，用于比对主表与分片的数据是否一致。
// 只取参与业务判断的列；时间统一格式化，避免精度/时区差异造成假阳性。
func orderFingerprint(o *entity.Orders) string {
	if o == nil {
		return "<nil>"
	}
	return strings.Join([]string{
		o.OrderNo,
		fmt.Sprint(o.UserId),
		fmt.Sprintf("%d/%d/%d/%d", o.TotalAmount, o.DiscountAmount, o.ShippingFee, o.PayAmount),
		o.Status,
		o.PaymentStatus,
		gtimeText(o.CreatedAt),
		gtimeText(o.UpdatedAt),
		gtimeText(o.DeletedAt),
	}, "|")
}

func gtimeText(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("Y-m-d H:i:s.u")
}

// shouldShadow 按订单号做**确定性**抽样：同一个订单每次都命中同一样本，
// 便于复现与排查（不用随机数，避免同一条订单时而比对时而不比对）。
func shouldShadow(orderNo string, percent int) bool {
	if percent >= 100 {
		return true
	}
	return int(crc32.ChecksumIEEE([]byte(orderNo))%100) < percent
}

// shadowCompareOrder 影子读：主读走主表时，抽样再读一次分片做对账。
//
// 只在「主读==主表」时有意义（monthly 模式主读已在分片）；任何失败都只记日志与计数，
// 绝不影响主流程返回值。
func shadowCompareOrder(ctx context.Context, orderNo string, primary shard, expected *entity.Orders) {
	if !primary.isZero() {
		return
	}
	percent := shadowReadPercent(ctx)
	if percent <= 0 || !shouldShadow(orderNo, percent) {
		return
	}
	sh, ok := shardOfOrderNo(orderNo)
	if !ok {
		return
	}

	var got *entity.Orders
	if err := model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().OrderNo, orderNo).
		Scan(&got); err != nil {
		shardShadowCounter(ctx, "error")
		g.Log().Debugf(ctx, "订单影子读失败（已忽略）: order_no=%s shard=%s err=%v",
			orderNo, sh.suffix, err)
		return
	}

	shardShadowCounter(ctx, "total")
	if a, b := orderFingerprint(expected), orderFingerprint(got); a != b {
		shardShadowCounter(ctx, "diff")
		g.Log().Warningf(ctx, "订单影子读不一致: order_no=%s 主表=%s 分片=%s", orderNo, a, b)
	}
}

// shardMirrorCounter / shardShadowCounter 用 Redis 计数，便于直接观察双写失败率与影子读差异率：
//
//	redis-cli MGET order:shard:mirror:mirror_ok order:shard:mirror:mirror_failed
//	redis-cli MGET order:shard:shadow:total order:shard:shadow:diff order:shard:shadow:error
//
// 计数失败（如 Redis 不可用）不影响主流程。
func shardMirrorCounter(ctx context.Context, kind string) {
	incrCounter(ctx, "order:shard:mirror:"+kind)
}

// incrCounter 尽力累加一个计数。
//
// Redis 不可用、或进程没注册 redis 驱动（例如单测二进制）时**静默跳过**：
// g.Redis() 在驱动缺失时会直接 panic，那会把「只该记个数的观测点」
// 变成主链路上的故障点。这里用 recover 兜住，兑现「计数失败不影响主流程」。
func incrCounter(ctx context.Context, key string) {
	defer func() { _ = recover() }()
	_, _ = g.Redis().Do(ctx, "INCR", key)
}

func shardShadowCounter(ctx context.Context, kind string) {
	incrCounter(ctx, "order:shard:shadow:"+kind)
}

// shardFallbackCounter 累计「分片读回落主表」的次数：灰度期间这个数应该很小且可解释，
// 长期不为 0 说明分片缺数据或分片表缺失。
func shardFallbackCounter(ctx context.Context, kind string) {
	incrCounter(ctx, "order:shard:fallback:"+kind)
}

// compareShardAgainstMain 反向影子读：主读走了分片时，再读一次主表做对账。
//
// 只在「主表仍是真相源」（single + 影子读抽样开启）时做：monthly 下主表已停写，
// 拿它当基准只会得到假差异。
func compareShardAgainstMain(ctx context.Context, orderNo string, shardOrder *entity.Orders) {
	if !allowShardFallback(ctx) || shadowReadPercent(ctx) <= 0 {
		return
	}
	var mainOrder *entity.Orders
	if err := model(ctx, shard{}, tableOrders).
		Where(dao.Orders.Columns().OrderNo, orderNo).
		Scan(&mainOrder); err != nil {
		shardShadowCounter(ctx, "error")
		return
	}
	shardShadowCounter(ctx, "total")
	if a, b := orderFingerprint(mainOrder), orderFingerprint(shardOrder); a != b {
		shardShadowCounter(ctx, "diff")
		g.Log().Warningf(ctx, "订单切读对账不一致: order_no=%s 主表=%s 分片=%s", orderNo, a, b)
	}
}
