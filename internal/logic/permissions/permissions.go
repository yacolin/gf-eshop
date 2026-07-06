package permissions

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"gf-eshop/api/permissions/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sPermissions struct{}

func init() {
	service.RegisterPermissions(&sPermissions{})
}

func (s *sPermissions) List(ctx context.Context, req *v1.PermissionListReq) (res *v1.PermissionListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.Permissions.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.Resource != "" {
		m = m.Where(dao.Permissions.Columns().Resource, req.Resource)
	}
	if req.Action != "" {
		m = m.Where(dao.Permissions.Columns().Action, req.Action)
	}
	if req.Category != "" {
		m = m.Where(dao.Permissions.Columns().Category, req.Category)
	}
	if req.Status > 0 {
		m = m.Where(dao.Permissions.Columns().Status, req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.PermissionListRes{
			List:  make([]*entity.Permissions, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.Permissions
	err = m.Page(page, size).OrderAsc(dao.Permissions.Columns().SortOrder).OrderDesc(dao.Permissions.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.PermissionListRes{List: list, Total: total}, nil
}

func (s *sPermissions) Detail(ctx context.Context, req *v1.PermissionDetailReq) (res *v1.PermissionDetailRes, err error) {
	var entity *entity.Permissions
	err = dao.Permissions.Ctx(ctx).Where(dao.Permissions.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrPermissionNotFound
	}
	return &v1.PermissionDetailRes{Permissions: entity}, nil
}

func (s *sPermissions) checkAdmin(ctx context.Context) error {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return errcode.ErrUnauthorized
	}
	isAdmin, err := service.Roles().IsAdmin(ctx, claims.StaffId)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errcode.ErrInsufficientPermissions
	}
	return nil
}

func (s *sPermissions) Create(ctx context.Context, req *v1.PermissionCreateReq) (res *v1.PermissionCreateRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	result, err := dao.Permissions.Ctx(ctx).Insert(do.Permissions{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.PermissionCreateRes{Id: id}, nil
}

func (s *sPermissions) Update(ctx context.Context, req *v1.PermissionUpdateReq) (res *v1.PermissionUpdateRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	count, err := dao.Permissions.Ctx(ctx).Where(dao.Permissions.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrPermissionNotFound
	}
	_, err = dao.Permissions.Ctx(ctx).Data(do.Permissions{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}).Where(dao.Permissions.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.PermissionUpdateRes{}, nil
}

func (s *sPermissions) Delete(ctx context.Context, req *v1.PermissionDeleteReq) (res *v1.PermissionDeleteRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	count, err := dao.Permissions.Ctx(ctx).Where(dao.Permissions.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrPermissionNotFound
	}
	_, err = dao.Permissions.Ctx(ctx).Where(dao.Permissions.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.PermissionDeleteRes{}, nil
}

func (s *sPermissions) Check(ctx context.Context, req *v1.PermissionCheckReq) (res *v1.PermissionCheckRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return &v1.PermissionCheckRes{HasPermission: false}, nil
	}
	has, err := service.Roles().HasPermission(ctx, claims.StaffId, req.Permission)
	if err != nil {
		return nil, err
	}
	return &v1.PermissionCheckRes{HasPermission: has}, nil
}

func (s *sPermissions) RolePermissionList(ctx context.Context, req *v1.RolePermissionListReq) (res *v1.RolePermissionListRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	var list []*entity.Permissions
	err = dao.Permissions.Ctx(ctx).
		Fields("sys_permissions.*").
		InnerJoin("sys_role_permissions", "sys_role_permissions.permission_id = sys_permissions.id").
		Where("sys_role_permissions.role_id", req.RoleId).
		OrderAsc(dao.Permissions.Columns().SortOrder).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.Permissions, 0)
	}
	return &v1.RolePermissionListRes{List: list}, nil
}

func (s *sPermissions) RolePermissionUpdate(ctx context.Context, req *v1.RolePermissionUpdateReq) (res *v1.RolePermissionUpdateRes, err error) {
	if err := s.checkAdmin(ctx); err != nil {
		return nil, err
	}
	roleCount, err := dao.Roles.Ctx(ctx).Where(dao.Roles.Columns().Id, req.RoleId).Count()
	if err != nil {
		return nil, err
	}
	if roleCount == 0 {
		return nil, errcode.ErrRoleNotFound
	}
	err = dao.RolePermissions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.RolePermissions.Ctx(ctx).TX(tx).Where(dao.RolePermissions.Columns().RoleId, req.RoleId).Delete()
		if err != nil {
			return err
		}
		for _, pid := range req.PermissionIds {
			_, err = tx.Model("sys_role_permissions").Insert(do.RolePermissions{
				RoleId:       req.RoleId,
				PermissionId: pid,
				ScopeType:    "platform",
				ScopeId:      0,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.RolePermissionUpdateRes{}, nil
}
