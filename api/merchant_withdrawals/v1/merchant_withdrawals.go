package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type ListReq struct {
	g.Meta `path:"/merchant-withdrawals" tags:"MerchantWithdrawals" method:"get" summary:"提现列表"`

	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	MerchantId int64 `json:"merchant_id"`
	Status     int   `json:"status"`
}

type ListRes struct {
	List  []*entity.MerchantWithdrawals `json:"list"`
	Total int                           `json:"total"`
}

// ---------- Detail ----------
type DetailReq struct {
	g.Meta `path:"/merchant-withdrawals/{id}" tags:"MerchantWithdrawals" method:"get" summary:"提现详情"`
	Id     int64 `json:"id"`
}

type DetailRes struct {
	*entity.MerchantWithdrawals
}

// ---------- Apply ----------
type ApplyReq struct {
	g.Meta `path:"/merchant-withdrawals" tags:"MerchantWithdrawals" method:"post" summary:"申请提现"`

	MerchantId    int64 `json:"merchant_id"      v:"required"`
	Amount        int64 `json:"amount"           v:"required|min:1"`
	BankAccountId int64 `json:"bank_account_id"  v:"required"`
}

type ApplyRes struct {
	Id int64 `json:"id"`
}

// ---------- Approve ----------
type ApproveReq struct {
	g.Meta `path:"/merchant-withdrawals/{id}/approve" tags:"MerchantWithdrawals" method:"put" summary:"审核通过"`
	Id     int64 `json:"id" v:"required"`
}

type ApproveRes struct{}

// ---------- Reject ----------
type RejectReq struct {
	g.Meta      `path:"/merchant-withdrawals/{id}/reject" tags:"MerchantWithdrawals" method:"put" summary:"审核拒绝"`
	Id          int64  `json:"id"           v:"required"`
	AuditRemark string `json:"audit_remark" v:""`
}

type RejectRes struct{}
