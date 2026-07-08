package cmd

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gf-eshop/internal/errcode"

	"gf-eshop/internal/controller/brands"
	"gf-eshop/internal/controller/categories"
	categoryBrands "gf-eshop/internal/controller/category_brands"
	"gf-eshop/internal/controller/permissions"
	"gf-eshop/internal/controller/notification"
	"gf-eshop/internal/controller/roles"
	staffCtrl "gf-eshop/internal/controller/staff"
	wsCtrl "gf-eshop/internal/controller/ws"

	"gf-eshop/internal/controller/attributes"
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
	_ "gf-eshop/internal/logic/dashboard"
	marketingLogic "gf-eshop/internal/logic/marketing"
	merchantsLogic "gf-eshop/internal/logic/merchants"
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
	reviewsCtrl "gf-eshop/internal/controller/reviews"
	"gf-eshop/internal/controller/merchants"
	"gf-eshop/internal/service"
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
						brandsLogic.Warmup(ctx)
						return 0, nil
					}),
					productsLogic.NewFuncStage("categories", func(ctx context.Context) (int, error) {
						categoriesLogic.Warmup(ctx)
						return 0, nil
					}),
					productsLogic.NewFuncStage("products", func(ctx context.Context) (int, error) {
						return productsLogic.Warmup(ctx)
					}),
				productsLogic.NewFuncStage("marketing", func(ctx context.Context) (int, error) {
					marketingLogic.Warmup(ctx)
					return 0, nil
				}),
				productsLogic.NewFuncStage("merchants", func(ctx context.Context) (int, error) {
					merchantsLogic.Warmup(ctx)
					return 0, nil
				}),
			)
				pipeline.Run(ctx)
				// 启动仪表盘定时刷新
				service.Dashboard().StartPeriodicRefresh(ctx)

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
					attributes.NewV1(),
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

				dashboardCtrl.NewV1(),
					reviewsCtrl.NewV1(),
				)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(authMiddleware)
				group.Bind(
					staffCtrl.NewV1(),
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
