package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gf-eshop/internal/errcode"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

func genTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("t-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// ErrorHandler 全局错误处理中间件，替换 ghttp.MiddlewareHandlerResponse
func ErrorHandler(r *ghttp.Request) {
	traceID := genTraceID()
	r.SetCtxVar("trace_id", traceID)
	r.Response.Header().Set("X-Trace-Id", traceID)
	r.Response.Header().Set("X-Request-Id", traceID)

	defer func() {
		if rec := recover(); rec != nil {
			g.Log().Errorf(r.Context(), "panic recovered: %v, stack: %s", rec, debug.Stack())
			r.Response.WriteStatus(500, &APIResponse{
				Code:    500,
				Message: "服务器内部错误",
				TraceID: traceID,
			})
		}
	}()

	r.Middleware.Next()

	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
		msg = "success"
		cd  = 0
	)

	if err != nil {
		cd = errcode.CodeOf(err)
		if cd == 0 {
			cd = 500
		}
		msg = err.Error()
	}

	httpStatus := mapErrorToHTTPStatus(cd)
	body := &APIResponse{
		Code:    cd,
		Message: msg,
		Data:    res,
		TraceID: traceID,
	}

	// 时间字段统一带毫秒（GoFrame 的 gtime 序列化写死了秒级，见 time_millis.go）。
	// 只有**真的含时间字段**时才走自编码；否则保持原路径，输出逐字节不变。
	if data, changed, err := marshalWithMillis(body); err != nil {
		g.Log().Warningf(r.Context(), "时间字段毫秒化失败，回退标准路径: %v", err)
	} else if changed {
		r.Response.WriteHeader(httpStatus)
		r.Response.WriteJson(json.RawMessage(data))
		return
	}
	r.Response.WriteStatus(httpStatus, body)
}

var errorStatusMap = map[int]int{
	// GoFrame 内置码（框架级错误保留 HTTP 语义）
	50: 404, // CodeNotFound
	54: 401, // CodeNotAuthorized
	51: 422, // CodeValidationFailed
	53: 400, // CodeInvalidParameter
	57: 500, // CodeOperationFailed
	// 自定义业务码（仅认证/鉴权类保留 HTTP 语义）
	1004: 401, // ErrUnauthorized
	1016: 401, // ErrInvalidToken
	2002: 403, // ErrInsufficientPermissions
}

func mapErrorToHTTPStatus(code int) int {
	if s, ok := errorStatusMap[code]; ok {
		return s
	}
	// 所有业务错误统一返回 200，通过 code 字段区分
	return 200
}
