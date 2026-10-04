package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// 命名约定的契约：配置键 → 环境变量名（大写蛇形，"." 与驼峰边界都变 "_"）。
func TestConfigKeyToEnvName(t *testing.T) {
	cases := map[string]string{
		"server.address":                 "SERVER_ADDRESS",
		"database.default.link":          "DATABASE_DEFAULT_LINK",
		"orderShard.mode":                "ORDER_SHARD_MODE",
		"orderShard.dualWrite":           "ORDER_SHARD_DUAL_WRITE",
		"orderShard.defaultWindowMonths": "ORDER_SHARD_DEFAULT_WINDOW_MONTHS",
		"email.dev_mode":                 "EMAIL_DEV_MODE",
		"elasticsearch.addresses":        "ELASTICSEARCH_ADDRESSES",
		"server.accessLogEnabled":        "SERVER_ACCESS_LOG_ENABLED",
	}
	for key, want := range cases {
		if got := ConfigKeyToEnvName(key); got != want {
			t.Errorf("ConfigKeyToEnvName(%q) = %q, want %q", key, got, want)
		}
	}
}

// ApplyEnvOverrides 的语义：只覆盖"设置过的"键，未设置的保持 yaml 原值；
// 空字符串也算"设置过"；大小写不敏感；清单外的变量一律忽略。
func TestApplyEnvOverridesSemantics(t *testing.T) {
	ctx := t.Context()
	before := configString(t, "server.address")

	// 1) 未设置：不改动
	applied := ApplyEnvOverrides(ctx, func(string) (string, bool) { return "", false })
	if len(applied) != 0 {
		t.Errorf("没有任何环境变量时不应有覆盖，实际 %v", applied)
	}
	if got := configString(t, "server.address"); got != before {
		t.Errorf("未设置时配置被改动了：%q → %q", before, got)
	}

	// 2) 设置（小写变量名也认）+ 空串也算设置 + 清单外的忽略
	// 假环境：键统一小写存放（模拟大小写不敏感的查询）
	env := map[string]string{
		"server_address":     ":19999", // 配置里写的是 server.address，环境里写小写也认
		"order_shard_mode":   "",       // 空串 = 显式设置
		"not_in_whitelist_x": "should-be-ignored",
	}
	applied = ApplyEnvOverrides(ctx, func(name string) (string, bool) {
		v, ok := env[strings.ToLower(name)]
		return v, ok
	})
	t.Logf("被覆盖的键: %v", applied)

	if got := configString(t, "server.address"); got != ":19999" {
		t.Errorf("server.address 未被覆盖，实际 %q", got)
	}
	if got := configString(t, "orderShard.mode"); got != "" {
		t.Errorf("空串应被视为显式设置，实际 %q", got)
	}
	found := false
	for _, k := range applied {
		if k == "server.address" || k == "orderShard.mode" {
			found = true
		}
		if strings.Contains(k, "not_in_whitelist") || strings.Contains(k, "NOT_IN_WHITELIST") {
			t.Errorf("清单外的变量不应被应用：%s", k)
		}
	}
	if !found {
		t.Errorf("返回值里应列出被覆盖的键，实际 %v", applied)
	}

	// 复原，避免影响同包其它用例
	_ = setConfig(t, "server.address", before)
	_ = setConfig(t, "orderShard.mode", "monthly")
}

// 防漂移：deploy/env.example 里必须出现清单中的**每一个**变量名
// （可以是 `KEY=` 也可以是注释 `# KEY=`），否则运维照着模板配会漏项。
func TestEnvExampleCoversWhitelist(t *testing.T) {
	path := filepath.Join("..", "..", "deploy", "env.example")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("读不到 %s（跳过）: %v", path, err)
	}
	text := string(raw)
	for _, key := range EnvOverridableKeys {
		name := ConfigKeyToEnvName(key)
		if !strings.Contains(text, name) {
			t.Errorf("deploy/env.example 缺少 %s（对应配置项 %s）—— 新增可覆盖项时要同步模板", name, key)
		}
	}
}

// ── 测试用小助手 ──────────────────────────────────────────────────────────

func configString(t *testing.T, key string) string {
	t.Helper()
	v, err := g.Cfg().Get(t.Context(), key)
	if err != nil {
		t.Fatalf("读配置 %s 失败: %v", key, err)
	}
	return v.String()
}

func setConfig(t *testing.T, key, value string) error {
	t.Helper()
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		t.Fatalf("配置适配器不是 *gcfg.AdapterFile")
	}
	return adapter.Set(key, value)
}
