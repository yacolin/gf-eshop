package orders

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/errcode"
)

// 本文件是订单域的**分片路由层**，完整方案见 docs/order-sharding-design.md。
//
// 订单域四张表（tx_orders / tx_sub_orders / tx_order_items / tx_order_logs）**同键同片**：
// 分片键是订单的 created_at 月份，物理表名形如 tx_orders_202608。
//
// Phase 0 只提供「接缝」：
//   - 默认 orderShard.mode=single，路由恒返回零值分片，物理表名不加后缀，
//     SQL 与分表前逐字一致（model() 与 dao.X.Ctx(ctx) 等价，见下方注释）。
//   - 切成 monthly 后，**能由单号推导出分片**的查询直接命中对应月表；
//     **推导不出分片**的查询（列表、看板聚合）当前会**明确报错**。
//
// 为什么推导不出就报错，而不是回落主表：主表在 Phase 5 之后会变成只读的旧表，
// 静默回落会读到过期数据 —— 这种错比一个明确的错误码难查得多。

const (
	tableOrders     = "tx_orders"
	tableSubOrders  = "tx_sub_orders"
	tableOrderItems = "tx_order_items"
	tableOrderLogs  = "tx_order_logs"
)

// 分表模式（配置项 orderShard.mode）。
const (
	// shardModeSingle 不分表：物理表名不加后缀。**默认值**，不配置即维持现状。
	shardModeSingle = "single"
	// shardModeMonthly 按月分表：tx_orders_202608。需要 Phase 2~4 全部落地后才可用。
	shardModeMonthly = "monthly"
)

// 单号内嵌时间的位置："ORD" + YYYYMMDDHHMMSS(14 位) + 序号。
// 历史单号后 4 位是随机数，Phase 1 起换成 6 位序列 —— 两者时间段位置相同，
// 所以**只认时间段**就能同时兼容新旧格式（实测历史 2000 条全部满足）。
// 相关常量（orderNoTimeStart/End、orderNoTimeGoFmt）与生成器一起放在 identity.go。

// shardSuffixGoFmt 分片后缀格式（如 202608）。
const shardSuffixGoFmt = "200601"

var shardWarnOnce sync.Once

// shard 标识一个订单分片；零值代表「不分表」。
type shard struct {
	suffix string // 形如 202608；零值为空串
}

// isZero 报告当前是否为不分表的零值分片。
func (s shard) isZero() bool { return s.suffix == "" }

// table 把逻辑表名解析成物理表名。
func (s shard) table(base string) string {
	if s.isZero() {
		return base
	}
	return base + "_" + s.suffix
}

// model 返回逻辑表 base 在当前分片下的 ORM 模型。
//
// 与 dao.X.Ctx(ctx) 完全等价（同样是 Model(table).Safe().Ctx(ctx)），差异只有表名：
//   - .Safe()  —— Update/Delete 无 WHERE 时报错，防止误伤全表；
//   - 软删除   —— GoFrame 在 SELECT 时按表字段自动补 deleted_at IS NULL（gdb_model_time.go）；
//   - 事务     —— ctx 里若已带事务，模型自动加入（gdb_core_underlying.go 的 TXFromCtx）。
//
// 表名在创建模型时就指定为最终物理表，因此软删除探测的是分片表自身的字段
// （分片表与原表结构一致，结论相同）。
func model(ctx context.Context, sh shard, base string) *gdb.Model {
	return g.DB().Model(sh.table(base)).Safe().Ctx(ctx)
}

// shardMode 读取 orderShard.mode，缺省 single（即不配置就维持现状）。
//
// 非法值按 single 处理并告警：一个拼写错误不该让整条订单链路不可用。
func shardMode(ctx context.Context) string {
	mode := g.Cfg().MustGet(ctx, "orderShard.mode", shardModeSingle).String()
	switch mode {
	case shardModeSingle:
		return shardModeSingle
	case shardModeMonthly:
		warnMonthlyNotReady(ctx)
		return shardModeMonthly
	default:
		shardWarnOnce.Do(func() {
			g.Log().Warningf(ctx,
				"orderShard.mode=%q 不是合法值（single/monthly），已按 single 处理", mode)
		})
		return shardModeSingle
	}
}

// warnMonthlyNotReady 在 monthly 模式下给出一次性告警：跨片能力尚未实现。
func warnMonthlyNotReady(ctx context.Context) {
	shardWarnOnce.Do(func() {
		g.Log().Warningf(ctx,
			"orderShard.mode=monthly：分片读写尚在 Phase 2~4 落地中，"+
				"当前只有「按单号/按创建时间」可正确路由，"+
				"列表与看板聚合会返回 %d；详见 docs/order-sharding-design.md",
			errcode.CodeOrderShardNotReady)
	})
}

// shardFromCreatedAt 写入路径的分片推导：按订单创建时间。
//
// 调用方必须传入「单一时间基准」——即同时用于生成 order_no、雪花主键与 created_at 的那一个 now。
// 否则月末边界会出现单号、主键、created_at 指向不同月份，读路径与写路径落到不同分片。
func shardFromCreatedAt(ctx context.Context, t *gtime.Time) shard {
	if shardMode(ctx) != shardModeMonthly {
		return shard{}
	}
	if t == nil {
		t = gtime.Now()
	}
	return shard{suffix: t.Format("Ym")}
}

// shardFromOrderNo 读取路径的分片推导：解析单号内嵌的 14 位时间戳。
func shardFromOrderNo(ctx context.Context, orderNo string) (shard, error) {
	if shardMode(ctx) != shardModeMonthly {
		return shard{}, nil
	}
	t, ok := parseOrderNoTime(orderNo)
	if !ok {
		return shard{}, errcode.Newf(errcode.CodeInvalidParams,
			"订单号 %q 无法解析出创建时间，无法定位分片", orderNo)
	}
	return shard{suffix: t.Format(shardSuffixGoFmt)}, nil
}

// shardForScan 用于无法由单号定位分片的查询（订单列表、看板聚合）。
//
// Phase 0 未实现跨片 fan-out 与日汇总表，monthly 下直接报错；
// 实现路径见 docs/order-sharding-design.md §5.5 / §5.6。
func shardForScan(ctx context.Context, op string) (shard, error) {
	if shardMode(ctx) != shardModeMonthly {
		return shard{}, nil
	}
	return shard{}, errcode.Newf(errcode.CodeOrderShardNotReady,
		"%s 需要跨分片查询，当前尚未实现（Phase 2/4）；"+
			"详见 docs/order-sharding-design.md §5.5/§5.6", op)
}

// parseOrderNoTime 解析单号中第 [3,17) 位的 14 位时间戳。
//
// 只校验时间段的合法性，不校验前缀（ORD / SUB 皆可），因此新旧单号格式通吃。
func parseOrderNoTime(orderNo string) (time.Time, bool) {
	if len(orderNo) < orderNoTimeEnd {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation(orderNoTimeGoFmt, orderNo[orderNoTimeStart:orderNoTimeEnd], time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
