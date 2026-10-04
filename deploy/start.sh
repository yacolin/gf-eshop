#!/usr/bin/env bash
# gf-eshop 启动脚本
#
# 为什么要包一层脚本（而不是让 systemd 直接跑二进制）：
#   1. 必须在「应用目录」启动 —— GoFrame 是按工作目录找 manifest/config/config.yaml 的，
#      工作目录不对就会读不到配置（报错或回落到默认值）。systemd/nohup 的 cwd 都不确定。
#   2. 启动前做一次配置体检，缺文件就明确报错退出，而不是让进程带着空配置跑起来。
#   3. 给小内存机器留一个统一的环境变量入口（GOGC/GOMEMLIMIT 等）。
#
# 用法：
#   前台： ./start.sh
#   后台： nohup ./start.sh >/dev/null 2>&1 &     （日志已由 logger.stdout 输出到 stdout）
#   指定环境变量文件： GFESHOP_ENV_FILE=/path/to.env ./start.sh
set -euo pipefail

APP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$APP_DIR"

CONF="$APP_DIR/manifest/config/config.yaml"
if [ ! -f "$CONF" ]; then
  echo "[start.sh] 找不到配置文件：$CONF" >&2
  echo "[start.sh] GoFrame 从工作目录读配置，请确认解包完整（manifest/config/config.yaml 存在）" >&2
  exit 1
fi

# 环境变量文件：这里的变量会**覆盖 manifest/config/config.yaml 的同名配置**
# （命名规则见 deploy/env.example：配置键 → 大写蛇形）。
# 为什么必须在这里 source：systemd 不会读 ~/.bashrc，而解包会覆盖 yaml，
# 所以线上真实配置（DB 密码、JWT 密钥…）只放在 /etc/gf-eshop.env（600）。
# 也可以用它给小内存机器设 GOMEMLIMIT / GOGC。
for f in "${GFESHOP_ENV_FILE:-}" /etc/gf-eshop.env; do
  if [ -n "$f" ] && [ -f "$f" ]; then
    set -a
    # shellcheck disable=SC1090
    . "$f"
    set +a
    echo "[start.sh] 已加载环境变量文件: $f"
  fi
done

# 内存调优：systemd 单元里已设 GOMEMLIMIT=400MiB（配合 MemoryMax=512M）。
# 若不用 systemd（nohup 手动跑），在 2G 机器上建议放开下面两行：
# : "${GOGC:=50}"; export GOGC
# : "${GOMEMLIMIT:=400MiB}"; export GOMEMLIMIT

echo "[start.sh] 启动 gf-eshop：dir=$APP_DIR"
exec "$APP_DIR/gf-eshop" "$@"
