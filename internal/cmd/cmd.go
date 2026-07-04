package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"gf-eshop/internal/controller/brands"
	"gf-eshop/internal/controller/categories"
	categoryBrands "gf-eshop/internal/controller/category_brands"

	"gf-eshop/internal/controller/hello"
	brandsLogic "gf-eshop/internal/logic/brands"
	categoriesLogic "gf-eshop/internal/logic/categories"
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
			})
			s.Run()
			return nil
		},
	}
)
