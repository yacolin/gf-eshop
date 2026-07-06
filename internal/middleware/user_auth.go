package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"gf-eshop/internal/errcode"
	"gf-eshop/utility"
)

func UserAuthMiddleware(r *ghttp.Request) {
	publicPaths := map[string]bool{
		"/api/v1/user/auth/login":    true,
		"/api/v1/user/auth/register": true,
		"/api/v1/user/auth/refresh":  true,
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

	claims, err := utility.ParseUserToken(r.Context(), tokenStr)
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
	r.SetCtxVar("user_claims", claims)
	r.Middleware.Next()
}
