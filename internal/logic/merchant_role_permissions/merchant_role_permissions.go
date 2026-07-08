package merchantRolePermissions

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"gf-eshop/api/merchant_role_permissions/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sMerchantRolePermissions struct{}

func init() {
	service.RegisterMerchantRolePermissions(&sMerchantRolePermissions{})
}

func (s *sMerchantRolePermissions) Create(ctx context.Context, req *v1.MerchantRolePermissionsCreateReq) (res *v1.MerchantRolePermissionsCreateRes, err error) {
	_, err = dao.MerchantRolePermissions.Ctx(ctx).Insert(do.MerchantRolePermissions{
		MerchantId:     req.MerchantId,
		RoleId:         req.RoleId,
		PermissionName: req.PermissionName,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolePermissionsCreateRes{}, nil
}

func (s *sMerchantRolePermissions) Delete(ctx context.Context, req *v1.MerchantRolePermissionsDeleteReq) (res *v1.MerchantRolePermissionsDeleteRes, err error) {
	_, err = dao.MerchantRolePermissions.Ctx(ctx).
		Where(dao.MerchantRolePermissions.Columns().Id, req.Id).
		Delete()
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolePermissionsDeleteRes{}, nil
}

func (s *sMerchantRolePermissions) List(ctx context.Context, req *v1.MerchantRolePermissionsListReq) (res *v1.MerchantRolePermissionsListRes, err error) {
	var list []*entity.MerchantRolePermissions
	err = dao.MerchantRolePermissions.Ctx(ctx).
		Where(dao.MerchantRolePermissions.Columns().RoleId, req.RoleId).
		OrderAsc(dao.MerchantRolePermissions.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.MerchantRolePermissionsListRes{List: list}, nil
}

func (s *sMerchantRolePermissions) Update(ctx context.Context, req *v1.MerchantRolePermissionsUpdateReq) (res *v1.MerchantRolePermissionsUpdateRes, err error) {
	err = dao.MerchantRolePermissions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先删旧权限
		_, err := dao.MerchantRolePermissions.Ctx(ctx).TX(tx).
			Where(dao.MerchantRolePermissions.Columns().RoleId, req.RoleId).
			Delete()
		if err != nil {
			return err
		}
		// 再插新权限
		for _, name := range req.PermissionNames {
			_, err := dao.MerchantRolePermissions.Ctx(ctx).TX(tx).Insert(do.MerchantRolePermissions{
				MerchantId:     req.MerchantId,
				RoleId:         req.RoleId,
				PermissionName: name,
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
	return &v1.MerchantRolePermissionsUpdateRes{}, nil
}
