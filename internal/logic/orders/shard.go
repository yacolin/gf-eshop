package orders

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

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
		warnMonthlyMode(ctx)
		return shardModeMonthly
	default:
		shardWarnOnce.Do(func() {
			g.Log().Warningf(ctx,
				"orderShard.mode=%q 不是合法值（single/monthly），已按 single 处理", mode)
		})
		return shardModeSingle
	}
}

// warnMonthlyMode 在 monthly 模式下给出一次性提示（Phase 4/5 之后语义已确定）。
//
// 这里说的都是**当前真实行为**，不再是「尚未实现」：跨片列表与看板已经可用，
// 仍然会失败的是「没有 order_no 的 offset 分页」，运维动作里被禁的是 migrate。
func warnMonthlyMode(ctx context.Context) {
	shardWarnOnce.Do(func() {
		g.Log().Warningf(ctx,
			"orderShard.mode=monthly：写只落分片、读只走分片（主表已停写，分片读**不回落**主表）。"+
				"列表请用游标分页（无 order_no 的 offset 分页返回 %d）；"+
				"migrate 已被禁止（会拿过期主表覆盖分片）；"+
				"详见 docs/order-sharding-design.md §7 Phase 4/5",
			errcode.CodeOrderShardNotReady)
	})
}

// ── 纯计算：分片只由「时间」决定，与运行模式无关 ────────────────────────────
//
// 模式决定的是「读哪里、写哪里、要不要镜像」，不应该影响分片怎么算。
// Phase 3 的双写需要「主表 + 分片」同时写，因此把计算与模式解耦。

// shardOfCreatedAt 按创建时间算分片。
func shardOfCreatedAt(t *gtime.Time) shard {
	if t == nil {
		t = gtime.Now()
	}
	return shard{suffix: t.Format("Ym")}
}

// shardOfOrderNo 按单号内嵌的 14 位时间戳算分片；解析失败返回 ok=false。
func shardOfOrderNo(orderNo string) (shard, bool) {
	t, ok := parseOrderNoTime(orderNo)
	if !ok {
		return shard{}, false
	}
	return shard{suffix: t.Format(shardSuffixGoFmt)}, true
}

// shardOfEncodedID 按「编码主键」算分片；老自增主键返回 ok=false（需查映射表）。
func shardOfEncodedID(id int64) (shard, bool) {
	t, ok := decodeOrderID(id)
	if !ok {
		return shard{}, false
	}
	return shard{suffix: t.Format(shardSuffixGoFmt)}, true
}

// ── 模式化入口：决定本次读写落在主表还是分片 ──────────────────────────────

// shardFromCreatedAt 主写入路径的分片推导（monthly 才落到分片）。
//
// 调用方必须传入「单一时间基准」——即同时用于生成 order_no、雪花主键与 created_at 的那一个 now。
// 否则月末边界会出现单号、主键、created_at 指向不同月份，读路径与写路径落到不同分片。
func shardFromCreatedAt(ctx context.Context, t *gtime.Time) shard {
	if shardMode(ctx) != shardModeMonthly {
		return shard{}
	}
	return shardOfCreatedAt(t)
}

// shardFromOrderNo 读取路径的分片推导：解析单号内嵌的 14 位时间戳。
func shardFromOrderNo(ctx context.Context, orderNo string) (shard, error) {
	if shardMode(ctx) != shardModeMonthly {
		return shard{}, nil
	}
	sh, ok := shardOfOrderNo(orderNo)
	if !ok {
		return shard{}, errcode.Newf(errcode.CodeInvalidParams,
			"订单号 %q 无法解析出创建时间，无法定位分片", orderNo)
	}
	return sh, nil
}

// ── Phase 3：双写与影子读的开关 ───────────────────────────────────────────

// dualWriteEnabled 是否开启双写（写主表的同时把同一订单镜像进分片）。
// 缺省 false：不配置就维持现状。
func dualWriteEnabled(ctx context.Context) bool {
	return g.Cfg().MustGet(ctx, "orderShard.dualWrite", false).Bool()
}

// mirrorShardOfCreatedAt 返回双写模式下需要额外镜像的分片；不需要时为零值。
// monthly 模式下主写本身就在分片，故不再镜像。
func mirrorShardOfCreatedAt(ctx context.Context, t *gtime.Time) shard {
	if !dualWriteEnabled(ctx) || shardMode(ctx) == shardModeMonthly {
		return shard{}
	}
	return shardOfCreatedAt(t)
}

// mirrorShardOfOrderNo 同上，按单号定位（用于改状态、支付回写这类已有单号的写入）。
func mirrorShardOfOrderNo(ctx context.Context, orderNo string) shard {
	if !dualWriteEnabled(ctx) || shardMode(ctx) == shardModeMonthly {
		return shard{}
	}
	sh, ok := shardOfOrderNo(orderNo)
	if !ok {
		return shard{}
	}
	return sh
}

// shadowReadPercent 影子读抽样比例（0~100），缺省 0 = 关闭。
func shadowReadPercent(ctx context.Context) int {
	p := g.Cfg().MustGet(ctx, "orderShard.shadowReadPercent", 0).Int()
	switch {
	case p <= 0:
		return 0
	case p > 100:
		return 100
	default:
		return p
	}
}

// ── Phase 4：读灰度 ───────────────────────────────────────────────────────
//
// 读路径按比例切到分片，是「切读」的灰度阶段：
//   - mode=single（双写阶段）：orderShard.readShardsPercent 控制比例，缺省 0；
//     命中分片但**没读到**时回落主表（此时主表仍是真相源），并计数；
//   - mode=monthly（终态）：恒定读分片、**不回落** —— 主表这时已停写，
//     回落到它只会读到过期数据（比如已支付却显示未支付）。

// readShardsPercent 读侧灰度比例（0~100）。
func readShardsPercent(ctx context.Context) int {
	p := g.Cfg().MustGet(ctx, "orderShard.readShardsPercent", 0).Int()
	switch {
	case p <= 0:
		return 0
	case p > 100:
		return 100
	default:
		return p
	}
}

// readFromShardByKey 点查是否走分片：monthly 恒真；single 下按比例做**确定性**抽样
// （复用 shouldShadow 的 crc32(key)%100，同一订单每次命中同一侧，便于复现）。
func readFromShardByKey(ctx context.Context, key string) bool {
	if shardMode(ctx) == shardModeMonthly {
		return true
	}
	p := readShardsPercent(ctx)
	if p <= 0 {
		return false
	}
	return shouldShadow(key, p)
}

// readShardsForList 列表/看板是否走分片。列表无法按单条订单做确定性抽样，
// 故 monthly 恒真，single 下按比例**每请求**抽样。
func readShardsForList(ctx context.Context) bool {
	if shardMode(ctx) == shardModeMonthly {
		return true
	}
	p := readShardsPercent(ctx)
	if p <= 0 {
		return false
	}
	if p >= 100 {
		return true
	}
	return grand.Intn(100) < p
}

// allowShardFallback 分片读未命中/失败时是否允许回落主表。
// 只有主表仍是真相源（single）时才允许；monthly 下主表已停写，回落会读到过期数据。
func allowShardFallback(ctx context.Context) bool {
	return shardMode(ctx) != shardModeMonthly
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
