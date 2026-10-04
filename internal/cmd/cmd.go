package cmd

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gf-eshop/internal/errcode"

	"gf-eshop/internal/controller/brands"
	"gf-eshop/internal/controller/categories"
	categoryAttributesCtrl "gf-eshop/internal/controller/category_attributes"
	categoryBrands "gf-eshop/internal/controller/category_brands"
	"gf-eshop/internal/controller/permissions"
	"gf-eshop/internal/controller/notification"
	"gf-eshop/internal/controller/roles"
	staffCtrl "gf-eshop/internal/controller/staff"
	wsCtrl "gf-eshop/internal/controller/ws"

	attributesCtrl "gf-eshop/internal/controller/attributes"
	attributeValuesCtrl "gf-eshop/internal/controller/attribute_values"
	departmentsCtrl "gf-eshop/internal/controller/departments"
	"gf-eshop/internal/controller/hello"
	inventoriesCtrl "gf-eshop/internal/controller/inventories"
	productVersionsCtrl "gf-eshop/internal/controller/product_versions"
	inventoryLogsCtrl "gf-eshop/internal/controller/inventory_logs"
	productAttributesCtrl "gf-eshop/internal/controller/product_attributes"
	productDescriptionsCtrl "gf-eshop/internal/controller/product_descriptions"
	productsCtrl "gf-eshop/internal/controller/products"
	skusCtrl "gf-eshop/internal/controller/skus"
	dashboardCtrl "gf-eshop/internal/controller/dashboard"
	brandsLogic "gf-eshop/internal/logic/brands"
	categoriesLogic "gf-eshop/internal/logic/categories"
	productsLogic "gf-eshop/internal/logic/products"
	ordersLogic "gf-eshop/internal/logic/orders"
	_ "gf-eshop/internal/logic/dashboard"
	marketingLogic "gf-eshop/internal/logic/marketing"
	merchantsLogic "gf-eshop/internal/logic/merchants"
	departmentsLogic "gf-eshop/internal/logic/departments"
	_ "gf-eshop/internal/logic/merchant_roles"
	_ "gf-eshop/internal/logic/merchant_role_permissions"
	_ "gf-eshop/internal/logic/merchant_withdrawals"
	_ "gf-eshop/internal/logic/merchant_balances"
	marketingCtrl "gf-eshop/internal/controller/marketing"
	"gf-eshop/internal/controller/address"
	"gf-eshop/internal/controller/user"
	"gf-eshop/internal/controller/user_admin"
	"gf-eshop/internal/controller/user_levels"
	"gf-eshop/internal/controller/user_points"
	pointsRulesCtrl "gf-eshop/internal/controller/points_rules"
	levelRulesCtrl "gf-eshop/internal/controller/level_rules"

	"gf-eshop/internal/controller/user_auth"
	"gf-eshop/internal/middleware"
	ordersCtrl "gf-eshop/internal/controller/orders"
	paymentsCtrl "gf-eshop/internal/controller/payments"
	cartsCtrl "gf-eshop/internal/controller/carts"

	merchantBankAccountsCtrl "gf-eshop/internal/controller/merchant_bank_accounts"
	merchantContactsCtrl "gf-eshop/internal/controller/merchant_contacts"
	merchantQualificationsCtrl "gf-eshop/internal/controller/merchant_qualifications"
	merchantWithdrawalsCtrl "gf-eshop/internal/controller/merchant_withdrawals"
	merchantBalancesCtrl "gf-eshop/internal/controller/merchant_balances"
	merchantRolesCtrl "gf-eshop/internal/controller/merchant_roles"
	merchantRolePermissionsCtrl "gf-eshop/internal/controller/merchant_role_permissions"
	reviewsCtrl "gf-eshop/internal/controller/reviews"
	"gf-eshop/internal/controller/merchants"
	"gf-eshop/internal/service"
	"gf-eshop/internal/verifycode"
	"gf-eshop/internal/ws"
	"gf-eshop/utility"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
				// 启动时缓存预热（并行管线）
				pipeline := productsLogic.NewWarmupPipeline(
					productsLogic.NewFuncStage("brands", func(ctx context.Context) (int, error) {
						return brandsLogic.Warmup(ctx)
					}),
					productsLogic.NewFuncStage("brands_es", func(ctx context.Context) (int, error) {
						return brandsLogic.WarmupES(ctx)
					}),
					productsLogic.NewFuncStage("categories", func(ctx context.Context) (int, error) {
						return categoriesLogic.Warmup(ctx)
					}),
					productsLogic.NewFuncStage("products", func(ctx context.Context) (int, error) {
						return productsLogic.Warmup(ctx)
					}),
					productsLogic.NewFuncStage("products_es", func(ctx context.Context) (int, error) {
						return productsLogic.WarmupES(ctx)
					}),
					productsLogic.NewFuncStage("marketing", func(ctx context.Context) (int, error) {
						marketingLogic.Warmup(ctx)
						return 0, nil
					}),
					productsLogic.NewFuncStage("merchants", func(ctx context.Context) (int, error) {
						merchantsLogic.Warmup(ctx)
						return 0, nil
					}),
					productsLogic.NewFuncStage("departments", func(ctx context.Context) (int, error) {
						departmentsLogic.Warmup(ctx)
						return 0, nil
					}),
			)
				pipeline.Run(ctx)
				// 启动仪表盘定时刷新
				service.Dashboard().StartPeriodicRefresh(ctx)
				// 验证码渠道配置自检：只告警不阻断启动
				warnVerifyConfigIfNeeded(ctx)

				// 订单主表自检：Phase 6 归档后若还停在 single 模式，订单接口会持续失败
				ordersLogic.CheckBaseTablesForMode(ctx)

			// 创建并启动 WebSocket Hub
			wsHub := ws.NewHub()
			go wsHub.Run()
			service.RegisterWsHub(wsHub)

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.ErrorHandler)
				group.Bind(
					hello.NewV1(),
				)
			})
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.ErrorHandler)
				group.Bind(
					brands.NewV1(),
					categories.NewV1(),
					categoryBrands.NewV1(),
					productsCtrl.NewV1(),
					skusCtrl.NewV1(),
					attributesCtrl.NewV1(),
					attributeValuesCtrl.NewV1(),
					categoryAttributesCtrl.NewV1(),
					productAttributesCtrl.NewV1(),
					productDescriptionsCtrl.NewV1(),
					inventoriesCtrl.NewV1(),
					inventoriesCtrl.NewWarehousesV1(),
					inventoryLogsCtrl.NewV1(),
					productVersionsCtrl.NewV1(),

					merchants.NewV1(),
					merchantBankAccountsCtrl.NewV1(),
					merchantContactsCtrl.NewV1(),
					merchantQualificationsCtrl.NewV1(),
					merchantWithdrawalsCtrl.NewV1(),
					merchantBalancesCtrl.NewV1(),
					merchantRolesCtrl.NewV1(),
					merchantRolePermissionsCtrl.NewV1(),

				dashboardCtrl.NewV1(),
					reviewsCtrl.NewV1(),
				)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(authMiddleware)
				group.Bind(
					staffCtrl.NewV1(),
						departmentsCtrl.NewV1(),
					permissions.NewV1(),
					notification.NewV1(),
					wsCtrl.NewV1(),
					user_admin.NewV1(),
					user_levels.NewV1(),
					user_points.NewV1(),	
				)
			})
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(authMiddleware, middleware.RequireAdmin)
				group.Bind(
					roles.NewV1(),
						pointsRulesCtrl.NewV1(),
						levelRulesCtrl.NewV1(),
				)
			})
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.UserAuthMiddleware)
				group.Bind(
					user.NewV1(),
					user_auth.NewV1(),
					address.NewV1(),
				)
			})
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.UserAuthMiddleware)
				group.Bind(
					cartsCtrl.NewV1(),
					ordersCtrl.NewV1(),
					paymentsCtrl.NewV1(),
					marketingCtrl.NewV1(),
				)
			})
		})
			// WS 升级路由（不使用 MiddlewareHandlerResponse，token 从查询参数获取）
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.GET("/ws", func(r *ghttp.Request) {
					tokenStr := r.GetQuery("token").String()
					if tokenStr == "" {
						r.Response.WriteStatus(401)
						return
					}
					claims, err := utility.ParseStaffToken(r.Context(), tokenStr)
					if err != nil {
						r.Response.WriteStatus(401)
						return
					}
					if claims.TokenType != utility.TokenTypeAccess {
						r.Response.WriteStatus(401)
						return
					}
					wsConn, err := r.WebSocket()
					if err != nil {
						return
					}
					lastSeqStr := r.GetQuery("last_seq").String()
					var lastSeq int64
					if lastSeqStr != "" {
						lastSeq, _ = strconv.ParseInt(lastSeqStr, 10, 64)
					}
					client := ws.NewClient(wsHub, wsConn, claims.StaffId, lastSeq)
					wsHub.Register() <- client
					go client.WritePump()
					client.ReadPump()
				})
			})
			s.Run()
			return nil
		},
	}
)

func init() {
	if err := Main.AddCommand(&gcmd.Command{
		Name:  "shard",
		Usage: "shard --action=create|migrate|restore|verify|archive|purge --from=2026-08 --to=2026-09 [--suffix=202610] [--batch=5000] [--force]",
		Brief: "订单分表运维：建分片 / 迁移历史 / 回灌主表(回滚) / 三重对账（见 docs/order-sharding-design.md §7）",
		Func: func(ctx context.Context, parser *gcmd.Parser) error {
			// 同 reindex：GoFrame 的 gcmd 会把多余的位置参数当成多级命令名，只能用选项传参
			var (
				action = strings.TrimSpace(parser.GetOpt("action", "").String())
				from   = strings.TrimSpace(parser.GetOpt("from", "").String())
				to     = strings.TrimSpace(parser.GetOpt("to", "").String())
				suffix string
				force  bool
				batch  = parser.GetOpt("batch", 5000).Int()
			)
			suffix = strings.TrimSpace(parser.GetOpt("suffix", "").String())
			force = parser.GetOpt("force", false).Bool()

			// archive / purge 用后缀定位归档表，不需要月份区间
			if action != "archive" && action != "purge" {
				if from == "" || to == "" {
					return gerror.New("必须指定 --from 与 --to（形如 2026-08）")
				}
			}

			switch action {
			case "create":
				tables, err := ordersLogic.CreateShards(ctx, from, to)
				if err != nil {
					return err
				}
				for _, t := range tables {
					g.Log().Infof(ctx, "分片表就绪: %s", t)
				}
				g.Log().Infof(ctx, "共 %d 张分片表就绪（%s ~ %s）", len(tables), from, to)
				return nil

			case "migrate":
				report, err := ordersLogic.MigrateShards(ctx, from, to, batch)
				if err != nil {
					return err
				}
				for _, m := range report.Months {
					g.Log().Infof(ctx, "迁移 %s %s：新增 %d 行", m.Month, m.Table, m.Rows)
				}
				g.Log().Infof(ctx, "老主键映射登记 %d 行，日汇总回填 %d 行",
					report.ShardMapRows, report.StatsRows)
				return nil

			case "restore":
				// Phase 5 回滚：把分片数据回灌主表。方向与 migrate 相反，故两者互斥。
				report, err := ordersLogic.RestoreToBase(ctx, from, to, batch)
				if err != nil {
					return err
				}
				var total int64
				for _, m := range report.Months {
					g.Log().Infof(ctx, "回灌 %s %s：影响 %d 行", m.Month, m.Table, m.Rows)
					total += m.Rows
				}
				g.Log().Infof(ctx, "回灌完成，共影响 %d 行；建议接着跑 --action=verify 确认两边一致，"+
					"确认后再把 orderShard.mode 改回 single", total)
				return nil

			case "verify":
				report, err := ordersLogic.VerifyShards(ctx, from, to)
				if err != nil {
					return err
				}
				for _, c := range report.Checks {
					if c.Pass {
						g.Log().Infof(ctx, "✓ %s %s：行数 %d，校验和 %d，金额 %s",
							c.Month, c.Table, c.RowsA, c.HashA, c.SumA)
						continue
					}
					g.Log().Errorf(ctx, "✗ %s %s：%s", c.Month, c.Table, c.Detail)
				}
				if !report.AllPass {
					return gerror.Newf("对账失败 %d 项", report.Failed)
				}
				g.Log().Infof(ctx, "对账全部通过（%d 项）", len(report.Checks))
				return nil

			case "archive":
				// Phase 6：把四张主表改名为归档表（无损性预检通过才动手）
				report, err := ordersLogic.ArchiveBaseTables(ctx, suffix)
				if err != nil {
					return err
				}
				for _, p := range report.Renamed {
					g.Log().Infof(ctx, "归档 %s → %s", p[0], p[1])
				}
				g.Log().Infof(ctx, "归档完成（后缀 %s）。观察窗口内可这样回滚：%s",
					report.Suffix, report.RollbackSQL())
				return nil

			case "purge":
				// 清理归档表：不加 --force 只预览 DROP 语句
				report, err := ordersLogic.PurgeLegacyTables(ctx, suffix, force, time.Now())
				if err != nil {
					return err
				}
				for _, s := range report.Statements {
					g.Log().Infof(ctx, "%s", s)
				}
				if len(report.Dropped) == 0 {
					g.Log().Infof(ctx, "以上为预览（未删除）。观察窗口最早结束于 %s；确认无需回滚后加 --force 执行",
						report.AllowedAt.Format("2006-01-02"))
					return nil
				}
				g.Log().Infof(ctx, "已删除 %d 张归档表", len(report.Dropped))
				return nil

			default:
				return gerror.Newf("未知 --action=%q，可选：create、migrate、restore、verify、archive、purge", action)
			}
		},
	}); err != nil {
		panic(err)
	}
}

// reindexTargets 已接入 ES 的业务实体 → 全量重建函数。
// 后续 products / skus 接入时在此登记，即可自动获得 `main reindex <entity>` 能力。
var reindexTargets = map[string]func(context.Context) (int, error){
	"brands":   brandsLogic.ReindexBrands,
	"products": productsLogic.ReindexProducts,
}

// reindexEntityNames 返回稳定排序的实体名，保证命令输出可预期。
func reindexEntityNames() []string {
	names := make([]string, 0, len(reindexTargets))
	for name := range reindexTargets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func init() {
	if err := Main.AddCommand(&gcmd.Command{
		Name:  "reindex",
		Usage: "reindex [--entity=brands]",
		Brief: "全量重建 Elasticsearch 索引（省略 --entity 时重建全部已接入实体）",
		Func: func(ctx context.Context, parser *gcmd.Parser) error {
			// 注意：GoFrame 的 gcmd 会把多余的位置参数当成多级命令名，
			// 因此实体名只能用选项传入，不能写成 `main reindex brands`。
			entity := strings.TrimSpace(parser.GetOpt("entity", "").String())

			if entity != "" {
				fn, ok := reindexTargets[entity]
				if !ok {
					return gerror.Newf("未知实体 %q，已接入的实体：%s",
						entity, strings.Join(reindexEntityNames(), ", "))
				}
				n, err := fn(ctx)
				if err != nil {
					return err
				}
				g.Log().Infof(ctx, "实体 %s 索引重建完成：%d 条", entity, n)
				return nil
			}

			for _, name := range reindexEntityNames() {
				n, err := reindexTargets[name](ctx)
				if err != nil {
					// 单个实体失败不阻断其他实体
					g.Log().Errorf(ctx, "实体 %s 索引重建失败：%v", name, err)
					continue
				}
				g.Log().Infof(ctx, "实体 %s 索引重建完成：%d 条", name, n)
			}
			return nil
		},
	}); err != nil {
		panic(err)
	}
}

func authMiddleware(r *ghttp.Request) {
	publicPaths := map[string]bool{
		"/api/v1/staff/login":   true,
		"/api/v1/staff/refresh": true,
	}

	if publicPaths[r.URL.Path] {
		r.Middleware.Next()
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		r.SetError(errcode.ErrUnauthorized)
		r.Exit()
		return
	}

	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	} else {
		r.SetError(errcode.ErrUnauthorized)
		r.Exit()
		return
	}

	claims, err := utility.ParseStaffToken(r.Context(), tokenStr)
	if err != nil {
		r.SetError(errcode.ErrUnauthorized)
		r.Exit()
		return
	}
	if claims.TokenType != utility.TokenTypeAccess {
		r.SetError(errcode.ErrUnauthorized)
		r.Exit()
		return
	}
	r.SetCtxVar("staff_claims", claims)
	r.Middleware.Next()
}

// warnVerifyConfigIfNeeded 启动时逐渠道自检验证码配置。
// 只在日志里告警，不阻断启动：验证码属于可选能力，
// 某渠道没配好只会让使用该渠道的接口返回 errcode.ErrVerifyChannelNotReady。
func warnVerifyConfigIfNeeded(ctx context.Context) {
	verifycode.WarnNotReady(ctx)
}
