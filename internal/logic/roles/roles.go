package roles

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/roles/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sRoles struct{}

func init() {
	service.RegisterRoles(&sRoles{})
}

func (s *sRoles) List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.Roles.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.Name != "" {
		m = m.WhereLike(dao.Roles.Columns().Name, "%"+req.Name+"%")
	}
	if req.RoleType != "" {
		m = m.Where(dao.Roles.Columns().RoleType, req.RoleType)
	}
	if req.Status > 0 {
		m = m.Where(dao.Roles.Columns().Status, req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.RoleListRes{
			List:  make([]*entity.Roles, 0),
			Total: 0,
		}, nil
	}
	var list []*entity.Roles
	err = m.Page(page, size).OrderAsc(dao.Roles.Columns().SortOrder).OrderDesc(dao.Roles.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.RoleListRes{List: list, Total: total}, nil
}

func (s *sRoles) Detail(ctx context.Context, req *v1.RoleDetailReq) (res *v1.RoleDetailRes, err error) {
	var entity *entity.Roles
	err = dao.Roles.Ctx(ctx).Where(dao.Roles.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "角色不存在")
	}
	return &v1.RoleDetailRes{Roles: entity}, nil
}

func (s *sRoles) Create(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	result, err := dao.Roles.Ctx(ctx).Insert(do.Roles{
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
	return &v1.RoleCreateRes{Id: id}, nil
}

func (s *sRoles) Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	count, err := dao.Roles.Ctx(ctx).Where(dao.Roles.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "角色不存在")
	}
	_, err = dao.Roles.Ctx(ctx).Data(do.Roles{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		RoleType:    req.RoleType,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}).Where(dao.Roles.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}

func (s *sRoles) Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	count, err := dao.Roles.Ctx(ctx).Where(dao.Roles.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "角色不存在")
	}
	_, err = dao.Roles.Ctx(ctx).Where(dao.Roles.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}

func (s *sRoles) IsAdmin(ctx context.Context, staffId int64) (bool, error) {
	count, err := dao.Roles.Ctx(ctx).
		InnerJoin("sys_staff_roles", "sys_staff_roles.role_id = sys_roles.id").
		Where("sys_staff_roles.staff_id", staffId).
		Where(dao.Roles.Columns().RoleType, "builtin").
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sRoles) HasPermission(ctx context.Context, staffId int64, permissionName string) (bool, error) {
	var perm entity.Permissions
	err := dao.Permissions.Ctx(ctx).
		Where(dao.Permissions.Columns().Name, permissionName).
		Scan(&perm)
	if err != nil || perm.Id == 0 {
		return false, nil
	}
	count, err := dao.RolePermissions.Ctx(ctx).
		InnerJoin("sys_staff_roles", "sys_staff_roles.role_id = sys_role_permissions.role_id").
		Where(dao.RolePermissions.Columns().PermissionId, perm.Id).
		Where("sys_staff_roles.staff_id", staffId).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sRoles) GetPermissions(ctx context.Context, staffId int64) (permissions []string, roleNames []string, err error) {
	roleVars, err := dao.Roles.Ctx(ctx).
		InnerJoin("sys_staff_roles", "sys_staff_roles.role_id = sys_roles.id").
		Where("sys_staff_roles.staff_id", staffId).
		Array(dao.Roles.Columns().Name)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range roleVars {
		roleNames = append(roleNames, v.String())
	}
	permVars, err := dao.Permissions.Ctx(ctx).
		InnerJoin("sys_role_permissions", "sys_role_permissions.permission_id = sys_permissions.id").
		InnerJoin("sys_staff_roles", "sys_staff_roles.role_id = sys_role_permissions.role_id").
		Where("sys_staff_roles.staff_id", staffId).
		Group(dao.Permissions.Columns().Name).
		Array(dao.Permissions.Columns().Name)
	if err != nil {
		return nil, roleNames, err
	}
	for _, v := range permVars {
		permissions = append(permissions, v.String())
	}
	if permissions == nil {
		permissions = make([]string, 0)
	}
	if roleNames == nil {
		roleNames = make([]string, 0)
	}
	return
}
