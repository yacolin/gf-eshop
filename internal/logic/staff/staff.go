package staff

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"gf-eshop/api/staff/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

func (s *sStaff) logAudit(ctx context.Context, operation, resource, resourceId, detail string, result int, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return
	}
	failureReason := ""
	if err != nil {
		failureReason = err.Error()
	}
	_ = service.OperationLogs().Log(ctx, &service.OperationLogInput{
		StaffId:       claims.StaffId,
		StaffName:     claims.RealName,
		Operation:     operation,
		Resource:      resource,
		ResourceId:    resourceId,
		Detail:        detail,
		Result:        result,
		FailureReason: failureReason,
	})
}

type sStaff struct{}

func init() {
	service.RegisterStaff(&sStaff{})
}

func refreshRedisKey(staffId int64, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", utility.RefreshRedisKey, staffId, tokenId)
}

func saveRefreshToken(ctx context.Context, staffId int64, tokenId string, expire time.Duration) {
	key := refreshRedisKey(staffId, tokenId)
	g.Redis().Do(ctx, "SETEX", key, int(expire.Seconds()), "1")
}

func deleteRefreshToken(ctx context.Context, staffId int64, tokenId string) {
	g.Redis().Do(ctx, "DEL", refreshRedisKey(staffId, tokenId))
}

func deleteAllRefreshTokens(ctx context.Context, staffId int64) {
	pattern := fmt.Sprintf("%s:%d:*", utility.RefreshRedisKey, staffId)
	keys, _ := g.Redis().Do(ctx, "KEYS", pattern)
	if !keys.IsNil() {
		for _, key := range keys.Vars() {
			g.Redis().Do(ctx, "DEL", key.String())
		}
	}
}

func (s *sStaff) Login(ctx context.Context, req *v1.StaffLoginReq) (res *v1.StaffLoginRes, err error) {
	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Username, req.Username).
		Scan(&staff)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errcode.ErrInvalidCredentials
	}
	if staff.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errcode.ErrInvalidCredentials
	}

	now := time.Now()
	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	ip := ""
	device := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
		device = r.Header.Get("User-Agent")
		if len(device) > 100 {
			device = device[:100]
		}
	}
	_, _ = dao.Staff.Ctx(ctx).Where(dao.Staff.Columns().Id, staff.Id).Update(do.Staff{
		LastLoginIp: ip,
		LastLoginAt: gtime.New(now),
	})
	if _, err := dao.SysLoginHistories.Ctx(ctx).Insert(do.SysLoginHistories{
		StaffId:     staff.Id,
		LoginIp:     ip,
		LoginDevice: device,
		LoginMethod: "password",
		LoginStatus: 1,
	}); err != nil {
		g.Log().Warning(ctx, "insert login history failed: %v", err)
	}

	refreshClaims, _ := utility.ParseStaffToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveRefreshToken(ctx, staff.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.StaffLoginRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
		StaffId:      staff.Id,
		Username:     staff.Username,
		RealName:     staff.RealName,
	}, nil
}

func (s *sStaff) RefreshToken(ctx context.Context, req *v1.StaffRefreshTokenReq) (res *v1.StaffRefreshTokenRes, err error) {
	claims, err := utility.ParseStaffToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	if claims.TokenType != utility.TokenTypeRefresh {
		return nil, errcode.ErrInvalidParams
	}

	key := refreshRedisKey(claims.StaffId, claims.TokenId)
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil || v.IsNil() {
		return nil, errcode.ErrInvalidToken
	}

	deleteRefreshToken(ctx, claims.StaffId, claims.TokenId)

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil || staff == nil {
		return nil, errcode.ErrUserNotFound
	}
	if staff.Status != 1 {
		return nil, errcode.ErrAccountDisabled
	}

	pair, err := utility.GenerateTokenPair(ctx, staff.Id, staff.Username, staff.RealName)
	if err != nil {
		return nil, gerror.NewCode(errcode.Code(57), "生成Token失败")
	}

	refreshClaims, _ := utility.ParseStaffToken(ctx, pair.RefreshToken)
	if refreshClaims != nil {
		saveRefreshToken(ctx, staff.Id, refreshClaims.TokenId, utility.JwtRefreshExpire(ctx))
	}

	return &v1.StaffRefreshTokenRes{
		AccessToken:  pair.AccessToken,
		ExpireIn:     pair.ExpireIn,
		RefreshToken: pair.RefreshToken,
		RefreshIn:    pair.RefreshIn,
	}, nil
}

func (s *sStaff) Logout(ctx context.Context, req *v1.StaffLogoutReq) (res *v1.StaffLogoutRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims != nil {
		deleteAllRefreshTokens(ctx, claims.StaffId)
	}
	return &v1.StaffLogoutRes{}, nil
}

func (s *sStaff) Profile(ctx context.Context, req *v1.StaffProfileReq) (res *v1.StaffProfileRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}

	var staff *entity.Staff
	err = dao.Staff.Ctx(ctx).
		Where(dao.Staff.Columns().Id, claims.StaffId).
		Scan(&staff)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errcode.ErrUserNotFound
	}

	deptIds, deptNames := s.getStaffDepartments(ctx, staff.Id)
	return &v1.StaffProfileRes{
		Id:              staff.Id,
		Username:        staff.Username,
		RealName:        staff.RealName,
		Email:           staff.Email,
		Phone:           staff.Phone,
		Avatar:          staff.Avatar,
		Status:          staff.Status,
		LastLoginIp:     staff.LastLoginIp,
		DepartmentIds:   deptIds,
		DepartmentNames: deptNames,
	}, nil
}

func (s *sStaff) Permissions(ctx context.Context, req *v1.StaffPermissionsReq) (res *v1.StaffPermissionsRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return nil, errcode.ErrUnauthorized
	}
	perms, roles, err := service.Roles().GetPermissions(ctx, claims.StaffId)
	if err != nil {
		return nil, err
	}
	return &v1.StaffPermissionsRes{
		Roles:       roles,
		Permissions: perms,
	}, nil
}

func (s *sStaff) getStaffDepartments(ctx context.Context, staffId int64) (ids []int64, names []string) {
	type DeptInfo struct {
		DepartmentId   int64
		DepartmentName string
	}
	var depts []DeptInfo
	err := dao.StaffDepartments.Ctx(ctx).
		InnerJoin("sys_departments", "sys_departments.id = sys_staff_departments.department_id").
		Where("sys_staff_departments.staff_id", staffId).
		Fields("sys_staff_departments.department_id", "sys_departments.name AS department_name").
		Scan(&depts)
	if err != nil || len(depts) == 0 {
		return []int64{}, []string{}
	}
	for _, d := range depts {
		ids = append(ids, d.DepartmentId)
		names = append(names, d.DepartmentName)
	}
	return
}

func (s *sStaff) List(ctx context.Context, req *v1.StaffListReq) (res *v1.StaffListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
		m    = dao.Staff.Ctx(ctx)
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if req.Keyword != "" {
		m = m.WhereOrLike(dao.Staff.Columns().Username, "%"+req.Keyword+"%").
			WhereOrLike(dao.Staff.Columns().RealName, "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.Staff.Columns().Status, *req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.StaffListRes{
			List:  make([]*v1.StaffListItem, 0),
			Total: 0,
		}, nil
	}
	var staffList []*entity.Staff
	err = m.Page(page, size).OrderDesc(dao.Staff.Columns().Id).Scan(&staffList)
	if err != nil {
		return nil, err
	}

	staffIds := make([]int64, 0, len(staffList))
	for _, s := range staffList {
		staffIds = append(staffIds, s.Id)
	}

	type StaffRole struct {
		StaffId  int64
		RoleId   int64
		RoleName string
	}
	var staffRoles []StaffRole
	err = dao.StaffRoles.Ctx(ctx).
		InnerJoin("sys_roles", "sys_roles.id = sys_staff_roles.role_id").
		WhereIn("sys_staff_roles.staff_id", staffIds).
		Fields("sys_staff_roles.staff_id", "sys_staff_roles.role_id", "sys_roles.name AS role_name").
		Scan(&staffRoles)
	if err != nil {
		return nil, err
	}

	roleMap := make(map[int64]*v1.StaffListItem)
	for _, s := range staffList {
		roleMap[s.Id] = &v1.StaffListItem{
			Id:              s.Id,
			Username:        s.Username,
			RealName:        s.RealName,
			Email:           s.Email,
			Phone:           s.Phone,
			Avatar:          s.Avatar,
			Status:          s.Status,
			LastLoginIp:     s.LastLoginIp,
			LastLoginAt:     s.LastLoginAt,
			CreatedAt:       s.CreatedAt,
			RoleIds:         make([]int64, 0),
			RoleNames:       make([]string, 0),
			DepartmentIds:   make([]int64, 0),
			DepartmentNames: make([]string, 0),
		}
	}
	for _, sr := range staffRoles {
		if item, ok := roleMap[sr.StaffId]; ok {
			item.RoleIds = append(item.RoleIds, sr.RoleId)
			item.RoleNames = append(item.RoleNames, sr.RoleName)
		}
	}

	type StaffDept struct {
		StaffId        int64
		DepartmentId   int64
		DepartmentName string
	}
	var staffDepts []StaffDept
	err = dao.StaffDepartments.Ctx(ctx).
		InnerJoin("sys_departments", "sys_departments.id = sys_staff_departments.department_id").
		WhereIn("sys_staff_departments.staff_id", staffIds).
		Fields("sys_staff_departments.staff_id", "sys_staff_departments.department_id", "sys_departments.name AS department_name").
		Scan(&staffDepts)
	if err == nil {
		for _, sd := range staffDepts {
			if item, ok := roleMap[sd.StaffId]; ok {
				item.DepartmentIds = append(item.DepartmentIds, sd.DepartmentId)
				item.DepartmentNames = append(item.DepartmentNames, sd.DepartmentName)
			}
		}
	}

	list := make([]*v1.StaffListItem, 0, len(staffList))
	for _, s := range staffList {
		list = append(list, roleMap[s.Id])
	}
	return &v1.StaffListRes{List: list, Total: total}, nil
}

func (s *sStaff) AssignRoles(ctx context.Context, req *v1.StaffAssignRolesReq) (res *v1.StaffAssignRolesRes, err error) {
	count, err := dao.Staff.Ctx(ctx).Where(dao.Staff.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		s.logAudit(ctx, "assign_roles", "staff", fmt.Sprintf("%d", req.Id), "staff_not_found", 0, errcode.ErrUserNotFound)
		return nil, errcode.ErrUserNotFound
	}

	err = dao.StaffRoles.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.StaffRoles.Ctx(ctx).TX(tx).Where(dao.StaffRoles.Columns().StaffId, req.Id).Delete()
		if err != nil {
			return err
		}
		for _, roleId := range req.RoleIds {
			_, err = tx.Model("sys_staff_roles").Insert(do.StaffRoles{
				StaffId: req.Id,
				RoleId:  roleId,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.logAudit(ctx, "assign_roles", "staff", fmt.Sprintf("%d", req.Id), "", 0, err)
		return nil, err
	}
	s.logAudit(ctx, "assign_roles", "staff", fmt.Sprintf("%d", req.Id),
		fmt.Sprintf("role_ids=%v", req.RoleIds), 1, nil)
	return &v1.StaffAssignRolesRes{}, nil
}


