package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// New 创建一个业务错误码，code 为稳定标识（一旦发布即永久保留）
func New(code int, msg string) error {
	return gerror.NewCode(gcode.New(code, "", nil), msg)
}

// CodeOf 提取错误码，未设置时返回 0
func CodeOf(err error) int {
	if err == nil {
		return 0
	}
	c := gerror.Code(err)
	if c == gcode.CodeNil {
		return 0
	}
	return c.Code()
}

// ── 通用域 1001-1999 ──

var (
	ErrInvalidParams       = New(1002, "参数错误")
	ErrUnauthorized        = New(1004, "未授权，请先登录")
	ErrUserNotFound        = New(1005, "用户不存在")
	ErrOrderNotFound       = New(1006, "订单不存在")
	ErrPaymentFailed       = New(1008, "支付失败")
	ErrInvalidCredentials  = New(1009, "用户名或密码错误")
	ErrNotFound            = New(1010, "资源不存在")
	ErrAccountDisabled     = New(1011, "账号已被禁用")
	ErrInvalidToken        = New(1016, "Token 无效或已过期")
	ErrInsufficientStock   = New(1024, "库存不足")
)

// ── 权限域 2001-2999 ──

var (
	ErrPermissionNotFound      = New(2001, "权限不存在")
	ErrInsufficientPermissions = New(2002, "无权限，需要管理员角色")
	ErrCannotModifySystemRole  = New(2003, "不能修改系统角色")
	ErrCannotDeleteSystemRole  = New(2004, "不能删除系统角色")
)

// ── 品牌/商品域 4001-4099 ──

var (
	ErrBrandNotFound    = New(4001, "品牌不存在")
	ErrCategoryNotFound = New(4010, "类目不存在")
	ErrProductNotFound  = New(4030, "产品不存在")
	ErrSKUNotFound      = New(4031, "SKU 不存在")
	ErrAttributeNotFound = New(4020, "属性不存在")
)

// ── 库存域 5001-5099 ──

var (
	ErrInventoryNotFound  = New(5001, "库存记录不存在")
	ErrInvalidStockChange = New(5003, "无效的库存变动")
)

// ── 交易域 6001-6099 ──

var (
	ErrPaymentNotFound = New(6001, "支付记录不存在")
	ErrRefundNotFound  = New(6010, "退款记录不存在")
	ErrRefundFailed    = New(6011, "退款失败")
)

// ── 订单域 7001-7099 ──

var (
	ErrInvalidOrderStatus = New(7002, "无效的订单状态变更")
	ErrOrderItemNotFound  = New(7003, "订单项不存在")
)
