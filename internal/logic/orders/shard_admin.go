package orders

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/errcode"
)

// 本文件是订单分表 Phase 2 的运维能力：建分片 → 迁移历史 → 三重对账。
// 方案与步骤见 docs/order-sharding-design.md §7；
// DDL（去自增 + 两张辅助表）归 schema 源仓库 std-eshop-db 管理：
// 基线 sql/tx_p5.sql，存量库迁移 sql/migrations/V001__tx_order_sharding.sql。
//
// 设计取向：
//   - 只做「增量、可重复执行」的事 —— 建表用 IF NOT EXISTS，写数据用
//     ON DUPLICATE KEY UPDATE，因此脚本中断后可以直接重跑。
//   - 不碰主表：迁移期间主表仍是唯一真相源，读路径完全不依赖分片。
//   - 对账不只看行数，还要比对金额合计与全字段校验和（行数相同但内容不同是最危险的）。

// legacyOrderIDMax 老数据自增主键的上界。
//
// 新主键是「分钟(25)|秒(6)|序列(20)」的编码：2026 年建的单约 8.7e13，
// 而 1<<40 ≈ 1.1e12 —— 凡是小于该值的 id 都不可能是编码主键，只可能是迁移前的自增主键。
// 反解时用它区分，避免把老 id（1..2000）误当成「2024-01-01 的单」而路由到错误的分片。
const legacyOrderIDMax = int64(1) << 40

// 分片运维涉及的四张表（子表必须与父订单同片）。
var shardedTables = []string{tableOrders, tableSubOrders, tableOrderItems, tableOrderLogs}

// ── 月份区间 ──────────────────────────────────────────────────────────────

// monthSpan 一个月的边界：[Start, End)。
type monthSpan struct {
	YM    string // 202608
	Start time.Time
	End   time.Time
}

// parseMonthSpans 解析 "2026-08" / "202608" 形式的闭区间月份。
func parseMonthSpans(from, to string) ([]monthSpan, error) {
	start, err := parseMonth(from)
	if err != nil {
		return nil, err
	}
	end, err := parseMonth(to)
	if err != nil {
		return nil, err
	}
	if end.Before(start) {
		return nil, errcode.Newf(errcode.CodeInvalidParams, "起始月份 %s 晚于结束月份 %s", from, to)
	}
	var spans []monthSpan
	for cur := start; !cur.After(end); cur = cur.AddDate(0, 1, 0) {
		spans = append(spans, monthSpan{
			YM:    cur.Format("200601"),
			Start: cur,
			End:   cur.AddDate(0, 1, 0),
		})
	}
	return spans, nil
}

func parseMonth(v string) (time.Time, error) {
	v = strings.TrimSpace(strings.ReplaceAll(v, "-", ""))
	t, err := time.ParseInLocation("200601", v, time.Local)
	if err != nil {
		return time.Time{}, errcode.Newf(errcode.CodeInvalidParams,
			"月份 %q 格式非法，应为 2026-08 或 202608", v)
	}
	return t, nil
}

// ── 建分片 ────────────────────────────────────────────────────────────────

// CreateShards 为 [from, to] 每个月创建四张分片表（幂等：已存在则跳过）。
func CreateShards(ctx context.Context, from, to string) ([]string, error) {
	spans, err := parseMonthSpans(from, to)
	if err != nil {
		return nil, err
	}
	var created []string
	for _, span := range spans {
		sh := shard{suffix: span.YM}
		for _, base := range shardedTables {
			table := sh.table(base)
			// 用 LIKE 复制主表结构，保证分片与主表结构永远一致
			if _, err = g.DB().Exec(ctx, "CREATE TABLE IF NOT EXISTS `"+table+"` LIKE `"+base+"`"); err != nil {
				return created, fmt.Errorf("创建分片表 %s 失败: %w", table, err)
			}
			created = append(created, table)
		}
	}
	return created, nil
}

// ── 迁移 ──────────────────────────────────────────────────────────────────

// MigrationReport 迁移结果汇总。
type MigrationReport struct {
	Months       []MonthMigration
	ShardMapRows int64
	StatsRows    int64
}

// MonthMigration 单月单表的迁移结果。
type MonthMigration struct {
	Month string
	Table string
	Rows  int64
}

// MigrateShards 把 [from, to] 的历史数据从主表复制进对应月份的分片表。
// batch 是单批行数（<=0 时取 5000），按主键区间推进，避免大事务与主从延迟。
func MigrateShards(ctx context.Context, from, to string, batch int) (*MigrationReport, error) {
	if batch <= 0 {
		batch = 5000
	}
	spans, err := parseMonthSpans(from, to)
	if err != nil {
		return nil, err
	}

	report := &MigrationReport{}
	for _, span := range spans {
		sh := shard{suffix: span.YM}
		for _, base := range shardedTables {
			n, err := migrateTable(ctx, base, sh, span, batch)
			if err != nil {
				return report, err
			}
			report.Months = append(report.Months, MonthMigration{Month: span.YM, Table: base, Rows: n})
		}
		// 老数据的主键不含时间位，必须登记映射才能按 id 定位分片
		n, err := recordLegacyShardMap(ctx, sh)
		if err != nil {
			return report, err
		}
		report.ShardMapRows += n
		// 日汇总从该月分片回填（每天必然落在唯一一个月分片里）
		m, err := backfillDailyStats(ctx, sh, span)
		if err != nil {
			return report, err
		}
		report.StatsRows += m
	}
	return report, nil
}

// migrateTable 按主键区间分批把一个月的行复制进分片表。
//
// 用 ON DUPLICATE KEY UPDATE 做幂等，而不是 INSERT IGNORE —— 后者会把
// CHECK 约束冲突降级成告警并丢行，等于静默丢数据。
// 并且这里是**全列覆盖**（见 upsertAllColumns），所以重跑 migrate 不只是「不重复」，
// 而是把主表的当前内容重新同步过去 —— 这正是设计文档 §7 Phase 3 说的
// 「迁移期间落在历史月份的零星更新」的重放方式。
func migrateTable(ctx context.Context, base string, sh shard, span monthSpan, batch int) (int64, error) {
	upsert, err := upsertAllColumns(ctx, base)
	if err != nil {
		return 0, err
	}
	var (
		target = sh.table(base)
		lastID int64
		total  int64
	)
	for {
		// 先取本批的最大 id，用它推进游标（不依赖分片表里已有的数据）
		v, err := g.DB().GetValue(ctx,
			"SELECT MAX(id) FROM (SELECT id FROM `"+base+"` "+
				"WHERE created_at >= ? AND created_at < ? AND id > ? ORDER BY id LIMIT ?) t",
			span.Start, span.End, lastID, batch)
		if err != nil {
			return total, fmt.Errorf("读取 %s(%s) 批次边界失败: %w", base, span.YM, err)
		}
		if v.IsNil() || v.Int64() == 0 {
			return total, nil
		}
		batchMax := v.Int64()

		res, err := g.DB().Exec(ctx,
			"INSERT INTO `"+target+"` SELECT * FROM `"+base+"` "+
				"WHERE created_at >= ? AND created_at < ? AND id > ? AND id <= ?"+upsert,
			span.Start, span.End, lastID, batchMax)
		if err != nil {
			return total, fmt.Errorf("迁移 %s(%s) 区间 (%d, %d] 失败: %w", base, span.YM, lastID, batchMax, err)
		}
		if affected, err := res.RowsAffected(); err == nil {
			total += affected
		}
		lastID = batchMax
	}
}

// recordLegacyShardMap 登记老主键 → 分片的映射（新主键可反解，不入表）。
func recordLegacyShardMap(ctx context.Context, sh shard) (int64, error) {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `tx_order_shard_map` (order_id, order_no, shard) "+
			"SELECT id, order_no, ? FROM `"+sh.table(tableOrders)+"` WHERE id < ? "+
			"ON DUPLICATE KEY UPDATE shard = VALUES(shard), order_no = VALUES(order_no)",
		sh.suffix, legacyOrderIDMax)
	if err != nil {
		return 0, fmt.Errorf("登记 %s 的老主键映射失败: %w", sh.suffix, err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// backfillDailyStats 从分片回填该月的订单日汇总（幂等覆盖式写入）。
func backfillDailyStats(ctx context.Context, sh shard, span monthSpan) (int64, error) {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `tx_order_daily_stats` "+
			"(stat_date, order_cnt, paid_cnt, cancelled_cnt, gmv, paid_amount) "+
			"SELECT DATE(created_at), COUNT(*), "+
			"COALESCE(SUM(payment_status = 'paid'), 0), "+
			"COALESCE(SUM(status = 'cancelled'), 0), "+
			"COALESCE(SUM(total_amount), 0), "+
			"COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN pay_amount ELSE 0 END), 0) "+
			"FROM `"+sh.table(tableOrders)+"` "+
			"WHERE created_at >= ? AND created_at < ? GROUP BY DATE(created_at) "+
			"ON DUPLICATE KEY UPDATE order_cnt = VALUES(order_cnt), paid_cnt = VALUES(paid_cnt), "+
			"cancelled_cnt = VALUES(cancelled_cnt), gmv = VALUES(gmv), paid_amount = VALUES(paid_amount)",
		span.Start, span.End)
	if err != nil {
		return 0, fmt.Errorf("回填 %s 日汇总失败: %w", sh.suffix, err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// ── 对账 ──────────────────────────────────────────────────────────────────

// VerificationReport 三重对账结果。
type VerificationReport struct {
	Checks  []VerifyCheck
	Failed  int
	AllPass bool
}

// VerifyCheck 单月单表的一项对账。
type VerifyCheck struct {
	Month  string
	Table  string
	RowsA  int64  // 主表
	RowsB  int64  // 分片
	SumA   string // 金额合计（主表）
	SumB   string // 金额合计（分片）
	SumCol string // 参与比对的金额列；多项时用 + 连接
	HashA  int64  // 全字段校验和（主表）
	HashB  int64  // 全字段校验和（分片）
	Pass   bool
	Detail string
}

// VerifyShards 对 [from, to] 的每个分片做三重对账：行数、金额合计、全字段校验和。
// 行数相同但内容不同是最危险的迁移事故，所以一定要比全字段校验和。
func VerifyShards(ctx context.Context, from, to string) (*VerificationReport, error) {
	spans, err := parseMonthSpans(from, to)
	if err != nil {
		return nil, err
	}

	report := &VerificationReport{AllPass: true}
	for _, span := range spans {
		sh := shard{suffix: span.YM}
		for _, base := range shardedTables {
			check, err := verifyTable(ctx, base, sh, span)
			if err != nil {
				return report, err
			}
			if !check.Pass {
				report.Failed++
				report.AllPass = false
			}
			report.Checks = append(report.Checks, check)
		}
	}
	return report, nil
}

func verifyTable(ctx context.Context, base string, sh shard, span monthSpan) (VerifyCheck, error) {
	check := VerifyCheck{Month: span.YM, Table: base}

	cols, err := tableColumns(ctx, base)
	if err != nil {
		return check, err
	}
	// 全字段校验和：两边都对同一批列做 CRC32 再求和
	hashExpr := "COALESCE(SUM(CRC32(CONCAT_WS('|', " + strings.Join(quoteAll(cols), ", ") + "))), 0)"

	baseWhere := " WHERE created_at >= ? AND created_at < ?"
	target := sh.table(base)

	if check.RowsA, err = scalarInt(ctx, "SELECT COUNT(*) FROM `"+base+"`"+baseWhere, span.Start, span.End); err != nil {
		return check, err
	}
	if check.RowsB, err = scalarInt(ctx, "SELECT COUNT(*) FROM `"+target+"`"); err != nil {
		return check, err
	}
	if check.HashA, err = scalarInt(ctx, "SELECT "+hashExpr+" FROM `"+base+"`"+baseWhere, span.Start, span.End); err != nil {
		return check, err
	}
	if check.HashB, err = scalarInt(ctx, "SELECT "+hashExpr+" FROM `"+target+"`"); err != nil {
		return check, err
	}

	// 金额合计（订单/子订单/明细各取语义上有意义的列）
	switch base {
	case tableOrders:
		check.SumCol = "pay_amount+total_amount+discount_amount"
	case tableSubOrders:
		check.SumCol = "pay_amount+total_amount"
	case tableOrderItems:
		check.SumCol = "subtotal+refund_amount"
	}
	if check.SumCol != "" {
		sumExpr := "COALESCE(SUM(" + check.SumCol + "), 0)"
		a, err := scalarInt(ctx, "SELECT "+sumExpr+" FROM `"+base+"`"+baseWhere, span.Start, span.End)
		if err != nil {
			return check, err
		}
		b, err := scalarInt(ctx, "SELECT "+sumExpr+" FROM `"+target+"`")
		if err != nil {
			return check, err
		}
		check.SumA, check.SumB = fmt.Sprint(a), fmt.Sprint(b)
	}

	check.Pass = check.RowsA == check.RowsB && check.HashA == check.HashB &&
		(check.SumCol == "" || check.SumA == check.SumB)
	if !check.Pass {
		check.Detail = fmt.Sprintf("行数 %d/%d，校验和 %d/%d，金额 %s/%s",
			check.RowsA, check.RowsB, check.HashA, check.HashB, check.SumA, check.SumB)
	}
	return check, nil
}

// tableColumns 取表的列名（按定义顺序），使对账能自动适配后续加列。
func tableColumns(ctx context.Context, table string) ([]string, error) {
	values, err := g.DB().GetArray(ctx,
		"SELECT COLUMN_NAME FROM information_schema.COLUMNS "+
			"WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", table)
	if err != nil {
		return nil, fmt.Errorf("读取 %s 列名失败: %w", table, err)
	}
	cols := make([]string, 0, len(values))
	for _, v := range values {
		cols = append(cols, v.String())
	}
	if len(cols) == 0 {
		return nil, errcode.Newf(errcode.CodeNotFound, "表 %s 不存在或没有列", table)
	}
	return cols, nil
}

// upsertClauseCache 缓存 ON DUPLICATE KEY UPDATE 子句：双写每次都要用，
// 不该每次都查 information_schema（4 张表 × 每条订单写入）。
var upsertClauseCache sync.Map // 表名 -> 子句

// upsertAllColumns 生成「全列覆盖」的 ON DUPLICATE KEY UPDATE 子句。
//
// 全列覆盖是刻意的：复制/重放必须让分片行与主表行**完全一致**，
// 若只写 `id = VALUES(id)` 之类的空更新，重跑就只是「不重复」而不会修正已漂移的行。
// 列名从 information_schema 动态取，因此加列后自动跟上。
func upsertAllColumns(ctx context.Context, table string) (string, error) {
	if v, ok := upsertClauseCache.Load(table); ok {
		return v.(string), nil
	}
	cols, err := tableColumns(ctx, table)
	if err != nil {
		return "", err
	}
	parts := make([]string, 0, len(cols))
	for _, c := range cols {
		// VALUES(col) 指「本该插入的那个值」；直接写 col = col 在 INSERT...SELECT
		// 下会被判为列名歧义（Error 1052）
		parts = append(parts, fmt.Sprintf("`%s` = VALUES(`%s`)", c, c))
	}
	clause := " ON DUPLICATE KEY UPDATE " + strings.Join(parts, ", ")
	upsertClauseCache.Store(table, clause)
	return clause, nil
}

func quoteAll(cols []string) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = "`" + c + "`"
	}
	return out
}

func scalarInt(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	v, err := g.DB().GetValue(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	if v.IsNil() {
		return 0, nil
	}
	return v.Int64(), nil
}

// ShardMapShardOf 查老主键落在哪个分片；新主键不需要查表（可反解）。
// 供读路径在 monthly 模式下定位老数据使用。
func ShardMapShardOf(ctx context.Context, orderID int64) (string, error) {
	v, err := g.DB().GetValue(ctx,
		"SELECT shard FROM `tx_order_shard_map` WHERE order_id = ?", orderID)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

// ShardOfOrderID 统一的「主键 → 分片」入口：新主键反解，老主键查映射表。
func ShardOfOrderID(ctx context.Context, orderID int64) (shard, error) {
	if orderID >= legacyOrderIDMax {
		if t, ok := decodeOrderID(orderID); ok {
			return shard{suffix: t.Format(shardSuffixGoFmt)}, nil
		}
	}
	ym, err := ShardMapShardOf(ctx, orderID)
	if err != nil {
		return shard{}, err
	}
	if ym == "" {
		return shard{}, errcode.Newf(errcode.CodeOrderShardNotReady,
			"订单主键 %d 既不是编码主键，也不在 tx_order_shard_map 中", orderID)
	}
	return shard{suffix: ym}, nil
}

// EnsureDailyStatsTable 保证日汇总表存在（幂等），供迁移脚本独立调用。
func EnsureDailyStatsTable(ctx context.Context) error {
	_, err := g.DB().Exec(ctx, `CREATE TABLE IF NOT EXISTS tx_order_daily_stats (
  stat_date     date     NOT NULL COMMENT '统计日期（Asia/Shanghai）',
  order_cnt     bigint   NOT NULL DEFAULT 0 COMMENT '下单数',
  paid_cnt      bigint   NOT NULL DEFAULT 0 COMMENT '支付笔数',
  refund_cnt    bigint   NOT NULL DEFAULT 0 COMMENT '退款笔数',
  cancelled_cnt bigint   NOT NULL DEFAULT 0 COMMENT '取消数',
  gmv           bigint   NOT NULL DEFAULT 0 COMMENT '下单金额（分）',
  paid_amount   bigint   NOT NULL DEFAULT 0 COMMENT '实收金额（分）',
  refund_amount bigint   NOT NULL DEFAULT 0 COMMENT '退款金额（分）',
  updated_at    datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (stat_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单日汇总（看板数据源）'`)
	return err
}

