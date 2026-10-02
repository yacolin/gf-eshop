package errcode

import (
	"testing"
)

func TestNoDuplicateCodes(t *testing.T) {
	all := []error{
		ErrInvalidParams,
		ErrUnauthorized,
		ErrUserNotFound,
		ErrOrderNotFound,
		ErrPaymentFailed,
		ErrInvalidCredentials,
		ErrNotFound,
		ErrAccountDisabled,
		ErrInvalidToken,
		ErrInsufficientStock,
		ErrPermissionNotFound,
		ErrInsufficientPermissions,
		ErrCannotModifySystemRole,
		ErrCannotDeleteSystemRole,
		ErrBrandNotFound,
		ErrCategoryNotFound,
		ErrProductNotFound,
		ErrSKUNotFound,
		ErrAttributeNotFound,
		ErrInventoryNotFound,
		ErrInvalidStockChange,
		ErrPaymentNotFound,
		ErrRefundNotFound,
		ErrRefundFailed,
		ErrInvalidOrderStatus,
		ErrOrderItemNotFound,
		ErrMerchantsNotFound,
		ErrMerchantBankAccountNotFound,
		ErrMerchantContactNotFound,
		ErrMerchantQualificationNotFound,
		ErrMerchantWithdrawalNotFound,
		ErrReviewNotFound,
		ErrUsernameAlreadyExists,
		ErrAddressNotFound,
		ErrAddressLimit,
		ErrVerifyCodeInvalid,
		ErrVerifyCodeTooFrequent,
		ErrVerifyCodeAttemptsExceed,
		ErrVerifyCodeDailyLimit,
		ErrVerifyChannelNotReady,
		ErrVerifySendFailed,
		ErrEmailAlreadyExists,
	}

	seen := make(map[int]string)
	for _, e := range all {
		code := CodeOf(e)
		if name, ok := seen[code]; ok {
			t.Errorf("duplicate code %d: %q conflicts with existing %q", code, e.Error(), name)
		}
		seen[code] = e.Error()
	}
}
