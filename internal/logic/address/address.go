package address

import (
	"context"

	"gf-eshop/api/address/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

const maxAddressPerUser = 20

type sAddress struct{}

func init() {
	service.RegisterAddress(&sAddress{})
}

func getUserId(ctx context.Context) (int64, error) {
	claims := utility.GetUserClaims(ctx)
	if claims == nil {
		return 0, errcode.ErrUnauthorized
	}
	return claims.UserId, nil
}

func (s *sAddress) Create(ctx context.Context, req *v1.AddressCreateReq) (res *v1.AddressCreateRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	count, err := dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().UserId, userId).Count()
	if err != nil {
		return nil, err
	}
	if count >= maxAddressPerUser {
		return nil, errcode.ErrAddressLimit
	}

	if req.IsDefault == 1 {
		_, _ = dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().UserId, userId).
			Update(do.Addresses{IsDefault: 0})
	}

	result, err := dao.Addresses.Ctx(ctx).Insert(do.Addresses{
		UserId:    userId,
		Consignee: req.Consignee,
		Phone:     req.Phone,
		Country:   req.Country,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		ZipCode:   req.ZipCode,
		Tag:       req.Tag,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.AddressCreateRes{Id: id}, nil
}

func (s *sAddress) List(ctx context.Context, req *v1.AddressListReq) (res *v1.AddressListRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	var list []*entity.Addresses
	m := dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().UserId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.AddressListRes{List: make([]*entity.Addresses, 0), Total: 0}, nil
	}

	err = m.OrderDesc(dao.Addresses.Columns().IsDefault).OrderDesc(dao.Addresses.Columns().UpdatedAt).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.AddressListRes{List: list, Total: total}, nil
}

func (s *sAddress) Detail(ctx context.Context, req *v1.AddressDetailReq) (res *v1.AddressDetailRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	var addr *entity.Addresses
	err = dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().Id, req.Id).
		Where(dao.Addresses.Columns().UserId, userId).Scan(&addr)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errcode.ErrAddressNotFound
	}
	return &v1.AddressDetailRes{Addresses: addr}, nil
}

func (s *sAddress) GetDefault(ctx context.Context, req *v1.AddressGetDefaultReq) (res *v1.AddressGetDefaultRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	var addr *entity.Addresses
	err = dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().UserId, userId).
		Where(dao.Addresses.Columns().IsDefault, 1).Scan(&addr)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errcode.ErrAddressNotFound
	}
	return &v1.AddressGetDefaultRes{Addresses: addr}, nil
}

func (s *sAddress) Update(ctx context.Context, req *v1.AddressUpdateReq) (res *v1.AddressUpdateRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	count, err := dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().Id, req.Id).
		Where(dao.Addresses.Columns().UserId, userId).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrAddressNotFound
	}

	if req.IsDefault != nil && *req.IsDefault == 1 {
		_, _ = dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().UserId, userId).
			Update(do.Addresses{IsDefault: 0})
	}

	updateData := do.Addresses{}
	if req.Consignee != nil {
		updateData.Consignee = *req.Consignee
	}
	if req.Phone != nil {
		updateData.Phone = *req.Phone
	}
	if req.Country != nil {
		updateData.Country = *req.Country
	}
	if req.Province != nil {
		updateData.Province = *req.Province
	}
	if req.City != nil {
		updateData.City = *req.City
	}
	if req.District != nil {
		updateData.District = *req.District
	}
	if req.Detail != nil {
		updateData.Detail = *req.Detail
	}
	if req.ZipCode != nil {
		updateData.ZipCode = *req.ZipCode
	}
	if req.Tag != nil {
		updateData.Tag = *req.Tag
	}
	if req.IsDefault != nil {
		updateData.IsDefault = *req.IsDefault
	}

	_, err = dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().Id, req.Id).
		Where(dao.Addresses.Columns().UserId, userId).Update(updateData)
	if err != nil {
		return nil, err
	}
	return &v1.AddressUpdateRes{}, nil
}

func (s *sAddress) Delete(ctx context.Context, req *v1.AddressDeleteReq) (res *v1.AddressDeleteRes, err error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}

	result, err := dao.Addresses.Ctx(ctx).Where(dao.Addresses.Columns().Id, req.Id).
		Where(dao.Addresses.Columns().UserId, userId).Delete()
	if err != nil {
		return nil, err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return nil, errcode.ErrAddressNotFound
	}
	return &v1.AddressDeleteRes{}, nil
}
