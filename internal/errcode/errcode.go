package errcode

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// New 创建一个业务错误，code 为稳定标识（一旦发布即永久保留）
func New(code int, msg string) error {
	return gerror.NewCode(gcode.New(code, "", nil), msg)
}

// Newf 创建带格式化消息的业务错误
func Newf(code int, format string, args ...interface{}) error {
	return gerror.NewCode(gcode.New(code, "", nil), fmt.Sprintf(format, args...))
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

// Code 返回一个 gcode.Code 值，用于 gerror.NewCode(code, msg) 动态消息
func Code(n int) gcode.Code {
	return gcode.New(n, "", nil)
}

// ============================================================================
// 错误码常量（用于动态消息的 errcode.Newf 调用）
// ============================================================================

const (
	CodeInvalidParams       = 1002
	CodeUnauthorized        = 1004
	CodeUserNotFound        = 1005
	CodeOrderNotFound       = 1006
	CodePaymentFailed       = 1008
	CodeInvalidCredentials  = 1009
	CodeNotFound            = 1010
	CodeAccountDisabled     = 1011
	CodeInvalidToken        = 1016
	CodeInsufficientStock   = 1024
	CodePermissionNotFound  = 2001
	CodeInsufficientPermissions = 2002
	CodeCannotModifySystemRole  = 2003
	CodeCannotDeleteSystemRole  = 2004
	CodeRoleNotFound            = 2005
	CodeAddressNotFound    = 1012
	CodeAddressLimit       = 1013
	CodeUsernameExists     = 1014

	CodeBrandNotFound       = 4001
	CodeCategoryNotFound    = 4010
	CodeAttributeNotFound   = 4020
	CodeProductNotFound     = 4030
	CodeSKUNotFound              = 4031
	CodeProductVersionNotFound   = 4032
	CodeProductDescriptionNotFound = 4033
	CodeInventoryNotFound   = 5001
	CodeInvalidStockChange   = 5003
	CodeInventoryLogNotFound = 5005
	CodePaymentNotFound      = 6001
	CodeCartNotFound          = 8001
	CodeRefundNotFound      = 6010
	CodeRefundFailed        = 6011
	CodeNotificationNotFound         = 9001
	CodeNotificationTemplateNotFound = 9002
	CodeInvalidOrderStatus           = 7002
	CodeOrderItemNotFound   = 7003

	CodeMerchantsNotFound           = 10001
	CodeMerchantBankAccountNotFound  = 10002
	CodeMerchantContactNotFound      = 10003
	CodeMerchantQualificationNotFound = 10004
	CodeMerchantWithdrawalNotFound   = 10005
	CodeReviewNotFound              = 10006
	CodeDepartmentNotFound          = 11001
)

// ============================================================================
// 错误码：通用域 1001-1999
// ============================================================================

var (
	ErrInvalidParams        = New(CodeInvalidParams, "参数错误")
	ErrUnauthorized         = New(CodeUnauthorized, "未授权，请先登录")
	ErrUserNotFound         = New(CodeUserNotFound, "用户不存在")
	ErrAddressNotFound      = New(CodeAddressNotFound, "地址不存在")
	ErrAddressLimit         = New(CodeAddressLimit, "地址数量已达上限")
	ErrUsernameAlreadyExists = New(CodeUsernameExists, "用户名已存在")
	ErrOrderNotFound       = New(CodeOrderNotFound, "订单不存在")
	ErrPaymentFailed       = New(CodePaymentFailed, "支付失败")
	ErrInvalidCredentials  = New(CodeInvalidCredentials, "用户名或密码错误")
	ErrNotFound            = New(CodeNotFound, "资源不存在")
	ErrAccountDisabled     = New(CodeAccountDisabled, "账号已被禁用")
	ErrInvalidToken        = New(CodeInvalidToken, "Token 无效或已过期")
	ErrInsufficientStock   = New(CodeInsufficientStock, "库存不足")
)

// ============================================================================
// 错误码：权限域 2001-2999
// ============================================================================

var (
	ErrPermissionNotFound      = New(CodePermissionNotFound, "权限不存在")
	ErrInsufficientPermissions = New(CodeInsufficientPermissions, "无权限，需要管理员角色")
	ErrCannotModifySystemRole  = New(CodeCannotModifySystemRole, "不能修改系统角色")
	ErrCannotDeleteSystemRole  = New(CodeCannotDeleteSystemRole, "不能删除系统角色")
	ErrRoleNotFound            = New(CodeRoleNotFound, "角色不存在")
)

// ============================================================================
// 错误码：品牌/商品域 4001-4099
// ============================================================================

var (
	ErrBrandNotFound    = New(CodeBrandNotFound, "品牌不存在")
	ErrCategoryNotFound = New(CodeCategoryNotFound, "类目不存在")
	ErrAttributeNotFound = New(CodeAttributeNotFound, "属性不存在")
	ErrProductNotFound           = New(CodeProductNotFound, "产品不存在")
	ErrSKUNotFound               = New(CodeSKUNotFound, "SKU 不存在")
	ErrProductVersionNotFound    = New(CodeProductVersionNotFound, "商品版本不存在")
	ErrProductDescriptionNotFound = New(CodeProductDescriptionNotFound, "商品描述不存在")
)

// ============================================================================
// 错误码：库存域 5001-5099
// ============================================================================

var (
	ErrInventoryNotFound  = New(CodeInventoryNotFound, "库存记录不存在")
	ErrInvalidStockChange  = New(CodeInvalidStockChange, "无效的库存变动")
	ErrInventoryLogNotFound = New(CodeInventoryLogNotFound, "库存日志不存在")
)

// ============================================================================
// 错误码：购物车域 8001-8099
// ============================================================================

var (
	ErrCartNotFound = New(CodeCartNotFound, "购物车不存在")
)

// ============================================================================
// 错误码：交易域 6001-6099
// ============================================================================

var (
	ErrPaymentNotFound = New(CodePaymentNotFound, "支付记录不存在")
	ErrRefundNotFound  = New(CodeRefundNotFound, "退款记录不存在")
	ErrRefundFailed    = New(CodeRefundFailed, "退款失败")
)

// ============================================================================
// 错误码：订单域 7001-7099
// ============================================================================

var (
	ErrInvalidOrderStatus = New(CodeInvalidOrderStatus, "无效的订单状态变更")
	ErrOrderItemNotFound  = New(CodeOrderItemNotFound, "订单项不存在")
)

// ============================================================================
// 错误码：通知域 9001-9099
// ============================================================================

var (
	ErrNotificationNotFound         = New(CodeNotificationNotFound, "通知不存在")
	ErrNotificationTemplateNotFound = New(CodeNotificationTemplateNotFound, "通知模板不存在")
)


// ============================================================================
// 错误码：商家域 10001-10099
// ============================================================================

var (
	ErrMerchantsNotFound            = New(CodeMerchantsNotFound, "商家不存在")
	ErrMerchantBankAccountNotFound  = New(CodeMerchantBankAccountNotFound, "银行账户不存在")
	ErrMerchantContactNotFound      = New(CodeMerchantContactNotFound, "联系人不存在")
	ErrMerchantQualificationNotFound = New(CodeMerchantQualificationNotFound, "资质不存在")
	ErrMerchantWithdrawalNotFound   = New(CodeMerchantWithdrawalNotFound, "提现记录不存在")
	ErrReviewNotFound               = New(CodeReviewNotFound, "评价不存在")
	ErrDepartmentNotFound           = New(CodeDepartmentNotFound, "部门不存在")
)