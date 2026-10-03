#!/usr/bin/env bash
# ES 在小内存机器上的「能不能撑住」观测脚本（Linux）
#
# 用途：试装 Elasticsearch 后，用它持续采样，判断机器是真撑住了还是在 swap 上硬扛。
# 只读，不改任何配置、不写任何文件（除你指定的输出文件），随时 Ctrl-C 退出。
#
# 用法：
#   ./es-watch.sh                 # 每 5 秒采一次，前台打印
#   ./es-watch.sh 60 10           # 每 10 秒采一次，共 60 次（约 10 分钟）
#   ./es-watch.sh 120 5 /tmp/es-watch.log
#
# 判读标准（脚本最后会给结论）：
#   ✅ 通过：si/so 全为 0，ES 堆稳定，无 OOM kill，无降级日志
#   ⚠️ 临界：偶发 so>0，堆峰值 >80%（还有余量但不能再加负载）
#   ❌ 不通过：持续 so>0（在拿 swap 当内存用）、JVM GC 时间快速增长、出现 OOM killer、应用开始降级
set -uo pipefail

ROUNDS="${1:-60}"        # 采样次数
INTERVAL="${2:-5}"       # 采样间隔（秒）
OUTFILE="${3:-}"         # 可选：同时写日志文件
ES_URL="${ES_URL:-http://127.0.0.1:9200}"
APP_UNIT="${APP_UNIT:-gf-eshop}"

LOG_LINES=0
FAIL_SWAP=0
FAIL_OOM=0
FAIL_FALLBACK=0
WARN_HEAP=0
MAX_HEAP_PCT=0

say() { printf '%s\n' "$*"; [ -n "$OUTFILE" ] && printf '%s\n' "$*" >>"$OUTFILE"; }

# 依赖检查
command -v free >/dev/null || { echo "需要 procps（free 命令）；这不是 Linux 环境？" >&2; exit 1; }
command -v curl >/dev/null || { echo "需要 curl" >&2; exit 1; }

say "== ES 观测开始：$(date '+%F %T')  轮次=$ROUNDS 间隔=${INTERVAL}s =="
say "机器内存概览："
free -m | sed 's/^/  /' | while read -r l; do say "$l"; done
say ""

printf '%-8s %-28s %-16s %-10s %-8s %s\n' "时间" "内存(used/free/swap_used MB)" "ES堆(used/max)" "堆占比" "si/so" "ES健康" | tee -a "${OUTFILE:-/dev/null}"
say ""

for i in $(seq 1 "$ROUNDS"); do
  TS=$(date '+%H:%M:%S')

  # 1) 内存与 swap
  read -r MEM_USED MEM_FREE SWAP_USED < <(free -m | awk '/^Mem:/{u=$3; f=$4} /^Swap:/{s=$3} END{print u, f, s}')
  # 2) 单位时间内的换入/换出（用 vmstat 两次采样的差值；so>0 表示正在把内存压到 swap）
  read -r SI SO < <(vmstat 1 2 | awk 'NR==4{print $7, $8}')

  # 3) ES JVM 堆
  ES_JSON=$(curl -s -m 3 "$ES_URL/_nodes/stats/jvm" 2>/dev/null)
  HEAP_USED=$(printf '%s' "$ES_JSON" | grep -o '"heap_used_in_bytes":[0-9]*' | head -1 | cut -d: -f2)
  HEAP_MAX=$(printf '%s' "$ES_JSON" | grep -o '"heap_max_in_bytes":[0-9]*' | head -1 | cut -d: -f2)
  if [ -n "${HEAP_USED:-}" ] && [ -n "${HEAP_MAX:-}" ] && [ "$HEAP_MAX" -gt 0 ]; then
    HEAP_PCT=$(( HEAP_USED * 100 / HEAP_MAX ))
    HEAP_STR="$(( HEAP_USED / 1048576 ))M/$(( HEAP_MAX / 1048576 ))M"
  else
    HEAP_PCT=0; HEAP_STR="不可用"
  fi
  [ "$HEAP_PCT" -gt "$MAX_HEAP_PCT" ] && MAX_HEAP_PCT=$HEAP_PCT
  [ "$HEAP_PCT" -gt 80 ] && WARN_HEAP=1

  # 4) 集群/索引健康
  HEALTH=$(curl -s -m 3 "$ES_URL/_cluster/health" 2>/dev/null | grep -o '"status":"[a-z]*"' | head -1 | cut -d'"' -f4)
  [ -z "$HEALTH" ] && HEALTH="无响应"

  # 5) 内核是否 OOM kill 过进程
  if command -v journalctl >/dev/null; then
    OOM=$(journalctl -k --since "${INTERVAL} seconds ago" 2>/dev/null | grep -ci "killed process" || true)
  else
    OOM=0
  fi
  [ "${OOM:-0}" -gt 0 ] && { FAIL_OOM=1; say "  !! 检测到 OOM kill："; journalctl -k --since "2 minutes ago" 2>/dev/null | grep -i "killed process" | tail -3 | sed 's/^/     /' | while read -r l; do say "$l"; done; }

  # 6) 应用是否已经开始降级（降级说明 ES 响应超时/熔断）
  if command -v journalctl >/dev/null; then
    NEW_FALLBACK=$(journalctl -u "$APP_UNIT" --since "${INTERVAL} seconds ago" --no-pager 2>/dev/null | grep -c "ES 检索失败，降级 DB 查询" || true)
  else
    NEW_FALLBACK=0
  fi
  if [ "${NEW_FALLBACK:-0}" -gt 0 ]; then
    FAIL_FALLBACK=1
    LOG_LINES=$(( LOG_LINES + NEW_FALLBACK ))
  fi

  # si/so 判读：so 持续 >0 视为在用 swap 当内存
  if [ "${SO:-0}" -gt 0 ]; then FAIL_SWAP=$(( FAIL_SWAP + 1 )); fi

  printf '%-8s %-28s %-16s %-10s %-8s %s\n' \
    "$TS" "${MEM_USED}/${MEM_FREE}/${SWAP_USED}" "$HEAP_STR" "${HEAP_PCT}%" "${SI}/${SO}" "$HEALTH" \
    | tee -a "${OUTFILE:-/dev/null}"

  [ "$i" -lt "$ROUNDS" ] && sleep "$INTERVAL"
done

say ""
say "================== 结论 =================="
if [ "$FAIL_OOM" -eq 1 ]; then
  say "❌ 出现过 OOM kill —— 内存不够，不要在这台机器上跑 ES（或把堆调更小再试）。"
elif [ "$FAIL_SWAP" -gt $(( ROUNDS / 3 )) ]; then
  say "❌ 有 $FAIL_SWAP/$ROUNDS 次采样检测到 swap 换出 —— 在拿 swap 当内存用，"
  say "   ES 的 GC 停顿会到秒级，应用侧会频繁降级。结论：2G 撑 ES 太勉强。"
elif [ "$FAIL_FALLBACK" -eq 1 ]; then
  say "❌ 应用出现 $LOG_LINES 次「ES 检索失败，降级 DB 查询」—— ES 已经不可用了（健康≠可用）。"
elif [ "$WARN_HEAP" -eq 1 ]; then
  say "⚠️ 没有 swap 换出，但 ES 堆峰值到过 ${MAX_HEAP_PCT}%（>80%）。当前能用，但余量不多："
  say "   要么把堆保持在 512m 并接受这个水位，要么加内存。"
elif [ "$FAIL_SWAP" -gt 0 ]; then
  say "⚠️ 偶发 swap 换出（$FAIL_SWAP/$ROUNDS 次），目前还算稳；建议再看一段时间的 GC 时间趋势。"
else
  say "✅ 全程没有 swap 换出、没有 OOM、没有降级，ES 堆峰值 ${MAX_HEAP_PCT}%。"
  say "   结论：这台 2G 机器跑「应用 + Redis + ES」在当前数据量下是撑得住的。"
fi
say "提示：真正的压力在「并发检索 + 段合并」同时发生时，建议配合 docs/perf-workflow.md 做一轮压测再定。"
say "========================================="
