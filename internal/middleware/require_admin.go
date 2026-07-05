package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

func RequireAdmin(r *ghttp.Request) {
	claims := utility.GetStaffClaims(r.Context())
	if claims == nil {
		r.Response.WriteJson(g.Map{"code": 401, "message": "未登录"})
		r.Exit()
		return
	}
	isAdmin, err := service.Roles().IsAdmin(r.Context(), claims.StaffId)
	if err != nil || !isAdmin {
		r.Response.WriteJson(g.Map{"code": 403, "message": "无权限，需要管理员角色"})
		r.Exit()
		return
	}
	r.Middleware.Next()
}
