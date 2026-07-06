package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"gf-eshop/internal/model/entity"
)

type PaymentsCreateReq struct {
	g.Meta        `path:"/payments" tags:"Payments" method:"post" summary:"创建支付"`
	OrderNo       string `json:"order_no"        v:"required" description:"订单号"`
	Amount        int64  `json:"amount"          v:"required|min:1" description:"支付金额"`
	PaymentMethod string `json:"payment_method"  v:"required" description:"支付方式"`
	Channel       string `json:"channel"         description:"支付渠道"`
}

type PaymentsCreateRes struct {
	*entity.Payments
}

type PaymentsCallbackReq struct {
	g.Meta        `path:"/payments/callback" tags:"Payments" method:"post" summary:"支付回调"`
	PaymentNo     string `json:"payment_no"     v:"required" description:"支付单号"`
	TransactionID string `json:"transaction_id" v:"required" description:"渠道交易ID"`
	Channel       string `json:"channel"        description:"渠道"`
	Status        string `json:"status"         v:"required|in:success,failed" description:"状态"`
	FailureReason string `json:"failure_reason" description:"失败原因"`
	RawBody       string `json:"raw_body"       description:"原始回调数据"`
}

type PaymentsCallbackRes struct {
	*entity.Payments
}

type PaymentsGetReq struct {
	g.Meta    `path:"/payments" tags:"Payments" method:"get" summary:"查询支付"`
	PaymentNo string `json:"payment_no" description:"支付单号"`
	OrderNo   string `json:"order_no"   description:"订单号"`
}

type PaymentsGetRes struct {
	*entity.Payments
}

type RefundsCreateReq struct {
	g.Meta    `path:"/payments/refunds" tags:"Payments" method:"post" summary:"创建退款"`
	PaymentNo string `json:"payment_no" v:"required" description:"支付单号"`
	Amount    int64  `json:"amount"     v:"required|min:1" description:"退款金额"`
	Reason    string `json:"reason"     description:"退款原因"`
}

type RefundsCreateRes struct {
	*entity.Refunds
}

type RefundsListReq struct {
	g.Meta        `path:"/payments/refunds" tags:"Refunds" method:"get" summary:"退款列表"`
	Page          int    `json:"page"        description:"页码"`
	PageSize      int    `json:"page_size"   description:"每页条数"`
	PaymentNo     string `json:"payment_no"  description:"支付单号"`
	OrderNo       string `json:"order_no"    description:"订单号"`
	RefundNo      string `json:"refund_no"   description:"退款单号"`
	Status        string `json:"status"      description:"退款状态"`
}

type RefundsListRes struct {
	List  []*entity.Refunds `json:"list"`
	Total int               `json:"total"`
}

type RefundsDetailReq struct {
	g.Meta   `path:"/payments/refunds/{refund_no}" tags:"Refunds" method:"get" summary:"退款详情"`
	RefundNo string `json:"refund_no"`
}

type RefundsDetailRes struct {
	*entity.Refunds
}