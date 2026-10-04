package orders

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
)

// 本文件是 tx_orders 主表的数据访问层：读取（单条 / 条件列表）与写入（新增 / 更新 / 聚合统计）。
// 上层 orders.go 只做编排，不直接触碰 dao 与 g.DB()；
// 其他模块（payments / dashboard）一律经 service.Orders() 调用，保证订单表只有一个访问出口。
//
// 所有模型都经 model(ctx, sh, ...) 构造（见 shard.go），因此：
//   - single 模式下 SQL 与分表前逐字一致；
//   - monthly 模式下落到 tx_orders_YYYYMM。

// ── 事务 ────────────────────────────────────────────────────────────────

// withOrdersTx 在订单库事务中执行 fn。
// 跨模块调用方（如 payments）自己开事务，repo 内的方法通过 ctx 自动加入同一事务。
func withOrdersTx(ctx context.Context, fn func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Orders.Transaction(ctx, fn)
}

// ── 读取 ────────────────────────────────────────────────────────────────

// findOrderWithShardByOrderNo 按业务单号查询订单，并把**实际读取源**的分片一并返回。
//
// 返回值里的 shard 有两个用途：① 同一条链路上的子表查询必须复用同一个源，
// 避免「父单读主表、子表读分片」这种混源；② 它**不代表写目标** ——
// 写路径另有 shardFromOrderNo（由 mode 决定），两者在灰度期间可能不同。
//
// Phase 4 读灰度：按比例命中分片时先读分片；读不到或读失败时，只要主表仍是真相源
// （single）就回落主表并计数；monthly 下不回落（主表已停写，回落会读到过期数据）。
func findOrderWithShardByOrderNo(ctx context.Context, orderNo string) (*entity.Orders, shard, error) {
	sh, shardOK := shardOfOrderNo(orderNo)
	if !shardOK && shardMode(ctx) == shardModeMonthly {
		return nil, shard{}, errcode.Newf(errcode.CodeInvalidParams,
			"订单号 %q 无法解析出创建时间，无法定位分片", orderNo)
	}

	if shardOK && readFromShardByKey(ctx, orderNo) {
		o, err := readOrderIn(ctx, sh, orderNo)
		switch {
		case err != nil:
			shardFallbackCounter(ctx, "read_order_error")
			g.Log().Warningf(ctx, "分片读订单失败: order_no=%s shard=%s err=%v", orderNo, sh.suffix, err)
			if !allowShardFallback(ctx) {
				return nil, sh, err
			}
		case o != nil:
			// 切读阶段与主表对账（主表仍是真相源时才有意义）
			compareShardAgainstMain(ctx, orderNo, o)
			return o, sh, nil
		default:
			// 分片里没有这条：可能是镜像还没覆盖到。monthly 下这就是「不存在」
			shardFallbackCounter(ctx, "read_order_miss")
			if !allowShardFallback(ctx) {
				return nil, sh, nil
			}
		}
	}

	o, err := readOrderIn(ctx, shard{}, orderNo)
	if err != nil {
		return nil, shard{}, err
	}
	// Phase 3 影子读：主读走主表时抽样比对分片（只记日志/计数）
	shadowCompareOrder(ctx, orderNo, shard{}, o)
	return o, shard{}, nil
}

// readOrderIn 在指定分片（零值=主表）按单号读订单；不存在时返回 (nil, nil)。
func readOrderIn(ctx context.Context, sh shard, orderNo string) (*entity.Orders, error) {
	var o *entity.Orders
	if err := model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().OrderNo, orderNo).
		Scan(&o); err != nil {
		return nil, err
	}
	return o, nil
}

// findOrderByOrderNo 按业务单号查询订单；不存在时返回 (nil, nil)。
func findOrderByOrderNo(ctx context.Context, orderNo string) (*entity.Orders, error) {
	o, _, err := findOrderWithShardByOrderNo(ctx, orderNo)
	return o, err
}

// findOrderByIDIn 在**已知分片**内按主键查询订单；不存在时返回 (nil, nil)。
//
// 分片由调用方给出（例如建单后回读：分片刚由同一个 now 推导出来），
// 不做二次推导 —— 自增 ID / 雪花 ID 都无法可靠反推月份。
func findOrderByIDIn(ctx context.Context, sh shard, id int64) (*entity.Orders, error) {
	var o *entity.Orders
	err := model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().Id, id).
		Scan(&o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// ── 列表（分页 / 筛选） ──────────────────────────────────────────────────

// orderListFilter 订单列表筛选条件。约定与改造前一致：零值代表「该条件未传」。
type orderListFilter struct {
	UserID        int64  // >0 生效
	Status        string // != "" 生效
	PaymentStatus string // != "" 生效
	OrderNo       string // != "" 生效
	// Window 列表的时间窗口（分表后必须有：否则要跨全部活跃分片 fan-out）。
	// 由 resolveListWindow 解析，见 list_window.go。
	Window orderListWindow
}

// countOrdersByFilter 按筛选条件统计订单数。
//
// 只服务于 offset 分页（page/page_size）的 total，因此与 pageOrders 一样只在主表上做：
// 分片模式下 offset 分页本身无法正确跨片归并（见 offsetListShard）。
func countOrdersByFilter(ctx context.Context, f orderListFilter) (int64, error) {
	src, err := offsetListSource(ctx, f)
	if err != nil {
		return 0, err
	}
	switch {
	case src.Main:
		n, err := listOrdersModel(ctx, shard{}, f).Count()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case len(src.Shards) == 0:
		return 0, nil // 窗口内没有分片：空结果（不回落主表，否则会返回无关月份的行）
	default:
		n, err := listOrdersModel(ctx, src.Shards[0], f).Count()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	}
}

// listShardForPointLookup 列表条件里带了 order_no 时，它其实等价于点查：
// 直接定位到那一个月，单分片内的 offset / keyset 都精确，不需要跨片归并。
// 返回 false 表示这不是点查（或当前模式仍应读主表）。
func listShardForPointLookup(ctx context.Context, f orderListFilter) (shard, bool) {
	if f.OrderNo == "" || shardMode(ctx) != shardModeMonthly {
		return shard{}, false
	}
	return shardOfOrderNo(f.OrderNo)
}

// offsetListSource 决定 offset 分页（page/page_size）读哪里。
//
// offset 分页在**多分片**下无法正确归并（每片各取 offset 再合并是错的），但：
//   - 带 order_no 的查询等价于点查 → 单分片，精确；
//   - **带 month / 区间且落在单个分片内** → 也精确 —— 这是新增月份参数顺带恢复的能力：
//     前端只要带上月份，老的 page/page_size 页面就能继续用，不必立刻改游标；
//   - 多分片：主表仍是真相源（single）时回落主表并计数；monthly 下明确报错并提示带月份。
func offsetListSource(ctx context.Context, f orderListFilter) (listSource, error) {
	src, err := resolveListSource(ctx, f)
	if err != nil {
		return listSource{}, err
	}
	if src.Main || len(src.Shards) <= 1 {
		return src, nil
	}
	return listSource{}, errcode.Newf(errcode.CodeOrderShardNotReady,
		"offset 分页在跨分片时无法正确归并（本次命中 %d 个分片）："+
			"请改用游标分页（cursor + size），或带上 month（如 month=%s）把范围收窄到单个月；"+
			"若已知订单号，直接带 order_no 查询即可精确定位",
		len(src.Shards), f.Window.Months[0])
}

// pageOrders 按筛选条件 offset 分页查询订单（固定 id 倒序）。
func pageOrders(ctx context.Context, f orderListFilter, page, size int) ([]*entity.Orders, error) {
	src, err := offsetListSource(ctx, f)
	if err != nil {
		return nil, err
	}
	switch {
	case src.Main:
		return listOrdersPage(ctx, shard{}, f, page, size, 0)
	case len(src.Shards) == 0:
		return nil, nil
	default:
		return listOrdersPage(ctx, src.Shards[0], f, page, size, 0)
	}
}

// pageOrdersByCursor 走 keyset 分页：取 id 小于 beforeID 的一页（固定 id 倒序）。
//
// 分片模式下每片各取一页再按 id 归并即可 —— keyset 的正确性来自「全局前 N 条
// 必然出现在某片的局部前 N 条里」，所以每片取 size 条就够了，不需要取全量。
func pageOrdersByCursor(ctx context.Context, f orderListFilter, beforeID int64, size int) ([]*entity.Orders, error) {
	src, err := resolveListSource(ctx, f)
	if err != nil {
		return nil, err
	}
	switch {
	case src.Main:
		return listOrdersPage(ctx, shard{}, f, 0, size, beforeID)
	case len(src.Shards) == 0:
		return nil, nil // 窗口内没有分片：空页
	case len(src.Shards) == 1:
		// 只命中一个分片：不需要归并（这也是「带 month 更省」的原因）
		return listOrdersPage(ctx, src.Shards[0], f, 0, size, beforeID)
	}

	var merged []*entity.Orders
	for _, sh := range src.Shards {
		part, err := listOrdersPage(ctx, sh, f, 0, size, beforeID)
		if err != nil {
			if allowShardFallback(ctx) {
				shardFallbackCounter(ctx, "list_shard_error")
				g.Log().Warningf(ctx, "分片列表查询失败，本次回落主表: shard=%s err=%v", sh.suffix, err)
				return listOrdersPage(ctx, shard{}, f, 0, size, beforeID)
			}
			return nil, err
		}
		merged = append(merged, part...)
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i].Id > merged[j].Id })
	if len(merged) > size {
		merged = merged[:size]
	}
	return merged, nil
}

// listOrdersPage 在指定分片上取一页（offset 或 keyset 二选一）。
func listOrdersPage(ctx context.Context, sh shard, f orderListFilter, page, size int, beforeID int64) ([]*entity.Orders, error) {
	m := listOrdersModel(ctx, sh, f)
	if beforeID > 0 {
		m = m.WhereLT(dao.Orders.Columns().Id, beforeID)
	}
	if page > 0 {
		m = m.Page(page, size)
	} else {
		m = m.Limit(size)
	}
	var list []*entity.Orders
	if err := m.OrderDesc(dao.Orders.Columns().Id).Scan(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// listOrdersModel 组装列表查询的公共条件，保证 Count 与分页两条 SQL 的 WHERE 完全一致。
func listOrdersModel(ctx context.Context, sh shard, f orderListFilter) *gdb.Model {
	m := model(ctx, sh, tableOrders)
	if f.UserID > 0 {
		m = m.Where(dao.Orders.Columns().UserId, f.UserID)
	}
	if f.Status != "" {
		m = m.Where(dao.Orders.Columns().Status, f.Status)
	}
	if f.PaymentStatus != "" {
		m = m.Where(dao.Orders.Columns().PaymentStatus, f.PaymentStatus)
	}
	if f.OrderNo != "" {
		m = m.Where(dao.Orders.Columns().OrderNo, f.OrderNo)
	}
	// 闭开区间：created_at >= Start AND created_at < End
	// （"到某日"在解析时已 +1 天，避免时间部分造成的边界漏单）
	//
	// 例外：带 order_no 的点查**不套时间窗** —— 客服拿单号找单时不该被
	// 「默认只看近 N 个月」挡住（单号里已有时间，路由本来就能定位到单分片）。
	if f.OrderNo == "" && f.Window.Start != nil {
		m = m.WhereGTE(dao.Orders.Columns().CreatedAt, f.Window.Start)
	}
	if f.OrderNo == "" && f.Window.End != nil {
		m = m.WhereLT(dao.Orders.Columns().CreatedAt, f.Window.End)
	}
	return m
}

// ── 写入 ────────────────────────────────────────────────────────────────

// insertOrder 写入主订单。
// 主键由调用方在 data 里显式给出（全局唯一 ID，见 identity.go），不再依赖自增。
func insertOrder(ctx context.Context, sh shard, data g.Map) error {
	_, err := model(ctx, sh, tableOrders).Insert(data)
	return err
}

// updateOrderByID 在已知分片内按主键更新主订单。
func updateOrderByID(ctx context.Context, sh shard, id int64, data g.Map) error {
	_, err := model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().Id, id).
		Data(data).
		Update()
	return err
}

// updateOrderByOrderNo 在已知分片内按业务单号更新主订单。
func updateOrderByOrderNo(ctx context.Context, sh shard, orderNo string, data g.Map) error {
	_, err := model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().OrderNo, orderNo).
		Data(data).
		Update()
	return err
}

// ── 看板聚合（Phase 4：跨片 fan-out） ─────────────────────────────────────
//
// 分片模式下看板的聚合必须跨片归并。这里没有走 tx_order_daily_stats：
// 该表只覆盖「订单数 / 营收」两项（趋势的金额口径、状态分布、热销商品都没有对应列），
// 而且目前只有 Phase 2 的回填、没有事件驱动的累加，直接用会读到过期数据。
// 详见 docs/order-sharding-design.md §5.6 与 §7 Phase 4。

// aggregate 在「本次读路径选定的源」上执行聚合并归并：
// 分片模式逐片查询后 merge，否则只查主表。
//
// 任一分片查询失败且允许回落时，**整段重做主表**（而不是把已累加的部分再叠加全量，
// 那会重复计数）；monthly 下不允许回落，直接报错。
func aggregate[T any](
	ctx context.Context, op string,
	query func(ctx context.Context, sh shard) (T, error),
	merge func(acc *T, part T),
) (T, error) {
	var zero T
	sources, err := readSources(ctx, op)
	if err != nil {
		return zero, err
	}
	var acc T
	for _, sh := range sources {
		part, err := query(ctx, sh)
		if err != nil {
			if sh.isZero() || !allowShardFallback(ctx) {
				return zero, err
			}
			shardFallbackCounter(ctx, op)
			g.Log().Warningf(ctx, "%s 在分片 %s 上失败，回落主表重算: %v", op, sh.suffix, err)
			return query(ctx, shard{})
		}
		merge(&acc, part)
	}
	return acc, nil
}

// readSources 返回本次查询要读的源：分片模式返回全部活跃分片，否则返回 [主表]。
// 一个分片表都没有时退回主表并计数（避免灰度把看板读成 0）。
func readSources(ctx context.Context, op string) ([]shard, error) {
	if !readShardsForList(ctx) {
		return []shard{{}}, nil
	}
	shards, err := activeShards(ctx)
	if err != nil {
		return nil, err
	}
	if len(shards) == 0 {
		shardFallbackCounter(ctx, op+"_no_shard")
		return []shard{{}}, nil
	}
	return shards, nil
}

// countAllOrders 统计订单总数（不含已软删除）。
func countAllOrders(ctx context.Context) (int64, error) {
	return aggregate(ctx, "count_orders",
		func(ctx context.Context, sh shard) (int64, error) {
			n, err := model(ctx, sh, tableOrders).Count()
			return int64(n), err
		},
		func(acc *int64, part int64) { *acc += part })
}

// sumPaidAmount 统计已支付订单的实付金额合计（分）。
func sumPaidAmount(ctx context.Context) (int64, error) {
	return aggregate(ctx, "sum_paid_amount",
		func(ctx context.Context, sh shard) (int64, error) {
			var row struct {
				Value int64 `orm:"value"`
			}
			err := model(ctx, sh, tableOrders).
				Fields("COALESCE(SUM(pay_amount), 0) AS value").
				Where(dao.Orders.Columns().PaymentStatus, "paid").
				Where("deleted_at IS NULL").
				Scan(&row)
			return row.Value, err
		},
		func(acc *int64, part int64) { *acc += part })
}

// orderDayTrendRow 按日聚合的原始行。
type orderDayTrendRow struct {
	Date   string `orm:"date"`
	Count  int64  `orm:"count"`
	Amount int64  `orm:"amount"`
}

// dailyOrderTrend 统计 since（含）之后的按日订单数与金额，date 形如 08-20。
func dailyOrderTrend(ctx context.Context, since string) ([]orderDayTrendRow, error) {
	merged, err := aggregate(ctx, "daily_order_trend",
		func(ctx context.Context, sh shard) (map[string]orderDayTrendRow, error) {
			var rows []orderDayTrendRow
			err := model(ctx, sh, tableOrders).
				Fields(
					"DATE_FORMAT(created_at, '%m-%d') AS date",
					"COUNT(*) AS count",
					"COALESCE(SUM(pay_amount), 0) AS amount",
				).
				Where("created_at >= ?", since).
				Where("deleted_at IS NULL").
				Group("date").
				Order("date ASC").
				Scan(&rows)
			if err != nil {
				return nil, err
			}
			byDate := make(map[string]orderDayTrendRow, len(rows))
			for _, r := range rows {
				byDate[r.Date] = r
			}
			return byDate, nil
		},
		func(acc *map[string]orderDayTrendRow, part map[string]orderDayTrendRow) {
			if *acc == nil {
				*acc = make(map[string]orderDayTrendRow, len(part))
			}
			for date, r := range part {
				cur := (*acc)[date]
				cur.Date = date
				cur.Count += r.Count
				cur.Amount += r.Amount
				(*acc)[date] = cur
			}
		})
	if err != nil {
		return nil, err
	}
	dates := make([]string, 0, len(merged))
	for date := range merged {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	out := make([]orderDayTrendRow, 0, len(dates))
	for _, date := range dates {
		out = append(out, merged[date])
	}
	return out, nil
}

// orderStatusCountRow 按状态聚合的原始行。
type orderStatusCountRow struct {
	Status string `orm:"status"`
	Value  int64  `orm:"value"`
}

// orderStatusDistribution 按状态统计订单数。
func orderStatusDistribution(ctx context.Context) ([]orderStatusCountRow, error) {
	merged, err := aggregate(ctx, "order_status_dist",
		func(ctx context.Context, sh shard) (map[string]int64, error) {
			var rows []orderStatusCountRow
			err := model(ctx, sh, tableOrders).
				Fields(dao.Orders.Columns().Status, "COUNT(*) AS value").
				Group(dao.Orders.Columns().Status).
				Scan(&rows)
			if err != nil {
				return nil, err
			}
			byStatus := make(map[string]int64, len(rows))
			for _, r := range rows {
				byStatus[r.Status] += r.Value
			}
			return byStatus, nil
		},
		func(acc *map[string]int64, part map[string]int64) {
			if *acc == nil {
				*acc = make(map[string]int64, len(part))
			}
			for status, v := range part {
				(*acc)[status] += v
			}
		})
	if err != nil {
		return nil, err
	}
	out := make([]orderStatusCountRow, 0, len(merged))
	for status, v := range merged {
		out = append(out, orderStatusCountRow{Status: status, Value: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Status < out[j].Status })
	return out, nil
}

// ── 供跨模块调用的复合写入 ────────────────────────────────────────────────

// markOrderPaidByOrderNo 支付成功后回写主订单与子订单的支付状态（按业务单号定位分片）。
//
// 与改造前的差异：四次要写入的 paid_at / updated_at 现在共用同一个 now，
// 而不是各自取一次时间（改造前主订单与子订单的 paid_at 会差几微秒）。
// orderID 仅用于双写镜像 tx_order_logs（该表按 order_id 定位，没有 order_no 索引）。
func markOrderPaidByOrderNo(ctx context.Context, orderNo string, orderID int64) error {
	sh, err := shardFromOrderNo(ctx, orderNo)
	if err != nil {
		return err
	}
	now := gtime.Now()
	err = updateOrderByOrderNo(ctx, sh, orderNo, g.Map{
		"payment_status": "paid",
		"status":         "paid",
		"paid_at":        now,
		"updated_at":     now,
	})
	if err != nil {
		return err
	}
	if err = updateSubOrdersByParentOrderNo(ctx, sh, orderNo, g.Map{
		"status":     "paid",
		"paid_at":    now,
		"updated_at": now,
	}); err != nil {
		return err
	}
	// Phase 3 双写：把这一单的最新行镜像进分片（同事务，失败只记日志）
	mirrorOrderBestEffort(ctx, orderNo, orderID, mirrorShardOfOrderNo(ctx, orderNo))
	return nil
}
