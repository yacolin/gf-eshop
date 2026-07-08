package merchantRoles

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/merchant_roles/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantRoles struct{}

func init() {
	service.RegisterMerchantRoles(&sMerchantRoles{})
}

func (s *sMerchantRoles) List(ctx context.Context, req *v1.MerchantRolesListReq) (res *v1.MerchantRolesListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.MerchantRoles.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.MerchantId > 0 {
		m = m.Where(dao.MerchantRoles.Columns().MerchantId, req.MerchantId)
	}
	if req.Name != "" {
		m = m.WhereLike(dao.MerchantRoles.Columns().Name, "%"+req.Name+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.MerchantRoles.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.MerchantRolesListRes{
			List:  make([]*entity.MerchantRoles, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.MerchantRoles
	err = m.Page(page, size).OrderAsc(dao.MerchantRoles.Columns().SortOrder).OrderDesc(dao.MerchantRoles.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolesListRes{List: list, Total: total}, nil
}

func (s *sMerchantRoles) Detail(ctx context.Context, req *v1.MerchantRolesDetailReq) (res *v1.MerchantRolesDetailRes, err error) {
	var entity *entity.MerchantRoles
	err = dao.MerchantRoles.Ctx(ctx).Where(dao.MerchantRoles.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "角色不存在")
	}
	return &v1.MerchantRolesDetailRes{MerchantRoles: entity}, nil
}

func (s *sMerchantRoles) Create(ctx context.Context, req *v1.MerchantRolesCreateReq) (res *v1.MerchantRolesCreateRes, err error) {
	result, err := dao.MerchantRoles.Ctx(ctx).Insert(do.MerchantRoles{
		MerchantId:  req.MerchantId,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		RoleType:    req.RoleType,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MerchantRolesCreateRes{Id: id}, nil
}

func (s *sMerchantRoles) Update(ctx context.Context, req *v1.MerchantRolesUpdateReq) (res *v1.MerchantRolesUpdateRes, err error) {
	count, err := dao.MerchantRoles.Ctx(ctx).Where(dao.MerchantRoles.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "角色不存在")
	}
	_, err = dao.MerchantRoles.Ctx(ctx).Data(do.MerchantRoles{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		RoleType:    req.RoleType,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}).Where(dao.MerchantRoles.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolesUpdateRes{}, nil
}

func (s *sMerchantRoles) Delete(ctx context.Context, req *v1.MerchantRolesDeleteReq) (res *v1.MerchantRolesDeleteRes, err error) {
	_, err = dao.MerchantRoles.Ctx(ctx).Where(dao.MerchantRoles.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolesDeleteRes{}, nil
}
