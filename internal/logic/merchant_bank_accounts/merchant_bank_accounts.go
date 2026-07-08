package merchantBankAccounts

import (
	"context"

	"gf-eshop/api/merchant_bank_accounts/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantBankAccounts struct{}

func init() {
	service.RegisterMerchantBankAccounts(&sMerchantBankAccounts{})
}

func (s *sMerchantBankAccounts) List(ctx context.Context, req *v1.MerchantBankAccountsListReq) (res *v1.MerchantBankAccountsListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.MerchantBankAccounts.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.MerchantBankAccounts.Columns().MerchantId, req.MerchantId)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.MerchantBankAccountsListRes{
			List:  make([]*entity.MerchantBankAccounts, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.MerchantBankAccounts
	err = m.Page(page, size).OrderDesc(dao.MerchantBankAccounts.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantBankAccountsListRes{List: list, Total: total}, nil
}

func (s *sMerchantBankAccounts) Detail(ctx context.Context, req *v1.MerchantBankAccountsDetailReq) (res *v1.MerchantBankAccountsDetailRes, err error) {
	var entity *entity.MerchantBankAccounts
	err = dao.MerchantBankAccounts.Ctx(ctx).Where(dao.MerchantBankAccounts.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrMerchantBankAccountNotFound
	}
	return &v1.MerchantBankAccountsDetailRes{MerchantBankAccounts: entity}, nil
}

func (s *sMerchantBankAccounts) Create(ctx context.Context, req *v1.MerchantBankAccountsCreateReq) (res *v1.MerchantBankAccountsCreateRes, err error) {
	result, err := dao.MerchantBankAccounts.Ctx(ctx).Insert(do.MerchantBankAccounts{
		MerchantId:  req.MerchantId,
		BankName:    req.BankName,
		BankBranch:  req.BankBranch,
		AccountName: req.AccountName,
		AccountNo:   req.AccountNo,
		AccountType: req.AccountType,
		IsDefault:   req.IsDefault,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MerchantBankAccountsCreateRes{Id: id}, nil
}

func (s *sMerchantBankAccounts) Update(ctx context.Context, req *v1.MerchantBankAccountsUpdateReq) (res *v1.MerchantBankAccountsUpdateRes, err error) {
	count, err := dao.MerchantBankAccounts.Ctx(ctx).Where(dao.MerchantBankAccounts.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrMerchantBankAccountNotFound
	}
	_, err = dao.MerchantBankAccounts.Ctx(ctx).Data(do.MerchantBankAccounts{
		BankName:    req.BankName,
		BankBranch:  req.BankBranch,
		AccountName: req.AccountName,
		AccountNo:   req.AccountNo,
		AccountType: req.AccountType,
		IsDefault:   req.IsDefault,
		Status:      req.Status,
	}).Where(dao.MerchantBankAccounts.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantBankAccountsUpdateRes{}, nil
}

func (s *sMerchantBankAccounts) Delete(ctx context.Context, req *v1.MerchantBankAccountsDeleteReq) (res *v1.MerchantBankAccountsDeleteRes, err error) {
	_, err = dao.MerchantBankAccounts.Ctx(ctx).Where(dao.MerchantBankAccounts.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantBankAccountsDeleteRes{}, nil
}
