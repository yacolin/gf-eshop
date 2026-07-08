package merchantWithdrawals

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/api/merchant_withdrawals/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantWithdrawals struct{}

func init() {
	service.RegisterMerchantWithdrawals(&sMerchantWithdrawals{})
}

func (s *sMerchantWithdrawals) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.MerchantWithdrawals.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.MerchantWithdrawals.Columns().MerchantId, req.MerchantId)
	}
	if req.Status > 0 {
		m = m.Where(dao.MerchantWithdrawals.Columns().Status, req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ListRes{
			List:  make([]*entity.MerchantWithdrawals, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.MerchantWithdrawals
	err = m.Page(page, size).OrderDesc(dao.MerchantWithdrawals.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ListRes{List: list, Total: total}, nil
}

func (s *sMerchantWithdrawals) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	var entity *entity.MerchantWithdrawals
	err = dao.MerchantWithdrawals.Ctx(ctx).Where(dao.MerchantWithdrawals.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrMerchantWithdrawalNotFound
	}
	return &v1.DetailRes{MerchantWithdrawals: entity}, nil
}

func (s *sMerchantWithdrawals) Apply(ctx context.Context, req *v1.ApplyReq) (res *v1.ApplyRes, err error) {
	withdrawNo := fmt.Sprintf("W%s%s", gtime.Now().Format("YmdHis"), grand.Letters(4))
	result, err := dao.MerchantWithdrawals.Ctx(ctx).Insert(do.MerchantWithdrawals{
		MerchantId:    req.MerchantId,
		WithdrawNo:    withdrawNo,
		Amount:        req.Amount,
		BankAccountId: req.BankAccountId,
		Status:        0,
		AppliedAt:     gtime.Now(),
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.ApplyRes{Id: id}, nil
}

func (s *sMerchantWithdrawals) Approve(ctx context.Context, req *v1.ApproveReq) (res *v1.ApproveRes, err error) {
	count, err := dao.MerchantWithdrawals.Ctx(ctx).Where(dao.MerchantWithdrawals.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantWithdrawalNotFound
	}
	_, err = dao.MerchantWithdrawals.Ctx(ctx).Data(do.MerchantWithdrawals{
		Status:     1,
		ApprovedAt: gtime.Now(),
	}).Where(dao.MerchantWithdrawals.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.ApproveRes{}, nil
}

func (s *sMerchantWithdrawals) Reject(ctx context.Context, req *v1.RejectReq) (res *v1.RejectRes, err error) {
	count, err := dao.MerchantWithdrawals.Ctx(ctx).Where(dao.MerchantWithdrawals.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantWithdrawalNotFound
	}
	_, err = dao.MerchantWithdrawals.Ctx(ctx).Data(do.MerchantWithdrawals{
		Status:      3,
		AuditRemark: req.AuditRemark,
		ApprovedAt:  gtime.Now(),
	}).Where(dao.MerchantWithdrawals.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.RejectRes{}, nil
}
