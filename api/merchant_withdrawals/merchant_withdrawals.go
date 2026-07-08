package merchantWithdrawals

import (
	"context"

	"gf-eshop/api/merchant_withdrawals/v1"
)

type IMerchantWithdrawalsV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Apply(ctx context.Context, req *v1.ApplyReq) (res *v1.ApplyRes, err error)
	Approve(ctx context.Context, req *v1.ApproveReq) (res *v1.ApproveRes, err error)
	Reject(ctx context.Context, req *v1.RejectReq) (res *v1.RejectRes, err error)
}
