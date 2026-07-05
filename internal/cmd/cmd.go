package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gf-eshop/internal/controller/brands"
	"gf-eshop/internal/controller/categories"
	categoryBrands "gf-eshop/internal/controller/category_brands"
	staffCtrl "gf-eshop/internal/controller/staff"

	"gf-eshop/internal/controller/hello"
	brandsLogic "gf-eshop/internal/logic/brands"
	categoriesLogic "gf-eshop/internal/logic/categories"
	"gf-eshop/utility"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 启动时缓存预热
			brandsLogic.Warmup(ctx)
			categoriesLogic.Warmup(ctx)

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
				)
			})
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					brands.NewV1(),
					categories.NewV1(),
					categoryBrands.NewV1(),
				)
				group.Group("/staff", func(group *ghttp.RouterGroup) {
					group.Middleware(authMiddleware)
					group.Bind(
						staffCtrl.NewV1(),
					)
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

	if publicPaths[r.URL.Path] && r.Method == "POST" {
		r.Middleware.Next()
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		r.Response.WriteJson(g.Map{"code": 401, "message": "未登录"})
		r.Exit()
		return
	}

	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	} else {
		r.Response.WriteJson(g.Map{"code": 401, "message": "无效的认证格式"})
		r.Exit()
		return
	}

	claims, err := utility.ParseStaffToken(r.Context(), tokenStr)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 401, "message": "Token无效或已过期"})
		r.Exit()
		return
	}
	if claims.TokenType != utility.TokenTypeAccess {
		r.Response.WriteJson(g.Map{"code": 401, "message": "无效的Token类型"})
		r.Exit()
		return
	}
	r.SetCtxVar("staff_claims", claims)
	r.Middleware.Next()
}
