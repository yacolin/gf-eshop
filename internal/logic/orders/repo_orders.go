package orders

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/dao"
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

// findOrderWithShardByOrderNo 按业务单号查询订单，并把订单所在分片一并返回。
// 同一条链路上的子表查询必须复用这个分片，避免各自推导导致读到不同月表。
// 订单不存在时返回 (nil, shard{}, nil)。
func findOrderWithShardByOrderNo(ctx context.Context, orderNo string) (*entity.Orders, shard, error) {
	sh, err := shardFromOrderNo(ctx, orderNo)
	if err != nil {
		return nil, shard{}, err
	}
	var o *entity.Orders
	err = model(ctx, sh, tableOrders).
		Where(dao.Orders.Columns().OrderNo, orderNo).
		Scan(&o)
	if err != nil {
		return nil, shard{}, err
	}
	// Phase 3 影子读：主读走主表时按比例再读分片做对账（只记日志/计数，不影响返回）
	shadowCompareOrder(ctx, orderNo, sh, o)
	return o, sh, nil
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
}

// countOrdersByFilter 按筛选条件统计订单数。
func countOrdersByFilter(ctx context.Context, f orderListFilter) (int64, error) {
	sh, err := shardForScan(ctx, "订单列表统计")
	if err != nil {
		return 0, err
	}
	n, err := listOrdersModel(ctx, sh, f).Count()
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

// pageOrders 按筛选条件分页查询订单（固定 id 倒序）。
func pageOrders(ctx context.Context, f orderListFilter, page, size int) ([]*entity.Orders, error) {
	sh, err := shardForScan(ctx, "订单列表分页")
	if err != nil {
		return nil, err
	}
	var list []*entity.Orders
	err = listOrdersModel(ctx, sh, f).
		Page(page, size).
		OrderDesc(dao.Orders.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// pageOrdersByCursor 走 keyset 分页：取 id 小于 beforeID 的一页（固定 id 倒序）。
// beforeID 为 0 表示首页。分表后每片各取一页再归并即可，不需要 offset。
func pageOrdersByCursor(ctx context.Context, f orderListFilter, beforeID int64, size int) ([]*entity.Orders, error) {
	sh, err := shardForScan(ctx, "订单游标分页")
	if err != nil {
		return nil, err
	}
	m := listOrdersModel(ctx, sh, f)
	if beforeID > 0 {
		m = m.WhereLT(dao.Orders.Columns().Id, beforeID)
	}
	var list []*entity.Orders
	err = m.
		OrderDesc(dao.Orders.Columns().Id).
		Limit(size).
		Scan(&list)
	if err != nil {
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

// ── 看板聚合（跨片能力见 shardForScan） ───────────────────────────────────

// countAllOrders 统计订单总数（不含已软删除）。
func countAllOrders(ctx context.Context) (int64, error) {
	sh, err := shardForScan(ctx, "订单总数统计")
	if err != nil {
		return 0, err
	}
	n, err := model(ctx, sh, tableOrders).Count()
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

// sumPaidAmount 统计已支付订单的实付金额合计（分）。
func sumPaidAmount(ctx context.Context) (int64, error) {
	sh, err := shardForScan(ctx, "营收合计")
	if err != nil {
		return 0, err
	}
	var row struct {
		Value int64 `orm:"value"`
	}
	err = model(ctx, sh, tableOrders).
		Fields("COALESCE(SUM(pay_amount), 0) AS value").
		Where(dao.Orders.Columns().PaymentStatus, "paid").
		Where("deleted_at IS NULL").
		Scan(&row)
	if err != nil {
		return 0, err
	}
	return row.Value, nil
}

// orderDayTrendRow 按日聚合的原始行。
type orderDayTrendRow struct {
	Date   string `orm:"date"`
	Count  int64  `orm:"count"`
	Amount int64  `orm:"amount"`
}

// dailyOrderTrend 统计 since（含）之后的按日订单数与金额，date 形如 08-20。
func dailyOrderTrend(ctx context.Context, since string) ([]orderDayTrendRow, error) {
	sh, err := shardForScan(ctx, "订单趋势统计")
	if err != nil {
		return nil, err
	}
	var rows []orderDayTrendRow
	err = model(ctx, sh, tableOrders).
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
	return rows, nil
}

// orderStatusCountRow 按状态聚合的原始行。
type orderStatusCountRow struct {
	Status string `orm:"status"`
	Value  int64  `orm:"value"`
}

// orderStatusDistribution 按状态统计订单数。
func orderStatusDistribution(ctx context.Context) ([]orderStatusCountRow, error) {
	sh, err := shardForScan(ctx, "订单状态分布统计")
	if err != nil {
		return nil, err
	}
	var rows []orderStatusCountRow
	err = model(ctx, sh, tableOrders).
		Fields(dao.Orders.Columns().Status, "COUNT(*) AS value").
		Group(dao.Orders.Columns().Status).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
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
