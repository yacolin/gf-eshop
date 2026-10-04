package orders

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/errcode"
)

// 本文件是订单分表 Phase 6 的能力：**归档旧主表**。
// 方案见 docs/order-sharding-design.md §7 Phase 6。
//
// 切到 monthly（Phase 5）之后主表不再被写入，它只剩两个用途：回滚依据、以及历史快照。
// Phase 6 把它改名成 tx_orders_legacy_<YYYYMM> 留观 30 天再删：
//   - 改名而不是立刻删 ⇒ 30 天内一条 RENAME 就能回滚；
//   - 留着旧名会让人（和工具）误以为它还是真相源 ⇒ 所以必须改名。
//
// 这里刻意不做「自动删」：清理是独立的动作，且有 30 天窗口守卫。

// legacyPrefix 归档表名的后缀标识：tx_orders → tx_orders_legacy_202610
const legacyPrefix = "_legacy_"

// legacyTableName 返回某张主表归档后的表名。
func legacyTableName(base, suffix string) string {
	return base + legacyPrefix + suffix
}

// parseLegacySuffix 从归档表名解析出后缀（如 tx_orders_legacy_202610 → 202610）。
func parseLegacySuffix(table string) (string, bool) {
	i := strings.LastIndex(table, legacyPrefix)
	if i < 0 {
		return "", false
	}
	suffix := table[i+len(legacyPrefix):]
	if !validMonthSuffix(suffix) {
		return "", false
	}
	return suffix, true
}

func validMonthSuffix(s string) bool {
	if len(s) != 6 || !isAllDigits(s) {
		return false
	}
	_, err := time.ParseInLocation("200601", s, time.Local)
	return err == nil
}

// legacyPurgeAllowedAt 返回某个归档后缀最早允许被删除的时间。
//
// 设计文档要求「观察 30 天」，但 RENAME 的时刻无法从 schema 反推
// （RENAME 不改变 CREATE_TIME），所以规则定为**归档后缀所在月的下一个月再过完一整月**：
// legacy_202610 最早 2026-12-01 才能清理 —— 比 30 天更保守，且只看表名就能判断。
func legacyPurgeAllowedAt(suffix string) (time.Time, error) {
	if !validMonthSuffix(suffix) {
		return time.Time{}, errcode.Newf(errcode.CodeInvalidParams, "归档后缀 %q 不是 YYYYMM", suffix)
	}
	start, err := time.ParseInLocation("200601", suffix, time.Local)
	if err != nil {
		return time.Time{}, errcode.Newf(errcode.CodeInvalidParams, "归档后缀 %q 无法解析: %v", suffix, err)
	}
	return start.AddDate(0, 2, 0), nil
}

// ArchiveReport 归档结果：改名对照表 + 回滚命令。
type ArchiveReport struct {
	Suffix  string
	Renamed [][2]string // {旧名, 新名}
}

// RollbackSQL 返回一条可直接执行的回滚语句（反 RENAME）。
func (r *ArchiveReport) RollbackSQL() string {
	var pairs []string
	for _, p := range r.Renamed {
		pairs = append(pairs, fmt.Sprintf("`%s` TO `%s`", p[1], p[0]))
	}
	if len(pairs) == 0 {
		return ""
	}
	return "RENAME TABLE " + strings.Join(pairs, ", ") + ";"
}

// ArchiveBaseTables 把四张主表改名为归档表（Phase 6）。
//
// 前置检查（缺一不可，都是「改名之后会疼」的地方）：
//  1. `orderShard.mode=monthly` —— 主表必须已经停写，否则改名后写入直接失败；
//  2. 主表本身还在（没被归档过）；
//  3. **无损性**：主表里的每一行都必须已存在于某个分片里 ——
//     否则改名后这些行就只剩归档表能查（而应用只查分片）。
//     这一条用 LEFT JOIN 逐表确认，比「行数相等」严格（行数相等也可能内容不同）。
//
// 四张表用**一条** RENAME 语句改名：MySQL 对多表 RENAME 是原子的，
// 不会出现「订单表改完了、明细表还没改」的中间态。
func ArchiveBaseTables(ctx context.Context, suffix string) (*ArchiveReport, error) {
	if shardMode(ctx) != shardModeMonthly {
		return nil, errcode.Newf(errcode.CodeShardMaintenanceRefused,
			"归档只能在 orderShard.mode=monthly 时执行：主表还在被写的话，"+
				"改名会让写入直接失败。请先完成 Phase 5（写只落分片）")
	}
	if !validMonthSuffix(suffix) {
		return nil, errcode.Newf(errcode.CodeInvalidParams, "归档后缀 %q 不是 YYYYMM", suffix)
	}

	var renamed [][2]string
	for _, base := range shardedTables {
		exists, err := tableExists(ctx, base)
		if err != nil {
			return nil, err
		}
		if !exists {
			legacy := legacyTableName(base, suffix)
			if ok, _ := tableExists(ctx, legacy); ok {
				return nil, errcode.Newf(errcode.CodeShardMaintenanceRefused,
					"主表 %s 已被归档为 %s，无需重复归档", base, legacy)
			}
			return nil, errcode.Newf(errcode.CodeNotFound, "主表 %s 不存在，无法归档", base)
		}
		missing, err := baseRowsMissingInShards(ctx, base)
		if err != nil {
			return nil, err
		}
		if missing > 0 {
			return nil, errcode.Newf(errcode.CodeShardMaintenanceRefused,
				"%s 有 %d 行只存在于主表、分片里没有：改名后这些行将查不到。"+
					"请先迁移/回灌再归档", base, missing)
		}
		renamed = append(renamed, [2]string{base, legacyTableName(base, suffix)})
	}
	if len(renamed) == 0 {
		return nil, errcode.Newf(errcode.CodeNotFound, "没有可归档的主表")
	}

	pairs := make([]string, 0, len(renamed))
	for _, p := range renamed {
		pairs = append(pairs, fmt.Sprintf("`%s` TO `%s`", p[0], p[1]))
	}
	if _, err := g.DB().Exec(ctx, "RENAME TABLE "+strings.Join(pairs, ", ")); err != nil {
		return nil, fmt.Errorf("归档改名失败: %w", err)
	}
	return &ArchiveReport{Suffix: suffix, Renamed: renamed}, nil
}

// PurgeReport 清理（DROP）结果。
type PurgeReport struct {
	Suffix     string
	AllowedAt  time.Time
	Dropped    []string
	Statements []string // 未 drop 时（或提前预览）给出可直接执行的语句
}

// PurgeLegacyTables 删除归档表。安全设计：
//   - 不传 force 只**打印** DROP 语句（预览），不动数据；
//   - 传 force 也要等过 30 天窗口（见 legacyPurgeAllowedAt），早于窗口直接拒绝。
func PurgeLegacyTables(ctx context.Context, suffix string, force bool, now time.Time) (*PurgeReport, error) {
	allowedAt, err := legacyPurgeAllowedAt(suffix)
	if err != nil {
		return nil, err
	}
	report := &PurgeReport{Suffix: suffix, AllowedAt: allowedAt}

	var found []string
	for _, base := range shardedTables {
		legacy := legacyTableName(base, suffix)
		ok, err := tableExists(ctx, legacy)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		found = append(found, legacy)
		report.Statements = append(report.Statements, "DROP TABLE `"+legacy+"`;")
	}
	if len(found) == 0 {
		return nil, errcode.Newf(errcode.CodeNotFound, "没有找到后缀 %s 的归档表", suffix)
	}
	if !force {
		return report, nil
	}
	if now.Before(allowedAt) {
		return nil, errcode.Newf(errcode.CodeShardMaintenanceRefused,
			"归档 %s 的观察窗口未满：最早可清理时间 %s（现在 %s）。"+
				"确认无需回滚后再传 --force",
			suffix, allowedAt.Format("2006-01-02"), now.Format("2006-01-02"))
	}
	for _, table := range found {
		if _, err := g.DB().Exec(ctx, "DROP TABLE `"+table+"`"); err != nil {
			return report, fmt.Errorf("删除归档表 %s 失败: %w", table, err)
		}
		report.Dropped = append(report.Dropped, table)
	}
	return report, nil
}

// ── 内部工具 ──────────────────────────────────────────────────────────────

func tableExists(ctx context.Context, table string) (bool, error) {
	v, err := g.DB().GetValue(ctx,
		"SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?",
		table)
	if err != nil {
		return false, fmt.Errorf("检查表 %s 是否存在失败: %w", table, err)
	}
	return v.Int64() > 0, nil
}

// baseRowsMissingInShards 统计「主表里有、但所有分片里都没有」的行数（按主键判断）。
// 这是归档的无损性条件：只要有一行缺失，改名后它就查不到了。
func baseRowsMissingInShards(ctx context.Context, base string) (int64, error) {
	shards, err := activeShards(ctx)
	if err != nil {
		return 0, err
	}
	if len(shards) == 0 {
		v, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM `"+base+"`")
		if err != nil {
			return 0, fmt.Errorf("统计 %s 行数失败: %w", base, err)
		}
		return v.Int64(), nil // 一个分片都没有：全部行都会丢
	}
	unions := make([]string, 0, len(shards))
	for _, sh := range shards {
		unions = append(unions, "SELECT id FROM `"+sh.table(base)+"`")
	}
	// LEFT JOIN + IS NULL 比 NOT IN 更稳（后者在子查询含 NULL 时会整体为假）
	sqlStr := "SELECT COUNT(*) FROM `" + base + "` b LEFT JOIN (" +
		strings.Join(unions, " UNION ALL ") + ") s ON s.id = b.id WHERE s.id IS NULL"
	v, err := g.DB().GetValue(ctx, sqlStr)
	if err != nil {
		return 0, fmt.Errorf("归档无损性预检失败（%s）: %w", base, err)
	}
	return v.Int64(), nil
}

// assertBaseTablesPresent 在运维动作前确认主表还在 —— 归档之后要给一句人话，
// 而不是让调用方看到「表不存在」这种底层报错。
func assertBaseTablesPresent(ctx context.Context, op string) error {
	var archived []string
	for _, base := range shardedTables {
		ok, err := tableExists(ctx, base)
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		// 主表不在：看看是不是已经归档了
		if suffix, found := findLegacySuffix(ctx, base); found {
			archived = append(archived, base+" → "+legacyTableName(base, suffix))
		} else {
			return errcode.Newf(errcode.CodeNotFound, "%s 需要的表 %s 不存在", op, base)
		}
	}
	if len(archived) == 0 {
		return nil
	}
	sort.Strings(archived)
	return errcode.Newf(errcode.CodeShardMaintenanceRefused,
		"%s 已随 Phase 6 归档，不再适用：%s。如需执行请先用 RENAME 把归档表改回原名（30 天窗口内）",
		op, strings.Join(archived, ", "))
}

// findLegacySuffix 找出某张主表是否存在归档版本。
func findLegacySuffix(ctx context.Context, base string) (string, bool) {
	values, err := g.DB().GetArray(ctx,
		"SELECT TABLE_NAME FROM information_schema.TABLES "+
			"WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME LIKE ?", base+legacyPrefix+"%")
	if err != nil {
		return "", false
	}
	var suffixes []string
	for _, v := range values {
		if s, ok := parseLegacySuffix(v.String()); ok {
			suffixes = append(suffixes, s)
		}
	}
	if len(suffixes) == 0 {
		return "", false
	}
	sort.Strings(suffixes)
	return suffixes[len(suffixes)-1], true
}

// CheckBaseTablesForMode 启动自检：single 模式依赖主表，而主表可能已被 Phase 6 归档。
// 这种组合下每个订单请求都会失败，所以在启动时就明确告警（只告警不阻断：
// 运维可能正打算 RENAME 回来，或马上要切 monthly）。
func CheckBaseTablesForMode(ctx context.Context) {
	missing := make([]string, 0, len(shardedTables))
	for _, base := range shardedTables {
		ok, err := tableExists(ctx, base)
		if err != nil {
			g.Log().Warningf(ctx, "订单主表自检失败（已忽略）: %v", err)
			return
		}
		if !ok {
			missing = append(missing, base)
		}
	}
	if len(missing) == 0 {
		return
	}

	mode := shardMode(ctx)
	if mode == shardModeMonthly {
		g.Log().Infof(ctx,
			"订单主表已归档（%s），当前 mode=monthly 只读写分片，符合 Phase 6 的预期状态",
			strings.Join(missing, ", "))
		return
	}
	g.Log().Warningf(ctx,
		"orderShard.mode=%s 但订单主表已归档（%s）：订单相关接口会持续失败。"+
			"要么把归档表 RENAME 回原名（30 天窗口内），要么把 mode 改成 monthly",
		mode, strings.Join(missing, ", "))
}
