// Package middleware 提供自定义 Gin 风格访问日志中间件。
// 默认未启用（由 GoFrame 内置 accessLog 替代），需要时在
// cmd.go 中 s.Use(middleware.AccessLog) 并关闭内置 accessLogEnabled。
package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// AccessLog 是 Gin 风格的请求日志中间件。
// 使用方式：s.Use(middleware.AccessLog)
// 输出格式：200 |    GET | 2.3ms | /api/v1/xxx
func AccessLog(r *ghttp.Request) {
	start := time.Now()
	r.Middleware.Next()
	cost := time.Since(start)

	g.Log().Printf(r.Context(), "%3d | %6s | %s | %s",
		r.Response.Status, r.Method, cost, r.URL.RequestURI(),
	)
}
